package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"kitops/pkg/artifact"
	"kitops/pkg/lib/constants"
	"path/filepath"
	"regexp"
	"strings"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/registry"
)

var (
	validTagRegex = regexp.MustCompile(`^[a-zA-Z0-9_][a-zA-Z0-9._-]{0,127}$`)
)

// ParseReference parses a reference string into a Reference struct. If the
// reference does not include a registry (e.g. myrepo:mytag), the placeholder
// 'localhost' is used. Additional tags can be specified in a comma-separated
// list (e.g. myrepo:tag1,tag2,tag3)
func ParseReference(refString string) (ref *registry.Reference, extraTags []string, err error) {
	// References _must_ contain host; use localhost to mark local-only
	if !strings.Contains(refString, "/") {
		refString = fmt.Sprintf("localhost/%s", refString)
	}

	refAndTags := strings.Split(refString, ",")
	baseRef, err := registry.ParseReference(refAndTags[0])
	if err != nil {
		return nil, nil, err
	}
	return &baseRef, refAndTags[1:], nil
}

func RepoPath(storagePath string, ref *registry.Reference) string {
	return filepath.Join(storagePath, ref.Registry, ref.Repository)
}

func GetManifestAndConfig(ctx context.Context, store content.Storage, manifestDesc ocispec.Descriptor) (*ocispec.Manifest, *artifact.KitFile, error) {
	manifest, err := GetManifest(ctx, store, manifestDesc)
	if err != nil {
		return nil, nil, err
	}
	config, err := GetConfig(ctx, store, manifest.Config)
	if err != nil {
		return nil, nil, err
	}
	return manifest, config, nil
}

func GetManifest(ctx context.Context, store content.Storage, manifestDesc ocispec.Descriptor) (*ocispec.Manifest, error) {
	manifestBytes, err := content.FetchAll(ctx, store, manifestDesc)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest %s: %w", manifestDesc.Digest, err)
	}
	manifest := &ocispec.Manifest{}
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest %s: %w", manifestDesc.Digest, err)
	}
	if manifest.Config.MediaType != constants.ModelConfigMediaType {
		return nil, fmt.Errorf("reference exists but is not a model")
	}

	return manifest, nil
}

func GetConfig(ctx context.Context, store content.Storage, configDesc ocispec.Descriptor) (*artifact.KitFile, error) {
	configBytes, err := content.FetchAll(ctx, store, configDesc)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	config := &artifact.KitFile{}
	if err := json.Unmarshal(configBytes, config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return config, nil
}

func GetTagsForDescriptor(ctx context.Context, store LocalStorage, desc ocispec.Descriptor) ([]string, error) {
	index, err := store.GetIndex()
	if err != nil {
		return nil, err
	}
	var tags []string
	for _, manifest := range index.Manifests {
		if manifest.Digest == desc.Digest {
			tags = append(tags, manifest.Annotations[ocispec.AnnotationRefName])
		}
	}
	return tags, nil
}

func ValidateTag(tag string) error {
	if !validTagRegex.MatchString(tag) {
		return fmt.Errorf("invalid tag")
	}
	return nil
}
