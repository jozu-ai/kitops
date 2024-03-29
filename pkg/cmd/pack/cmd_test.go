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

package pack

import (
	"context"
	"kitops/pkg/lib/constants"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackOptions_Complete(t *testing.T) {
	options := &packOptions{}
	ctx := context.WithValue(context.Background(), constants.ConfigKey{}, "/home/user/.kitops")
	args := []string{"arg1"}

	err := options.complete(ctx, args)

	assert.NoError(t, err)
	assert.Equal(t, args[0], options.contextDir)
	assert.Equal(t, filepath.Join(args[0], constants.DefaultKitFileName), options.modelFile)
}

func TestPackOptions_RunPack(t *testing.T) {
	t.Skip("Skipping test for now")
	options := &packOptions{
		modelFile:  "Kitfile",
		contextDir: "/path/to/context",
	}

	err := runPack(context.Background(), options)

	assert.NoError(t, err)
}
