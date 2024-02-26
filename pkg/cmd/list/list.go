/*
Copyright © 2024 Jozu.com
*/
package list

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"kitops/pkg/artifact"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/storage"
	"math"
	"os"
	"path/filepath"
	"strings"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

const (
	listTableHeader = "REPOSITORY\tTAG\tMAINTAINER\tNAME\tSIZE\tDIGEST"
	listTableFmt    = "%s\t%s\t%s\t%s\t%s\t%s\t"
)

func listLocalKits(ctx context.Context, opts *listOptions) ([]string, error) {
	storageRoot := constants.StoragePath(opts.configHome)
	storeDirs, err := findRepos(storageRoot)
	if err != nil {
		return nil, err
	}

	var allInfoLines []string
	for _, storeDir := range storeDirs {
		store := storage.NewLocalStore(storageRoot, storeDir)

		infolines, err := listKits(ctx, store)
		if err != nil {
			return nil, err
		}
		allInfoLines = append(allInfoLines, infolines...)
	}
	return allInfoLines, nil
}

func listKits(ctx context.Context, store storage.Store) ([]string, error) {
	index, err := store.ParseIndexJson()
	if err != nil {
		return nil, err
	}

	var infolines []string
	for _, manifestDesc := range index.Manifests {
		manifest, err := getManifest(ctx, store, manifestDesc)
		if err != nil {
			return nil, err
		}
		if manifest.Config.MediaType != constants.ModelConfigMediaType {
			continue
		}
		manifestConf, err := readManifestConfig(ctx, store, manifest)
		if err != nil {
			return nil, err
		}
		infoline := getManifestInfoLine(store.GetRepository(), manifestDesc, manifest, manifestConf)
		infolines = append(infolines, infoline)
	}

	return infolines, nil
}

func getManifest(ctx context.Context, store storage.Store, manifestDesc ocispec.Descriptor) (*ocispec.Manifest, error) {
	manifestBytes, err := store.Fetch(ctx, manifestDesc)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest %s: %w", manifestDesc.Digest, err)
	}
	manifest := &ocispec.Manifest{}
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest %s: %w", manifestDesc.Digest, err)
	}
	return manifest, nil
}

func readManifestConfig(ctx context.Context, store storage.Store, manifest *ocispec.Manifest) (*artifact.KitFile, error) {
	configBytes, err := store.Fetch(ctx, manifest.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	config := &artifact.KitFile{}
	if err := json.Unmarshal(configBytes, config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return config, nil
}

func getManifestInfoLine(repo string, desc ocispec.Descriptor, manifest *ocispec.Manifest, config *artifact.KitFile) string {
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
	var author string
	if len(config.Kit.Authors) > 0 {
		author = config.Kit.Authors[0]
	} else {
		author = "<none>"
	}

	info := fmt.Sprintf(listTableFmt, repo, ref, author, config.Kit.Name, sizeStr, desc.Digest)
	return info
}

func findRepos(storePath string) ([]string, error) {
	var indexPaths []string
	err := filepath.WalkDir(storePath, func(file string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "index.json" && !info.IsDir() {
			dir := filepath.Dir(file)
			relDir, err := filepath.Rel(storePath, dir)
			if err != nil {
				return err
			}
			if relDir == "." {
				relDir = ""
			}
			indexPaths = append(indexPaths, relDir)
		}
		return nil
	})
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read local storage: %w", err)
	}
	return indexPaths, nil
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
