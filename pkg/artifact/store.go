package artifact

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	_ "crypto/sha256"
	_ "crypto/sha512"

	"github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/specs-go"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/content/oci"
)

type Store struct {
	Storage *oci.Store
}

func NewArtifactStore() *Store {

	store, error := oci.New(".jozuStore")
	if error != nil {
		panic(error)
	}

	return &Store{
		Storage: store,
	}
}

func (store *Store) saveContentLayer(layer *ModelLayer) (*v1.Descriptor, error) {
	ctx := context.Background()

	buf := &bytes.Buffer{}
	err := layer.Apply(buf)
	if err != nil {
		return nil, err
	}

	// Create a descriptor for the layer
	desc := v1.Descriptor{
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

func (store *Store) saveConfigFile(model *JozuFile) (*v1.Descriptor, error)  {
	ctx := context.Background()
	buf, err := model.MarshalToJSON()
	if err != nil {
		return nil, err
	}
	desc := v1.Descriptor{
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

func (store *Store) saveModelManifest(model *Model, config *v1.Descriptor) (*v1.Manifest, error) {
	ctx := context.Background()
	manifest := v1.Manifest{
		Versioned: specs.Versioned{SchemaVersion: 2},
		Config: *config,
		Layers: []v1.Descriptor{},
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
	desc := v1.Descriptor{
		MediaType: v1.MediaTypeImageManifest,
		Digest:    digest.FromBytes(manifestBytes),
		Size:      int64(len(manifestBytes)),
	}
	err = store.Storage.Push(ctx, desc, bytes.NewReader(manifestBytes))
	if err != nil {
		return nil, err
	}
	return &manifest, nil
}

func (store *Store) SaveModel(model *Model) (*v1.Manifest, error) {
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