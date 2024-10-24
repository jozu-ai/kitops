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

// Root config command.
func ConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Short:   shortDesc,
		Long:    longDesc,
		Example: example,
	}

	// Add subcommands to the root config command.
	cmd.AddCommand(setCmd())
	cmd.AddCommand(getCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(resetCmd())

	return cmd
}

// Subcommand for 'set'
func setCmd() *cobra.Command {
	opts := &configOptions{}
	cmd := &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Set a configuration value",
		Args:  cobra.ExactArgs(2), // Ensure exactly 2 arguments: key and value.
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			opts.key, opts.value = args[0], args[1]
			if err := opts.complete(ctx); err != nil {
				return fmt.Errorf("failed to complete options: %w", err)
			}
			if err := setConfig(ctx, opts); err != nil {
				return fmt.Errorf("failed to set config: %w", err)
			}
			output.Infof("Configuration key '%s' set to '%s'", opts.key, opts.value)
			return nil
		},
	}

	return cmd
}

// Subcommand for 'get'
func getCmd() *cobra.Command {
	opts := &configOptions{}
	cmd := &cobra.Command{
		Use:   "get [key]",
		Short: "Get a configuration value",
		Args:  cobra.ExactArgs(1), // Ensure exactly 1 argument: key.
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			opts.key = args[0]
			if err := opts.complete(ctx); err != nil {
				return fmt.Errorf("failed to complete options: %w", err)
			}
			value, err := getConfig(ctx, opts)
			if err != nil {
				return fmt.Errorf("failed to get config: %w", err)
			}
			output.Infof("Configuration key '%s': '%s'", opts.key, value)
			return nil
		},
	}

	return cmd
}

// Subcommand for 'list'
func listCmd() *cobra.Command {
	opts := &configOptions{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all configuration values",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if err := opts.complete(ctx); err != nil {
				return fmt.Errorf("failed to complete options: %w", err)
			}
			if err := listConfig(ctx, opts); err != nil {
				return fmt.Errorf("failed to list configs: %w", err)
			}
			return nil
		},
	}

	return cmd
}

// Subcommand for 'reset'
func resetCmd() *cobra.Command {
	opts := &configOptions{}
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset configuration to default values",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if err := opts.complete(ctx); err != nil {
				return fmt.Errorf("failed to complete options: %w", err)
			}
			if err := resetConfig(ctx, opts); err != nil {
				return fmt.Errorf("failed to reset config: %w", err)
			}
			output.Infof("Configuration reset to default values")
			return nil
		},
	}

	return cmd
}

// complete populates configOptions fields.
func (opts *configOptions) complete(ctx context.Context) error {
	configHome, ok := ctx.Value(constants.ConfigKey{}).(string)
	if !ok {
		return fmt.Errorf("default config path not set on command context")
	}
	opts.configHome = configHome
	return nil
}
