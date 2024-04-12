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

package filesystem

import (
	"fmt"
	"io/fs"
	"kitops/pkg/lib/constants"
	"os"
	"path/filepath"
	"strings"
)

// VerifySubpath checks that filepath.Join(context, subDir) is a subdirectory of context, following
// symlinks if present.
func VerifySubpath(context, subDir string) (absPath, relPath string, err error) {
	if filepath.IsAbs(subDir) {
		return "", "", fmt.Errorf("absolute paths are not supported (%s)", subDir)
	}

	// Get absolute path for context
	absContext, err := filepath.Abs(context)
	if err != nil {
		return "", "", fmt.Errorf("failed to resolve absolute path for %s: %w", context, err)
	}
	if _, exists := PathExists(absContext); exists {
		res, err := filepath.EvalSymlinks(absContext)
		if err != nil {
			return "", "", fmt.Errorf("error resolving %s: %w", absContext, err)
		}
		absContext = res
	}

	// Get absolute path for context + subpath
	fullPath := filepath.Clean(filepath.Join(absContext, subDir))
	if _, exists := PathExists(fullPath); exists {
		res, err := filepath.EvalSymlinks(fullPath)
		if err != nil {
			return "", "", fmt.Errorf("error resolving %s: %w", absContext, err)
		}
		fullPath = res
	}

	// Get relative path between context and the full path to check if the
	// actual full, absolute path is a subdirectory of context
	relPath, err = filepath.Rel(absContext, fullPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to get relative path: %w", err)
	}
	if strings.Contains(relPath, "..") {
		return "", "", fmt.Errorf("paths must be within context directory")
	}

	return fullPath, relPath, nil
}

func PathExists(path string) (fs.FileInfo, bool) {
	fi, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return nil, false
	}
	return fi, true
}

// Searches for a kit file in the given context directory.
// It checks for accepted kitfile names and returns the absolute path for the first found kitfile.
// If no kitfile is found, returns error
func FindKitfileInPath(contextDir string) (string, error) {
	var defaultKitFileNames = constants.DefaultKitfileNames()
	for _, fileName := range defaultKitFileNames {
		if _, exists := PathExists(filepath.Join(contextDir, fileName)); exists {
			absPath, err := filepath.Abs(filepath.Join(contextDir, fileName))
			if err != nil {
				return "", fmt.Errorf("Failed to find Kitfile: %w", err)
			}
			return absPath, nil
		}
	}
	return "", fmt.Errorf("No Kitfile found in %s. Consider using the -f flag to specify its path", contextDir)
}
