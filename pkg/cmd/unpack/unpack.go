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

package unpack

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"kitops/pkg/artifact"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/filesystem"
	"kitops/pkg/lib/repo/util"
	"kitops/pkg/output"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/content"
)

// runUnpack fetches and unpacks a *registry.Reference from an oras.Target. It returns an error if
// unpacking fails, or if any path specified in the modelkit is not a subdirectory of the current
// unpack target directory.
func runUnpack(ctx context.Context, opts *unpackOptions) error {
	return runUnpackRecursive(ctx, opts, []string{})
}

func runUnpackRecursive(ctx context.Context, opts *unpackOptions, visitedRefs []string) error {
	if len(visitedRefs) > constants.MaxModelRefChain {
		return fmt.Errorf("reached maximum number of model references: [%s]", strings.Join(visitedRefs, "=>"))
	}

	ref := opts.modelRef
	store, err := getStoreForRef(ctx, opts)
	if err != nil {
		ref := util.FormatRepositoryForDisplay(opts.modelRef.String())
		return fmt.Errorf("failed to find reference %s: %s", ref, err)
	}
	manifestDesc, err := store.Resolve(ctx, ref.Reference)
	if err != nil {
		return fmt.Errorf("failed to resolve reference: %w", err)
	}
	manifest, config, err := util.GetManifestAndConfig(ctx, store, manifestDesc)
	if err != nil {
		return fmt.Errorf("failed to read model: %s", err)
	}
	if config.Model != nil && util.IsModelKitReference(config.Model.Path) {
		output.Infof("Unpacking referenced modelkit %s", config.Model.Path)
		if err := unpackParent(ctx, config.Model.Path, opts, visitedRefs); err != nil {
			return err
		}
	}

	if shouldUnpackLayer(config, opts.filterConfs) {
		if err := unpackConfig(config, opts.unpackDir, opts.overwrite); err != nil {
			return err
		}
	}

	// Since there might be multiple datasets, etc. we need to synchronously iterate
	// through the config's relevant field to get the correct path for unpacking
	var modelPartIdx, codeIdx, datasetIdx, docsIdx int
	for _, layerDesc := range manifest.Layers {
		var relPath string
		mediaType := constants.ParseMediaType(layerDesc.MediaType)
		switch mediaType.BaseType {
		case constants.ModelType:
			if !shouldUnpackLayer(config.Model, opts.filterConfs) {
				continue
			}
			_, relPath, err = filesystem.VerifySubpath(opts.unpackDir, config.Model.Path)
			if err != nil {
				return fmt.Errorf("error resolving model path: %w", err)
			}
			output.Infof("Unpacking model %s to %s", config.Model.Name, relPath)

		case constants.ModelPartType:
			part := config.Model.Parts[modelPartIdx]
			if !shouldUnpackLayer(part, opts.filterConfs) {
				modelPartIdx += 1
				continue
			}
			_, relPath, err = filesystem.VerifySubpath(opts.unpackDir, part.Path)
			if err != nil {
				return fmt.Errorf("error resolving code path: %w", err)
			}
			output.Infof("Unpacking model part %s to %s", part.Name, relPath)
			modelPartIdx += 1

		case constants.CodeType:
			codeEntry := config.Code[codeIdx]
			if !shouldUnpackLayer(codeEntry, opts.filterConfs) {
				codeIdx += 1
				continue
			}
			_, relPath, err = filesystem.VerifySubpath(opts.unpackDir, codeEntry.Path)
			if err != nil {
				return fmt.Errorf("error resolving code path: %w", err)
			}
			output.Infof("Unpacking code to %s", relPath)
			codeIdx += 1

		case constants.DatasetType:
			datasetEntry := config.DataSets[datasetIdx]
			if !shouldUnpackLayer(datasetEntry, opts.filterConfs) {
				datasetIdx += 1
				continue
			}
			_, relPath, err = filesystem.VerifySubpath(opts.unpackDir, datasetEntry.Path)
			if err != nil {
				return fmt.Errorf("error resolving dataset path for dataset %s: %w", datasetEntry.Name, err)
			}
			output.Infof("Unpacking dataset %s to %s", datasetEntry.Name, relPath)
			datasetIdx += 1

		case constants.DocsType:
			docsEntry := config.Docs[docsIdx]
			if !shouldUnpackLayer(docsEntry, opts.filterConfs) {
				docsIdx += 1
				continue
			}
			_, relPath, err = filesystem.VerifySubpath(opts.unpackDir, docsEntry.Path)
			if err != nil {
				return fmt.Errorf("error resolving path %s for docs: %w", docsEntry.Path, err)
			}
			output.Infof("Unpacking docs to %s", docsEntry.Path)
			docsIdx += 1
		}

		if err := unpackLayer(ctx, store, layerDesc, relPath, opts.overwrite, mediaType.Compression); err != nil {
			return fmt.Errorf("failed to unpack: %w", err)
		}
	}
	output.Debugf("Unpacked %d model part layers", modelPartIdx)
	output.Debugf("Unpacked %d code layers", codeIdx)
	output.Debugf("Unpacked %d dataset layers", datasetIdx)
	output.Debugf("Unpacked %d docs layers", docsIdx)

	return nil
}

