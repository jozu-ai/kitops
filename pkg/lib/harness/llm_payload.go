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

//go:build embed_harness
// +build embed_harness

package harness

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

//go:embed llamafile*.tar.gz
var serverEmbed embed.FS

//go:embed ui.tar.gz
var uiEmbed embed.FS

func extractServer(harnessHome string) error {
	// Create the harnessHome directory once before extracting files
	if err := os.MkdirAll(harnessHome, os.FileMode(0755)); err != nil {
		return fmt.Errorf("error creating directory %s: %w", harnessHome, err)
	}
	if err := extractFile(serverEmbed, "llamafile.tar.gz", harnessHome); err != nil {
		return fmt.Errorf("error extracting file %s to %s: %w", "llamafile.tar.gz", harnessHome, err)
	}

	// Set executable permissions and rename on Windows
	llamaFilePath := filepath.Join(harnessHome, "llamafile")
	if runtime.GOOS == "windows" {
		llamaExePath := filepath.Join(harnessHome, "llamafile.exe")
		if err := os.Rename(llamaFilePath, llamaExePath); err != nil {
			return fmt.Errorf("error renaming file to executable: %w", err)
		}
	} else {
		if err := os.Chmod(llamaFilePath, os.FileMode(0755)); err != nil {
			return fmt.Errorf("error setting executable permission: %w", err)
		}
	}

	return nil
}

func extractUI(harnessHome string) error {
	uiHome := filepath.Join(harnessHome, "ui")
	if err := os.MkdirAll(uiHome, os.FileMode(0755)); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", uiHome, err)
	}
	return extractFile(uiEmbed, "ui.tar.gz", uiHome)
}
