package pull

import (
	"fmt"
	"jmm/pkg/lib/storage"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"oras.land/oras-go/v2/content/oci"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
)

const (
	shortDesc = `Pull model from registry`
	longDesc  = `Pull model from registry TODO`
)

var (
	flags *PullFlags
	opts  *PullOptions
)

type PullFlags struct {
	UseHTTP bool
}

type PullOptions struct {
	usehttp     bool
	configHome  string
	storageHome string
	modelRef    *registry.Reference
}

func (opts *PullOptions) complete(args []string) error {
	opts.configHome = viper.GetString("config")
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

func (opts *PullOptions) validate() error {
	return nil
}

func PullCommand() *cobra.Command {
	opts = &PullOptions{}
	flags = &PullFlags{}

	cmd := &cobra.Command{
		Use:   "pull",
		Short: shortDesc,
		Long:  longDesc,
		Run:   runCommand(opts),
	}

	cmd.Args = cobra.ExactArgs(1)
	cmd.Flags().BoolVar(&flags.UseHTTP, "http", false, "Use plain HTTP when connecting to remote registries")
	return cmd
}

func runCommand(opts *PullOptions) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := opts.complete(args); err != nil {
			fmt.Printf("Failed to process arguments: %s", err)
		}
		err := opts.validate()
		if err != nil {
			fmt.Println(err)
			return
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
		desc, err := doPull(cmd.Context(), remoteRegistry, localStore, opts.modelRef)
		if err != nil {
			fmt.Printf("Failed to pull: %s\n", err)
			return
		}
		fmt.Printf("Pulled %s\n", desc.Digest)
	}
}
