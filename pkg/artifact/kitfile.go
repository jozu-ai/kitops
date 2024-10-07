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

package artifact

import (
	"bytes"
	"encoding/json"
	"io"

	"gopkg.in/yaml.v3"
)

type (
	KitFile struct {
		ManifestVersion string    `json:"manifestVersion" yaml:"manifestVersion"`
		Package         Package   `json:"package,omitempty" yaml:"package,omitempty"`
		Model           *Model    `json:"model,omitempty" yaml:"model,omitempty"`
		Code            []Code    `json:"code,omitempty" yaml:"code,omitempty"`
		DataSets        []DataSet `json:"datasets,omitempty" yaml:"datasets,omitempty"`
		Docs            []Docs    `json:"docs,omitempty" yaml:"docs,omitempty"`
	}

	Docs struct {
		Path        string `json:"path" yaml:"path"`
		Description string `json:"description" yaml:"description"`
	}

	Package struct {
		Name        string   `json:"name,omitempty" yaml:"name,omitempty"`
		Version     string   `json:"version,omitempty" yaml:"version,omitempty"`
		Description string   `json:"description,omitempty" yaml:"description,omitempty"`
		License     string   `json:"license,omitempty" yaml:"license,omitempty"`
		Authors     []string `json:"authors,omitempty" yaml:"authors,omitempty,flow"`
	}

	Model struct {
		Name        string      `json:"name,omitempty" yaml:"name,omitempty"`
		Path        string      `json:"path,omitempty" yaml:"path,omitempty"`
		License     string      `json:"license,omitempty" yaml:"license,omitempty"`
		Framework   string      `json:"framework,omitempty" yaml:"framework,omitempty"`
		Version     string      `json:"version,omitempty" yaml:"version,omitempty"`
		Description string      `json:"description,omitempty" yaml:"description,omitempty"`
		Parts       []ModelPart `json:"parts,omitempty" yaml:"parts,omitempty"`
		// Parameters is an arbitrary section of yaml that can be used to store any additional
		// data that may be relevant to the current model, with a few caveats:
		//  * Only a json-compatible subset of yaml is supported
		//  * Strings will be serialized without flow parameters, etc.
		//  * Numbers will be converted to decimal representations (0xFF -> 255, 1.2e+3 -> 1200)
		//  * Maps will be sorted alphabetically by key
		Parameters any `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	}

	ModelPart struct {
		Name    string `json:"name,omitempty" yaml:"name,omitempty"`
		Path    string `json:"path,omitempty" yaml:"path,omitempty"`
		License string `json:"license,omitempty" yaml:"license,omitempty"`
		Type    string `json:"type,omitempty" yaml:"type,omitempty"`
	}

	Code struct {
		Path        string `json:"path,omitempty" yaml:"path,omitempty"`
		Description string `json:"description,omitempty" yaml:"description,omitempty"`
		License     string `json:"license,omitempty" yaml:"license,omitempty"`
	}

	DataSet struct {
		Name        string `json:"name,omitempty" yaml:"name,omitempty"`
		Path        string `json:"path,omitempty" yaml:"path,omitempty"`
		Description string `json:"description,omitempty" yaml:"description,omitempty"`
		License     string `json:"license,omitempty" yaml:"license,omitempty"`
		// Parameters is an arbitrary section of yaml that can be used to store any additional
		// metadata relevant to the dataset, with a few caveats:
		//  * Only a json-compatible subset of yaml is supported
		//  * Strings will be serialized without flow parameters, etc.
		//  * Numbers will be converted to decimal representations
		//  * Maps will be sorted alphabetically by key
		//  * It's recommended to store metadata like preprocessing steps, formats, etc.
		Parameters any `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	}
)

func (kf *KitFile) LoadModel(kitfileContent io.ReadCloser) error {
	decoder := yaml.NewDecoder(kitfileContent)
	decoder.KnownFields(true)
	if err := decoder.Decode(kf); err != nil {
		return err
	}
	return nil
}

func (kf *KitFile) MarshalToJSON() ([]byte, error) {
	jsonData, err := json.Marshal(kf)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func (kf *KitFile) MarshalToYAML() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := yaml.NewEncoder(buf)
	enc.SetIndent(2)
	if err := enc.Encode(kf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
