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

package info

import (
	"context"
	"errors"
	"fmt"
	"kitops/pkg/cmd/options"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/repo"
	"kitops/pkg/output"
	"strings"

	"github.com/spf13/cobra"
	"oras.land/oras-go/v2/errdef"
	"oras.land/oras-go/v2/registry"
)

const (
	shortDesc = `Show the configuration for a modelkit`
	longDesc  = `Print the contents of a modelkit config to the screen.

By default, kit will check local storage for the specified modelkit. To see
the configuration for a modelkit stored on a remote registry, use the
--remote flag.`
	example = `# See configuration for a local modelkit:
kit info mymodel:mytag

# See configuration for a local modelkit by digest:
kit info mymodel@sha256:44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a

# See configuration for a remote modelkit:
kit info --remote registry.example.com/my-model:1.0.0`
)

type infoOptions struct {
	options.NetworkOptions
	configHome  string
	checkRemote bool
	modelRef    *registry.Reference
}

func InfoCommand() *cobra.Command {
	opts := &infoOptions{}

	cmd := &cobra.Command{
		Use:     "info [flags] MODELKIT",
		Short:   shortDesc,
		Long:    longDesc,
		Example: example,
		RunE:    runCommand(opts),
		Args:    cobra.ExactArgs(1),
	}

	opts.AddNetworkFlags(cmd)
	cmd.Flags().BoolVarP(&opts.checkRemote, "remote", "r", false, "Check remote registry instead of local storage")
	return cmd
}

func runCommand(opts *infoOptions) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if err := opts.complete(cmd.Context(), args); err != nil {
			return output.Fatalf("Invalid arguments: %s", err)
		}
		config, err := getInfo(cmd.Context(), opts)
		if err != nil {
			if errors.Is(err, errdef.ErrNotFound) {
				return output.Fatalf("Could not find modelkit %s", repo.FormatRepositoryForDisplay(opts.modelRef.String()))
			}
			return output.Fatalf("Error resolving modelkit: %s", err)
		}
		yamlBytes, err := config.MarshalToYAML()
		if err != nil {
			return output.Fatalf("Error formatting manifest: %w", err)
		}
		fmt.Println(string(yamlBytes))
		return nil
	}
}

func (opts *infoOptions) complete(ctx context.Context, args []string) error {
	configHome, ok := ctx.Value(constants.ConfigKey{}).(string)
	if !ok {
		return fmt.Errorf("default config path not set on command context")
	}
	opts.configHome = configHome

	ref, extraTags, err := repo.ParseReference(args[0])
	if err != nil {
		return err
	}
	if len(extraTags) > 0 {
		return fmt.Errorf("invalid reference format: extra tags are not supported: %s", strings.Join(extraTags, ", "))
	}
	opts.modelRef = ref

	if opts.modelRef.Registry == repo.DefaultRegistry && opts.checkRemote {
		return fmt.Errorf("can not check remote: %s does not contain registry", repo.FormatRepositoryForDisplay(opts.modelRef.String()))
	}

	return nil
}
