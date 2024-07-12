/*
Copyright © 2024 Jozu.com
*/
package list

import (
	"context"
	"errors"
	"sort"

	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/repo/local"
	"kitops/pkg/lib/repo/util"
)

func listLocalKits(ctx context.Context, opts *listOptions) ([]string, error) {
	storageRoot := constants.StoragePath(opts.configHome)

	localRepos, err := local.GetAllLocalRepos(storageRoot)
	if err != nil {
		return nil, err
	}
	var allInfoLines []string
	for _, repo := range localRepos {
		infolines, err := readInfoFromRepo(ctx, repo)
		if err != nil {
			return nil, err
		}
		allInfoLines = append(allInfoLines, infolines...)
	}

	return allInfoLines, nil
}

func readInfoFromRepo(ctx context.Context, repo local.LocalRepo) ([]string, error) {
	var infolines []string
	manifestDescs := repo.GetAllModels()
	for _, manifestDesc := range manifestDescs {
		manifest, config, err := util.GetManifestAndConfig(ctx, repo, manifestDesc)
		if err != nil && !errors.Is(err, util.ErrNotAModelKit) {
			return nil, err
		}
		tags := repo.GetTags(manifestDesc)
		// Strip localhost from repo if present, since we added it
		repository := util.FormatRepositoryForDisplay(repo.GetRepoName())
		if repository == "" {
			repository = "<none>"
		}
		info := modelInfo{
			repo:   repository,
			digest: string(manifestDesc.Digest),
			tags:   tags,
		}
		info.fill(manifest, config)

		infolines = append(infolines, info.format()...)
	}

	sort.Strings(infolines)
	return infolines, nil
}
