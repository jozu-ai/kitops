// Copyright 2024 The KitOps Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"kitops/pkg/artifact"
	"kitops/pkg/cmd/version"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/filesystem"
	kfutils "kitops/pkg/lib/kitfile"
	"kitops/pkg/lib/repo"
	"kitops/pkg/output"

	"github.com/opencontainers/go-digest"
	specs "github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
)

// SaveModel saves an *artifact.Model to the provided oras.Target, compressing layers. It attempts to block
// modelkits that include paths that leave the base context directory, allowing only subdirectories of the root
// context to be included in the modelkit.
func SaveModel(ctx context.Context, store repo.LocalStorage, kitfile *artifact.KitFile, ignore filesystem.IgnorePaths, compression string) (*ocispec.Descriptor, error) {
	configDesc, err := saveConfig(ctx, store, kitfile)
	if err != nil {
		return nil, err
	}
	layerDescs, err := saveKitfileLayers(ctx, store, kitfile, ignore, compression)
	if err != nil {
		return nil, err
	}

	manifest := CreateManifest(configDesc, layerDescs)

	manifestDesc, err := saveModelManifest(ctx, store, manifest)
	if err != nil {
		return nil, err
	}
	return manifestDesc, nil
}

func saveConfig(ctx context.Context, store repo.LocalStorage, kitfile *artifact.KitFile) (ocispec.Descriptor, error) {
	modelBytes, err := kitfile.MarshalToJSON()
	if err != nil {
		return ocispec.DescriptorEmptyJSON, err
	}
	desc := ocispec.Descriptor{
		MediaType: constants.ModelConfigMediaType.String(),
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

func saveKitfileLayers(ctx context.Context, store repo.LocalStorage, kitfile *artifact.KitFile, ignore filesystem.IgnorePaths, compression string) ([]ocispec.Descriptor, error) {
	var layers []ocispec.Descriptor
	if kitfile.Model != nil {
		if kitfile.Model.Path != "" && !kfutils.IsModelKitReference(kitfile.Model.Path) {
			mediaType := constants.MediaType{
				BaseType:    constants.ModelType,
				Compression: compression,
			}
			layer, err := saveContentLayer(ctx, store, kitfile.Model.Path, mediaType, ignore)
			if err != nil {
				return nil, err
			}
			layers = append(layers, layer)
		}
		for _, part := range kitfile.Model.Parts {
			mediaType := constants.MediaType{
				BaseType:    constants.ModelPartType,
				Compression: compression,
			}
			layer, err := saveContentLayer(ctx, store, part.Path, mediaType, ignore)
			if err != nil {
				return nil, err
			}
			layers = append(layers, layer)
		}
	}
	for _, code := range kitfile.Code {
		mediaType := constants.MediaType{
			BaseType:    constants.CodeType,
			Compression: compression,
		}
		layer, err := saveContentLayer(ctx, store, code.Path, mediaType, ignore)
		if err != nil {
			return nil, err
		}
		layers = append(layers, layer)
	}
	for _, dataset := range kitfile.DataSets {
		mediaType := constants.MediaType{
			BaseType:    constants.DatasetType,
			Compression: compression,
		}
		layer, err := saveContentLayer(ctx, store, dataset.Path, mediaType, ignore)
		if err != nil {
			return nil, err
		}
		layers = append(layers, layer)
	}
	return layers, nil
}

func saveContentLayer(ctx context.Context, store repo.LocalStorage, path string, mediaType constants.MediaType, ignore filesystem.IgnorePaths) (ocispec.Descriptor, error) {
	// We want to store a gzipped tar file in store, but to do so we need a descriptor, so we have to compress
	// to a temporary file. Ideally, we'd also add this to the internal store by moving the file to avoid
	// copying if possible.
	tempPath, desc, err := compressLayer(path, mediaType, ignore)
	if err != nil {
		return ocispec.DescriptorEmptyJSON, err
	}
	defer func() {
		if err := os.Remove(tempPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			output.Errorf("Failed to remove temporary file %s: %s", tempPath, err)
		}
	}()

	if exists, err := store.Exists(ctx, desc); err != nil {
		return ocispec.DescriptorEmptyJSON, err
	} else if exists {
		output.Infof("Already saved %s layer: %s", mediaType.BaseType, desc.Digest)
		return desc, nil
	}

	// Workaround to avoid copying a potentially very large file: move it to the expected path
	// and verify that it exists afterwards.
	blobPath := repo.BlobPathForManifest(store, desc)
	if err := os.Rename(tempPath, blobPath); err != nil {
		// This may fail on some systems (e.g. linux where / and /home are different partitions)
		// Fallback to regular push which is basically a copy
		output.Debugf("Failed to move temp file into storage: %s", err)
		file, err := os.Open(tempPath)
		if err != nil {
			return ocispec.DescriptorEmptyJSON, fmt.Errorf("failed to open temporary file: %w", err)
		}
		defer file.Close()
		if err := store.Push(ctx, desc, file); err != nil {
			return ocispec.DescriptorEmptyJSON, fmt.Errorf("failed to add layer to storage: %w", err)
		}
	}

	// Verify blob is in store now
	exists, err := store.Exists(ctx, desc)
	if err != nil {
		return ocispec.DescriptorEmptyJSON, err
	}
	if !exists {
		return ocispec.DescriptorEmptyJSON, fmt.Errorf("failed to move layer to storage: file is not stored")
	}

	output.Infof("Saved %s layer: %s", mediaType.BaseType, desc.Digest)
	return desc, nil
}

func saveModelManifest(ctx context.Context, store oras.Target, manifest ocispec.Manifest) (*ocispec.Descriptor, error) {
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
	return &desc, nil
}

func CreateManifest(configDesc ocispec.Descriptor, layerDescs []ocispec.Descriptor) ocispec.Manifest {
	manifest := ocispec.Manifest{
		Versioned: specs.Versioned{SchemaVersion: 2},
		Config:    configDesc,
		Layers:    layerDescs,
		Annotations: map[string]string{
			constants.CliVersionAnnotation: version.Version,
		},
	}

	return manifest
}
