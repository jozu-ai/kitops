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

package list

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		input  int64
		output string
	}{
		{input: 0, output: "0 B"},
		{input: 500, output: "500 B"},
		{input: 1<<10 - 1, output: "1023 B"},
		{input: 1 << 10, output: "1.0 KiB"},
		{input: 4.5 * (1 << 10), output: "4.5 KiB"},
		{input: 1<<20 - 1, output: "1023.9 KiB"},
		{input: 1 << 20, output: "1.0 MiB"},
		{input: 6.5 * (1 << 20), output: "6.5 MiB"},
		{input: 1<<30 - 1, output: "1023.9 MiB"},
		{input: 1 << 30, output: "1.0 GiB"},
		{input: 1 << 40, output: "1.0 TiB"},
		{input: 1 << 50, output: "1.0 PiB"},
		{input: 500 * (1 << 50), output: "500.0 PiB"},
		{input: 1 << 60, output: "1024.0 PiB"},
	}
	for idx, tt := range tests {
		t.Run(fmt.Sprintf("test %d", idx), func(t *testing.T) {
			output := formatBytes(tt.input)
			assert.Equalf(t, tt.output, output, "Should convert %d to %s", tt.input, tt.output)
		})
	}
}
