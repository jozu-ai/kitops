/*
Copyright © 2024 Jozu.com
*/
package build

import (
	"fmt"
	"os"
	"path"

	"kitops/pkg/artifact"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/filesystem"
	"kitops/pkg/lib/storage"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"oras.land/oras-go/v2/registry"
)

var (
	shortDesc = `Build a model`
	longDesc  = `Build a model TODO`
)

type BuildFlags struct {
	ModelFile  string
	FullTagRef string
}

type BuildOptions struct {
	ModelFile   string
	ContextDir  string
	configHome  string
	storageHome string
	modelRef    *registry.Reference
	extraRefs   []string
}

func NewCmdBuild() *cobra.Command {
	buildFlags := NewBuildFlags()

	cmd := &cobra.Command{
		Use:   "build",
		Short: shortDesc,
		Long:  longDesc,
		Run: func(cmd *cobra.Command, args []string) {
			options, err := buildFlags.ToOptions()
			if err != nil {
				fmt.Println(err)
				return
			}
			err = options.Complete(cmd, args)
			if err != nil {
				fmt.Println(err)
				return
			}
			err = options.Validate()
			if err != nil {
				fmt.Println(err)
				return
			}
			err = options.RunBuild()
			if err != nil {
				fmt.Println(err)
				return
			}
		},
	}
	buildFlags.AddFlags(cmd)
	return cmd
}

func (options *BuildOptions) Complete(cmd *cobra.Command, argsIn []string) error {
	options.ContextDir = argsIn[0]
	if options.ModelFile == "" {
		options.ModelFile = path.Join(options.ContextDir, constants.DefaultModelFileName)
	}
	options.configHome = viper.GetString("config")
	fmt.Println("config: ", options.configHome)
	options.storageHome = storage.StorageHome(options.configHome)
	return nil
}

func (o *BuildOptions) Validate() error {
	return nil
}

func (options *BuildOptions) RunBuild() error {
	fmt.Println("build called")
	// 1. Read the model file
	modelfile, err := os.Open(options.ModelFile)
	if err != nil {
		return err
	}
	defer modelfile.Close()
	kitfile := &artifact.KitFile{}
	if err = kitfile.LoadModel(modelfile); err != nil {
		fmt.Println(err)
		return err
	}

	model := &artifact.Model{}
	model.Config = kitfile

	// 2. package the Code
	for _, code := range kitfile.Code {
		codePath, err := filesystem.VerifySubpath(options.ContextDir, code.Path)
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
		datasetPath, err := filesystem.VerifySubpath(options.ContextDir, dataset.Path)
		if err != nil {
			return err
		}
		layer := &artifact.ModelLayer{
			BaseDir:   datasetPath,
			MediaType: constants.DataSetLayerMediaType,
		}
		model.Layers = append(model.Layers, *layer)
	}

	// 4. package the TrainedModels
	for _, trainedModel := range kitfile.Models {
		modelPath, err := filesystem.VerifySubpath(options.ContextDir, trainedModel.Path)
		if err != nil {
			return err
		}
		layer := &artifact.ModelLayer{
			BaseDir:   modelPath,
			MediaType: constants.ModelLayerMediaType,
		}
		model.Layers = append(model.Layers, *layer)
	}

	modelStorePath := options.storageHome
	repo := ""
	tag := ""
	if options.modelRef != nil {
		repo = path.Join(options.modelRef.Registry, options.modelRef.Repository)
		tag = options.modelRef.Reference
	}
	store := storage.NewLocalStore(modelStorePath, repo)
	manifestDesc, err := store.SaveModel(model, tag)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, tag := range options.extraRefs {
		if err := store.TagModel(*manifestDesc, tag); err != nil {
			return err
		}
	}

	fmt.Println("Model saved: ", manifestDesc.Digest)

	return nil
}

func (o *BuildFlags) ToOptions() (*BuildOptions, error) {
	options := &BuildOptions{}
	if o.ModelFile != "" {
		options.ModelFile = o.ModelFile
	}
	if o.FullTagRef != "" {
		modelRef, extraRefs, err := storage.ParseReference(o.FullTagRef)
		if err != nil {
			return nil, err
		}
		options.modelRef = modelRef
		options.extraRefs = extraRefs
	}
	return options, nil
}

func (flags *BuildFlags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&flags.ModelFile, "file", "f", "", "Path to the model file")
	cmd.Flags().StringVarP(&flags.FullTagRef, "tag", "t", "", "Tag for the model. Example: -t registry/repository:tag1,tag2")
	cmd.Args = cobra.ExactArgs(1)
}

func NewBuildFlags() *BuildFlags {
	return &BuildFlags{}
}
