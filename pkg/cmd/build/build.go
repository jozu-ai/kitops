/*
Copyright © 2024 Jozu.com
*/
package build

import (
	"fmt"
	"os"

	"jmm/pkg/artifact"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const DEFAULT_MODEL_FILE = "Jozufile"

var (
	shortDesc = `Build a model`
	longDesc  = `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`
)

type BuildFlags struct {
	ModelFile string
}

type BuildOptions struct {
	ModelFile  string
	ContextDir string
	JozuHome   string
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
	fmt.Println("config: ", viper.GetString("config"))
	options.JozuHome = viper.GetString("config")
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

	// 2. Run the build steps from the model file

	// 3. Tar the build context and push to local registry
	layer := artifact.NewLayer(options.ContextDir)
	model := &artifact.Model{}
	model.Layers = append(model.Layers, layer)
	model.Config = jozufile

	store := artifact.NewArtifactStore(options.JozuHome)
	var manifest *v1.Manifest
	manifest, err = store.SaveModel(model)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Model saved: ", manifest.Config.Digest)
	return nil
}

func (o *BuildFlags) ToOptions() (*BuildOptions, error) {
	options := &BuildOptions{}
	if o.ModelFile != "" {
		options.ModelFile = o.ModelFile
	}
	return options, nil
}

func (flags *BuildFlags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&flags.ModelFile, "file", "f", "", "Path to the model file")
	cmd.Args = cobra.ExactArgs(1)

}

func NewBuildFlags() *BuildFlags {
	return &BuildFlags{}
}
