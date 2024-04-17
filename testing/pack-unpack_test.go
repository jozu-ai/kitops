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

package testing

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"kitops/pkg/lib/constants"
)

// TestPackUnpack tests kit functionality by generating a file tree, packing it,
// unpacking it, and verifying that the unpacked contents match expectations.
// We work in a new temporary directory for each test to avoid interaction between
// tests.
func TestPackUnpack(t *testing.T) {
	tests := loadAllTestCasesOrPanic(t, "testdata/pack-unpack")
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s (%s)", tt.Name, tt.Description), func(t *testing.T) {
			// Set up temporary directory for work
			tmpDir, removeTmp := setupTempDir(t)
			defer removeTmp()

			// Set up paths to use for test
			modelKitPath, unpackPath, contextPath := setupTestDirs(t, tmpDir)
			t.Setenv("KITOPS_HOME", contextPath)

			// Create Kitfile
			kitfilePath := filepath.Join(modelKitPath, constants.DefaultKitfileName)
			if err := os.WriteFile(kitfilePath, []byte(tt.Kitfile), 0644); err != nil {
				t.Fatal(err)
			}
			// Create .kitignore, if it exists
			if tt.Kitignore != "" {
				ignorePath := filepath.Join(modelKitPath, constants.IgnoreFileName)
				if err := os.WriteFile(ignorePath, []byte(tt.Kitignore), 0644); err != nil {
					t.Fatal(err)
				}
			}
			// Create files for test case
			setupFiles(t, modelKitPath, append(tt.Files, tt.IgnoredFiles...))

			runCommand(t, "pack", modelKitPath, "-t", modelKitTag, "-v")
			runCommand(t, "list")
			runCommand(t, "unpack", modelKitTag, "-d", unpackPath, "-v")

			checkFilesExist(t, unpackPath, tt.Files)
			checkFilesDoNotExist(t, unpackPath, append(tt.IgnoredFiles, ".kitignore"))
		})
	}
}
