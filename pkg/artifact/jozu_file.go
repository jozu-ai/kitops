package artifact

import (
	"encoding/json"
	"io"
	"os"

	"sigs.k8s.io/yaml"
)

type (
	JozuFile struct {
		ManifestVersion string         `json:"manifestVersion"`
		Package         Package        `json:"package,omitempty"`
		Code            []Code         `json:"code,omitempty"`
		DataSets        []DataSet      `json:"datasets,omitempty"`
		Models          []TrainedModel `json:"models,omitempty"`
	}

	Package struct {
		Name        string   `json:"name,omitempty"`
		Version     string   `json:"version,omitempty"`
		Description string   `json:"description,omitempty"`
		License     string   `json:"license,omitempty"`
		Authors     []string `json:"authors,omitempty"`
	}

	Code struct {
		Path        string `json:"path,omitempty"`
		License     string `json:"license,omitempty"`
		Description string `json:"description,omitempty"`
	}

	DataSet struct {
		Name          string `json:"name,omitempty"`
		Path          string `json:"path,omitempty"`
		Description   string `json:"description,omitempty"`
		License       string `json:"license,omitempty"`
		Preprocessing string `json:"preprocessing,omitempty"`
	}

	TrainedModel struct {
		Name        string      `json:"name,omitempty"`
		Path        string      `json:"path,omitempty"`
		Framework   string      `json:"framework,omitempty"`
		Version     string      `json:"version,omitempty"`
		Description string      `json:"description,omitempty"`
		License     string      `json:"license,omitempty"`
		Training    *Training   `json:"training,omitempty"`
		Validation  *Validation `json:"validation,omitempty"`
	}

	Training struct {
		DataSet    string                 `json:"dataset,omitempty"`
		Parameters map[string]interface{} `json:"parameters,omitempty"`
	}

	Validation struct {
		DataSet string                 `json:"dataset,omitempty"`
		Metrics map[string]interface{} `json:"metrics,omitempty"`
	}
)

func (jf *JozuFile) LoadModel(file *os.File) error {
	// Read the file
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, jf)
	if err != nil {
		return err
	}
	return nil
}

func (jf *JozuFile) MarshalToJSON() ([]byte, error) {
	// Marshal the JozuFile to JSON
	jsonData, err := json.Marshal(jf)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
