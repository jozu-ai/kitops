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

	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/registry"
)

const (
	DefaultRegistry   = "localhost"
	DefaultRepository = "_"
)

var (
	validTagRegex = regexp.MustCompile(`^[a-zA-Z0-9_][a-zA-Z0-9._-]{0,127}$`)
)

// ParseReference parses a reference string into a Reference struct. It attempts to make
// references conform to an expected structure, with a defined registry and repository by filling
// default values for registry and repository where appropriate. Where the first part of a reference
// doesn't look like a registry URL, the default registry is used, turning e.g. testorg/testrepo into
// localhost/testorg/testrepo.
func ParseReference(refString string) (ref *registry.Reference, extraTags []string, err error) {
	// Check if provided input is a plain digest
	if _, err := digest.Parse(refString); err == nil {
		ref := &registry.Reference{
			Registry:   DefaultRegistry,
			Repository: DefaultRepository,
			Reference:  refString,
		}
		return ref, []string{}, nil
	}

	// Handle registry, which may or may not be specified; if unspecified, use a default value for registry
	refParts := strings.Split(refString, "/")
	if len(refParts) == 1 {
		// Just a repo, need to add default registry
		refString = fmt.Sprintf("%s/%s", DefaultRegistry, refString)
	} else {
		// Check if registry part "looks" like a URL; we're trying to distinguish between cases:
		// a) testorg/testrepo --> should be localhost/testorg/testrepo
		// b) registry.io/testrepo --> should be registry.io/testrepo
		// c) localhost:5000/testrepo --> should be localhost:5000/testrepo
		registry := refParts[0]
		if !strings.Contains(registry, ":") && !strings.Contains(registry, ".") {
			refString = fmt.Sprintf("%s/%s", DefaultRegistry, refString)
		}
	}

	// Split off extra tags (e.g. repo:tag1,tag2,tag3)
	refAndTags := strings.Split(refString, ",")
	baseRef, err := registry.ParseReference(refAndTags[0])
	if err != nil {
		return nil, nil, err
	}
	return &baseRef, refAndTags[1:], nil
}

// DefaultReference returns a reference that can be used when no reference is supplied. It uses
// the default registry and repository
func DefaultReference() *registry.Reference {
	return &registry.Reference{
		Registry:   DefaultRegistry,
		Repository: DefaultRepository,
	}
}

// FormatRepositoryForDisplay removes default values from a repository string to avoid surfacing defaulted fields
// when displaying references, which may be confusing.
func FormatRepositoryForDisplay(repo string) string {
	repo = strings.TrimPrefix(repo, DefaultRegistry+"/")
	repo = strings.TrimPrefix(repo, DefaultRepository)
	return repo
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
