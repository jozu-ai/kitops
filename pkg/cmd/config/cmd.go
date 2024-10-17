// Copyright 2024 The KitOps Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0
package config

import (
	"context"
	"fmt"
	"kitops/pkg/lib/constants"
	"kitops/pkg/output"

	"github.com/spf13/cobra"
)

const (
	shortDesc = `Manage configuration for KitOps CLI`
	longDesc  = `Allows setting, getting, listing, and resetting configuration options for the KitOps CLI.

This command provides functionality to manage configuration settings such as
storage paths, credentials file location, CLI version, and update notification preferences.
The configuration values can be set using specific keys, retrieved for inspection, listed,
or reset to default values.`

	example = `# Set a configuration option
kit config set storageSubpath /path/to/storage

# Get a configuration option
kit config get storageSubpath

# List all configuration options
kit config list

# Reset configuration to default values
kit config reset`
)

func (opts *configOptions) complete(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no configuration key provided")
	}
	opts.key = args[0]

	if len(args) > 1 {
		opts.value = args[1]
	}

	configHome, ok := ctx.Value(constants.ConfigKey{}).(string)
	if !ok {
		return fmt.Errorf("default config path not set on command context")
	}
	opts.configHome = configHome

	return nil
}

// ConfigCommand represents the config command
func ConfigCommand() *cobra.Command {
	opts := &configOptions{}

	cmd := &cobra.Command{
		Use:     "config [set|get|list|reset] KEY [VALUE]",
		Short:   shortDesc,
		Long:    longDesc,
		Example: example,
		RunE:    runCommand(opts),
	}

	cmd.Args = cobra.MinimumNArgs(1)
	cmd.Flags().SortFlags = false

	return cmd
}

func runCommand(opts *configOptions) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		// Handle each command and its required/optional arguments
		switch args[0] {
		case "set":
			if len(args) < 3 {
				return output.Fatalf("Missing key or value for 'set'. Usage: kit config set <key> <value>")
			}
			opts.key, opts.value = args[1], args[2]
			if err := opts.complete(ctx, args); err != nil {
				return output.Fatalf("Invalid arguments: %s", err)
			}
			if err := setConfig(ctx, opts); err != nil {
				return output.Fatalf("Failed to set config: %s", err)
			}
			output.Infof("Configuration key '%s' set to '%s'", opts.key, opts.value)

		case "get":
			if len(args) < 2 {
				return output.Fatalf("Missing key for 'get'. Usage: kit config get <key>")
			}
			opts.key = args[1]
			if err := opts.complete(ctx, args); err != nil {
				return output.Fatalf("Invalid arguments: %s", err)
			}
			value, err := getConfig(ctx, opts)
			if err != nil {
				return output.Fatalf("Failed to get config: %s", err)
			}
			output.Infof("Configuration key '%s': '%s'", opts.key, value)

		case "list":
			// No key required for 'list'
			if err := opts.complete(ctx, args); err != nil {
				return output.Fatalf("Invalid arguments: %s", err)
			}
			if err := listConfig(ctx, opts); err != nil {
				return output.Fatalf("Failed to list configs: %s", err)
			}

		case "reset":
			// No key required for 'reset'
			if err := opts.complete(ctx, args); err != nil {
				return output.Fatalf("Invalid arguments: %s", err)
			}
			if err := resetConfig(ctx, opts); err != nil {
				return output.Fatalf("Failed to reset config: %s", err)
			}
			output.Infof("Configuration reset to default values")

		default:
			return output.Fatalf("Unknown command: %s. Available commands are: set, get, list, reset", args[0])
		}

		return nil
	}
}
