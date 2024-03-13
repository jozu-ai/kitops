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

package options

import (
	"github.com/spf13/cobra"
)

// NetworkOptions represent common networking-related flags that are used by multiple commands.
// The flags should be added to the command via AddNetworkFlags before running.
type NetworkOptions struct {
	PlainHTTP bool
	TlsVerify bool
}

func (o *NetworkOptions) AddNetworkFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.PlainHTTP, "plain-http", false, "Use plain HTTP when connecting to remote registries")
	cmd.Flags().BoolVar(&o.TlsVerify, "tls-verify", true, "Require TLS and verify certificates when connecting to remote registries")
}
