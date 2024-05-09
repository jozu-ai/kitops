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

package dev

const (
	devShortDesc = "Run models locally (experimental)"
	devLongDesc  = "Start a local server and interact with a model in the browser"
	devExample   = "kit dev start"

	devStartShortDesc = "Start development server (experimental)"
	devStartLongDesc  = `Start development server (experimental) from a modelkit

Start a development server for an unpacked modelkit, using a context directory
that includes the model and a kitfile.
`
	devStartExample = `# Serve the model located in the current directory
kit dev start

# Serve the modelkit in ./my-model on port 8080
kit dev start ./my-model --port 8080
`

	devStopShortDesc = "Stop development server"
	devStopLongDesc  = "Stop the development server if it is running"

	devLogsShortDesc = "View logs for development server"
	devLogsLongDesc  = `Print any logs output by the development server.

If the development server is currently running, the logs for this server will
be printed. If it is stopped, the logs for the previous run of the server, if
available, will be printed instead.
`
)
