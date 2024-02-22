package pull

import (
	"context"
	"fmt"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/storage"
	"kitops/pkg/output"

	"github.com/spf13/cobra"
	"oras.land/oras-go/v2/content/oci"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
)

const (
	shortDesc = `Pull model from registry`
	longDesc  = `Pull model from registry TODO`
)

type pullFlags struct {
	useHTTP bool
}

type pullOptions struct {
	usehttp     bool
	configHome  string
	storageHome string
	modelRef    *registry.Reference
}

func (opts *pullOptions) complete(ctx context.Context, flags *pullFlags, args []string) error {
	configHome, ok := ctx.Value(constants.ConfigKey{}).(string)
	if !ok {
		return fmt.Errorf("default config path not set on command context")
	}
	opts.configHome = configHome
	opts.storageHome = storage.StorageHome(opts.configHome)

	modelRef, extraTags, err := storage.ParseReference(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse reference %s: %w", modelRef, err)
	}
	if modelRef.Registry == "localhost" {
		return fmt.Errorf("registry is required when pulling")
	}
	if len(extraTags) > 0 {
		return fmt.Errorf("reference cannot include multiple tags")
	}
	opts.modelRef = modelRef
	opts.usehttp = flags.useHTTP

	printConfig(opts)
	return nil
}

func PullCommand() *cobra.Command {
	flags := &pullFlags{}

	cmd := &cobra.Command{
		Use:   "pull",
		Short: shortDesc,
		Long:  longDesc,
		Run:   runCommand(flags),
	}

	cmd.Args = cobra.ExactArgs(1)
	cmd.Flags().BoolVar(&flags.useHTTP, "http", false, "Use plain HTTP when connecting to remote registries")
	return cmd
}

func runCommand(flags *pullFlags) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		opts := &pullOptions{}
		if err := opts.complete(cmd.Context(), flags, args); err != nil {
			output.Fatalf("Failed to process arguments: %s", err)
		}
		remoteRegistry, err := remote.NewRegistry(opts.modelRef.Registry)
		if err != nil {
			output.Fatalln(err)
		}
		if opts.usehttp {
			remoteRegistry.PlainHTTP = true
		}

		localStorePath := storage.LocalStorePath(opts.storageHome, opts.modelRef)
		localStore, err := oci.New(localStorePath)
		if err != nil {
			output.Fatalln(err)
		}

		output.Infof("Pulling %s", opts.modelRef.String())
		desc, err := pullModel(cmd.Context(), remoteRegistry, localStore, opts.modelRef)
		if err != nil {
			output.Fatalf("Failed to pull: %s", err)
			return
		}
		output.Infof("Pulled %s", desc.Digest)
	}
}

func printConfig(opts *pullOptions) {
	output.Debugf("Using storage path: %s", opts.storageHome)
}
