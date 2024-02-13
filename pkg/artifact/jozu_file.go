package artifact

import (
	"encoding/json"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type (
	JozuFile struct {
		ManifestVersion string         `yaml:"manifestVersion"`
		Package         Package        `yaml:"package"`
		Code            []Code         `yaml:"code"`
		DataSets        []DataSet      `yaml:"datasets"`
		Models          []TrainedModel `yaml:"models"`
	}

	Package struct {
		Name        string   `yaml:"name"`
		Version     string   `yaml:"version"`
		Description string   `yaml:"description"`
		Authors     []string `yaml:"authors"`
	}

	Code struct {
		Path        string `yaml:"path"`
		License     string `yaml:"license"`
		Description string `yaml:"description"`
	}

	DataSet struct {
		Name          string `yaml:"name"`
		Path          string `yaml:"path"`
		Description   string `yaml:"description"`
		License       string `yaml:"license"`
		Preprocessing string `yaml:"preprocessing"`
	}

	TrainedModel struct {
		Name        string     `yaml:"name"`
		Path        string     `yaml:"path"`
		Framework   string     `yaml:"framework"`
		Version     string     `yaml:"version"`
		Description string     `yaml:"description"`
		License     string     `yaml:"license"`
		Training    Training   `yaml:"training"`
		Validation  Validation `yaml:"validation"`
	}

	Training struct {
		DataSet    string                 `yaml:"dataset"`
		Parameters map[string]interface{} `yaml:"parameters"`
	}

	Validation struct {
		DataSet string                 `yaml:"dataset"`
		Metrics map[string]interface{} `yaml:"metrics"`
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
