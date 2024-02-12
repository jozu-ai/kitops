/*
Copyright © 2024 Jozu.com
*/
package models

import (
	"context"
	"encoding/json"
	"fmt"
	"jmm/pkg/artifact"
	"jmm/pkg/lib/storage"
	"math"
	"strings"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

const (
	ModelsTableHeader = "REPOSITORY\tTAG\tMAINTAINER\tMODEL FORMAT\tSIZE\tDIGEST"
	ModelsTableFmt    = "%s\t%s\t%s\t%s\t%s\t%s\t"
)

func listModels(store storage.Store) ([]string, error) {
	index, err := store.ParseIndexJson()
	if err != nil {
		return nil, err
	}

	var infolines []string
	for _, manifestDesc := range index.Manifests {
		manifest, err := getManifest(store, manifestDesc)
		if err != nil {
			return nil, err
		}
		manifestConf, err := readManifestConfig(store, manifest)
		if err != nil {
			return nil, err
		}
		// TODO: filter list for our manifests only, ignore other artifacts
		infoline, err := getManifestInfoLine(store.GetRepository(), manifestDesc, manifest, manifestConf)
		if err != nil {
			return nil, err
		}
		infolines = append(infolines, infoline)
	}

	return infolines, nil
}

func getManifest(store storage.Store, manifestDesc ocispec.Descriptor) (*ocispec.Manifest, error) {
	manifestBytes, err := store.Fetch(context.Background(), manifestDesc)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest %s: %w", manifestDesc.Digest, err)
	}
	manifest := &ocispec.Manifest{}
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest %s: %w", manifestDesc.Digest, err)
	}
	return manifest, nil
}

func readManifestConfig(store storage.Store, manifest *ocispec.Manifest) (*artifact.JozuFile, error) {
	configBytes, err := store.Fetch(context.Background(), manifest.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	config := &artifact.JozuFile{}
	if err := json.Unmarshal(configBytes, config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return config, nil
}

func getManifestInfoLine(repo string, desc ocispec.Descriptor, manifest *ocispec.Manifest, config *artifact.JozuFile) (string, error) {
	ref := desc.Annotations[ocispec.AnnotationRefName]
	if ref == "" {
		ref = "<none>"
	}

	// Strip localhost from repo if present, since we added it
	repo = strings.TrimPrefix(repo, "localhost/")
	if repo == "" {
		repo = "<none>"
	}

	var size int64
	for _, layer := range manifest.Layers {
		size += layer.Size
	}
	sizeStr := formatBytes(size)

	info := fmt.Sprintf(ModelsTableFmt, repo, ref, config.Maintainer, config.ModelFormat, sizeStr, desc.Digest)
	return info, nil
}

func formatBytes(i int64) string {
	if i == 0 {
		return "0 B"
	}

	if i < 1024 {
		// Catch bytes to avoid printing fractional amounts of bytes e.g. 123.0 bytes
		return fmt.Sprintf("%d B", i)
	}

	suffixes := []string{"KiB", "MiB", "GiB", "TiB"}
	unit := float64(1024)

	size := float64(i) / unit
	for _, suffix := range suffixes {
		if size < unit {
			// Round down to the nearest tenth of a unit to avoid e.g. 1MiB - 1B = 1024KiB
			niceSize := math.Floor(size*10) / 10
			return fmt.Sprintf("%.1f %s", niceSize, suffix)
		}
		size = size / unit
	}

	// Fall back to printing whatever's left as PiB
	return fmt.Sprintf("%.1f PiB", size)
}