func unpackParent(ctx context.Context, ref string, optsIn *unpackOptions, visitedRefs []string) error {
	if idx := getIndex(visitedRefs, ref); idx != -1 {
		cycleStr := fmt.Sprintf("[%s=>%s]", strings.Join(visitedRefs[idx:], "=>"), ref)
		return fmt.Errorf("found cycle in modelkit references: %s", cycleStr)
	}

	parentRef, _, err := util.ParseReference(ref)
	if err != nil {
		return err
	}
	opts := *optsIn
	opts.modelRef = parentRef
	// Unpack only model, ignore code/datasets
	modelFilter, err := parseFilter("model")
	if err != nil {
		// Shouldn't happen, ever
		return fmt.Errorf("failed to parse filter for parent modelkit: %w", err)
	}
	opts.filterConfs = []filterConf{*modelFilter}

	return runUnpackRecursive(ctx, &opts, append(visitedRefs, ref))
}

func unpackConfig(config *artifact.KitFile, unpackDir string, overwrite bool) error {
	configPath := filepath.Join(unpackDir, constants.DefaultKitfileName)
	if fi, exists := filesystem.PathExists(configPath); exists {
		if !overwrite {
			return fmt.Errorf("failed to unpack config: path %s already exists", configPath)
		} else if !fi.Mode().IsRegular() {
			return fmt.Errorf("failed to unpack config: path %s exists and is not a regular file", configPath)
		}
	}

	configBytes, err := config.MarshalToYAML()
	if err != nil {
		return fmt.Errorf("failed to unpack config: %w", err)
	}

	output.Infof("Unpacking config to %s", configPath)
	if err := os.WriteFile(configPath, configBytes, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

func unpackLayer(ctx context.Context, store content.Storage, desc ocispec.Descriptor, unpackPath string, overwrite bool, compression string) error {
	rc, err := store.Fetch(ctx, desc)
	if err != nil {
		return fmt.Errorf("failed get layer %s: %w", desc.Digest, err)
	}
	var logger *output.ProgressLogger
	rc, logger = output.WrapReadCloser(desc.Size, rc)
	defer rc.Close()

	var cr io.ReadCloser
	var cErr error
	switch compression {
	case constants.GzipCompression, constants.GzipFastestCompression:
		cr, cErr = gzip.NewReader(rc)
	case constants.NoneCompression:
		cr = rc
	}
	if cErr != nil {
		return fmt.Errorf("error setting up decompress: %w", cErr)
	}
	defer cr.Close()
	tr := tar.NewReader(cr)

	unpackDir := filepath.Dir(unpackPath)
	if err := os.MkdirAll(unpackDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", unpackDir, err)
	}

	if err := extractTar(tr, unpackDir, overwrite, logger); err != nil {
		return err
	}
	logger.Wait()
	return nil
}

func extractTar(tr *tar.Reader, dir string, overwrite bool, logger *output.ProgressLogger) error {
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
				logger.Debugf("Path %s already exists", outPath)
			} else {
				logger.Debugf("Creating directory %s", outPath)
				if err := os.MkdirAll(outPath, header.FileInfo().Mode()); err != nil {
					return fmt.Errorf("failed to create directory %s: %w", outPath, err)
				}
			}

		case tar.TypeReg:
			if fi, exists := filesystem.PathExists(outPath); exists {
				if !overwrite {
					return fmt.Errorf("path '%s' already exists", outPath)
				}
				if !fi.Mode().IsRegular() {
					return fmt.Errorf("path '%s' already exists and is not a regular file", outPath)
				}
			}
			logger.Debugf("Unpacking file %s", outPath)
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
			return fmt.Errorf("unrecognized type in archive: %s", header.Name)
		}
	}
	return nil
}

func getIndex(list []string, s string) int {
	for idx, item := range list {
		if s == item {
			return idx
		}
	}
	return -1
}
