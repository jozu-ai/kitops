package export

import (
	"context"
	"errors"
	"fmt"
	"kitops/pkg/cmd/options"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/repo"
	"kitops/pkg/output"

	"github.com/spf13/cobra"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/oci"
	"oras.land/oras-go/v2/errdef"
	"oras.land/oras-go/v2/registry"
)

const (
	shortDesc = `Produce the components from a modelkit on the local filesystem`
	longDesc  = `Produces all or selected components of a modelkit on the local filesystem.
	
This command exports of a modelkit's components, including models, code, datasets, and configuration files,
to a specified directory on the local filesystem. By default, it attempts to find the modelkit in local storage;
if not found, it searches the remote registry and retrieves it. This process ensures that the necessary components
are always available for export, optimizing for efficiency by fetching only specified components from the remote 
registry when necessary`

	example = `  # Export all components of a modelkit to the current directory
  kit export myrepo/my-model:latest -d /path/to/export

  # Export only the model and datasets of a modelkit to a specified directory
  kit export myrepo/my-model:latest --model --datasets -d /path/to/export

  # Export a modelkit from a remote registry with overwrite enabled
  kit export registry.example.com/myrepo/my-model:latest -o -d /path/to/export`
)

type exportOptions struct {
	options.NetworkOptions
	configHome string
	exportDir  string
	exportConf exportConf
	modelRef   *registry.Reference
	overwrite  bool
}

type exportConf struct {
	exportConfig   bool
	exportModels   bool
	exportCode     bool
	exportDatasets bool
}

func (opts *exportOptions) complete(ctx context.Context, args []string) error {
	configHome, ok := ctx.Value(constants.ConfigKey{}).(string)
	if !ok {
		return fmt.Errorf("default config path not set on command context")
	}
	opts.configHome = configHome
	modelRef, extraTags, err := repo.ParseReference(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse reference %s: %w", args[0], err)
	}
	if len(extraTags) > 0 {
		return fmt.Errorf("can not export multiple tags")
	}
	opts.modelRef = modelRef

	conf := opts.exportConf
	if !conf.exportConfig && !conf.exportModels && !conf.exportCode && !conf.exportDatasets {
		opts.exportConf.exportConfig = true
		opts.exportConf.exportModels = true
		opts.exportConf.exportCode = true
		opts.exportConf.exportDatasets = true
	}

	printConfig(opts)
	return nil
}

func ExportCommand() *cobra.Command {
	opts := &exportOptions{}

	cmd := &cobra.Command{
		Use:     "export [registry/]repository[:tag|@digest]",
		Short:   shortDesc,
		Long:    longDesc,
		Example: example,
		Run:     runCommand(opts),
	}

	cmd.Args = cobra.ExactArgs(1)
	cmd.Flags().StringVarP(&opts.exportDir, "dir", "d", "", "The target directory to export components into. This directory will be created if it does not exist")
	cmd.Flags().BoolVarP(&opts.overwrite, "overwrite", "o", false, "Overwrites existing files and directories in the target export directory without prompting")
	cmd.Flags().BoolVar(&opts.exportConf.exportConfig, "config", false, "Export only config file")
	cmd.Flags().BoolVar(&opts.exportConf.exportModels, "model", false, "Export only model")
	cmd.Flags().BoolVar(&opts.exportConf.exportCode, "code", false, "Export only code")
	cmd.Flags().BoolVar(&opts.exportConf.exportDatasets, "datasets", false, "Export only datasets")
	opts.AddNetworkFlags(cmd)

	return cmd
}

func runCommand(opts *exportOptions) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := opts.complete(cmd.Context(), args); err != nil {
			output.Fatalf("Failed to process arguments: %s", err)
		}

		if opts.modelRef.Reference == "" {
			output.Fatalf("Invalid reference: exporting requires a tag or digest")
		}
		store, err := getStoreForRef(cmd.Context(), opts)
		if err != nil {
			ref := repo.FormatRepositoryForDisplay(opts.modelRef.String())
			output.Fatalf("Failed to find reference %s: %s", ref, err)
		}

		exportTo := opts.exportDir
		if exportTo == "" {
			exportTo = "current directory"
		}
		output.Infof("Exporting to %s", exportTo)
		err = exportModel(cmd.Context(), store, opts.modelRef, opts)
		if err != nil {
			output.Fatalln(err)
		}
	}
}

func getStoreForRef(ctx context.Context, opts *exportOptions) (oras.Target, error) {
	storageHome := constants.StoragePath(opts.configHome)
	localStore, err := oci.New(repo.RepoPath(storageHome, opts.modelRef))
	if err != nil {
		return nil, fmt.Errorf("failed to read local storage: %s\n", err)
	}

	if _, err := localStore.Resolve(ctx, opts.modelRef.Reference); err == nil {
		// Reference is present in local storage
		return localStore, nil
	}

	if opts.modelRef.Registry == repo.DefaultRegistry {
		return nil, fmt.Errorf("not found")
	}
	// Not in local storage, check remote
	remoteRegistry, err := repo.NewRegistry(opts.modelRef.Registry, &repo.RegistryOptions{
		PlainHTTP:       opts.PlainHTTP,
		SkipTLSVerify:   !opts.TlsVerify,
		CredentialsPath: constants.CredentialsPath(opts.configHome),
	})
	if err != nil {
		return nil, fmt.Errorf("could not resolve registry %s: %w", opts.modelRef.Registry, err)
	}

	repo, err := remoteRegistry.Repository(ctx, opts.modelRef.Repository)
	if err != nil {
		return nil, fmt.Errorf("could not resolve repository %s in registry %s", opts.modelRef.Repository, opts.modelRef.Registry)
	}
	if _, err := repo.Resolve(ctx, opts.modelRef.Reference); err != nil {
		if errors.Is(err, errdef.ErrNotFound) {
			return nil, fmt.Errorf("reference %s is not present in local storage and could not be found in remote", opts.modelRef.String())
		}
		return nil, fmt.Errorf("unexpected error retrieving reference from remote: %w", err)
	}

	return repo, nil
}

func printConfig(opts *exportOptions) {
	output.Debugf("Using config path: %s", opts.configHome)
	output.Debugf("Overwrite: %t", opts.overwrite)
	output.Debugf("Exporting %s", opts.modelRef.String())
}
