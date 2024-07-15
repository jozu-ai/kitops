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

package harness

import (
	"fmt"
	"io"
	"kitops/pkg/lib/constants"
	"kitops/pkg/output"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

type LLMHarness struct {
	Host       string
	Port       int
	ConfigHome string
}

func (harness *LLMHarness) Init() error {
	err := extractServer(constants.HarnessPath(harness.ConfigHome), "llama.cpp/build/*/*/*/bin/**")
	if err != nil {
		return fmt.Errorf("failed to extract dev server files: %s", err)
	}
	err = extractUI(constants.HarnessPath(harness.ConfigHome))
	if err != nil {
		return fmt.Errorf("failed to extract dev UI files: %s", err)
	}
	return nil
}

func (harness *LLMHarness) Start(modelPath string) error {
	harnessPath := constants.HarnessPath(harness.ConfigHome)
	pidFile := filepath.Join(harnessPath, constants.HarnessProcessFile)
	logFile := filepath.Join(harnessPath, constants.HarnessLogFile)

	if _, err := os.Stat(pidFile); !os.IsNotExist(err) {
		// Attempt to read the PID from the file.
		pid, err := readPIDFromFile(pidFile)
		if err != nil {
			return fmt.Errorf("failed to read PID file: %w", err)
		}
		// Check if the process is still running.
		if isProcessRunning(pid) {
			return fmt.Errorf("a server process with PID %d is already running", pid)
		} else {
			output.Infoln("The process previously recorded is not running. Proceeding to start a new process.")
		}
	}

	uiHome := filepath.Join(harnessPath, "ui")
	cmd := exec.Command("./server",
		"--host", harness.Host,
		"--port", strconv.Itoa(harness.Port),
		"--model", modelPath,
		"--path", uiHome)
	cmd.Dir = harnessPath
	logs, err := os.OpenFile(logFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file for harness: %w", err)
	}
	defer logs.Close()
	output.Debugf("Saving server logs to %s", logFile)
	cmd.Stdout = logs
	cmd.Stderr = logs

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting llm harness: %s", err)
	}

	pid := cmd.Process.Pid
	if err := writePIDFile(pidFile, pid); err != nil {
		return fmt.Errorf("failed to write PID file: %w", err)
	}

	output.Debugf("Started harness with PID %d and saved to file.\n", pid)

	return nil
}

func (harness *LLMHarness) Stop() error {
	pidFile := filepath.Join(constants.HarnessPath(harness.ConfigHome), constants.HarnessProcessFile)

	pid, err := readPIDFromFile(pidFile)
	if os.IsNotExist(err) {
		return fmt.Errorf("no Running server found")
	}
	if err != nil {
		return err
	}

	// Check if the process is still running.
	if !isProcessRunning(pid) {
		return fmt.Errorf("no running process found with PID %d", pid)
	}

	// Kill the process using the PID.
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("error finding process: %s", err)
	}

	err = process.Signal(syscall.SIGTERM) // Try to kill it gently
	if err != nil {
		output.Infof("Error killing process %s", err)
		// If SIGTERM failed, kill it with SIGKILL
		err = process.Kill()
		if err != nil {
			return fmt.Errorf("error killing process: %s", err)
		}
	}

	output.Debugf("Process with PID %d has been killed.", pid)
	// Delete the PID file to clean up.
	err = os.Remove(pidFile)
	if err != nil {
		return fmt.Errorf("error removing PID file: %s", err)
	}

	return nil
}

func PrintLogs(configHome string, w io.Writer) error {
	harnessPath := constants.HarnessPath(configHome)
	logPath := filepath.Join(harnessPath, constants.HarnessLogFile)
	logFile, err := os.Open(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			output.Errorf("No log file found")
			return nil
		}
		return fmt.Errorf("Error reading log file: %w", err)
	}
	defer logFile.Close()
	if _, err = io.Copy(w, logFile); err != nil {
		return fmt.Errorf("Failed to print log file: %w", err)
	}
	return nil
}

func isProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// Sending signal 0 to a process does not affect it but can be used for error checking.
	// If an error is returned, the process does not exist.
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// ensures the directory for pidFilePath exists and writes the PID to the file.
func writePIDFile(pidFilePath string, pid int) error {
	// Ensure the directory for the pidFilePath exists.
	dir := filepath.Dir(pidFilePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", dir, err)
	}

	// Write the PID to the file.
	pidData := []byte(fmt.Sprintf("%d", pid))
	if err := os.WriteFile(pidFilePath, pidData, 0644); err != nil {
		return fmt.Errorf("error writing PID to file %s: %w", pidFilePath, err)
	}

	return nil
}

func readPIDFromFile(filePath string) (int, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0, err
	}
	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}
	return pid, nil
}
