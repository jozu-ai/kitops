// Copyright 2024 The KitOps Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0
package dev

import (
	"context"
	"kitops/pkg/artifact"
	"kitops/pkg/lib/filesystem"
	"kitops/pkg/lib/harness"
	"kitops/pkg/output"
	"os"
)

func runDev(ctx context.Context, options *DevOptions) error {

	kitfile := &artifact.KitFile{}
	if fileInfo, err := os.Stat(options.modelFile); err == nil && fileInfo.IsDir() {
		options.modelFile = filesystem.FindKitfileInPath(options.modelFile)
	}

	modelfile, err := os.Open(options.modelFile)
	if err != nil {
		return err
	}
	defer modelfile.Close()
	if err := kitfile.LoadModel(modelfile); err != nil {
		return err
	}
	output.Infof("Loaded Kitfile: %s", kitfile.Model.Path)
	modelPath, _, err := filesystem.VerifySubpath(options.contextDir, kitfile.Model.Path)
	if err != nil {
		return err
	}

	llmHarness := &harness.LLMHarness{}
	llmHarness.Port = options.port
	llmHarness.ConfigHome = options.configHome
	llmHarness.Init()

	if err := llmHarness.Start(modelPath); err != nil {
		output.Errorf("Error starting llm harness: %s", err)
		return err
	}

	output.Infof("Development server started at http://localhost:%d", options.port)

	return nil
}

func stopDev(ctx context.Context, options *DevOptions) error {

	llmHarness := &harness.LLMHarness{}
	llmHarness.ConfigHome = options.configHome
	llmHarness.Init()

	return llmHarness.Stop()
}
