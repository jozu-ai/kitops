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

package harness

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

func extractLibraries(harnessHome string, glob string) error {
	files, err := fs.Glob(serverEmbed, glob)
	if err != nil {
		return fmt.Errorf("error globbing files: %w", err)
	} else if len(files) == 0 {
		return fmt.Errorf("no files matched the glob pattern")
	}
	// Create the harnessHome directory once before extracting files
	if err := os.MkdirAll(harnessHome, 0o755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", harnessHome, err)
	}

	g := new(errgroup.Group)
	for _, file := range files {

		file := file
		g.Go(func() error {
			return extractFile(serverEmbed, file, harnessHome)
		})

	}

	return g.Wait()
}

func extractFile(fs embed.FS, file, harnessHome string) error {
	srcFile, err := fs.Open(file)
	if err != nil {
		return fmt.Errorf("read payload %s: %v", file, err)
	}
	defer srcFile.Close()

	destFile := filepath.Join(harnessHome, filepath.Base(file))
	dest, err := os.OpenFile(destFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o755) // Keep executable permissions
	if err != nil {
		return fmt.Errorf("write payload %s: %v", file, err)
	}
	defer dest.Close()

	if _, err := io.Copy(dest, srcFile); err != nil {
		return fmt.Errorf("copy payload %s: %v", file, err)
	}
	return nil
}
