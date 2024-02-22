package list

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kitops/pkg/artifact"
	"kitops/pkg/lib/constants"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
)

func listRemoteKits(ctx context.Context, remoteRef *registry.Reference, useHttp bool) ([]string, error) {
	remoteRegistry, err := remote.NewRegistry(remoteRef.Registry)
	if err != nil {
		return nil, fmt.Errorf("failed to read registry: %w", err)
	}
	remoteRegistry.PlainHTTP = useHttp

	repo, err := remoteRegistry.Repository(ctx, remoteRef.Repository)
	if err != nil {
		return nil, fmt.Errorf("failed to read repository: %w", err)
	}
	if remoteRef.Reference != "" {
		return listImageTag(ctx, repo, remoteRef)
	}
	return listTags(ctx, repo, remoteRef)
}

func listTags(ctx context.Context, repo registry.Repository, ref *registry.Reference) ([]string, error) {
	var tags []string
	err := repo.Tags(ctx, "", func(tagsPage []string) error {
		tags = append(tags, tagsPage...)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list tags on repostory: %w", err)
	}

	var allLines []string
	for _, tag := range tags {
		tagRef := &registry.Reference{
			Registry:   ref.Registry,
			Repository: ref.Repository,
			Reference:  tag,
		}
		infoLines, err := listImageTag(ctx, repo, tagRef)
		if err != nil {
			return nil, err
		}
		allLines = append(allLines, infoLines...)
	}

	return allLines, nil
}

func listImageTag(ctx context.Context, repo registry.Repository, ref *registry.Reference) ([]string, error) {
	manifestDesc, manifestReader, err := repo.FetchReference(ctx, ref.Reference)
	if err != nil {
		return nil, fmt.Errorf("failed to read reference: %w", err)
	}
	manifestBytes, err := io.ReadAll(manifestReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}
	manifest := &ocispec.Manifest{}
	if err := json.Unmarshal(manifestBytes, manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}
	if manifest.Config.MediaType != constants.ModelConfigMediaType {
		return nil, nil
	}

	configReader, err := repo.Fetch(ctx, manifest.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to read config reference: %w", err)
	}
	configBytes, err := io.ReadAll(configReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	config := &artifact.KitFile{}
	if err := json.Unmarshal(configBytes, config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Manifest descriptor may not have annotation for tag, add it here for safety
	if _, ok := manifestDesc.Annotations[ocispec.AnnotationRefName]; !ok {
		if manifestDesc.Annotations == nil {
			manifestDesc.Annotations = map[string]string{}
		}
		manifestDesc.Annotations[ocispec.AnnotationRefName] = ref.Reference
	}

	info := getManifestInfoLine(ref.Repository, manifestDesc, manifest, config)
	return []string{info}, nil
}
