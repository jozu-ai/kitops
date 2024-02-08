package artifact

import (
	"encoding/json"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type (
	JozuFile struct {
		Maintainer  string      `yaml:"maintainer"`
		ModelFormat string      `yaml:"modelFormat"`
		Inputs      []Parameter `yaml:"inputs"`
		Outputs     []Parameter `yaml:"outputs"`
	}

	Parameter struct {
		Name     string `yaml:"name"`
		Datatype string `yaml:"datatype"`
		Dims     []int  `yaml:"dims"`
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
