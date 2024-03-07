package pack

import (
	"context"
	"fmt"
	"kitops/pkg/artifact"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/filesystem"
	"kitops/pkg/lib/repo"
	"kitops/pkg/lib/storage"
	"kitops/pkg/output"
)

func RunPack(ctx context.Context, options *packOptions) error {
	// 1. Read the model file
	kitfile := &artifact.KitFile{}
	if err := kitfile.LoadModel(options.modelFile); err != nil {
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
			Path:      codePath,
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
			Path:      datasetPath,
			MediaType: constants.DataSetLayerMediaType,
		}
		model.Layers = append(model.Layers, *layer)
	}

	// 4. package the TrainedModel
	if kitfile.Model != nil {
		modelPath, err := filesystem.VerifySubpath(options.contextDir, kitfile.Model.Path)
		if err != nil {
			return err
		}
		layer := &artifact.ModelLayer{
			Path:      modelPath,
			MediaType: constants.ModelLayerMediaType,
		}
		model.Layers = append(model.Layers, *layer)
	}

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
