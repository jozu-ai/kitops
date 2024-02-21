package export

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"jmm/pkg/artifact"
	"jmm/pkg/lib/constants"
	"jmm/pkg/lib/filesystem"
	"jmm/pkg/lib/repo"
	"os"
	"path"
	"path/filepath"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/registry"
	"sigs.k8s.io/yaml"
)

func ExportModel(ctx context.Context, store oras.Target, ref *registry.Reference, options *ExportOptions) error {
	manifestDesc, err := store.Resolve(ctx, ref.Reference)
	if err != nil {
		return fmt.Errorf("failed to resolve local reference: %w", err)
	}
	manifest, config, err := repo.GetManifestAndConfig(ctx, store, manifestDesc)
	if err != nil {
		return fmt.Errorf("failed to read local model: %s", err)
	}

	if options.exportConf.ExportConfig {
		if err := ExportConfig(config, options.exportDir, options.overwrite); err != nil {
			return err
		}
	}

	// Since there might be multiple models, etc. we need to synchronously iterate
	// through the config's relevant field to get the correct path for exporting
	var modelIdx, codeIdx, datasetIdx int
	for _, layerDesc := range manifest.Layers {
		layerDir := ""
		switch layerDesc.MediaType {
		case constants.ModelLayerMediaType:
			if !options.exportConf.ExportModels {
				continue
			}
			modelEntry := config.Models[modelIdx]
			layerDir = filepath.Join(options.exportDir, modelEntry.Path)
			fmt.Printf("Exporting model %s to %s\n", modelEntry.Name, layerDir)
			modelIdx += 1

		case constants.CodeLayerMediaType:
			if !options.exportConf.ExportCode {
				continue
			}
			codeEntry := config.Code[codeIdx]
			layerDir = filepath.Join(options.exportDir, codeEntry.Path)
			fmt.Printf("Exporting code to %s\n", layerDir)
			codeIdx += 1

		case constants.DataSetLayerMediaType:
			if !options.exportConf.ExportDatasets {
				continue
			}
			datasetEntry := config.DataSets[datasetIdx]
			layerDir = filepath.Join(options.exportDir, datasetEntry.Path)
			fmt.Printf("Exporting dataset %s to %s\n", datasetEntry.Name, layerDir)
			datasetIdx += 1
		}
		if _, err := filesystem.VerifySubpath(options.exportDir, layerDir); err != nil {
			return err
		}
		if err := ExportLayer(ctx, store, layerDesc, layerDir, options.overwrite); err != nil {
			return err
		}
	}

	return nil
}

func ExportConfig(config *artifact.JozuFile, exportDir string, overwrite bool) error {
	configPath := path.Join(exportDir, constants.DefaultModelFileName)
	if fi, exists := filesystem.PathExists(configPath); exists {
		if !overwrite {
			return fmt.Errorf("failed to export config: path %s already exists", exportDir)
		} else if !fi.Mode().IsRegular() {
			return fmt.Errorf("failed to export config: path %s exists and is not a regular file", exportDir)
		}
	}

	configBytes, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to export config: %w", err)
	}

	fmt.Printf("Exporting config to %s\n", configPath)
	if err := os.WriteFile(configPath, configBytes, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

func ExportLayer(ctx context.Context, store content.Storage, desc ocispec.Descriptor, exportDir string, overwrite bool) error {
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

	if fi, exists := filesystem.PathExists(exportDir); exists {
		if !overwrite {
			return fmt.Errorf("failed to export: path %s already exists", exportDir)
		} else if !fi.IsDir() {
			return fmt.Errorf("failed to export: path %s exists and is not a directory", exportDir)
		}
	}
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
		outPath := path.Join(dir, header.Name)
		fmt.Printf("Extracting %s\n", outPath)

		switch header.Typeflag {
		case tar.TypeDir:
			if fi, exists := filesystem.PathExists(outPath); exists {
				if !overwrite {
					return fmt.Errorf("path '%s' already exists", outPath)
				}
				if !fi.IsDir() {
					return fmt.Errorf("path '%s' already exists and is not a directory", outPath)
				}
			}
			fmt.Printf("Creating directory %s\n", outPath)
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
			fmt.Printf("Extracting file %s\n", outPath)
			file, err := os.OpenFile(outPath, os.O_TRUNC|os.O_RDWR|os.O_EXCL, header.FileInfo().Mode())
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
