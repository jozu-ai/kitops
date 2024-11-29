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

package kitfile

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"kitops/pkg/artifact"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/filesystem"
	"kitops/pkg/output"

	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// compressLayer compresses an *artifact.ModelLayer to a gzipped tar file. In order to return
// a descriptor (including hash) for the compressed file, the layer is saved to a temporary file
// on disk and must be moved to an appropriate location. It is the responsibility of the caller
// to clean up the temporary file when it is no longer needed.
func compressLayer(path string, mediaType constants.MediaType, ignore filesystem.IgnorePaths) (tempFilePath string, desc ocispec.Descriptor, layerInfo *artifact.LayerInfo, err error) {
	// Clean path to ensure consistent format (./path vs path/ vs path)
	path = filepath.Clean(path)

	if layerIgnored, err := ignore.Matches(path, path); err != nil {
		return "", ocispec.DescriptorEmptyJSON, nil, err
	} else if layerIgnored {
		output.Errorf("Warning: layer path %s ignored by kitignore", path)
	}

	pathInfo, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", ocispec.DescriptorEmptyJSON, nil, fmt.Errorf("%s path %s does not exist", mediaType.BaseType, path)
		}
		return "", ocispec.DescriptorEmptyJSON, nil, err
	}
	totalSize, err := getTotalSize(path, pathInfo, ignore)
	if err != nil {
		return "", ocispec.DescriptorEmptyJSON, nil, fmt.Errorf("failed to get size of layer: %w", err)
	}

	tempFile, err := os.CreateTemp("", "kitops_layer_*")
	if err != nil {
		return "", ocispec.DescriptorEmptyJSON, nil, fmt.Errorf("failed to create temporary file: %w", err)
	}
	tempFileName := tempFile.Name()
	output.Debugf("Compressing layer to temporary file %s", tempFileName)

	digester := digest.Canonical.Digester()
	mw := io.MultiWriter(tempFile, digester.Hash())

	var cw io.WriteCloser
	var tw *tar.Writer
	switch mediaType.Compression {
	case constants.GzipCompression:
		cw = gzip.NewWriter(mw)
		tw = tar.NewWriter(cw)
	case constants.GzipFastestCompression:
		cw, err = gzip.NewWriterLevel(mw, gzip.BestSpeed)
		if err != nil {
			return "", ocispec.DescriptorEmptyJSON, nil, fmt.Errorf("failed to set up gzip compression: %w", err)
		}
		tw = tar.NewWriter(cw)
	case constants.NoneCompression:
		tw = tar.NewWriter(mw)
	}
	ptw, plog := output.TarProgress(totalSize, tw)

	// Wrapper function for closing writers before returning an error
	// Note: we have to close gzip writer before reading digest from digester as closing is what writes the GZIP footer
	handleErr := func(err error) (string, ocispec.Descriptor, *artifact.LayerInfo, error) {
		// Don't care about these errors since we'll be deleting the file anyways
		_ = ptw.Close()
		_ = tw.Close()
		if cw != nil {
			_ = cw.Close()
		}
		_ = tempFile.Close()
		removeTempFile(tempFileName)
		return "", ocispec.DescriptorEmptyJSON, nil, err
	}

	if pathInfo.Mode().IsRegular() {
		if err := writeHeaderToTar(pathInfo.Name(), pathInfo, ptw, plog); err != nil {
			return handleErr(err)
		}
		if err := writeFileToTar(path, pathInfo, ptw, plog); err != nil {
			return handleErr(err)
		}
	} else if pathInfo.IsDir() {
		if err := writeDirToTar(path, ignore, ptw, plog); err != nil {
			return handleErr(err)
		}
	} else {
		return handleErr(fmt.Errorf("path %s is neither a file nor a directory", path))
	}
	plog.Wait()

	callAndPrintError(ptw.Close, "Failed to close writer: %s")
	callAndPrintError(tw.Close, "Failed to close tar writer: %s")
	if cw != nil {
		callAndPrintError(cw.Close, "Failed to close compression writer: %s")
	}

	tempFileInfo, err := tempFile.Stat()
	if err != nil {
		removeTempFile(tempFileName)
		return "", ocispec.DescriptorEmptyJSON, nil, fmt.Errorf("failed to stat temporary file: %w", err)
	}
	callAndPrintError(tempFile.Close, "Failed to close temporary file: %s")

	desc = ocispec.Descriptor{
		MediaType: mediaType.String(),
		Digest:    digester.Digest(),
		Size:      tempFileInfo.Size(),
	}
	return tempFileName, desc, nil, nil
}

