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
	"kitops/pkg/lib/git"
	kfutils "kitops/pkg/lib/kitfile"
	"kitops/pkg/lib/repo/local"
	repoutils "kitops/pkg/lib/repo/util"
	"kitops/pkg/output"

	"oras.land/oras-go/v2/registry"
)

func doImport(ctx context.Context, opts *importOptions) error {
	tmpDir, err := os.MkdirTemp("", "kitops_import_tmp")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			output.Logf(output.LogLevelWarn, "Failed to remove temporary directory %s", tmpDir)
		}
	}()

	parsedRef, _, err := repoutils.ParseReference(opts.tag)
	if err != nil {
		return err
	}

	if err := cloneRepository(opts.repo, tmpDir); err != nil {
		return err
	}

	// Check on the off-chance a Kitfile already exists
	var kitfile *artifact.KitFile
	if kfpath, err := filesystem.FindKitfileInPath(tmpDir); err != nil {
		kf, err := generateKitfile(tmpDir)
		if err != nil {
			return err
		}
		kitfile = kf
	} else {
		kf, err := readExistingKitfile(kfpath)
		if err != nil {
			return err
		}
		kitfile = kf
	}

	output.Infof("Packing model to %s", opts.tag)
	if err := os.Chdir(tmpDir); err != nil {
		return fmt.Errorf("failed to use context path %s: %w", tmpDir, err)
	}
	if err := packDirectory(ctx, opts.configHome, tmpDir, kitfile, parsedRef); err != nil {
		return fmt.Errorf("failed to pack ModelKit: %w", err)
	}
	output.Infof("Model is packed as %s", opts.tag)
	return nil
}

func cloneRepository(repo, destDir string) error {
	hfRepo := fmt.Sprintf("https://huggingface.co/%s", repo)
	if err := git.CloneRepository(hfRepo, destDir); err != nil {
		return err
	}
	// Clean up git-related files, since we probably don't want those
	if err := git.CleanGitMetadata(destDir); err != nil {
		return err
	}
	return nil
}

func generateKitfile(contextDir string) (*artifact.KitFile, error) {
	kitfile, err := kfutils.GenerateKitfile(contextDir, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Kitfile: %w", err)
	}
	kitfileBytes, err := kitfile.MarshalToYAML()
	if err != nil {
		return nil, fmt.Errorf("failed to write Kitfile: %w", err)
	}
	kitfilePath := filepath.Join(contextDir, constants.DefaultKitfileName)
	if err := os.WriteFile(kitfilePath, kitfileBytes, 0644); err != nil {
		return nil, fmt.Errorf("failed to write Kitfile: %s", err)
	}
	output.Infof("Generated Kitfile:\n\n%s\n", string(kitfileBytes))
	return kitfile, nil
}

func packDirectory(ctx context.Context, configHome, contextDir string, kitfile *artifact.KitFile, ref *registry.Reference) error {
	localRepo, err := local.NewLocalRepo(constants.StoragePath(configHome), ref)
	if err != nil {
		return err
	}
	ignore, err := filesystem.NewIgnoreFromContext(contextDir, kitfile)
	if err != nil {
		return err
	}
	manifestDesc, err := kfutils.SaveModel(ctx, localRepo, kitfile, ignore, constants.NoneCompression)
	if err != nil {
		return err
	}
	if err := localRepo.Tag(ctx, *manifestDesc, ref.Reference); err != nil {
		return fmt.Errorf("failed to tag manifest: %w", err)
	}
	return nil
}

func readExistingKitfile(kfPath string) (*artifact.KitFile, error) {
	kfFile, err := os.Open(kfPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read existing Kitfile: %w", err)
	}
	defer kfFile.Close()
	kitfile := &artifact.KitFile{}
	if err := kitfile.LoadModel(kfFile); err != nil {
		return nil, fmt.Errorf("failed to load existing Kitfile: %w", err)
	}
	if err := kfutils.ValidateKitfile(kitfile); err != nil {
		return nil, fmt.Errorf("existing Kitfile is invalid: %w", err)
	}
	return kitfile, nil
}
