/*
Copyright © 2024 Jozu.com
*/
package pack

import (
	"context"
	"fmt"
	"strings"

	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/filesystem"
	"kitops/pkg/lib/repo"
	"kitops/pkg/output"

	"github.com/spf13/cobra"
	"oras.land/oras-go/v2/registry"
)

const (
	shortDesc = `Pack a modelkit`
	longDesc  = `Pack a modelkit from a kitfile using the given context directory.

The packing process involves taking the configuration and resources defined in
your kitfile and using them to create a modelkit. This modelkit is then stored
in your local registry, making it readily available for further actions such
as pushing to a remote registry for collaboration.

Unless a different location is specified, this command looks for the kitfile
at the root of the provided context directory. Any relative paths defined
within the kitfile are interpreted as being relative to this context
directory.`

	examples = `# Pack a modelkit using the kitfile in the current directory
kit pack .

# Pack a modelkit with a specific kitfile and tag
kit pack . -f /path/to/your/Kitfile -t registry/repository:modelv1`
)

type packOptions struct {
	modelFile   string
	contextDir  string
	configHome  string
	storageHome string
	fullTagRef  string
	modelRef    *registry.Reference
	extraRefs   []string
}

func PackCommand() *cobra.Command {
	opts := &packOptions{}

	cmd := &cobra.Command{
		Use:     "pack [flags] DIRECTORY",
		Short:   shortDesc,
		Long:    longDesc,
		Example: examples,
		Run:     runCommand(opts),
	}
	cmd.Flags().StringVarP(&opts.modelFile, "file", "f", "", "Specifies the path to the Kitfile if it's not located at the root of the context directory")
	cmd.Flags().StringVarP(&opts.fullTagRef, "tag", "t", "", "Assigns one or more tags to the built modelkit. Example: -t registry/repository:tag1,tag2")
	cmd.Args = cobra.ExactArgs(1)
	return cmd
}

func runCommand(opts *packOptions) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := opts.complete(cmd.Context(), args)
		if err != nil {
			output.Fatalf("Failed to process configuration: %s", err)
			return
		}
		err = runPack(cmd.Context(), opts)
		if err != nil {
			output.Fatalf("Failed to pack model kit: %s", err)
			return
		}
	}
}

func (opts *packOptions) complete(ctx context.Context, args []string) error {
	opts.contextDir = args[0]

	if opts.modelFile == "" {
		opts.modelFile = filesystem.FindKitfileInPath(opts.contextDir)
	}

	configHome, ok := ctx.Value(constants.ConfigKey{}).(string)
	if !ok {
		return fmt.Errorf("default config path not set on command context")
	}
	opts.configHome = configHome
	opts.storageHome = constants.StoragePath(opts.configHome)

	if opts.fullTagRef != "" {
		modelRef, extraRefs, err := repo.ParseReference(opts.fullTagRef)
		if err != nil {
			return fmt.Errorf("failed to parse reference %s: %w", opts.fullTagRef, err)
		}
		opts.modelRef = modelRef
		opts.extraRefs = extraRefs
	} else {
		opts.modelRef = repo.DefaultReference()
	}
	printConfig(opts)
	return nil
}

func printConfig(opts *packOptions) {
	output.Debugf("Using storage path: %s", opts.storageHome)
	output.Debugf("Context dir: %s", opts.contextDir)
	output.Debugf("Model file: %s", opts.modelFile)
	if opts.modelRef != nil {
		output.Debugf("Packing %s", opts.modelRef.String())
	} else {
		output.Debugln("No tag or reference specified")
	}
	if len(opts.extraRefs) > 0 {
		output.Debugf("Additional tags: %s", strings.Join(opts.extraRefs, ", "))
	}
}
