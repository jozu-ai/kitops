package push

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
	shortDesc = `Uploads modelkits to a specified registry`
	longDesc  = `This command pushes modelkits to a remote registry.

The modelkits should be tagged with the target registry and repository before
they can be pushed`

	example = `# Push the latest modelkits to a remote registry
kit push registry.example.com/my-model:latest

# Push a specific version of a modelkits using a tag:
kit push registry.example.com/my-model:1.0.0`
)

type pushOptions struct {
	options.NetworkOptions
	configHome string
	modelRef   *registry.Reference
}

func (opts *pushOptions) complete(ctx context.Context, args []string) error {
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
		return fmt.Errorf("registry is required when pushing")
	}
	if len(extraTags) > 0 {
		return fmt.Errorf("reference cannot include multiple tags")
	}
	opts.modelRef = modelRef

	printConfig(opts)
	return nil
}

func PushCommand() *cobra.Command {
	opts := &pushOptions{}
	cmd := &cobra.Command{
		Use:     "push registry/repository[:tag|@digest]",
		Short:   shortDesc,
		Long:    longDesc,
		Example: example,
		Run:     runCommand(opts),
	}

	cmd.Args = cobra.ExactArgs(1)
	opts.AddNetworkFlags(cmd)

	return cmd
}

func runCommand(opts *pushOptions) func(*cobra.Command, []string) {
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

		output.Infof("Pushing %s", opts.modelRef.String())
		desc, err := PushModel(cmd.Context(), localStore, remoteRegistry, opts.modelRef)
		if err != nil {
			output.Fatalf("Failed to push: %s", err)
			return
		}
		output.Infof("Pushed %s", desc.Digest)
	}
}

func printConfig(opts *pushOptions) {
	output.Debugf("Using config path: %s", opts.configHome)
}
