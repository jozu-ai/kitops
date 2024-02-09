/*
Copyright © 2024 Jozu.com
*/
package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"jmm/pkg/artifact"
	"os"
	"text/tabwriter"

	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

const (
	ModelsTableHeader = "DIGEST\tMAINTAINER\tMODEL FORMAT\tSIZE"
	ModelsTableFmt    = "%s\t%s\t%s\t%s\t"
)

func listModels(opts *ModelsOptions) error {
	store := artifact.NewArtifactStore(opts.configHome)
	index, err := store.ParseIndexJson()
	if err != nil {
		return err
	}

	manifests, err := manifestsFromIndex(index, store)
	if err != nil {
		return err
	}

	if err := printManifestsSummary(manifests, store); err != nil {
		return err
	}

	return nil
}

func manifestsFromIndex(index *ocispec.Index, store *artifact.Store) (map[digest.Digest]ocispec.Manifest, error) {
	manifests := map[digest.Digest]ocispec.Manifest{}
	for _, manifestDesc := range index.Manifests {
		manifestReader, err := store.Storage.Fetch(context.Background(), manifestDesc)
		if err != nil {
			return nil, fmt.Errorf("failed to get manifest %s: %w", manifestDesc.Digest, err)
		}
		manifestBytes, err := io.ReadAll(manifestReader)
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

func readManifestConfig(manifest *ocispec.Manifest, store *artifact.Store) (*artifact.JozuFile, error) {
	configReader, err := store.Storage.Fetch(context.Background(), manifest.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}
	configBytes, err := io.ReadAll(configReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	config := &artifact.JozuFile{}
	if err := json.Unmarshal(configBytes, config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return config, nil
}

func printManifestsSummary(manifests map[digest.Digest]ocispec.Manifest, store *artifact.Store) error {
	tw := tabwriter.NewWriter(os.Stdout, 0, 2, 3, ' ', 0)
	fmt.Fprintln(tw, ModelsTableHeader)
	for digest, manifest := range manifests {
		// TODO: filter this list for manifests we're interested in (build needs to set a manifest mediaType/artifactType)
		line, err := getManifestInfoLine(digest, &manifest, store)
		if err != nil {
			return err
		}
		fmt.Fprintln(tw, line)
	}
	tw.Flush()
	return nil
}

func getManifestInfoLine(digest digest.Digest, manifest *ocispec.Manifest, store *artifact.Store) (string, error) {
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

	suffixes := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"}
	unit := float64(1024)

	size := float64(i)
	for _, suffix := range suffixes {
		if size < unit {
			return fmt.Sprintf("%.1f %s", size, suffix)
		}
		size = size / unit
	}

	// Fall back to printing 1000's of PiB
	return fmt.Sprintf("%.1f PiB", size)
}
