package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"kitops/pkg/artifact"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/repo"
	"kitops/pkg/output"

	"github.com/opencontainers/go-digest"
	specs "github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
)

func SaveModel(ctx context.Context, store oras.Target, model *artifact.Model, tag string) (*ocispec.Descriptor, error) {
	configDesc, err := saveConfigFile(ctx, store, model.Config)
	if err != nil {
		return nil, err
	}
	var layerDescs []ocispec.Descriptor
	for _, layer := range model.Layers {
		layerDesc, err := saveContentLayer(ctx, store, &layer)
		if err != nil {
			return nil, err
		}
		layerDescs = append(layerDescs, layerDesc)
	}

	manifest := CreateManifest(configDesc, layerDescs)
	manifestDesc, err := saveModelManifest(ctx, store, manifest, tag)
	if err != nil {
		return nil, err
	}
	return manifestDesc, nil
}

func saveContentLayer(ctx context.Context, store oras.Target, layer *artifact.ModelLayer) (ocispec.Descriptor, error) {
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

	exists, err := store.Exists(ctx, desc)
	if err != nil {
		return ocispec.DescriptorEmptyJSON, err
	}
	if exists {
		output.Infof("Model layer already saved: %s", desc.Digest)
	} else {
		// Does not exist in storage, need to push
		err = store.Push(ctx, desc, buf)
		if err != nil {
			return ocispec.DescriptorEmptyJSON, err
		}
		output.Infof("Saved model layer: %s", desc.Digest)
	}

	return desc, nil
}

func saveConfigFile(ctx context.Context, store oras.Target, model *artifact.KitFile) (ocispec.Descriptor, error) {
	modelBytes, err := model.MarshalToJSON()
	if err != nil {
		return ocispec.DescriptorEmptyJSON, err
	}
	desc := ocispec.Descriptor{
		MediaType: constants.ModelConfigMediaType,
		Digest:    digest.FromBytes(modelBytes),
		Size:      int64(len(modelBytes)),
	}

	exists, err := store.Exists(ctx, desc)
	if err != nil {
		return ocispec.DescriptorEmptyJSON, err
	}
	if !exists {
		// Does not exist in storage, need to push
		err = store.Push(ctx, desc, bytes.NewReader(modelBytes))
		if err != nil {
			return ocispec.DescriptorEmptyJSON, err
		}
		output.Infof("Saved configuration: %s", desc.Digest)
	} else {
		output.Infof("Configuration already exists in storage: %s", desc.Digest)
	}

	return desc, nil
}

func saveModelManifest(ctx context.Context, store oras.Target, manifest ocispec.Manifest, tag string) (*ocispec.Descriptor, error) {
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

	if exists, err := store.Exists(ctx, desc); err != nil {
		return nil, err
	} else if !exists {
		// Does not exist in storage, need to push
		err = store.Push(ctx, desc, bytes.NewReader(manifestBytes))
		if err != nil {
			return nil, err
		}
		output.Infof("Saved manifest to storage: %s", desc.Digest)
	} else {
		output.Infof("Manifest already exists in storage: %s", desc.Digest)
	}

	if tag != "" {
		if err := repo.ValidateTag(tag); err != nil {
			return nil, err
		}
		if err := store.Tag(ctx, desc, tag); err != nil {
			return nil, fmt.Errorf("failed to tag manifest: %w", err)
		}
		output.Debugf("Added tag to manifest: %s", tag)
	}

	return &desc, nil
}

func CreateManifest(configDesc ocispec.Descriptor, layerDescs []ocispec.Descriptor) ocispec.Manifest {
	manifest := ocispec.Manifest{
		Versioned:   specs.Versioned{SchemaVersion: 2},
		Config:      configDesc,
		Layers:      layerDescs,
		Annotations: map[string]string{},
	}

	return manifest
}
