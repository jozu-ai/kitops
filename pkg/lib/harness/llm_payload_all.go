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
	"archive/tar"
	"compress/gzip"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"kitops/pkg/lib/filesystem"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

func extractServer(harnessHome string, glob string) error {
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

func extractUI(harnessHome string) error {
	uiHome := filepath.Join(harnessHome, "ui")
	if err := os.MkdirAll(uiHome, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", uiHome, err)
	}
	return extractFile(uiEmbed, "ui.tar.gz", uiHome)
}

func extractFile(fs embed.FS, file, harnessHome string) error {
	srcFile, err := fs.Open(file)
	if err != nil {
		return fmt.Errorf("read payload %s: %v", file, err)
	}
	defer srcFile.Close()

	srcReader := io.Reader(srcFile)
	if strings.HasSuffix(file, ".tar.gz") {
		gzr, err := gzip.NewReader(srcReader)
		if err != nil {
			return fmt.Errorf("error extracting gzipped file: %w", err)
		}
		defer gzr.Close()
		tarReader := tar.NewReader(gzr)
		return extractTar(tarReader, harnessHome)
	}

	if strings.HasSuffix(file, ".gz") {
		srcReader, err = gzip.NewReader(srcReader)
		if err != nil {
			return fmt.Errorf("failed to decompress payload %s: %v", file, err)
		}
		file = strings.TrimSuffix(file, ".gz")
	}

	destFile := filepath.Join(harnessHome, filepath.Base(file))
	dest, err := os.OpenFile(destFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o755) // Keep executable permissions
	if err != nil {
		return fmt.Errorf("write payload %s: %v", file, err)
	}
	defer dest.Close()

	if _, err := io.Copy(dest, srcReader); err != nil {
		return fmt.Errorf("copy payload %s: %v", file, err)
	}
	return nil
}

func extractTar(tr *tar.Reader, dir string) error {
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		outPath := filepath.Join(dir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if fi, exists := filesystem.PathExists(outPath); exists {
				if !fi.IsDir() {
					return fmt.Errorf("path '%s' already exists and is not a directory", outPath)
				}
			} else {
				if err := os.MkdirAll(outPath, header.FileInfo().Mode()); err != nil {
					return fmt.Errorf("failed to create directory %s: %w", outPath, err)
				}
			}

		case tar.TypeReg:
			if fi, exists := filesystem.PathExists(outPath); exists {
				if !fi.Mode().IsRegular() {
					return fmt.Errorf("path '%s' already exists and is not a regular file", outPath)
				}
			}
			file, err := os.OpenFile(outPath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, header.FileInfo().Mode())
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", outPath, err)
			}
			defer file.Close()

			written, err := io.Copy(file, tr)
			if err != nil {
				return fmt.Errorf("failed to write file %s: %w", outPath, err)
			}
			if written != header.Size {
				return fmt.Errorf("could not unpack file %s", outPath)
			}

		default:
			return fmt.Errorf("Unrecognized type in archive: %s", header.Name)
		}
	}
	return nil
}
