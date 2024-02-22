package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"kitops/pkg/artifact"
	"kitops/pkg/lib/constants"
	"kitops/pkg/output"
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
	repo      string
}

// Assert LocalStore implements the Store interface.
var _ Store = (*LocalStore)(nil)

func NewLocalStore(storeRoot, repo string) Store {
	storeHome := filepath.Join(storeRoot, repo)
	indexPath := filepath.Join(storeHome, "index.json")

	store, err := oci.New(storeHome)
	if err != nil {
		panic(err)
	}

	return &LocalStore{
		storage:   store,
		indexPath: indexPath,
		repo:      repo,
	}
}

func (store *LocalStore) SaveModel(model *artifact.Model, tag string) (*ocispec.Descriptor, error) {
	configDesc, err := store.saveConfigFile(model.Config)
	if err != nil {
		return nil, err
	}
	var layerDescs []ocispec.Descriptor
	for _, layer := range model.Layers {
		layerDesc, err := store.saveContentLayer(&layer)
		if err != nil {
			return nil, err
		}
		layerDescs = append(layerDescs, layerDesc)
	}

	manifestDesc, err := store.saveModelManifest(layerDescs, configDesc, tag)
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

func (store *LocalStore) GetRepository() string {
	return store.repo
}

func (store *LocalStore) saveContentLayer(layer *artifact.ModelLayer) (ocispec.Descriptor, error) {
	ctx := context.Background()

	buf := &bytes.Buffer{}
	err := layer.Apply(buf)
	if err != nil {
		return ocispec.DescriptorEmptyJSON, err
	}

	// Create a descriptor for the layer
	desc := ocispec.Descriptor{
		MediaType: layer.MediaType,
		Digest:    digest.FromBytes(buf.Bytes()),
		Size:      int64(buf.Len()),
	}

	exists, err := store.storage.Exists(ctx, desc)
	if err != nil {
		return ocispec.DescriptorEmptyJSON, err
	}
	if exists {
		output.Infof("Model layer already saved: %s", desc.Digest)
	} else {
		// Does not exist in storage, need to push
		err = store.storage.Push(ctx, desc, buf)
		if err != nil {
			return ocispec.DescriptorEmptyJSON, err
		}
		output.Infof("Saved model layer: %s", desc.Digest)
	}

	return desc, nil
}

func (store *LocalStore) saveConfigFile(model *artifact.KitFile) (ocispec.Descriptor, error) {
	ctx := context.Background()
	modelBytes, err := model.MarshalToJSON()
	if err != nil {
		return ocispec.DescriptorEmptyJSON, err
	}
	desc := ocispec.Descriptor{
		MediaType: constants.ModelConfigMediaType,
		Digest:    digest.FromBytes(modelBytes),
		Size:      int64(len(modelBytes)),
	}

	exists, err := store.storage.Exists(ctx, desc)
	if err != nil {
		return ocispec.DescriptorEmptyJSON, err
	}
	if !exists {
		// Does not exist in storage, need to push
		err = store.storage.Push(ctx, desc, bytes.NewReader(modelBytes))
		if err != nil {
			return ocispec.DescriptorEmptyJSON, err
		}
		output.Infof("Saved configuration: %s", desc.Digest)
	} else {
		output.Infof("Configuration already exists in storage: %s", desc.Digest)
	}

	return desc, nil
}

func (store *LocalStore) saveModelManifest(layerDescs []ocispec.Descriptor, config ocispec.Descriptor, tag string) (*ocispec.Descriptor, error) {
	ctx := context.Background()
	manifest := ocispec.Manifest{
		Versioned:   specs.Versioned{SchemaVersion: 2},
		Config:      config,
		Layers:      []ocispec.Descriptor{},
		Annotations: map[string]string{},
	}
	// Add the layers to the manifest
	manifest.Layers = append(manifest.Layers, layerDescs...)

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
		output.Infof("Saved manifest to storage: %s", desc.Digest)
	} else {
		output.Infof("Manifest already exists in storage: %s", desc.Digest)
	}

	if tag != "" {
		if err := validateTag(tag); err != nil {
			return nil, err
		}
		if err := store.storage.Tag(context.Background(), desc, tag); err != nil {
			return nil, fmt.Errorf("failed to tag manifest: %w", err)
		}
		output.Debugf("Added tag to manifest: %s", tag)
	}

	return &desc, nil
}
