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

package pull

import (
	"context"
	"fmt"
	"kitops/pkg/cmd/options"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/repo"
	"kitops/pkg/output"

	"github.com/spf13/cobra"
	"oras.land/oras-go/v2/registry"
)

const (
	shortDesc = `Retrieve modelkits from a remote registry to your local environment.`
	longDesc  = `Downloads modelkits from a specified registry. The downloaded modelkits
are stored in the local registry.`

	example = `# Pull the latest version of a modelkit from a remote registry
kit pull registry.example.com/my-model:latest`
)

type pullOptions struct {
	options.NetworkOptions
	configHome string
	modelRef   *registry.Reference
}

func (opts *pullOptions) complete(ctx context.Context, args []string) error {
	configHome, ok := ctx.Value(constants.ConfigKey{}).(string)
	if !ok {
		return fmt.Errorf("default config path not set on command context")
	}
	opts.configHome = configHome

	modelRef, extraTags, err := repo.ParseReference(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse reference: %w", err)
	}
	if modelRef.Registry == "localhost" {
		return fmt.Errorf("registry is required when pulling")
	}
	if len(extraTags) > 0 {
		return fmt.Errorf("reference cannot include multiple tags")
	}
	opts.modelRef = modelRef

	return nil
}

func PullCommand() *cobra.Command {
	opts := &pullOptions{}
	cmd := &cobra.Command{
		Use:     "pull [flags] registry/repository[:tag|@digest]",
		Short:   shortDesc,
		Long:    longDesc,
		Example: example,
		Run:     runCommand(opts),
	}

	cmd.Args = cobra.ExactArgs(1)
	opts.AddNetworkFlags(cmd)

	return cmd
}

func runCommand(opts *pullOptions) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := opts.complete(cmd.Context(), args); err != nil {
			output.Fatalf("Invalid arguments: %s", err)
		}

		output.Infof("Pulling %s", opts.modelRef.String())
		desc, err := runPull(cmd.Context(), opts)
		if err != nil {
			output.Fatalf("Failed to pull: %s", err)
			return
		}
		output.Infof("Pulled %s", desc.Digest)
	}
}
