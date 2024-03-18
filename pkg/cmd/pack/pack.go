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

package pack

import (
	"context"
	"fmt"
	"kitops/pkg/artifact"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/filesystem"
	"kitops/pkg/lib/repo"
	"kitops/pkg/lib/storage"
	"kitops/pkg/output"
)

// runPack compresses and stores a modelkit based on a Kitfile. Returns an error if packing
// fails for any reason, or if any path in the Kitfile is not a subdirectory of the current
// context directory.
//
// Packed modelkits are saved to the local on-disk cache. As OCI-spec indexes only support one
// registry/repository reference at a time, individual blobs may be duplicated on disk if stored
// under different references.
func runPack(ctx context.Context, options *packOptions) error {
	// 1. Read the model file
	kitfile := &artifact.KitFile{}
	if err := kitfile.LoadModel(options.modelFile); err != nil {
		return err
	}

	model := &artifact.Model{}
	model.Config = kitfile

	// 2. package the Code
	for _, code := range kitfile.Code {
		codePath, _, err := filesystem.VerifySubpath(options.contextDir, code.Path)
		if err != nil {
			return err
		}
		layer := &artifact.ModelLayer{
			Path:      codePath,
			MediaType: constants.CodeLayerMediaType,
		}
		model.Layers = append(model.Layers, *layer)
	}

	// 3. package the DataSets
	for _, dataset := range kitfile.DataSets {
		datasetPath, _, err := filesystem.VerifySubpath(options.contextDir, dataset.Path)
		if err != nil {
			return err
		}
		layer := &artifact.ModelLayer{
			Path:      datasetPath,
			MediaType: constants.DataSetLayerMediaType,
		}
		model.Layers = append(model.Layers, *layer)
	}

	// 4. package the TrainedModel
	if kitfile.Model != nil {
		modelPath, _, err := filesystem.VerifySubpath(options.contextDir, kitfile.Model.Path)
		if err != nil {
			return err
		}
		layer := &artifact.ModelLayer{
			Path:      modelPath,
			MediaType: constants.ModelLayerMediaType,
		}
		model.Layers = append(model.Layers, *layer)
	}

	tag := ""
	if options.modelRef != nil {
		tag = options.modelRef.Reference
	}
	storageHome := constants.StoragePath(options.configHome)
	localStore, err := repo.NewLocalStore(storageHome, options.modelRef)
	if err != nil {
		return fmt.Errorf("failed to open local storage: %w", err)
	}

	manifestDesc, err := storage.SaveModel(ctx, localStore, model, tag)
	if err != nil {
		return err
	}

	for _, tag := range options.extraRefs {
		if err := localStore.Tag(ctx, *manifestDesc, tag); err != nil {
			return err
		}
	}

	output.Infof("Model saved: %s", manifestDesc.Digest)

	return nil
}
