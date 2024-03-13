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

package version

import (
	"kitops/pkg/output"
	"runtime"

	"github.com/spf13/cobra"
)

const (
	shortDesc = `Display the version information for the CLI`
	longDesc  = `The version command prints detailed version information.

This information includes the current version of the tool, the Git commit that
the version was built from, the build time, and the version of Go it was
compiled with.`
)

// Default build-time variable.
// These values are overridden via ldflags
var (
	Version   = "unknown"
	GitCommit = "unknown"
	BuildTime = "unknown"
	GoVersion = runtime.Version()
)

func VersionCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "version",
		Short: shortDesc,
		Long:  longDesc,
		Run: func(cmd *cobra.Command, args []string) {
			output.Infof("Version: %s\nCommit: %s\nBuilt: %s\nGo version: %s\n", Version, GitCommit, BuildTime, GoVersion)
		},
	}
	return cmd
}

func init() {}
