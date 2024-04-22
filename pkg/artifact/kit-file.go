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
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

type (
	KitFile struct {
		ManifestVersion string        `json:"manifestVersion" yaml:"manifestVersion"`
		Kit             ModelKit      `json:"package,omitempty" yaml:"package,omitempty"`
		Model           *TrainedModel `json:"model,omitempty" yaml:"model,omitempty"`
		Code            []Code        `json:"code,omitempty" yaml:"code,omitempty"`
		DataSets        []DataSet     `json:"datasets,omitempty" yaml:"datasets,omitempty"`
	}

	ModelKit struct {
		Name        string   `json:"name,omitempty" yaml:"name,omitempty"`
		Version     string   `json:"version,omitempty" yaml:"version,omitempty"`
		Description string   `json:"description,omitempty" yaml:"description,omitempty"`
		License     string   `json:"license,omitempty" yaml:"license,omitempty"`
		Authors     []string `json:"authors,omitempty" yaml:"authors,omitempty,flow"`
	}

	Code struct {
		Path        string `json:"path,omitempty" yaml:"path,omitempty"`
		License     string `json:"license,omitempty" yaml:"license,omitempty"`
		Description string `json:"description,omitempty" yaml:"description,omitempty"`
	}

	DataSet struct {
		Name          string `json:"name,omitempty" yaml:"name,omitempty"`
		Path          string `json:"path,omitempty" yaml:"path,omitempty"`
		Description   string `json:"description,omitempty" yaml:"description,omitempty"`
		License       string `json:"license,omitempty" yaml:"license,omitempty"`
		Preprocessing string `json:"preprocessing,omitempty" yaml:"preprocessing,omitempty"`
	}

	TrainedModel struct {
		Name        string      `json:"name,omitempty" yaml:"name,omitempty"`
		Path        string      `json:"path,omitempty" yaml:"path,omitempty"`
		Framework   string      `json:"framework,omitempty" yaml:"framework,omitempty"`
		Version     string      `json:"version,omitempty" yaml:"version,omitempty"`
		Description string      `json:"description,omitempty" yaml:"description,omitempty"`
		License     string      `json:"license,omitempty" yaml:"license,omitempty"`
		Training    *Training   `json:"training,omitempty" yaml:"training,omitempty"`
		Validation  *Validation `json:"validation,omitempty" yaml:"validation,omitempty"`
		Parts       []ModelPart `json:"parts,omitempty" yaml:"parts,omitempty"`
	}

	ModelPart struct {
		Name        string `json:"name,omitempty" yaml:"name,omitempty"`
		Path        string `json:"path,omitempty" yaml:"path,omitempty"`
		Description string `json:"description,omitempty" yaml:"description,omitempty"`
	}

	Training struct {
		DataSet    string                 `json:"dataset,omitempty" yaml:"dataset,omitempty"`
		Parameters map[string]interface{} `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	}

	Validation struct {
		DataSet string                 `json:"dataset,omitempty" yaml:"dataset,omitempty"`
		Metrics map[string]interface{} `json:"metrics,omitempty" yaml:"metrics,omitempty"`
	}
)

func (kf *KitFile) LoadModel(kitfileContent io.ReadCloser) error {

	data, err := io.ReadAll(kitfileContent)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return fmt.Errorf("empty kitfile")
	}
	err = yaml.Unmarshal(data, kf)
	if err != nil {
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
