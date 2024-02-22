package push

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
	shortDesc = `Push model to registry`
	longDesc  = `Push model to registry TODO`
)

type pushFlags struct {
	UseHTTP bool
}

type pushOptions struct {
	usehttp     bool
	configHome  string
	storageHome string
	modelRef    *registry.Reference
}

func (opts *pushOptions) complete(ctx context.Context, flags *pushFlags, args []string) error {
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
	opts.usehttp = flags.UseHTTP
	return nil
}

func PushCommand() *cobra.Command {
	flags := &pushFlags{}

	cmd := &cobra.Command{
		Use:   "push",
		Short: shortDesc,
		Long:  longDesc,
		Run:   runCommand(flags),
	}

	cmd.Args = cobra.ExactArgs(1)
	cmd.Flags().BoolVar(&flags.UseHTTP, "http", false, "Use plain HTTP when connecting to remote registries")
	return cmd
}

func runCommand(flags *pushFlags) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		opts := &pushOptions{}
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

		fmt.Printf("Pushing %s\n", opts.modelRef.String())
		desc, err := PushModel(cmd.Context(), localStore, remoteRegistry, opts.modelRef)
		if err != nil {
			fmt.Printf("Failed to push: %s\n", err)
			return
		}
		fmt.Printf("Pushed %s\n", desc.Digest)
	}
}
