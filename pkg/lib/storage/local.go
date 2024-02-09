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
	Storage   *oci.Store
	indexPath string
}

// Assert LocalStore implements the Store interface.
var _ Store = (*LocalStore)(nil)

func NewLocalStore(jozuhome string) *LocalStore {
	storeHome := filepath.Join(jozuhome, ".jozuStore")
	indexPath := filepath.Join(storeHome, "index.json")

	store, err := oci.New(storeHome)
	if err != nil {
		panic(err)
	}

	return &LocalStore{
		Storage:   store,
		indexPath: indexPath,
	}
}

func (store *LocalStore) SaveModel(model *artifact.Model) (*ocispec.Manifest, error) {
	config, err := store.saveConfigFile(model.Config)
	if err != nil {
		return nil, err
	}
	for _, layer := range model.Layers {
		_, err = store.saveContentLayer(layer)
		if err != nil {
			return nil, err
		}
	}

	manifest, err := store.saveModelManifest(model, config)
	if err != nil {
		return nil, err
	}
	return manifest, nil
}

func (store *LocalStore) Fetch(ctx context.Context, desc ocispec.Descriptor) ([]byte, error) {
	bytes, err := content.FetchAll(ctx, store.Storage, desc)
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
	err = store.Storage.Push(ctx, desc, buf)
	layer.Descriptor = desc
	if err != nil {
		return nil, err
	}
	fmt.Println("Saved model layer: ", desc.Digest)
	return &desc, nil
}

func (store *LocalStore) saveConfigFile(model *artifact.JozuFile) (*ocispec.Descriptor, error) {
	ctx := context.Background()
	buf, err := model.MarshalToJSON()
	if err != nil {
		return nil, err
	}
	desc := ocispec.Descriptor{
		MediaType: constants.ModelConfigMediaType,
		Digest:    digest.FromBytes(buf),
		Size:      int64(len(buf)),
	}
	err = store.Storage.Push(ctx, desc, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	return &desc, nil
}

func (store *LocalStore) saveModelManifest(model *artifact.Model, config *ocispec.Descriptor) (*ocispec.Manifest, error) {
	ctx := context.Background()
	manifest := ocispec.Manifest{
		Versioned: specs.Versioned{SchemaVersion: 2},
		Config:    *config,
		Layers:    []ocispec.Descriptor{},
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
	err = store.Storage.Push(ctx, desc, bytes.NewReader(manifestBytes))
	if err != nil {
		return nil, err
	}
	return &manifest, nil
}
