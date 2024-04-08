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
	"fmt"
	"kitops/pkg/lib/constants"
	"os"
	"path/filepath"

	"github.com/moby/patternmatcher"
	"github.com/moby/patternmatcher/ignorefile"
)

func readIgnoreFile(contextDir string) (*patternmatcher.PatternMatcher, error) {
	ignorePath := filepath.Join(contextDir, constants.IgnoreFileName)
	ignoreFile, err := os.Open(ignorePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &patternmatcher.PatternMatcher{}, nil
		}
		return nil, err
	}
	patterns, err := ignorefile.ReadAll(ignoreFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s file: %w", constants.IgnoreFileName, err)
	}
	pm, err := patternmatcher.New(patterns)
	if err != nil {
		return nil, fmt.Errorf("invalid %s file: %w", constants.IgnoreFileName, err)
	}
	return pm, nil
}
