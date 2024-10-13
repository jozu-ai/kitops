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

package kitfile

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"kitops/pkg/lib/repo/local"
	"kitops/pkg/lib/repo/util"

	"kitops/pkg/artifact"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/filesystem"
	"kitops/pkg/output"

	"github.com/opencontainers/go-digest"
	specs "github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
)

// SaveModel saves an *artifact.Model to the provided oras.Target, compressing layers. It attempts to block
// modelkits that include paths that leave the base context directory, allowing only subdirectories of the root
// context to be included in the modelkit.
func SaveModel(ctx context.Context, localRepo local.LocalRepo, kitfile *artifact.KitFile, ignore filesystem.IgnorePaths, compression string) (*ocispec.Descriptor, error) {
	configDesc, err := saveConfig(ctx, localRepo, kitfile)
	if err != nil {
		return nil, err
	}
	layerDescs, err := saveKitfileLayers(ctx, localRepo, kitfile, ignore, compression)
	if err != nil {
		return nil, err
	}

	manifest := createManifest(configDesc, layerDescs)

	manifestDesc, err := saveModelManifest(ctx, localRepo, manifest)
	if err != nil {
		return nil, err
	}
	return manifestDesc, nil
}

func saveConfig(ctx context.Context, localRepo local.LocalRepo, kitfile *artifact.KitFile) (ocispec.Descriptor, error) {
	modelBytes, err := kitfile.MarshalToJSON()
	if err != nil {
		return ocispec.DescriptorEmptyJSON, err
	}
	desc := ocispec.Descriptor{
		MediaType: constants.ModelConfigMediaType.String(),
		Digest:    digest.FromBytes(modelBytes),
		Size:      int64(len(modelBytes)),
	}

	exists, err := localRepo.Exists(ctx, desc)
	if err != nil {
		return ocispec.DescriptorEmptyJSON, err
	}
	if !exists {
		// Does not exist in storage, need to push
		err = localRepo.Push(ctx, desc, bytes.NewReader(modelBytes))
		if err != nil {
			return ocispec.DescriptorEmptyJSON, err
		}
		output.Infof("Saved configuration: %s", desc.Digest)
	} else {
		output.Infof("Configuration already exists in storage: %s", desc.Digest)
	}

	return desc, nil
}

func saveKitfileLayers(ctx context.Context, localRepo local.LocalRepo, kitfile *artifact.KitFile, ignore filesystem.IgnorePaths, compression string) ([]ocispec.Descriptor, error) {

	modelPartsLen := 0
	if kitfile.Model != nil {
		modelPartsLen = len(kitfile.Model.Parts)
		if kitfile.Model.Path != "" && !util.IsModelKitReference(kitfile.Model.Path) {
			modelPartsLen++ // Account for the model path
		}
	}
	var layers = make([]ocispec.Descriptor, modelPartsLen+len(kitfile.Code)+len(kitfile.DataSets)+len(kitfile.Docs))
	var wg sync.WaitGroup
	errChan := make(chan error, len(layers))

	processLayer := func(index int, path string, mediaType constants.MediaType) {
		defer wg.Done()

		if ctx.Err() != nil {
			errChan <- ctx.Err()
			return
		}

		layer, err := saveContentLayer(ctx, localRepo, path, mediaType, ignore)
		if err != nil {
			errChan <- err
			return
		}
		// Place the layer in the correct position in the layers slice
		layers[index] = layer
	}

	// Counter to track index of each layer
	layerIndex := 0

	// Process model layers
	if kitfile.Model != nil {
		if kitfile.Model.Path != "" && !util.IsModelKitReference(kitfile.Model.Path) {
			wg.Add(1)
			go processLayer(layerIndex, kitfile.Model.Path, constants.MediaType{
				BaseType:    constants.ModelType,
				Compression: compression,
			})
			layerIndex++
		}
		for _, part := range kitfile.Model.Parts {
			wg.Add(1)
			go processLayer(layerIndex, part.Path, constants.MediaType{
				BaseType:    constants.ModelPartType,
				Compression: compression,
			})
			layerIndex++
		}
	}

	// Process code layers
	for _, code := range kitfile.Code {
		wg.Add(1)
		go processLayer(layerIndex, code.Path, constants.MediaType{
			BaseType:    constants.CodeType,
			Compression: compression,
		})
		layerIndex++
	}

	// Process dataset layers
	for _, dataset := range kitfile.DataSets {
		wg.Add(1)
		go processLayer(layerIndex, dataset.Path, constants.MediaType{
			BaseType:    constants.DatasetType,
			Compression: compression,
		})
		layerIndex++
	}

	// Process documentation layers
	for _, docs := range kitfile.Docs {
		wg.Add(1)
		go processLayer(layerIndex, docs.Path, constants.MediaType{
			BaseType:    constants.DocsType,
			Compression: compression,
		})
		layerIndex++
	}
	wg.Wait()
	close(errChan)

	var allErrors []error
	for err := range errChan {
		if err != nil {
			allErrors = append(allErrors, err)
		}
	}

	if len(allErrors) > 0 {
		return nil, errors.Join(allErrors...)
	}

	return layers, nil
}

func saveContentLayer(ctx context.Context, localRepo local.LocalRepo, path string, mediaType constants.MediaType, ignore filesystem.IgnorePaths) (ocispec.Descriptor, error) {
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

	if exists, err := localRepo.Exists(ctx, desc); err != nil {
		return ocispec.DescriptorEmptyJSON, err
	} else if exists {
		output.Infof("Already saved %s layer: %s", mediaType.BaseType, desc.Digest)
		return desc, nil
	}

	// Workaround to avoid copying a potentially very large file: move it to the expected path
	// and verify that it exists afterwards.
	blobPath := localRepo.BlobPath(desc)
	if err := os.Rename(tempPath, blobPath); err != nil {
		// This may fail on some systems (e.g. linux where / and /home are different partitions)
		// Fallback to regular push which is basically a copy
		output.Debugf("Failed to move temp file into storage (will copy instead): %s", err)
		file, err := os.Open(tempPath)
		if err != nil {
			return ocispec.DescriptorEmptyJSON, fmt.Errorf("failed to open temporary file: %w", err)
		}
		defer file.Close()
		if err := localRepo.Push(ctx, desc, file); err != nil {
			return ocispec.DescriptorEmptyJSON, fmt.Errorf("failed to add layer to storage: %w", err)
		}
	}

	// Verify blob is in store now
	exists, err := localRepo.Exists(ctx, desc)
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

func createManifest(configDesc ocispec.Descriptor, layerDescs []ocispec.Descriptor) ocispec.Manifest {
	manifest := ocispec.Manifest{
		Versioned: specs.Versioned{SchemaVersion: 2},
		Config:    configDesc,
		Layers:    layerDescs,
		Annotations: map[string]string{
			constants.CliVersionAnnotation: constants.Version,
		},
	}

	return manifest
}
