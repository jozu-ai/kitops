package harness

import (
	"fmt"
	"kitops/pkg/lib/constants"
	"kitops/pkg/output"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

type LLMHarness struct {
	Port       int
	ConfigHome string
}

func (harness *LLMHarness) Init() {
	extractLibraries(constants.HarnessPath(harness.ConfigHome), "llama.cpp/build/*/*/*/bin/**")
}

func (harness *LLMHarness) Start(modelPath string) error {
	pidFile := filepath.Join(constants.HarnessPath(harness.ConfigHome), constants.HarnessProcessFile)

	if _, err := os.Stat(pidFile); !os.IsNotExist(err) {
		// Attempt to read the PID from the file.
		pid, err := readPIDFromFile(pidFile)
		if err != nil {
			return fmt.Errorf("failed to read PID file: %w", err)
		}
		// Check if the process is still running.
		if isProcessRunning(pid) {
			return fmt.Errorf("a process with PID %d is already running", pid)
		} else {
			fmt.Println("The process previously recorded is not running. Proceeding to start a new process.")
		}
	}

	cmd := exec.Command("./server",
		"--port", strconv.Itoa(harness.Port),
		"--model", modelPath)
	cmd.Dir = constants.HarnessPath(harness.ConfigHome)

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting llm harness: %s\n", err)
		return err
	}

	pid := cmd.Process.Pid
	writePIDFile(pidFile, pid)

	output.Debugf("Started harness with PID %d and saved to file.\n", pid)

	return nil
}

func (harness *LLMHarness) Stop() error {
	pidFile := filepath.Join(constants.HarnessPath(harness.ConfigHome), constants.HarnessProcessFile)

	pid, err := readPIDFromFile(pidFile)
	if err != nil {
		fmt.Printf("Error reading PID file: %s\n", err)
		return err
	}

	// Check if the process is still running.
	if !isProcessRunning(pid) {
		fmt.Printf("No running process found with PID %d. Nothing to stop.\n", pid)
		return err
	}

	// Kill the process using the PID.
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Printf("Error finding process: %s\n", err)
		return err
	}

	err = process.Signal(syscall.SIGTERM) // Try to kill it gently
	if err != nil {
		fmt.Printf("Error killing process: %s\n", err)
		// If SIGTERM failed, kill it with SIGKILL
		err = process.Kill()
		if err != nil {
			fmt.Printf("Error killing process with SIGKILL: %s\n", err)
			return err
		}
	}

	fmt.Printf("Process with PID %d has been killed.\n", pid)
	// Delete the PID file to clean up.
	err = os.Remove(pidFile)
	if err != nil {
		fmt.Printf("Error removing PID file: %s\n", err)
		return nil
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
