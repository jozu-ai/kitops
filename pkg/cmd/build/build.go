package build

import (
	"context"
	"fmt"
	"kitops/pkg/artifact"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/filesystem"
	"kitops/pkg/lib/repo"
	"kitops/pkg/lib/storage"
	"kitops/pkg/output"
	"os"
)

func RunBuild(ctx context.Context, options *buildOptions) error {
	// 1. Read the model file
	modelfile, err := os.Open(options.modelFile)
	if err != nil {
		return err
	}
	defer modelfile.Close()
	kitfile := &artifact.KitFile{}
	if err = kitfile.LoadModel(modelfile); err != nil {
		return err
	}

	model := &artifact.Model{}
	model.Config = kitfile

	// 2. package the Code
	for _, code := range kitfile.Code {
		codePath, err := filesystem.VerifySubpath(options.contextDir, code.Path)
		if err != nil {
			return err
		}
		layer := &artifact.ModelLayer{
			BaseDir:   codePath,
			MediaType: constants.CodeLayerMediaType,
		}
		model.Layers = append(model.Layers, *layer)
	}

	// 3. package the DataSets
	for _, dataset := range kitfile.DataSets {
		datasetPath, err := filesystem.VerifySubpath(options.contextDir, dataset.Path)
		if err != nil {
			return err
		}
		layer := &artifact.ModelLayer{
			BaseDir:   datasetPath,
			MediaType: constants.DataSetLayerMediaType,
		}
		model.Layers = append(model.Layers, *layer)
	}

	// 4. package the TrainedModel
	modelPath, err := filesystem.VerifySubpath(options.contextDir, kitfile.Model.Path)
	if err != nil {
		return err
	}
	layer := &artifact.ModelLayer{
		BaseDir:   modelPath,
		MediaType: constants.ModelLayerMediaType,
	}
	model.Layers = append(model.Layers, *layer)

	tag := ""
	if options.modelRef != nil {
		tag = options.modelRef.Reference
	}
	storageHome := constants.StoragePath(options.configHome)
	localStore, err := repo.NewLocalStore(storageHome, options.modelRef)
	if err != nil {
		return fmt.Errorf("failed to open local storage: %w", err)
	}

	manifestDesc, err := storage.SaveModel(ctx, localStore, model, tag)
	if err != nil {
		return err
	}

	for _, tag := range options.extraRefs {
		if err := localStore.Tag(ctx, *manifestDesc, tag); err != nil {
			return err
		}
	}

	output.Infof("Model saved: %s", manifestDesc.Digest)

	return nil
}
