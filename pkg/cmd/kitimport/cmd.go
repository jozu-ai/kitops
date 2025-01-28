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
	"kitops/pkg/lib/constants"
	"kitops/pkg/output"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

const (
	shortDesc = `Import a model from HuggingFace`
	longDesc  = `Download a repository from HuggingFace and package it as a ModelKit.

The repository can be specified either via a repository (e.g. myorg/myrepo) or
with a full URL (https://huggingface.co/myorg/myrepo). The repository will be
downloaded to a temporary directory and be packaged using a generated Kitfile.

In interactive settings, this command will read the EDITOR environment variable
to determine which editor should be used for editing the Kitfile.

Note: importing repositories requires 'git' and 'git-lfs' to be installed.`

	example = `# Download repository myorg/myrepo and package it, using the default tag (myorg/myrepo:latest)
kit import myorg/myrepo

# Download repository and tag it 'myrepository:mytag'
kit import myorg/myrepo --tag myrepository:mytag

# Download repository and pack it using an existing Kitfile
kit import myorg/myrepo --file ./path/to/Kitfile`
)

var repoToTagRegexp = regexp.MustCompile(`^.*?([0-9A-Za-z_-]+/[0-9A-Za-z_-]+)[^/]*$`)

type importOptions struct {
	configHome  string
	repo        string
	tag         string
	token       string
	kitfilePath string
}

func ImportCommand() *cobra.Command {
	opts := &importOptions{}

	cmd := &cobra.Command{
		Use:     "import [flags] REPOSITORY",
		Short:   shortDesc,
		Long:    longDesc,
		Example: example,
		RunE:    runCommand(opts),
		Args:    cobra.ExactArgs(1),
	}

	cmd.Flags().StringVar(&opts.token, "token", "", "Token to use for authenticating with repository")
	cmd.Flags().StringVarP(&opts.tag, "tag", "t", "", "Tag for the ModelKit (default is '[repository]:latest')")
	cmd.Flags().StringVarP(&opts.kitfilePath, "file", "f", "", "Path to Kitfile to use for packing")
	cmd.Flags().SortFlags = false
	return cmd
}

func runCommand(opts *importOptions) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if err := opts.complete(cmd.Context(), args); err != nil {
			return output.Fatalf("Invalid arguments: %s", err)
		}

		if err := doImport(cmd.Context(), opts); err != nil {
			return output.Fatalln(err)
		}

		return nil
	}
}

func (opts *importOptions) complete(ctx context.Context, args []string) error {
	configHome, ok := ctx.Value(constants.ConfigKey{}).(string)
	if !ok {
		return fmt.Errorf("default config path not set on command context")
	}
	opts.configHome = configHome
	opts.repo = args[0]

	if opts.tag == "" {
		tag := repoToTagRegexp.ReplaceAllString(opts.repo, "${1}")
		tag = strings.ToLower(tag)
		opts.tag = fmt.Sprintf("%s:latest", tag)
		output.Infof("Using tag %s. Use flag --tag to override", opts.tag)
	}
	return nil
}
