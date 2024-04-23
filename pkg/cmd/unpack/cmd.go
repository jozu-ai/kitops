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
	"context"
	"fmt"
	"kitops/pkg/cmd/options"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/repo"
	"kitops/pkg/output"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"oras.land/oras-go/v2/registry"
)

const (
	shortDesc = `Produce the components from a modelkit on the local filesystem`
	longDesc  = `Produces all or selected components of a modelkit on the local filesystem.

This command unpacks a modelkit's components, including models, code,
datasets, and configuration files, to a specified directory on the local
filesystem. By default, it attempts to find the modelkit in local storage; if
not found, it searches the remote registry and retrieves it. This process
ensures that the necessary components are always available for unpacking,
optimizing for efficiency by fetching only specified components from the
remote registry when necessary`

	example = `# Unpack all components of a modelkit to the current directory
kit unpack myrepo/my-model:latest -d /path/to/unpacked

# Unpack only the model and datasets of a modelkit to a specified directory
kit unpack myrepo/my-model:latest --model --datasets -d /path/to/unpacked

# Unpack a modelkit from a remote registry with overwrite enabled
kit unpack registry.example.com/myrepo/my-model:latest -o -d /path/to/unpacked`
)

type unpackOptions struct {
	options.NetworkOptions
	configHome string
	unpackDir  string
	unpackConf unpackConf
	modelRef   *registry.Reference
	overwrite  bool
}

type unpackConf struct {
	unpackKitfile  bool
	unpackModels   bool
	unpackCode     bool
	unpackDatasets bool
}

func (opts *unpackOptions) complete(ctx context.Context, args []string) error {
	configHome, ok := ctx.Value(constants.ConfigKey{}).(string)
	if !ok {
		return fmt.Errorf("default config path not set on command context")
	}
	opts.configHome = configHome
	modelRef, extraTags, err := repo.ParseReference(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse reference: %w", err)
	}
	if len(extraTags) > 0 {
		return fmt.Errorf("can not unpack multiple tags")
	}
	opts.modelRef = modelRef

	conf := opts.unpackConf
	if !conf.unpackKitfile && !conf.unpackModels && !conf.unpackCode && !conf.unpackDatasets {
		opts.unpackConf.unpackKitfile = true
		opts.unpackConf.unpackModels = true
		opts.unpackConf.unpackCode = true
		opts.unpackConf.unpackDatasets = true
	}

	absDir, err := filepath.Abs(opts.unpackDir)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path %s: %w", opts.unpackDir, err)
	}
	opts.unpackDir = absDir

	printConfig(opts)
	return nil
}

func UnpackCommand() *cobra.Command {
	opts := &unpackOptions{}

	cmd := &cobra.Command{
		Use:     "unpack [flags] [registry/]repository[:tag|@digest]",
		Short:   shortDesc,
		Long:    longDesc,
		Example: example,
		Run:     runCommand(opts),
	}

	cmd.Args = cobra.ExactArgs(1)
	cmd.Flags().StringVarP(&opts.unpackDir, "dir", "d", "", "The target directory to unpack components into. This directory will be created if it does not exist")
	cmd.Flags().BoolVarP(&opts.overwrite, "overwrite", "o", false, "Overwrites existing files and directories in the target unpack directory without prompting")
	cmd.Flags().BoolVar(&opts.unpackConf.unpackKitfile, "kitfile", false, "Unpack only Kitfile")
	cmd.Flags().BoolVar(&opts.unpackConf.unpackModels, "model", false, "Unpack only model")
	cmd.Flags().BoolVar(&opts.unpackConf.unpackCode, "code", false, "Unpack only code")
	cmd.Flags().BoolVar(&opts.unpackConf.unpackDatasets, "datasets", false, "Unpack only datasets")
	opts.AddNetworkFlags(cmd)

	return cmd
}

func runCommand(opts *unpackOptions) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := opts.complete(cmd.Context(), args); err != nil {
			output.Fatalf("Invalid arguments: %s", err)
		}

		if opts.modelRef.Reference == "" {
			output.Fatalf("Invalid reference: unpacking requires a tag or digest")
		}

		unpackTo := opts.unpackDir
		if unpackTo == "" {
			unpackTo = "current directory"
		}
		// Make sure target directory exists, in case user is using the -d flag
		if err := os.MkdirAll(opts.unpackDir, 0755); err != nil {
			output.Fatalf("failed to create directory %s: %w", opts.unpackDir, err)
		}
		// Change working directory to context path to make sure relative paths within
		// tarballs are correct. This is the equivalent of using the -C parameter for tar
		if err := os.Chdir(opts.unpackDir); err != nil {
			output.Fatalf("Failed to use unpack path %s: %w", opts.unpackDir, err)
		}

		output.Infof("Unpacking to %s", unpackTo)
		err := runUnpack(cmd.Context(), opts)
		if err != nil {
			output.Fatalln(err)
		}
	}
}

func printConfig(opts *unpackOptions) {
	output.Debugf("Overwrite: %t", opts.overwrite)
	output.Debugf("Unpacking %s", opts.modelRef.String())
}
