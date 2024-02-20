package export

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"jmm/pkg/artifact"
	"jmm/pkg/lib/constants"
	"jmm/pkg/lib/repo"
	"os"
	"path"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/content/oci"
	"oras.land/oras-go/v2/registry"
	"sigs.k8s.io/yaml"
)

func ExportModel(ctx context.Context, localStore *oci.Store, ref *registry.Reference, exportDir string, conf ExportConf) error {
	manifestDesc, err := localStore.Resolve(ctx, ref.Reference)
	if err != nil {
		return fmt.Errorf("failed to resolve local reference: %w", err)
	}
	manifest, config, err := repo.GetManifestAndConfig(ctx, localStore, manifestDesc)
	if err != nil {
		return fmt.Errorf("failed to read local model: %s", err)
	}

	if conf.ExportConfig {
		if err := ExportConfig(config, exportDir); err != nil {
			return err
		}
	}

	// Since there might be multiple models, etc. we need to synchronously iterate
	// through the config's relevant field to get the correct path for exporting
	var modelIdx, codeIdx, datasetIdx int
	for _, layerDesc := range manifest.Layers {
		var layerExportErr error
		switch layerDesc.MediaType {
		case constants.ModelLayerMediaType:
			if !conf.ExportModels {
				continue
			}
			modelEntry := config.Models[modelIdx]
			layerDir := path.Join(exportDir, modelEntry.Path)
			fmt.Printf("Exporting model %s to %s\n", modelEntry.Name, layerDir)
			layerExportErr = ExportLayer(ctx, localStore, layerDesc, layerDir)
			modelIdx += 1

		case constants.CodeLayerMediaType:
			if !conf.ExportCode {
				continue
			}
			codeEntry := config.Code[codeIdx]
			layerDir := path.Join(exportDir, codeEntry.Path)
			fmt.Printf("Exporting code to %s\n", layerDir)
			layerExportErr = ExportLayer(ctx, localStore, layerDesc, layerDir)
			codeIdx += 1

		case constants.DataSetLayerMediaType:
			if !conf.ExportDatasets {
				continue
			}
			datasetEntry := config.DataSets[datasetIdx]
			layerDir := path.Join(exportDir, datasetEntry.Path)
			fmt.Printf("Exporting dataset %s to %s\n", datasetEntry.Name, layerDir)
			layerExportErr = ExportLayer(ctx, localStore, layerDesc, layerDir)
			datasetIdx += 1
		}
		if layerExportErr != nil {
			return layerExportErr
		}
	}

	return nil
}

func ExportConfig(config *artifact.JozuFile, exportDir string) error {
	configPath := path.Join(exportDir, constants.DefaultModelFileName)
	if pathExists(configPath) {
		return fmt.Errorf("failed to export config: path %s already exists", configPath)
	}

	configBytes, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to export config: %w", err)
	}

	fmt.Printf("Exporting config to %s\n", configPath)
	if err := os.WriteFile(configPath, configBytes, 0644); err != nil {
		return fmt.Errorf("failed to export config file: %w", err)
	}
	return nil
}

func ExportLayer(ctx context.Context, localStore content.Storage, desc ocispec.Descriptor, exportDir string) error {
	rc, err := localStore.Fetch(ctx, desc)
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

	if pathExists(exportDir) {
		return fmt.Errorf("failed to export: path %s already exists", exportDir)
	}
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", exportDir, err)
	}

	return extractTar(tr, exportDir)
}

func extractTar(tr *tar.Reader, dir string) error {
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
			if pathExists(outPath) {
				return fmt.Errorf("directory '%s' already exists", outPath)
			}
			fmt.Printf("Creating directory %s\n", outPath)
			if err := os.MkdirAll(outPath, header.FileInfo().Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", outPath, err)
			}

		case tar.TypeReg:
			fmt.Printf("Extracting file %s\n", outPath)
			file, err := os.OpenFile(outPath, os.O_CREATE|os.O_RDWR|os.O_EXCL, header.FileInfo().Mode())
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

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
