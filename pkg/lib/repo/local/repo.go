package local

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"kitops/pkg/lib/constants"

	"github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/content/oci"
	"oras.land/oras-go/v2/errdef"
	"oras.land/oras-go/v2/registry"
)

type LocalRepo interface {
	GetRepoName() string
	GetIndex() *ocispec.Index
	BlobPath(desc ocispec.Descriptor) string
	oras.Target
	content.Deleter
	content.Untagger
}

type localRepo struct {
	storagePath string
	nameRef     string
	localIndex  *localIndex
	*oci.Store
}

func NewLocalRepo(storagePath string, ref *registry.Reference) (LocalRepo, error) {
	repo := &localRepo{}
	repo.storagePath = storagePath
	repo.nameRef = path.Join(ref.Registry, ref.Repository)

	store, err := oci.New(storagePath)
	if err != nil {
		return nil, err
	}
	repo.Store = store

	// Initialize repo-specific index.json
	indexPath := constants.IndexJsonPathForRepo(storagePath, repo.nameRef)
	localIndex, err := parseIndex(indexPath)
	if err != nil {
		return nil, err
	}
	repo.localIndex = localIndex

	return repo, nil
}

func (r *localRepo) GetIndex() *ocispec.Index {
	return &r.localIndex.Index
}

// GetRepo returns the registry and repository for the current OCI store.
func (r *localRepo) GetRepoName() string {
	return r.nameRef
}

func (r *localRepo) BlobPath(desc ocispec.Descriptor) string {
	return filepath.Join(r.storagePath, "blobs", "sha256", desc.Digest.Encoded())
}

func (l *localRepo) Delete(ctx context.Context, target ocispec.Descriptor) error {
	err := l.Store.Delete(ctx, target)
	if err != nil {
		return err
	}
	if target.MediaType == ocispec.MediaTypeImageManifest {
		return l.localIndex.delete(ctx, target)
	}
	return nil
}

func (l *localRepo) Exists(ctx context.Context, target ocispec.Descriptor) (bool, error) {
	if target.MediaType == ocispec.MediaTypeImageManifest {
		return l.localIndex.exists(ctx, target)
	} else {
		return l.Store.Exists(ctx, target)
	}
}

func (l *localRepo) Fetch(ctx context.Context, target ocispec.Descriptor) (io.ReadCloser, error) {
	if target.MediaType == ocispec.MediaTypeImageManifest {
		if exists, err := l.localIndex.exists(ctx, target); err != nil {
			return nil, err
		} else if !exists {
			return nil, errdef.ErrNotFound
		}
	}
	return l.Store.Fetch(ctx, target)
}

func (l *localRepo) Push(ctx context.Context, expected ocispec.Descriptor, content io.Reader) error {
	if expected.MediaType == ocispec.MediaTypeImageManifest {
		// Attempting to push a manifest to oci.Store will return an error if it already exists.
		// Normally, clients check before pushing, but in our case, the manifest may exist in the
		// oci.Store but not the local index. As a result, we have to check if it exists before pushing.
		exists, err := l.Store.Exists(ctx, expected)
		if err != nil {
			return err
		}
		if !exists {
			if err := l.Store.Push(ctx, expected, content); err != nil {
				return err
			}
		}
		return l.localIndex.addManifest(expected)
	} else {
		return l.Store.Push(ctx, expected, content)
	}
}

func (l *localRepo) Resolve(ctx context.Context, reference string) (ocispec.Descriptor, error) {
	return l.localIndex.resolve(ctx, reference)
}

func (l *localRepo) Tag(ctx context.Context, desc ocispec.Descriptor, reference string) error {
	// TODO: should we tag it in the general index.json too?
	return l.localIndex.tag(ctx, desc, reference)
}

func (l *localRepo) Untag(ctx context.Context, reference string) error {
	return l.localIndex.untag(ctx, reference)
}

var _ LocalRepo = (*localRepo)(nil)

type localIndex struct {
	indexPath string
	ocispec.Index
}

func (li *localIndex) addManifest(manifestDesc ocispec.Descriptor) error {
	// TODO: consider using ORAS' tag resolver to make this a little cleaner
	curTag := manifestDesc.Annotations[ocispec.AnnotationRefName]
	for _, m := range li.Manifests {
		manifestTag := m.Annotations[ocispec.AnnotationRefName]
		if m.Digest == manifestDesc.Digest && manifestTag == curTag {
			// Already included
			return nil
		}
	}
	li.Manifests = append(li.Manifests, manifestDesc)
	return li.save()
}

func (li *localIndex) save() error {
	indexJson, err := json.Marshal(li.Index)
	if err != nil {
		return fmt.Errorf("failed to marshal index: %w", err)
	}
	return os.WriteFile(li.indexPath, indexJson, 0666)
}

func (li *localIndex) exists(_ context.Context, target ocispec.Descriptor) (bool, error) {
	for _, manifestDesc := range li.Manifests {
		if manifestDesc.Digest == target.Digest {
			return true, nil
		}
	}
	return false, nil
}

func (li *localIndex) delete(_ context.Context, target ocispec.Descriptor) error {
	var newManifests []ocispec.Descriptor
	for _, manifestDesc := range li.Manifests {
		if manifestDesc.Digest != target.Digest {
			newManifests = append(newManifests, manifestDesc)
		}
	}
	li.Manifests = newManifests
	return li.save()
}

func (li *localIndex) resolve(_ context.Context, reference string) (ocispec.Descriptor, error) {
	for _, manifestDesc := range li.Manifests {
		if manifestDesc.Annotations[ocispec.AnnotationRefName] == reference {
			return manifestDesc, nil
		}
	}
	return ocispec.DescriptorEmptyJSON, errdef.ErrNotFound
}

func (li *localIndex) tag(_ context.Context, desc ocispec.Descriptor, reference string) error {
	// TODO: should probably de-duplicate this (don't store a manifest without a tag)
	descExists := false
	for _, m := range li.Manifests {
		tag := m.Annotations[ocispec.AnnotationRefName]
		if m.Digest == desc.Digest {
			if tag == reference {
				return nil
			}
			descExists = true
		}
	}
	if !descExists {
		return fmt.Errorf("%s: %s: %w", desc.Digest, desc.MediaType, errdef.ErrNotFound)
	}
	if desc.Annotations == nil {
		desc.Annotations = map[string]string{}
	}
	desc.Annotations[ocispec.AnnotationRefName] = reference
	return li.addManifest(desc)
}

func (li *localIndex) untag(_ context.Context, reference string) error {
	deleted := false
	for idx, m := range li.Manifests {
		tag := m.Annotations[ocispec.AnnotationRefName]
		if tag == reference {
			delete(li.Manifests[idx].Annotations, ocispec.AnnotationRefName)
			deleted = true
		}
	}
	if !deleted {
		return errdef.ErrNotFound
	}
	return li.save()
}

// parseIndexJson parses an OCI index at specified path
func parseIndex(indexPath string) (*localIndex, error) {
	index := &localIndex{
		Index: ocispec.Index{
			Versioned: specs.Versioned{
				SchemaVersion: 2,
			},
		},
		indexPath: indexPath,
	}

	indexBytes, err := os.ReadFile(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return index, nil
		}
		return nil, fmt.Errorf("failed to read index: %w", err)
	}

	if err := json.Unmarshal(indexBytes, &index.Index); err != nil {
		return nil, fmt.Errorf("failed to parse index: %w", err)
	}

	return index, nil
}
