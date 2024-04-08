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

func getIgnoreMatcher(contextDir string) (*patternmatcher.PatternMatcher, error) {
	filePatterns, err := readIgnoreFile(contextDir)
	if err != nil {
		return nil, err
	}
	filePatterns = append(filePatterns, constants.DefaultKitfileNames()...)
	filePatterns = append(filePatterns, constants.IgnoreFileName)

	pm, err := patternmatcher.New(filePatterns)
	if err != nil {
		return nil, fmt.Errorf("invalid %s file: %w", constants.IgnoreFileName, err)
	}
	return pm, nil
}

func readIgnoreFile(contextDir string) ([]string, error) {
	ignorePath := filepath.Join(contextDir, constants.IgnoreFileName)
	ignoreFile, err := os.Open(ignorePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	patterns, err := ignorefile.ReadAll(ignoreFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s file: %w", constants.IgnoreFileName, err)
	}
	return patterns, nil
}
