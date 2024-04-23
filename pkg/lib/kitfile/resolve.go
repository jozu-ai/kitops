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
	"context"
	"fmt"
	"kitops/pkg/artifact"
	"kitops/pkg/cmd/info"
	"kitops/pkg/lib/constants"
	"strings"
)

// ResolveKitfile returns the Kitfile for a reference. Any references to other modelkits
// are fetched and included in the resolved Kitfile, giving the equivalent Kitfile including
// the model, datasets, and code from those referenced modelkits.
func ResolveKitfile(ctx context.Context, configHome, kitfileRef, baseRef string) (*artifact.KitFile, error) {
	resolved := &artifact.KitFile{}
	refChain := []string{baseRef, kitfileRef}
	for i := 0; i < constants.MaxModelRefChain; i++ {
		kitfile, err := info.GetKitfileForRefString(ctx, configHome, kitfileRef)
		if err != nil {
			return nil, err
		}
		resolved = mergeKitfiles(resolved, kitfile)
		if resolved.Model == nil || !IsModelKitReference(resolved.Model.Path) {
			return resolved, nil
		}
		if idx := getIndex(refChain, resolved.Model.Path); idx != -1 {
			cycleStr := fmt.Sprintf("[%s=>%s]", strings.Join(refChain[idx:], "=>"), resolved.Model.Path)
			return nil, fmt.Errorf("Found cycle in modelkit references: %s", cycleStr)
		}
		refChain = append(refChain, resolved.Model.Path)
		kitfileRef = resolved.Model.Path
	}
	return nil, fmt.Errorf("Reached maximum number of model references: [%s]", strings.Join(refChain, "=>"))
}

func mergeKitfiles(into, from *artifact.KitFile) *artifact.KitFile {
	firstNonEmpty := func(strs ...string) string {
		for _, s := range strs {
			if s != "" {
				return s
			}
		}
		return ""
	}

	result := &artifact.KitFile{}
	result.ManifestVersion = firstNonEmpty(into.ManifestVersion, from.ManifestVersion)
	result.Kit.Name = firstNonEmpty(into.Kit.Name, from.Kit.Name)
	result.Kit.Description = firstNonEmpty(into.Kit.Description, from.Kit.Description)
	result.Kit.License = firstNonEmpty(into.Kit.License, from.Kit.License)
	result.Kit.Version = firstNonEmpty(into.Kit.Version, from.Kit.Version)
	result.Kit.Authors = append(into.Kit.Authors, from.Kit.Authors...)

	if into.Model != nil || from.Model != nil {
		result.Model = &artifact.TrainedModel{}
		intoModel := into.Model
		fromModel := from.Model
		if intoModel == nil {
			intoModel = &artifact.TrainedModel{}
		}
		if fromModel == nil {
			fromModel = &artifact.TrainedModel{}
		}
		result.Model.Path = fromModel.Path
		result.Model.Name = firstNonEmpty(intoModel.Name, fromModel.Name)
		result.Model.Description = firstNonEmpty(intoModel.Description, fromModel.Description)
		result.Model.License = firstNonEmpty(intoModel.License, fromModel.License)
		result.Model.Framework = firstNonEmpty(intoModel.Framework, fromModel.Framework)
		result.Model.Version = firstNonEmpty(intoModel.Version, fromModel.Version)
		result.Model.Parts = append(intoModel.Parts, fromModel.Parts...)
	}

	result.Code = into.Code
	result.Code = append(result.Code, from.Code...)
	result.DataSets = into.DataSets
	result.DataSets = append(result.DataSets, from.DataSets...)
	return result
}

func getIndex(list []string, s string) int {
	for idx, item := range list {
		if s == item {
			return idx
		}
	}
	return -1
}