// writeDirToTar walks the filesystem at basePath, compressing contents via the *tar.Writer.
// Any non-regular files and directories (e.g. symlinks) are skipped.
func writeDirToTar(basePath string, ignore filesystem.IgnorePaths, ptw *output.ProgressTar, plog *output.ProgressLogger) error {
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
				plog.Debugf("Skipping directory %s: ignored", file)
				return filepath.SkipDir
			}
			plog.Debugf("Skipping file %s: ignored", file)
			return nil
		}

		relPath, err := filepath.Rel(trimPath, file)
		if err != nil {
			return fmt.Errorf("failed to find relative path for %s", file)
		}

		if err := writeHeaderToTar(relPath, fi, ptw, plog); err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		return writeFileToTar(file, fi, ptw, plog)
	})
}

func writeHeaderToTar(name string, fi os.FileInfo, ptw *output.ProgressTar, plog *output.ProgressLogger) error {
	header, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		return fmt.Errorf("failed to generate header for %s: %w", name, err)
	}
	header.Name = name
	sanitizeTarHeader(header)
	if err := ptw.WriteHeader(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	plog.Debugf("Wrote header %s to tar file", header.Name)
	return nil
}

func writeFileToTar(file string, fi os.FileInfo, ptw *output.ProgressTar, plog *output.ProgressLogger) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open file for archiving: %w", err)
	}
	defer f.Close()

	if written, err := io.Copy(ptw, f); err != nil {
		return fmt.Errorf("failed to add file to archive: %w", err)
	} else if written != fi.Size() {
		return fmt.Errorf("error writing file: %w", err)
	}
	plog.Debugf("Wrote file %s to tar file", file)
	return nil
}

func getTotalSize(basePath string, pathInfo fs.FileInfo, ignore filesystem.IgnorePaths) (int64, error) {
	if !output.ProgressEnabled() {
		// Won't use this information anyways, save the work.
		return 0, nil
	}
	if pathInfo.Mode().IsRegular() {
		return pathInfo.Size(), nil
	} else if pathInfo.IsDir() {
		var total int64
		err := filepath.WalkDir(basePath, func(file string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if shouldIgnore, err := ignore.Matches(file, basePath); err != nil {
				return fmt.Errorf("failed to match %s against ignore file: %w", file, err)
			} else if shouldIgnore {
				if !ignore.HasExclusions() && d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if d.Type().IsRegular() {
				fi, err := d.Info()
				if err != nil {
					return fmt.Errorf("failed to stat %s: %w", file, err)
				}
				total += fi.Size()
			}
			return nil
		})
		if err != nil {
			return 0, err
		}
		return total, nil
	} else {
		return 0, fmt.Errorf("path %s is neither a file nor a directory", basePath)
	}
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

func sanitizeTarHeader(header *tar.Header) {
	// On windows, store paths linux-style (forward slashes). This is a no-op if
	// filepath.Separator is '/'
	header.Name = filepath.ToSlash(header.Name)
	// Clear fields that break reproducible tars
	header.AccessTime = time.Time{}
	header.ModTime = time.Time{}
	header.ChangeTime = time.Time{}
	header.Uid = 0
	header.Gid = 0
	header.Uname = ""
	header.Gname = ""
}
