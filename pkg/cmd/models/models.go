/*
Copyright © 2024 Jozu.com
*/
package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"jmm/pkg/artifact"
	"jmm/pkg/lib/storage"
	"math"
	"text/tabwriter"

	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

const (
	ModelsTableHeader = "DIGEST\tMAINTAINER\tMODEL FORMAT\tSIZE"
	ModelsTableFmt    = "%s\t%s\t%s\t%s\t"
)

func listModels(store storage.Store) (string, error) {
	index, err := store.ParseIndexJson()
	if err != nil {
		return "", err
	}

	manifests, err := manifestsFromIndex(index, store)
	if err != nil {
		return "", err
	}

	summary, err := printManifestsSummary(manifests, store)
	if err != nil {
		return "", err
	}

	return summary, nil
}

func manifestsFromIndex(index *ocispec.Index, store storage.Store) (map[digest.Digest]ocispec.Manifest, error) {
	manifests := map[digest.Digest]ocispec.Manifest{}
	for _, manifestDesc := range index.Manifests {
		manifestBytes, err := store.Fetch(context.Background(), manifestDesc)
		if err != nil {
			return nil, fmt.Errorf("failed to read manifest %s: %w", manifestDesc.Digest, err)
		}
		manifest := ocispec.Manifest{}
		if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
			return nil, fmt.Errorf("failed to parse manifest %s: %w", manifestDesc.Digest, err)
		}
		manifests[manifestDesc.Digest] = manifest
	}
	return manifests, nil
}

func readManifestConfig(manifest *ocispec.Manifest, store storage.Store) (*artifact.JozuFile, error) {
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

func printManifestsSummary(manifests map[digest.Digest]ocispec.Manifest, store storage.Store) (string, error) {
	buf := bytes.Buffer{}
	tw := tabwriter.NewWriter(&buf, 0, 2, 3, ' ', 0)
	fmt.Fprintln(tw, ModelsTableHeader)
	for digest, manifest := range manifests {
		// TODO: filter this list for manifests we're interested in (build needs to set a manifest mediaType/artifactType)
		line, err := getManifestInfoLine(digest, &manifest, store)
		if err != nil {
			return "", err
		}
		fmt.Fprintln(tw, line)
	}
	tw.Flush()
	return buf.String(), nil
}

func getManifestInfoLine(digest digest.Digest, manifest *ocispec.Manifest, store storage.Store) (string, error) {
	config, err := readManifestConfig(manifest, store)
	if err != nil {
		return "", err
	}
	var size int64
	for _, layer := range manifest.Layers {
		size += layer.Size
	}
	sizeStr := formatBytes(size)

	info := fmt.Sprintf(ModelsTableFmt, digest, config.Maintainer, config.ModelFormat, sizeStr)
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
