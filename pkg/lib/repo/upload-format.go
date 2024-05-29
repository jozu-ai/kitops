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

package repo

import "kitops/pkg/output"

type uploadFormat int

const (
	uploadMonolithicPut uploadFormat = iota
	uploadChunkedPatch
	uploadUndefined
)

func getUploadFormat(registry string, size int64) uploadFormat {
	output.SafeDebugf("Getting upload format for: %s", registry)
	switch registry {
	case "ghcr.io":
		// ghcr.io returns 416 is a PATCH has Content-Length greater than 4.0 MiB for some reason
		// Transfer-Encoding: chunked is supported by the registry, but not implemented yet.
		return uploadMonolithicPut
	}
	// No matches above, use heuristic
	if size < 100<<20 {
		return uploadMonolithicPut
	} else {
		return uploadChunkedPatch
	}
}
