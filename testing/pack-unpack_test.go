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
	"time"

	"kitops/pkg/lib/constants"

	"github.com/stretchr/testify/assert"
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

func TestPackReproducibility(t *testing.T) {
	tmpDir, removeTmp := setupTempDir(t)
	defer removeTmp()

	modelKitPath, _, contextPath := setupTestDirs(t, tmpDir)
	t.Setenv("KITOPS_HOME", contextPath)

	testKitfile := `
manifestVersion: 1.0.0
package:
  name: test-repack
model:
  path: test-file.txt
dataset:
  - path: test-dir/test-subfile.txt
`
	kitfilePath := filepath.Join(modelKitPath, constants.DefaultKitfileName)
	if err := os.WriteFile(kitfilePath, []byte(testKitfile), 0644); err != nil {
		t.Fatal(err)
	}
	setupFiles(t, modelKitPath, []string{"test-file.txt", "test-dir/test-subfile.txt"})

	packOut := runCommand(t, "pack", modelKitPath, "-t", "test:repack1", "-v")
	digestOne := digestFromPack(t, packOut)

	// Change timestamps on file to simulate an unpacked modelkit at a future time
	futureTime := time.Now().Add(time.Hour)
	if err := os.Chtimes(filepath.Join(modelKitPath, "test-file.txt"), futureTime, futureTime); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(filepath.Join(modelKitPath, "test-dir"), futureTime, futureTime); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(filepath.Join(modelKitPath, "test-dir/test-subfile.txt"), futureTime, futureTime); err != nil {
		t.Fatal(err)
	}

	packOut = runCommand(t, "pack", modelKitPath, "-t", "test:repack2", "-v")
	digestTwo := digestFromPack(t, packOut)

	assert.Equal(t, digestOne, digestTwo, "Digests should be the same")
}
