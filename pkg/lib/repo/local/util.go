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

package local

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// parseIndexJson parses an OCI index at specified path
func parseIndex(indexPath string) (*ocispec.Index, error) {
	index := &ocispec.Index{
		Versioned: specs.Versioned{
			SchemaVersion: 2,
		},
	}
	indexBytes, err := os.ReadFile(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return index, nil
		}
		return nil, fmt.Errorf("failed to read index: %w", err)
	}
	if err := json.Unmarshal(indexBytes, index); err != nil {
		return nil, fmt.Errorf("failed to parse index: %w", err)
	}
	return index, nil
}

func parseTagsIndex(tagsIndexPath string) (*tagsIndex, error) {
	bytes, err := os.ReadFile(tagsIndexPath)
	tags := emptyTagsIndex(tagsIndexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return tags, nil
		}
		return nil, fmt.Errorf("failed to read tags index: %w", err)
	}
	if err := json.Unmarshal(bytes, &tags.tagToDigest); err != nil {
		return nil, fmt.Errorf("failed to parse tags index: %w", err)
	}
	return tags, nil
}
