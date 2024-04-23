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

package info

import (
	"context"
	"fmt"
	"kitops/pkg/artifact"
	"kitops/pkg/cmd/options"
	"kitops/pkg/lib/repo"
)

func GetKitfileForRef(ctx context.Context, configHome, ref string) (*artifact.KitFile, error) {
	modelRef, _, err := repo.ParseReference(ref)
	if err != nil {
		return nil, err
	}

	opts := &infoOptions{
		configHome: configHome,
		modelRef:   modelRef,
		NetworkOptions: options.NetworkOptions{
			PlainHTTP: false,
			TlsVerify: true,
		},
	}

	kitfile, err := getLocalConfig(ctx, opts)
	if err == nil {
		return kitfile, nil
	}
	kitfile, err = getRemoteConfig(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch Kitfile for %s: %w", ref, err)
	}
	return kitfile, err
}
