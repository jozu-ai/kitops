package artifact

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type parameterTestCase struct {
	Name        string
	Description string `yaml:"description"`
	KitfileYaml string `yaml:"kitfileYaml"`
	KitfileJson string `yaml:"kitfileJson"`
}

func (tc parameterTestCase) withName(name string) parameterTestCase {
	tc.Name = name
	return tc
}

func TestParameterMarshalUnmarshal(t *testing.T) {
	tests := loadAllTestCasesOrPanic[parameterTestCase](t, filepath.Join("testdata", "parameters"))
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s (%s)", tt.Name, tt.Description), func(t *testing.T) {
			kf := &KitFile{}
			rc := io.NopCloser(strings.NewReader(tt.KitfileYaml))
			err := kf.LoadModel(rc)
			if !assert.NoError(t, err) {
				return
			}

			unmarshalledYaml, err := kf.MarshalToYAML()
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, tt.KitfileYaml, string(unmarshalledYaml))

			unmarshalledJson, err := kf.MarshalToJSON()
			if !assert.NoError(t, err) {
				return
			}
			if tt.KitfileJson != "" {
				assert.Equal(t, tt.KitfileJson, string(unmarshalledJson))
			}
		})
	}
}

func loadAllTestCasesOrPanic[T interface{ withName(string) T }](t *testing.T, testsPath string) []T {
	files, err := os.ReadDir(testsPath)
	if err != nil {
		t.Fatal(err)
	}
	var tests []T
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		bytes, err := os.ReadFile(filepath.Join(testsPath, file.Name()))
		if err != nil {
			t.Fatal(err)
		}
		var testcase T
		if err := yaml.Unmarshal(bytes, &testcase); err != nil {
			t.Fatal(err)
		}
		testcase = testcase.withName(file.Name())
		tests = append(tests, testcase)
	}
	return tests
}
