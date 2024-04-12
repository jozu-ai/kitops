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
package dev

import (
	"context"
	"fmt"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/filesystem"
	"kitops/pkg/output"
	"net"
	"strconv"

	"github.com/spf13/cobra"
)

const (
	shortDesc = `Start development server (experimental)`
	longDesc  = `Start development server (experimental) with the specified context directory and kitfile`
	example   = `kit dev ./my-model --port 8080`
)

var (
	flags *DevFlags
	opts  *DevOptions
)

type DevFlags struct {
	Port      int
	ModelFile string
	Stop      bool
}

type DevOptions struct {
	port       int
	modelFile  string
	contextDir string
	configHome string
	stop       bool
}

func (opts *DevOptions) complete(ctx context.Context, args []string) error {
	opts.contextDir = ""
	if len(args) == 1 {
		opts.contextDir = args[0]
	}

	configHome, ok := ctx.Value(constants.ConfigKey{}).(string)
	if !ok {
		return fmt.Errorf("default config path not set on command context")
	}
	opts.configHome = configHome
	if flags.Stop {
		opts.stop = flags.Stop
	}

	opts.modelFile = flags.ModelFile
	if opts.modelFile == "" {
		opts.modelFile = filesystem.FindKitfileInPath(opts.contextDir)
	}
	opts.port = flags.Port
	if opts.port == 0 {
		availPort, err := findAvailablePort()
		if err != nil {
			output.Fatalf("failed to find available port: %v", err)
			return err
		}
		opts.port = availPort
	}
	return nil
}

func DevCommand() *cobra.Command {
	opts = &DevOptions{}
	flags = &DevFlags{}
	cmd := &cobra.Command{
		Use:     "dev <directory> [flags]",
		Short:   shortDesc,
		Long:    longDesc,
		Example: example,
		Run:     runCommand(opts),
	}
	cmd.Args = cobra.MaximumNArgs(1)
	cmd.Flags().StringVarP(&flags.ModelFile, "file", "f", "", "Path to the kitfile")
	cmd.Flags().IntVar(&flags.Port, "port", 0, "Port for development server to listen on")
	cmd.Flags().BoolVar(&flags.Stop, "stop", false, "Stop the development server")
	return cmd
}

func runCommand(opts *DevOptions) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := opts.complete(cmd.Context(), args); err != nil {
			output.Errorf("failed to complete options: %w", err)
		}
		if opts.stop {
			output.Infoln("Stopping development server...")
			err := stopDev(cmd.Context(), opts)
			if err != nil {
				output.Fatalf("Failed to stop dev server: %s", err)
				return
			}
			return
		}
		err := runDev(cmd.Context(), opts)
		if err != nil {
			output.Fatalf("Failed to start dev server: %s", err)
			return
		}
	}
}

func findAvailablePort() (int, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	_, portStr, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		return 0, err
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, err
	}
	return port, nil
}
