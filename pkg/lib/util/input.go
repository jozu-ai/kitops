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

package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// PromptForInput prints the provided prompt and reads a line of input from stdin.
// If isSensitive is true, input is treated as a password (i.e. not printed)
func PromptForInput(prompt string, isSensitive bool) (string, error) {
	var bytes []byte
	var err error
	if !term.IsTerminal(int(syscall.Stdin)) {
		return "", fmt.Errorf("attempting to read input from non-terminal")
	}

	fmt.Print(prompt)
	if isSensitive {
		bytes, err = term.ReadPassword(int(syscall.Stdin))
		fmt.Print("\n")
	} else {
		reader := bufio.NewReader(os.Stdin)
		bytes, err = reader.ReadBytes('\n')
	}
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}
	return strings.TrimSpace(string(bytes)), nil
}
