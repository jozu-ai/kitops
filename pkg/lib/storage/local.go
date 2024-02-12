package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"jmm/pkg/artifact"
	"jmm/pkg/lib/constants"
	"os"
	"path/filepath"

	"github.com/opencontainers/go-digest"
	specs "github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/content/oci"
)

type LocalStore struct {
	storage   *oci.Store
	indexPath string
}

// Assert LocalStore implements the Store interface.
var _ Store = (*LocalStore)(nil)

func NewLocalStore(storeHome string) Store {
	indexPath := filepath.Join(storeHome, "index.json")

	store, err := oci.New(storeHome)
	if err != nil {
		panic(err)
	}

	return &LocalStore{
		storage:   store,
		indexPath: indexPath,
	}
}

func (store *LocalStore) SaveModel(model *artifact.Model, tag string) (*ocispec.Descriptor, error) {
	config, err := store.saveConfigFile(model.Config)
	if err != nil {
		return nil, err
	}
	for _, layer := range model.Layers {
		_, err = store.saveContentLayer(&layer)
		if err != nil {
			return nil, err
		}
	}

	manifestDesc, err := store.saveModelManifest(model, config, tag)
	if err != nil {
		return nil, err
	}
	return manifestDesc, nil
}

func (store *LocalStore) TagModel(manifestDesc ocispec.Descriptor, tag string) error {
	if err := validateTag(tag); err != nil {
		return err
	}

	if err := store.storage.Tag(context.Background(), manifestDesc, tag); err != nil {
		return fmt.Errorf("failed to tag manifest: %w", err)
	}

	return nil
}

func (store *LocalStore) Fetch(ctx context.Context, desc ocispec.Descriptor) ([]byte, error) {
	bytes, err := content.FetchAll(ctx, store.storage, desc)
	return bytes, err
}

func (store *LocalStore) ParseIndexJson() (*ocispec.Index, error) {
	indexBytes, err := os.ReadFile(store.indexPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read index: %w", err)
	}

	index := &ocispec.Index{}
	if err := json.Unmarshal(indexBytes, index); err != nil {
		return nil, fmt.Errorf("failed to parse index: %w", err)
	}

	return index, nil
}

func (store *LocalStore) saveContentLayer(layer *artifact.ModelLayer) (*ocispec.Descriptor, error) {
	ctx := context.Background()

	buf := &bytes.Buffer{}
	err := layer.Apply(buf)
	if err != nil {
		return nil, err
	}

	// Create a descriptor for the layer
	desc := ocispec.Descriptor{
		MediaType: constants.ModelLayerMediaType,
		Digest:    digest.FromBytes(buf.Bytes()),
		Size:      int64(buf.Len()),
	}

	exists, err := store.storage.Exists(ctx, desc)
	if err != nil {
		return nil, err
	}
	if exists {
		fmt.Println("Model layer already saved: ", desc.Digest)
	} else {
		// Does not exist in storage, need to push
		err = store.storage.Push(ctx, desc, buf)
		if err != nil {
			return nil, err
		}
		fmt.Println("Saved model layer: ", desc.Digest)
	}

	return &desc, nil
}

func (store *LocalStore) saveConfigFile(model *artifact.JozuFile) (*ocispec.Descriptor, error) {
	ctx := context.Background()
	modelBytes, err := model.MarshalToJSON()
	if err != nil {
		return nil, err
	}
	desc := ocispec.Descriptor{
		MediaType: constants.ModelConfigMediaType,
		Digest:    digest.FromBytes(modelBytes),
		Size:      int64(len(modelBytes)),
	}

	exists, err := store.storage.Exists(ctx, desc)
	if err != nil {
		return nil, err
	}
	if !exists {
		// Does not exist in storage, need to push
		err = store.storage.Push(ctx, desc, bytes.NewReader(modelBytes))
		if err != nil {
			return nil, err
		}
	}

	return &desc, nil
}

func (store *LocalStore) saveModelManifest(model *artifact.Model, config *ocispec.Descriptor, tag string) (*ocispec.Descriptor, error) {
	ctx := context.Background()
	manifest := ocispec.Manifest{
		Versioned:   specs.Versioned{SchemaVersion: 2},
		Config:      *config,
		Layers:      []ocispec.Descriptor{},
		Annotations: map[string]string{},
	}
	// Add the layers to the manifest
	for _, layer := range model.Layers {
		manifest.Layers = append(manifest.Layers, layer.Descriptor)
	}

	manifestBytes, err := json.Marshal(manifest)
	if err != nil {
		return nil, err
	}
	// Push the manifest to the store
	desc := ocispec.Descriptor{
		MediaType: ocispec.MediaTypeImageManifest,
		Digest:    digest.FromBytes(manifestBytes),
		Size:      int64(len(manifestBytes)),
	}

	if exists, err := store.storage.Exists(ctx, desc); err != nil {
		return nil, err
	} else if !exists {
		// Does not exist in storage, need to push
		err = store.storage.Push(ctx, desc, bytes.NewReader(manifestBytes))
		if err != nil {
			return nil, err
		}
	}

	if tag != "" {
		if err := validateTag(tag); err != nil {
			return nil, err
		}
		if err := store.storage.Tag(context.Background(), desc, tag); err != nil {
			return nil, fmt.Errorf("failed to tag manifest: %w", err)
		}
	}

	return &desc, nil
}
