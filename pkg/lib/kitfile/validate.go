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

package kitfile

import (
	"fmt"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"kitops/pkg/artifact"
)

// ValidateKitfile returns an error if the parsed Kitfile is not valid. The error string
// is multiple lines, each consisting of an issue in the kitfile (e.g. duplicate path).
func ValidateKitfile(kf *artifact.KitFile) error {
	// Map of paths to the component that uses them; used to detect duplicate paths
	paths := map[string][]string{}
	addPath := func(path, source string) {
		if path == "" {
			path = "."
		}
		paths[path] = append(paths[path], source)
	}

	if kf.Model != nil {
		addPath(kf.Model.Path, fmt.Sprintf("model %s", kf.Model.Name))
		for _, part := range kf.Model.Parts {
			addPath(part.Path, fmt.Sprintf("modelpart %s", part.Name))
		}
	}
	for _, dataset := range kf.DataSets {
		addPath(dataset.Path, fmt.Sprintf("dataset %s", dataset.Name))
	}
	for idx, code := range kf.Code {
		addPath(code.Path, fmt.Sprintf("code layer %d", idx))
	}

	var errs []string
	for layerPath, layerIds := range paths {
		if len := len(layerIds); len > 1 {
			errMsg := fmt.Sprintf("  * %s and %s use the same path %s", strings.Join(layerIds[:len-1], ", "), layerIds[len-1], layerPath)
			errs = append(errs, errMsg)
		}
		if path.IsAbs(layerPath) || filepath.IsAbs(layerPath) {
			errMsg := fmt.Sprintf("  * absolute paths are not supported in a Kitfile (path %s in %s)", layerPath, layerIds[0])
			errs = append(errs, errMsg)
		}
	}
	if len(errs) > 0 {
		// Iterating through the paths map is random; sort to get a consistent message
		slices.Sort(errs)
		return fmt.Errorf("errors while validating Kitfile: \n%s", strings.Join(errs, "\n"))
	}

	return nil
}
