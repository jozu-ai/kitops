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

package testing

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"kitops/cmd"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

const modelKitTag = "test:test"

type testcase struct {
	Name         string
	Description  string   `yaml:"description"`
	Kitfile      string   `yaml:"kitfile"`
	Kitignore    string   `yaml:"kitignore"`
	Files        []string `yaml:"files"`
	IgnoredFiles []string `yaml:"ignored"`
}

// runCommand executes kit <args>, saving stdout/stderr output to a buffer
// that is then printed through the test interface. If the kit command
// calls `os.Exit`, this command will terminate without generating any logs.
// Returns the stdout and stderr output of the command.
func runCommand(t *testing.T, args ...string) string {
	t.Logf("Running command: kit %s", strings.Join(args, " "))
	runCmd := cmd.RunCommand()
	runCmd.SetArgs(args)

	// Set up buffer to capture command output
	outbuf := &bytes.Buffer{}
	runCmd.SetOut(outbuf)
	runCmd.SetErr(outbuf)

	err := runCmd.Execute()
	if !assert.NoError(t, err, "Command returned error") {
		return ""
	}

	outlog, err := io.ReadAll(outbuf)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Command output: \n%s", string(outlog))
	return string(outlog)
}

func setupTempDir(t *testing.T) (tmpDir string, removeTmpDir func()) {
	// Set up temporary directory for work
	tmpDir, err := os.MkdirTemp("", "kitops-testing-*")
	if !assert.NoError(t, err) {
		t.Fatalf("Could not create temporary directory: %s", err)
	}
	removeTmpDir = func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Error removing temp dir: %s", err)
		}
	}
	t.Logf("Using temp directory: %s", tmpDir)
	return tmpDir, removeTmpDir
}

// setupTestDirs generates the test directories used for storing $KIT_HOME, the original modelkit
// and the unpacked modelkit as subdirectories of tmpDir.
func setupTestDirs(t *testing.T, tmpDir string) (modelKitPath, unpackPath, contextPath string) {
	// Set up paths to use for test
	modelKitPath = filepath.Join(tmpDir, "test-modelkit-in")
	unpackPath = filepath.Join(tmpDir, "test-modelkit-out")
	contextPath = filepath.Join(tmpDir, ".kitops")
	for _, path := range []string{modelKitPath, unpackPath, contextPath} {
		if err := os.MkdirAll(path, 0755); err != nil {
			t.Fatal(err)
		}
	}
	return
}

// setupFiles ensures that all paths in files exist within tmpDir. Directories along the
// path are created if necessary, and files contain the text "testing: <filename>".
func setupFiles(t *testing.T, tmpDir string, files []string) {
	for _, file := range files {
		path := filepath.Join(tmpDir, file)
		dirName := filepath.Dir(path)
		if err := os.MkdirAll(dirName, 0755); err != nil {
			t.Fatal(err)
		}
		t.Logf("creating path %s", path)
		if err := os.WriteFile(path, []byte("testing: "+file), 0644); err != nil {
			t.Fatal(err)
		}
	}
}

// checkFilesExist tests that every path listed in files exists within tmpDir, failing the
// current test if not.
func checkFilesExist(t *testing.T, tmpDir string, files []string) {
	for _, file := range files {
		path := filepath.Join(tmpDir, file)
		stat, err := os.Stat(path)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				t.Errorf("File %s should exist", file)
			} else {
				t.Errorf("Unexpected error: %s", err)
			}
		} else {
			assert.True(t, stat.Mode().IsRegular(), "Path %s should be regular file", path)
		}
	}
}

// checkFilesDoNotExist checks that none of the paths listed in files exist within tmpDir.
func checkFilesDoNotExist(t *testing.T, tmpDir string, files []string) {
	for _, file := range files {
		path := filepath.Join(tmpDir, file)
		_, err := os.Stat(path)
		if err == nil {
			t.Errorf("File %s should not exist", file)
		} else if !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("Unexpected error: %s", err)
		}
	}
}

func loadAllTestCasesOrPanic(t *testing.T, testsPath string) []testcase {
	files, err := os.ReadDir(testsPath)
	if err != nil {
		t.Fatal(err)
	}
	var tests []testcase
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		bytes, err := os.ReadFile(filepath.Join(testsPath, file.Name()))
		if err != nil {
			t.Fatal(err)
		}
		testcase := testcase{}
		if err := yaml.Unmarshal(bytes, &testcase); err != nil {
			t.Fatal(err)
		}
		testcase.Name = file.Name()
		tests = append(tests, testcase)
	}
	return tests
}
