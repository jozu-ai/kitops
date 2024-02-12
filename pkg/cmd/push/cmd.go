package push

import (
	"fmt"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"oras.land/oras-go/v2/content/oci"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
)

const (
	shortDesc = `Push model to registry`
	longDesc  = `Push model to registry TODO`
)

var (
	flags *PushFlags
	opts  *PushOptions
)

type PushFlags struct {
	UseHTTP bool
}

type PushOptions struct {
	usehttp     bool
	configHome  string
	storageHome string
	modelRef    *registry.Reference
}

func (opts *PushOptions) complete(args []string) error {
	opts.configHome = viper.GetString("config")
	opts.storageHome = path.Join(opts.configHome, "storage")
	modelRef, err := registry.ParseReference(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse reference %s: %w", modelRef, err)
	}
	opts.modelRef = &modelRef
	opts.usehttp = flags.UseHTTP
	return nil
}

func (opts *PushOptions) validate() error {
	return nil
}

func PushCommand() *cobra.Command {
	opts = &PushOptions{}
	flags = &PushFlags{}

	cmd := &cobra.Command{
		Use:   "push",
		Short: shortDesc,
		Long:  longDesc,
		Run:   runCommand(opts),
	}

	cmd.Args = cobra.ExactArgs(1)
	cmd.Flags().BoolVar(&flags.UseHTTP, "http", false, "Push to http registry")
	return cmd
}

func runCommand(opts *PushOptions) func(*cobra.Command, []string) {
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

		localStorePath := path.Join(opts.storageHome, opts.modelRef.Registry, opts.modelRef.Repository)
		localStore, err := oci.New(localStorePath)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Pushing %s\n", opts.modelRef.String())
		desc, err := doPush(cmd.Context(), localStore, remoteRegistry, opts.modelRef)
		if err != nil {
			fmt.Printf("Failed to push: %s\n", err)
		}
		fmt.Printf("Pushed %s\n", desc.Digest)
	}
}
