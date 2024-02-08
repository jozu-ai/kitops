package artifact

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	_ "crypto/sha256"
	_ "crypto/sha512"

	"github.com/opencontainers/go-digest"
	specs "github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/content/oci"
)

type Store struct {
	Storage *oci.Store
}

func NewArtifactStore() *Store {

	store, err := oci.New(".jozuStore")
	if err != nil {
		panic(err)
	}

	return &Store{
		Storage: store,
	}
}

func (store *Store) saveContentLayer(layer *ModelLayer) (*ocispec.Descriptor, error) {
	ctx := context.Background()

	buf := &bytes.Buffer{}
	err := layer.Apply(buf)
	if err != nil {
		return nil, err
	}

	// Create a descriptor for the layer
	desc := ocispec.Descriptor{
		MediaType: ModelLayerMediaType,
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

func (store *Store) saveConfigFile(model *JozuFile) (*ocispec.Descriptor, error) {
	ctx := context.Background()
	buf, err := model.MarshalToJSON()
	if err != nil {
		return nil, err
	}
	desc := ocispec.Descriptor{
		MediaType: ModelConfigMediaType,
		Digest:    digest.FromBytes(buf),
		Size:      int64(len(buf)),
	}
	err = store.Storage.Push(ctx, desc, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	return &desc, nil
}

func (store *Store) saveModelManifest(model *Model, config *ocispec.Descriptor) (*ocispec.Manifest, error) {
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

func (store *Store) SaveModel(model *Model) (*ocispec.Manifest, error) {
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
