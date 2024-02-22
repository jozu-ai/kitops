package pull

import (
	"context"
	"fmt"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/storage"

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
			fmt.Printf("Failed to process arguments: %s", err)
		}
		remoteRegistry, err := remote.NewRegistry(opts.modelRef.Registry)
		if err != nil {
			fmt.Println(err)
			return
		}
		if opts.usehttp {
			remoteRegistry.PlainHTTP = true
		}

		localStorePath := storage.LocalStorePath(opts.storageHome, opts.modelRef)
		localStore, err := oci.New(localStorePath)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Pulling %s\n", opts.modelRef.String())
		desc, err := pullModel(cmd.Context(), remoteRegistry, localStore, opts.modelRef)
		if err != nil {
			fmt.Printf("Failed to pull: %s\n", err)
			return
		}
		fmt.Printf("Pulled %s\n", desc.Digest)
	}
}
