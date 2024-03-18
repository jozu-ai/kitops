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

package inspect

import (
	"context"
	"fmt"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/repo"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

func inspectReference(ctx context.Context, opts *inspectOptions) (*ocispec.Manifest, error) {
	if opts.checkRemote {
		return getRemoteManifest(ctx, opts)
	} else {
		return getLocalManifest(ctx, opts)
	}
}

func getLocalManifest(ctx context.Context, opts *inspectOptions) (*ocispec.Manifest, error) {
	storageRoot := constants.StoragePath(opts.configHome)
	store, err := repo.NewLocalStore(storageRoot, opts.modelRef)
	if err != nil {
		return nil, fmt.Errorf("failed to read local storage: %w", err)
	}
	return repo.ResolveManifest(ctx, store, opts.modelRef.Reference)
}

func getRemoteManifest(ctx context.Context, opts *inspectOptions) (*ocispec.Manifest, error) {
	repository, err := repo.NewRepository(ctx, opts.modelRef.Registry, opts.modelRef.Repository, &repo.RegistryOptions{
		PlainHTTP:       opts.PlainHTTP,
		SkipTLSVerify:   !opts.TlsVerify,
		CredentialsPath: constants.CredentialsPath(opts.configHome),
	})
	if err != nil {
		return nil, err
	}
	return repo.ResolveManifest(ctx, repository, opts.modelRef.Reference)
}
