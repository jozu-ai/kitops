// Copyright 2025 The KitOps Authors.
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

package kitimport

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"kitops/pkg/artifact"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/filesystem"
	"kitops/pkg/lib/hf"
	kfutils "kitops/pkg/lib/kitfile"
	kfgen "kitops/pkg/lib/kitfile/generate"
	repoutil "kitops/pkg/lib/repo/util"
	"kitops/pkg/lib/util"
	"kitops/pkg/output"
)

func importUsingHF(ctx context.Context, opts *importOptions) error {
	// TODO: Handle full huggingface URLs

	tmpDir, err := os.MkdirTemp("", "kitops_import_tmp")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	doCleanup := true
	defer func() {
		if doCleanup {
			if err := os.RemoveAll(tmpDir); err != nil {
				output.Logf(output.LogLevelWarn, "Failed to remove temporary directory %s: %s", tmpDir, err)
			}
		}
	}()

	dirListing, err := hf.ListFiles(ctx, opts.repo)
	if err != nil {
		return fmt.Errorf("failed to list files from HuggingFace API: %w", err)
	}

	var kitfile *artifact.KitFile
	if opts.kitfilePath == "-" {
		kitfile = &artifact.KitFile{}
		if err := kitfile.LoadModel(os.Stdin); err != nil {
			return fmt.Errorf("failed to read Kitfile from input: %w", err)
		}
		if err := kfutils.ValidateKitfile(kitfile); err != nil {
			return err
		}
	} else if opts.kitfilePath != "" {
		kf, err := readExistingKitfile(opts.kitfilePath)
		if err != nil {
			return err
		}
		kitfile = kf
	} else {
		kf, err := generateKitfile(dirListing, opts.repo, tmpDir)
		if err != nil {
			return err
		}
		kitfile = kf

		if util.IsInteractiveSession() {
			// If we hit an error here, we don't want to clean up files so that user
			// can manually edit them.
			doCleanup = false
			newKitfile, err := promptToEditKitfile(tmpDir, kf)
			if err != nil {
				kfPath := filepath.Join(tmpDir, constants.DefaultKitfileName)
				output.Logf(output.LogLevelWarn, "Could not determine default editor from $EDITOR environment variable")
				output.Logf(output.LogLevelWarn, "Please manually edit Kitfile at path")
				output.Logf(output.LogLevelWarn, "    %s", kfPath)
				output.Logf(output.LogLevelWarn, "and run command")
				output.Logf(output.LogLevelWarn, "    kit import %s -t %s -f %s", opts.repo, opts.tag, kfPath)
				output.Logf(output.LogLevelWarn, "to complete process")
				return err
			}
			kitfile = newKitfile
			doCleanup = true
		}
	}

	toDownload, err := filterListingForKitfile(dirListing, kitfile)
	if err != nil {
		return err
	}
	if err := hf.DownloadFiles(ctx, opts.repo, tmpDir, toDownload); err != nil {
		return fmt.Errorf("error downloading repository: %w", err)
	}

	output.Infof("Packing model to %s", opts.tag)
	if err := packDirectory(ctx, opts.configHome, tmpDir, kitfile, opts.modelKitRef); err != nil {
		return fmt.Errorf("failed to pack ModelKit: %w", err)
	}
	output.Infof("Model is packed as %s", opts.tag)
	return nil
}

func filterListingForKitfile(contents *kfgen.DirectoryListing, kitfile *artifact.KitFile) ([]string, error) {
	// Repurpose the ignore implementation to find which files we need to download and which ones we can skip.
	// This works because ignore is designed to _also_ ignore paths that are packed as part of another layer
	// instead of the current one.
	ignore, err := filesystem.NewIgnore(nil, kitfile)
	if err != nil {
		return nil, fmt.Errorf("failed to process Kitfile to get file list: %w", err)
	}

	hasCatchall := kitfileHasCatchallLayer(kitfile)
	var pathsToDownload []string
	var processDir func(dir *kfgen.DirectoryListing) error
	processDir = func(dir *kfgen.DirectoryListing) error {
		for _, file := range dir.Files {
			if hasCatchall {
				pathsToDownload = append(pathsToDownload, file.Path)
				continue
			}
			matches, err := ignore.Matches(file.Path, "")
			if err != nil {
				return fmt.Errorf("failed to process path %s: %w", file.Path, err)
			}
			if matches {
				pathsToDownload = append(pathsToDownload, file.Path)
			}
		}
		for _, subDir := range dir.Subdirs {
			if err := processDir(&subDir); err != nil {
				return err
			}
		}
		return nil
	}
	if err := processDir(contents); err != nil {
		return nil, err
	}

	return pathsToDownload, nil
}

func kitfileHasCatchallLayer(kitfile *artifact.KitFile) bool {
	layerPaths := repoutil.LayerPathsFromKitfile(kitfile)
	for _, path := range layerPaths {
		if path == "." {
			return true
		}
	}
	return false
}
