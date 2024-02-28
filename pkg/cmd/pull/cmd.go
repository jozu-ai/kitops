package pull

import (
	"context"
	"fmt"
	"kitops/pkg/cmd/options"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/repo"
	"kitops/pkg/output"

	"github.com/spf13/cobra"
	"oras.land/oras-go/v2/content/oci"
	"oras.land/oras-go/v2/registry"
)

const (
	shortDesc = `Pull model from registry`
	longDesc  = `Pull model from registry TODO`
)

type pullOptions struct {
	options.NetworkOptions
	configHome string
	modelRef   *registry.Reference
}

func (opts *pullOptions) complete(ctx context.Context, args []string) error {
	configHome, ok := ctx.Value(constants.ConfigKey{}).(string)
	if !ok {
		return fmt.Errorf("default config path not set on command context")
	}
	opts.configHome = configHome

	modelRef, extraTags, err := repo.ParseReference(args[0])
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

	printConfig(opts)
	return nil
}

func PullCommand() *cobra.Command {
	opts := &pullOptions{}
	cmd := &cobra.Command{
		Use:   "pull",
		Short: shortDesc,
		Long:  longDesc,
		Run:   runCommand(opts),
	}

	cmd.Args = cobra.ExactArgs(1)
	opts.AddNetworkFlags(cmd)

	return cmd
}

func runCommand(opts *pullOptions) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := opts.complete(cmd.Context(), args); err != nil {
			output.Fatalf("Failed to process arguments: %s", err)
		}
		remoteRegistry, err := repo.NewRegistry(opts.modelRef.Registry, &repo.RegistryOptions{
			PlainHTTP:       opts.PlainHTTP,
			SkipTLSVerify:   !opts.TlsVerify,
			CredentialsPath: constants.CredentialsPath(opts.configHome),
		})
		if err != nil {
			output.Fatalln(err)
		}

		storageHome := constants.StoragePath(opts.configHome)
		localStorePath := repo.RepoPath(storageHome, opts.modelRef)
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
	output.Debugf("Using config path: %s", opts.configHome)
}
