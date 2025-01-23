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
	"os/exec"
	"path/filepath"
	"strings"

	"kitops/pkg/artifact"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/filesystem"
	"kitops/pkg/lib/git"
	kfutils "kitops/pkg/lib/kitfile"
	"kitops/pkg/lib/repo/local"
	repoutils "kitops/pkg/lib/repo/util"
	"kitops/pkg/lib/util"
	"kitops/pkg/output"

	"oras.land/oras-go/v2/registry"
)

func doImport(ctx context.Context, opts *importOptions) error {
	tmpDir, err := os.MkdirTemp("", "kitops_import_tmp")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	curDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	doCleanup := true
	defer func() {
		if doCleanup {
			// Make sure we're not in tmpDir before trying to remove it; on Windows we
			// cannot remove the current working directory.
			if err := os.Chdir(curDir); err != nil {
				output.Logf(output.LogLevelWarn, "Failed to change directory to %s: %s", curDir, err)
			}
			if err := os.RemoveAll(tmpDir); err != nil {
				output.Logf(output.LogLevelWarn, "Failed to remove temporary directory %s: %s", tmpDir, err)
			}
		}
	}()

	parsedRef, _, err := repoutils.ParseReference(opts.tag)
	if err != nil {
		return err
	}

	if err := cloneRepository(opts.repo, tmpDir, opts.token); err != nil {
		return err
	}

	// Check on the off-chance a Kitfile already exists
	var kitfile *artifact.KitFile
	if kfpath, err := filesystem.FindKitfileInPath(tmpDir); err != nil {
		// Fill fields in package so that they're not empty in `kit list` later.
		sections := strings.Split(opts.repo, "/")
		var modelPackage *artifact.Package
		if len(sections) >= 2 {
			modelPackage = &artifact.Package{
				Name:    sections[len(sections)-1],
				Authors: []string{sections[len(sections)-2]},
			}
		}

		kf, err := generateKitfile(tmpDir, modelPackage)
		if err != nil {
			return err
		}
		kitfile = kf

		if util.IsInteractiveSession() {
			// If we hit an error here, we don't want to clean up files so that user
			// can manually edit them.
			doCleanup = false
			newKitfile, err := promptToEditKitfile(tmpDir, opts.tag, kf)
			if err != nil {
				return err
			}
			kitfile = newKitfile
			doCleanup = true
		}
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

func cloneRepository(repo, destDir, token string) error {
	fullRepo := repo
	if !strings.HasPrefix(fullRepo, "http") {
		fullRepo = fmt.Sprintf("https://huggingface.co/%s", repo)
	}
	if err := git.CloneRepository(fullRepo, destDir, token); err != nil {
		return err
	}
	// Clean up git-related files, since we probably don't want those
	if err := git.CleanGitMetadata(destDir); err != nil {
		return err
	}
	return nil
}

func generateKitfile(contextDir string, modelPackage *artifact.Package) (*artifact.KitFile, error) {
	kitfile, err := kfutils.GenerateKitfile(contextDir, modelPackage)
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

func promptToEditKitfile(contextDir, tag string, currentKitfile *artifact.KitFile) (*artifact.KitFile, error) {
	kitfilePath := filepath.Join(contextDir, constants.DefaultKitfileName)
	ans, err := util.PromptForInput("Would you like to edit Kitfile before packing? (y/N): ", false)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}
	if !strings.HasPrefix(strings.ToLower(ans), "y") {
		// Current one is fine!
		return currentKitfile, nil
	}
	editor := os.Getenv("EDITOR")
	if editor == "" {
		output.Logf(output.LogLevelWarn, "Could not determine default editor from $EDITOR environment variable")
		output.Logf(output.LogLevelWarn, "Please manually edit Kitfile at path")
		output.Logf(output.LogLevelWarn, "    %s", kitfilePath)
		output.Logf(output.LogLevelWarn, "and run command")
		output.Logf(output.LogLevelWarn, "    kit pack -t %s %s", tag, contextDir)
		output.Logf(output.LogLevelWarn, "to complete process")
		return nil, fmt.Errorf("no editor found")
	}
	editCmd := exec.Command(editor, kitfilePath)
	editCmd.Stdin = os.Stdin
	editCmd.Stdout = os.Stdout
	editCmd.Stderr = os.Stderr
	if err := editCmd.Run(); err != nil {
		return nil, fmt.Errorf("error running external editor: %w", err)
	}
	return readExistingKitfile(kitfilePath)
}
