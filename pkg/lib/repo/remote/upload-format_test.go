// Copyright 2025 The KitOps Authors.
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

package remote

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUploadFormat(t *testing.T) {
	tests := []struct {
		registry       string
		size           int64
		expectedFormat uploadFormat
	}{
		{
			registry:       "quay.io",
			size:           100,
			expectedFormat: uploadMonolithicPut,
		},
		{
			registry:       "docker.io",
			size:           500,
			expectedFormat: uploadMonolithicPut,
		},
		{
			registry:       "quay.io",
			size:           uploadChunkDefaultSize - 1,
			expectedFormat: uploadMonolithicPut,
		},
		{
			registry:       "docker.io",
			size:           uploadChunkDefaultSize - 1,
			expectedFormat: uploadMonolithicPut,
		},
		{
			registry:       "quay.io",
			size:           uploadChunkDefaultSize,
			expectedFormat: uploadChunkedPatch,
		},
		{
			registry:       "docker.io",
			size:           uploadChunkDefaultSize,
			expectedFormat: uploadChunkedPatch,
		},
		// GitHub container registry -- should always use monolithic
		{
			registry:       "ghcr.io",
			size:           100,
			expectedFormat: uploadMonolithicPut,
		},
		{
			registry:       "ghcr.io",
			size:           uploadChunkDefaultSize * 2,
			expectedFormat: uploadMonolithicPut,
		},
	}

	for _, tt := range tests {
		t.Run(tt.registry, func(t *testing.T) {
			actualFormat := getUploadFormat(tt.registry, tt.size)
			assert.Equal(t, tt.expectedFormat, actualFormat)
		})
	}
}

func TestGetUploadFormatGoogleArtifactRegistry(t *testing.T) {
	testRegistries := []string{
		".pkg.dev",
		"docker.pkg.dev",
		"gcr.io",
		"asia.gcr.io",
		"marketplace.gcr.io",
		"eu.gcr.io",
		"us.gcr.io",
		"region-docker.pkg.dev",
		"northamerica-northeast1-docker.pkg.dev",
		"us-central1-docker.pkg.dev",
		"us-east1-docker.pkg.dev",
	}

	for _, registry := range testRegistries {
		uploadFormatSmall := getUploadFormat(registry, 100)
		assert.Equal(t, uploadMonolithicPut, uploadFormatSmall, "Small layers should use monolithic put")
		uploadFormatLarge := getUploadFormat(registry, uploadChunkDefaultSize)
		assert.Equal(t, uploadMonolithicPut, uploadFormatLarge, "Large layers should use monolithic put")
	}
}
