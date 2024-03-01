/*
Copyright © 2024 Jozu.com
*/
package list

import (
	"context"
	"fmt"
	"kitops/pkg/artifact"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/repo"
	"math"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

const (
	listTableHeader = "REPOSITORY\tTAG\tMAINTAINER\tNAME\tSIZE\tDIGEST"
	listTableFmt    = "%s\t%s\t%s\t%s\t%s\t%s\t"
)

func listLocalKits(ctx context.Context, opts *listOptions) ([]string, error) {
	storageRoot := constants.StoragePath(opts.configHome)
	stores, err := repo.GetAllLocalStores(storageRoot)
	if err != nil {
		return nil, err
	}

	var allInfoLines []string
	for _, store := range stores {
		infolines, err := listKits(ctx, store)
		if err != nil {
			return nil, err
		}
		allInfoLines = append(allInfoLines, infolines...)
	}
	return allInfoLines, nil
}

func listKits(ctx context.Context, store repo.LocalStorage) ([]string, error) {
	index, err := store.GetIndex()
	if err != nil {
		return nil, err
	}

	var infolines []string
	for _, manifestDesc := range index.Manifests {
		manifest, config, err := repo.GetManifestAndConfig(ctx, store, manifestDesc)
		if err != nil {
			return nil, err
		}
		infoline := getManifestInfoLine(store.GetRepo(), manifestDesc, manifest, config)
		infolines = append(infolines, infoline)
	}

	return infolines, nil
}

func getManifestInfoLine(repository string, desc ocispec.Descriptor, manifest *ocispec.Manifest, config *artifact.KitFile) string {
	ref := desc.Annotations[ocispec.AnnotationRefName]
	if ref == "" {
		ref = "<none>"
	}

	// Strip localhost from repo if present, since we added it
	repository = repo.FormatRepositoryForDisplay(repository)
	if repository == "" {
		repository = "<none>"
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

	info := fmt.Sprintf(listTableFmt, repository, ref, author, config.Kit.Name, sizeStr, desc.Digest)
	return info
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
