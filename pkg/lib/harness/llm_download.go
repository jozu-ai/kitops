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

//go:build !embed_harness
// +build !embed_harness

package harness

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

const (
	llamafileDownloadURL = "http://downloads.jozu.ml?file=llamafile.tar.gz&version=" + LlamaFileVersion
	uiDownloadURL        = "http://downloads.jozu.ml?file=ui.tar.gz&version=" + LlamaFileVersion
)

func extractServer(harnessHome string) error {
	if err := os.MkdirAll(harnessHome, 0o755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", harnessHome, err)
	}
	tmpFolder := filepath.Join(harnessHome, "tmp")

	err := downloadFile(llamafileDownloadURL, tmpFolder, "llamafile.tar.gz")
	if err != nil {
		return fmt.Errorf("failed to extract llamafile: %w", err)
	}
	localFS := os.DirFS(tmpFolder)

	err = extractFile(localFS, "llamafile.tar.gz", harnessHome)
	if err != nil {
		return fmt.Errorf("failed to unpack llamafile: %w", err)
	}

	llamaFilePath := filepath.Join(harnessHome, "llamafile")
	if runtime.GOOS == "windows" {
		llamaExePath := filepath.Join(harnessHome, "llamafile.exe")
		if err := os.Rename(llamaFilePath, llamaExePath); err != nil {
			return fmt.Errorf("error renaming file to executable: %w", err)
		}
	} else {
		if err := os.Chmod(llamaFilePath, 0o755); err != nil {
			return fmt.Errorf("error setting executable permission: %w", err)
		}
	}

	return nil
}

func extractUI(harnessHome string) error {
	tmpFolder := filepath.Join(harnessHome, "tmp")

	err := downloadFile(uiDownloadURL, tmpFolder, "ui.tar.gz")
	if err != nil {
		return fmt.Errorf("failed to extract UI: %w", err)
	}

	uiHome := filepath.Join(harnessHome, "ui")
	if err := os.MkdirAll(uiHome, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", uiHome, err)
	}
	localFS := os.DirFS(tmpFolder)

	return extractFile(localFS, "ui.tar.gz", uiHome)
}

func downloadFile(url string, folder string, filename string) error {

	err := os.MkdirAll(folder, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create folder %s: %w", folder, err)
	}

	filePath := filepath.Join(folder, filename)

	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", folder, err)
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download from url %s failed: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status downloading file %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filePath, err)
	}

	return nil

}
