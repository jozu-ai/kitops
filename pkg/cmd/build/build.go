/*
Copyright © 2024 Jozu.com
*/
package build

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"jmm/pkg/artifact"
	"jmm/pkg/lib/constants"
	"jmm/pkg/lib/storage"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"oras.land/oras-go/v2/registry"
)

const (
	DEFAULT_MODEL_FILE = "Jozufile"
)

var (
	shortDesc = `Build a model`
	longDesc  = `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`
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
		options.ModelFile = options.ContextDir + "/" + DEFAULT_MODEL_FILE
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
	jozufile := &artifact.JozuFile{}
	if err = jozufile.LoadModel(modelfile); err != nil {
		fmt.Println(err)
		return err
	}

	model := &artifact.Model{}
	model.Config = jozufile

	// 2. package the Code
	for _, code := range jozufile.Code {
		codePath, err := toAbsPath(options.ContextDir, code.Path)
		if err != nil {
			return err
		}
		layer := artifact.NewLayer(codePath, constants.CodeLayerMediaType)
		model.Layers = append(model.Layers, *layer)
	}
	// 3. package the DataSets
	datasetPath := ""
	for _, dataset := range jozufile.DataSets {
		datasetPath, err = toAbsPath(options.ContextDir, dataset.Path)
		if err != nil {
			return err
		}
		layer := artifact.NewLayer(datasetPath, constants.DataSetLayerMediaType)
		model.Layers = append(model.Layers, *layer)
	}

	// 4. package the TrainedModels
	for _, trainedModel := range jozufile.Models {
		modelPath, err := toAbsPath(options.ContextDir, trainedModel.Path)
		if err != nil {
			return err
		}
		layer := artifact.NewLayer(modelPath, constants.ModelLayerMediaType)
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
func toAbsPath(context string, relativePath string) (string, error) {

	absContext, err := filepath.Abs(context)
	if err != nil {
		fmt.Println("Error resolving base path:", err)
		return "", err
	}
	combinedPath := filepath.Join(absContext, relativePath)

	cleanPath := filepath.Clean(combinedPath)
	return cleanPath, nil

}
