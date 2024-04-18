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

package storage

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"kitops/pkg/lib/filesystem"
	"kitops/pkg/output"
	"os"
	"path/filepath"
	"time"

	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// compressLayer compresses an *artifact.ModelLayer to a gzipped tar file. In order to return
// a descriptor (including hash) for the compressed file, the layer is saved to a temporary file
// on disk and must be moved to an appropriate location. It is the responsibility of the caller
// to clean up the temporary file when it is no longer needed.
func compressLayer(path, mediaType string, ignore filesystem.IgnorePaths) (tempFilePath string, desc ocispec.Descriptor, err error) {
	// Clean path to ensure consistent format (./path vs path/ vs path)
	path = filepath.Clean(path)

	if layerIgnored, err := ignore.Matches(path, path); err != nil {
		return "", ocispec.DescriptorEmptyJSON, err
	} else if layerIgnored {
		output.Errorf("Warning: layer path %s ignored by kitignore", path)
	}

	pathInfo, err := os.Stat(path)
	if err != nil {
		return "", ocispec.DescriptorEmptyJSON, err
	}
	tempFile, err := os.CreateTemp("", "kitops_layer_*")
	if err != nil {
		return "", ocispec.DescriptorEmptyJSON, fmt.Errorf("failed to create temporary file: %w", err)
	}
	tempFileName := tempFile.Name()
	output.Debugf("Compressing layer to temporary file %s", tempFileName)

	digester := digest.Canonical.Digester()
	mw := io.MultiWriter(tempFile, digester.Hash())

	// Note: we have to close gzip writer before reading digest from digester as closing is what writes the GZIP footer
	gzw := gzip.NewWriter(mw)
	tw := tar.NewWriter(gzw)

	// Wrapper function for closing writers before returning an error
	handleErr := func(err error) (string, ocispec.Descriptor, error) {
		// Don't care about these errors since we'll be deleting the file anyways
		_ = tw.Close()
		_ = gzw.Close()
		_ = tempFile.Close()
		removeTempFile(tempFileName)
		return "", ocispec.DescriptorEmptyJSON, err
	}

	if pathInfo.Mode().IsRegular() {
		if err := writeHeaderToTar(pathInfo.Name(), pathInfo, tw); err != nil {
			return handleErr(err)
		}
		if err := writeFileToTar(path, pathInfo, tw); err != nil {
			return handleErr(err)
		}
	} else if pathInfo.IsDir() {
		if err := writeDirToTar(path, ignore, tw); err != nil {
			return handleErr(err)
		}
	} else {
		return handleErr(fmt.Errorf("path %s is neither a file nor a directory", path))
	}

	callAndPrintError(tw.Close, "Failed to close tar writer: %s")
	callAndPrintError(gzw.Close, "Failed to close gzip writer: %s")

	tempFileInfo, err := tempFile.Stat()
	if err != nil {
		removeTempFile(tempFileName)
		return "", ocispec.DescriptorEmptyJSON, fmt.Errorf("failed to stat temporary file: %w", err)
	}
	callAndPrintError(tempFile.Close, "Failed to close temporary file: %s")

	desc = ocispec.Descriptor{
		MediaType: mediaType,
		Digest:    digester.Digest(),
		Size:      tempFileInfo.Size(),
	}
	return tempFileName, desc, nil
}

// writeDirToTar walks the filesystem at basePath, compressing contents via the *tar.Writer.
// Any non-regular files and directories (e.g. symlinks) are skipped.
func writeDirToTar(basePath string, ignore filesystem.IgnorePaths, tw *tar.Writer) error {
	// We'll want paths in the tarball to be relative to the *parent* of basePath since we want
	// to compress the directory pointed at by basePath
	trimPath := filepath.Dir(basePath)
	if trimPath == "." {
		// Avoid accidentally trimming leading `.` from filenames
		trimPath = ""
	}
	return filepath.Walk(basePath, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip adding an entry for the context directory to the tarball
		if file == "." {
			return nil
		}
		// Skip anything that's not a regular file or directory
		if !fi.Mode().IsRegular() && !fi.Mode().IsDir() {
			return nil
		}

		if shouldIgnore, err := ignore.Matches(file, basePath); err != nil {
			return fmt.Errorf("failed to match %s against ignore file: %w", file, err)
		} else if shouldIgnore {
			if !ignore.HasExclusions() && fi.IsDir() {
				output.Debugf("Skipping directory %s: ignored", file)
				return filepath.SkipDir
			}
			output.Debugf("Skipping file %s: ignored", file)
			return nil
		}

		relPath, err := filepath.Rel(trimPath, file)
		if err != nil {
			return fmt.Errorf("failed to find relative path for %s", file)
		}

		if err := writeHeaderToTar(relPath, fi, tw); err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		return writeFileToTar(file, fi, tw)
	})
}

func writeHeaderToTar(name string, fi os.FileInfo, tw *tar.Writer) error {
	header, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		return fmt.Errorf("failed to generate header for %s: %w", name, err)
	}
	header.Name = name
	header.AccessTime = time.Time{}
	header.ModTime = time.Time{}
	header.ChangeTime = time.Time{}
	if err := tw.WriteHeader(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	output.Debugf("Wrote header %s to tar file", header.Name)
	return nil
}

func writeFileToTar(file string, fi os.FileInfo, tw *tar.Writer) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open file for archiving: %w", err)
	}
	defer f.Close()

	if written, err := io.Copy(tw, f); err != nil {
		return fmt.Errorf("failed to add file to archive: %w", err)
	} else if written != fi.Size() {
		return fmt.Errorf("error writing file: %w", err)
	}
	output.Debugf("Wrote file %s to tar file", file)
	return nil
}

// callAndPrintError is a wrapper to print an error message for a function that
// may return an error. The error is printed and then discarded.
func callAndPrintError(f func() error, msg string) {
	if err := f(); err != nil {
		output.Errorf(msg, err)
	}
}

func removeTempFile(filepath string) {
	if err := os.Remove(filepath); err != nil && !os.IsNotExist(err) {
		output.Errorf("Failed to clean up temporary file %s: %s", filepath, err)
	}
}
