package export

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"kitops/pkg/artifact"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/filesystem"
	"kitops/pkg/lib/repo"
	"kitops/pkg/output"
	"os"
	"path/filepath"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/registry"
	"sigs.k8s.io/yaml"
)

func exportModel(ctx context.Context, store oras.Target, ref *registry.Reference, options *exportOptions) error {
	manifestDesc, err := store.Resolve(ctx, ref.Reference)
	if err != nil {
		return fmt.Errorf("failed to resolve local reference: %w", err)
	}
	manifest, config, err := repo.GetManifestAndConfig(ctx, store, manifestDesc)
	if err != nil {
		return fmt.Errorf("failed to read local model: %s", err)
	}

	if options.exportConf.exportConfig {
		if err := exportConfig(config, options.exportDir, options.overwrite); err != nil {
			return err
		}
	}

	// Since there might be multiple models, etc. we need to synchronously iterate
	// through the config's relevant field to get the correct path for exporting
	var codeIdx, datasetIdx int
	for _, layerDesc := range manifest.Layers {
		layerDir := ""
		switch layerDesc.MediaType {
		case constants.ModelLayerMediaType:
			if !options.exportConf.exportModels {
				continue
			}
			layerDir = filepath.Join(options.exportDir, config.Model.Path)
			output.Infof("Exporting model to %s", layerDir)

		case constants.CodeLayerMediaType:
			if !options.exportConf.exportCode {
				continue
			}
			codeEntry := config.Code[codeIdx]
			layerDir = filepath.Join(options.exportDir, codeEntry.Path)
			output.Infof("Exporting code to %s", layerDir)
			codeIdx += 1

		case constants.DataSetLayerMediaType:
			if !options.exportConf.exportDatasets {
				continue
			}
			datasetEntry := config.DataSets[datasetIdx]
			layerDir = filepath.Join(options.exportDir, datasetEntry.Path)
			output.Infof("Exporting dataset %s to %s", datasetEntry.Name, layerDir)
			datasetIdx += 1
		}
		if err := exportLayer(ctx, store, layerDesc, layerDir, options.overwrite); err != nil {
			return err
		}
	}
	output.Debugf("Exported %d code layers", codeIdx)
	output.Debugf("Exported %d dataset layers", datasetIdx)

	return nil
}

func exportConfig(config *artifact.KitFile, exportDir string, overwrite bool) error {
	configPath := filepath.Join(exportDir, constants.DefaultKitFileName)
	if fi, exists := filesystem.PathExists(configPath); exists {
		if !overwrite {
			return fmt.Errorf("failed to export config: path %s already exists", configPath)
		} else if !fi.Mode().IsRegular() {
			return fmt.Errorf("failed to export config: path %s exists and is not a regular file", configPath)
		}
	}

	configBytes, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to export config: %w", err)
	}

	output.Infof("Exporting config to %s", configPath)
	if err := os.WriteFile(configPath, configBytes, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

func exportLayer(ctx context.Context, store content.Storage, desc ocispec.Descriptor, exportPath string, overwrite bool) error {
	rc, err := store.Fetch(ctx, desc)
	if err != nil {
		return fmt.Errorf("failed get layer %s: %w", desc.Digest, err)
	}
	defer rc.Close()

	gzr, err := gzip.NewReader(rc)
	if err != nil {
		return fmt.Errorf("error extracting gzipped file: %w", err)
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)

	if _, exists := filesystem.PathExists(exportPath); exists {
		if !overwrite {
			return fmt.Errorf("failed to export: path %s already exists", exportPath)
		}
		output.Debugf("Directory %s already exists", exportPath)
	}
	exportDir := filepath.Dir(exportPath)
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", exportDir, err)
	}

	return extractTar(tr, exportDir, overwrite)
}

func extractTar(tr *tar.Reader, dir string, overwrite bool) error {
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		outPath := filepath.Join(dir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if fi, exists := filesystem.PathExists(outPath); exists {
				if !overwrite {
					return fmt.Errorf("path '%s' already exists", outPath)
				}
				if !fi.IsDir() {
					return fmt.Errorf("path '%s' already exists and is not a directory", outPath)
				}
				output.Debugf("Path %s already exists", outPath)
			}
			output.Debugf("Creating directory %s", outPath)
			if err := os.MkdirAll(outPath, header.FileInfo().Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", outPath, err)
			}

		case tar.TypeReg:
			if fi, exists := filesystem.PathExists(outPath); exists {
				if !overwrite {
					return fmt.Errorf("path '%s' already exists", outPath)
				}
				if !fi.Mode().IsRegular() {
					return fmt.Errorf("path '%s' already exists and is not a regular file", outPath)
				}
			}
			output.Debugf("Extracting file %s", outPath)
			file, err := os.OpenFile(outPath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, header.FileInfo().Mode())
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", outPath, err)
			}
			defer file.Close()

			written, err := io.Copy(file, tr)
			if err != nil {
				return fmt.Errorf("failed to write file %s: %w", outPath, err)
			}
			if written != header.Size {
				return fmt.Errorf("could not extract file %s", outPath)
			}

		default:
			return fmt.Errorf("Unrecognized type in archive: %s", header.Name)
		}
	}
	return nil
}
