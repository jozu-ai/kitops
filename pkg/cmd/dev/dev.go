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
