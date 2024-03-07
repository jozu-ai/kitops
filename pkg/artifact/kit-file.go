package artifact

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

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

func (kf *KitFile) LoadModel(filePath string) error {
	modelfile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer modelfile.Close()
	// Read the file
	data, err := io.ReadAll(modelfile)
	if err != nil {
		return err
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
