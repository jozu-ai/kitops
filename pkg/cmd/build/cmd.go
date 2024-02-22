/*
Copyright © 2024 Jozu.com
*/
package build

import (
	"context"
	"fmt"
	"path"

	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/storage"

	"github.com/spf13/cobra"
	"oras.land/oras-go/v2/registry"
)

var (
	shortDesc = `Build a model`
	longDesc  = `Build a model TODO`
)

type buildFlags struct {
	modelFile  string
	fullTagRef string
}

type buildOptions struct {
	modelFile   string
	contextDir  string
	configHome  string
	storageHome string
	modelRef    *registry.Reference
	extraRefs   []string
}

func BuildCommand() *cobra.Command {
	flags := &buildFlags{}

	cmd := &cobra.Command{
		Use:   "build",
		Short: shortDesc,
		Long:  longDesc,
		Run:   runCommand(flags),
	}

	cmd.Flags().StringVarP(&flags.modelFile, "file", "f", "", "Path to the model file")
	cmd.Flags().StringVarP(&flags.fullTagRef, "tag", "t", "", "Tag for the model. Example: -t registry/repository:tag1,tag2")
	cmd.Args = cobra.ExactArgs(1)
	return cmd
}

func runCommand(flags *buildFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		opts := &buildOptions{}
		err := opts.complete(cmd.Context(), flags, args)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = RunBuild(opts)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (opts *buildOptions) complete(ctx context.Context, flags *buildFlags, args []string) error {
	opts.contextDir = args[0]

	opts.modelFile = flags.modelFile
	if opts.modelFile == "" {
		opts.modelFile = path.Join(opts.contextDir, constants.DefaultModelFileName)
	}

	configHome, ok := ctx.Value(constants.ConfigKey{}).(string)
	if !ok {
		return fmt.Errorf("default config path not set on command context")
	}
	opts.configHome = configHome
	opts.storageHome = storage.StorageHome(opts.configHome)

	if flags.fullTagRef != "" {
		modelRef, extraRefs, err := storage.ParseReference(flags.fullTagRef)
		if err != nil {
			return fmt.Errorf("failed to parse reference %s: %w", flags.fullTagRef, err)
		}
		opts.modelRef = modelRef
		opts.extraRefs = extraRefs
	}
	return nil
}
