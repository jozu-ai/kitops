package config

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

import (
	"context"
	"fmt"
	"kitops/pkg/cmd/options"
	"kitops/pkg/lib/constants"
	"kitops/pkg/output"

	"github.com/spf13/cobra"
)

const (
	shortDesc = `Manage configuration for KitOps CLI`
	longDesc  = `Allows setting, getting, and resetting configuration options for the KitOps CLI.

This command provides functionality to manage the configuration settings such as
storage paths, credentials file location, CLI version, and update notification preferences.
The configuration values can be set using specific keys, retrieved for inspection, or reset to default values.`

	example = `# Set a configuration option
kit config set storagePath /path/to/storage

# Get a configuration option
kit config get storagePath

# Reset configuration to default values
kit config reset`
)

type configOptions struct {
	options.NetworkOptions
	configHome string 
	key   string
	value string
}

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

	if err := opts.NetworkOptions.Complete(ctx, args); err != nil {
		return err
	}

	return nil
}

// ConfigCommand represents the config command
func ConfigCommand() *cobra.Command {
	opts := &configOptions{}

	cmd := &cobra.Command{
		Use:     "config [set|get|reset] <key> [value]",
		Short:   shortDesc,
		Long:    longDesc,
		Example: example,
		RunE:    runCommand(opts),
	}

	cmd.Args = cobra.MinimumNArgs(1)
	opts.AddNetworkFlags(cmd)
	cmd.Flags().SortFlags = false

	return cmd
}

func runCommand(opts *configOptions) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if err := opts.complete(cmd.Context(), args); err != nil {
			return output.Fatalf("Invalid arguments: %s", err)
		}

		switch args[0] {
		case "set":
			if len(args) < 2 {
				return output.Fatalf("Missing value for key %s", args[1])
			}
			err := setConfig(cmd.Context(), opts)
			if err != nil {
				return output.Fatalf("Failed to set config: %s", err)
			}
			output.Infof("Configuration key %s set to %s", opts.key, opts.value)
		case "get":
			value, err := getConfig(cmd.Context(), opts)
			if err != nil {
				return output.Fatalf("Failed to get config: %s", err)
			}
			output.Infof("Configuration key %s: %s", opts.key, value)
		case "reset":
			err := resetConfig(cmd.Context(), opts)
			if err != nil {
				return output.Fatalf("Failed to reset config: %s", err)
			}
			output.Infof("Configuration reset to default values")
		default:
			return output.Fatalf("Unknown command %s", args[0])
		}

		return nil
	}
}


