package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"kitops/pkg/artifact"
	"kitops/pkg/lib/constants"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/content"
)

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
