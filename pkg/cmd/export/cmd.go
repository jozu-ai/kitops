package export

import (
	"context"
	"errors"
	"fmt"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/storage"
	"kitops/pkg/output"

	"github.com/spf13/cobra"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/oci"
	"oras.land/oras-go/v2/errdef"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
)

const (
	shortDesc = `Export model from registry`
	longDesc  = `Export model from registry TODO`
)

type exportFlags struct {
	overwrite      bool
	useHTTP        bool
	exportConfig   bool
	exportModels   bool
	exportDatasets bool
	exportCode     bool
	exportDir      string
}

type exportOptions struct {
	configHome  string
	storageHome string
	exportDir   string
	overwrite   bool
	exportConf  exportConf
	modelRef    *registry.Reference
	usehttp     bool
}

type exportConf struct {
	exportConfig   bool
	exportModels   bool
	exportCode     bool
	exportDatasets bool
}

func (opts *exportOptions) complete(ctx context.Context, flags *exportFlags, args []string) error {
	configHome, ok := ctx.Value(constants.ConfigKey{}).(string)
	if !ok {
		return fmt.Errorf("default config path not set on command context")
	}
	opts.configHome = configHome
	opts.storageHome = storage.StorageHome(opts.configHome)
	modelRef, extraTags, err := storage.ParseReference(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse reference %s: %w", args[0], err)
	}
	if len(extraTags) > 0 {
		return fmt.Errorf("can not export multiple tags")
	}
	opts.modelRef = modelRef
	opts.overwrite = flags.overwrite
	opts.usehttp = flags.useHTTP
	opts.exportDir = flags.exportDir

	if !flags.exportConfig && !flags.exportModels && !flags.exportCode && !flags.exportDatasets {
		opts.exportConf.exportConfig = true
		opts.exportConf.exportModels = true
		opts.exportConf.exportCode = true
		opts.exportConf.exportDatasets = true
	} else {
		opts.exportConf.exportConfig = flags.exportConfig
		opts.exportConf.exportModels = flags.exportModels
		opts.exportConf.exportCode = flags.exportCode
		opts.exportConf.exportDatasets = flags.exportDatasets
	}

	printConfig(opts)
	return nil
}

func ExportCommand() *cobra.Command {
	flags := &exportFlags{}

	cmd := &cobra.Command{
		Use:   "export",
		Short: shortDesc,
		Long:  longDesc,
		Run:   runCommand(flags),
	}

	cmd.Args = cobra.ExactArgs(1)
	cmd.Flags().StringVarP(&flags.exportDir, "dir", "d", "", "Directory to export into. Will be created if it does not exist")
	cmd.Flags().BoolVarP(&flags.overwrite, "overwrite", "o", false, "Overwrite existing files and directories in the export dir")
	cmd.Flags().BoolVar(&flags.exportConfig, "config", false, "Export only config file")
	cmd.Flags().BoolVar(&flags.exportModels, "models", false, "Export only models")
	cmd.Flags().BoolVar(&flags.exportCode, "code", false, "Export only code")
	cmd.Flags().BoolVar(&flags.exportDatasets, "datasets", false, "Export only datasets")
	cmd.Flags().BoolVar(&flags.useHTTP, "http", false, "Use plain HTTP when connecting to remote registries")

	return cmd
}

func runCommand(flags *exportFlags) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		opts := &exportOptions{}
		if err := opts.complete(cmd.Context(), flags, args); err != nil {
			output.Fatalf("Failed to process arguments: %s", err)
		}

		store, err := getStoreForRef(cmd.Context(), opts)
		if err != nil {
			output.Fatalln(err)
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
	localStore, err := oci.New(storage.LocalStorePath(opts.storageHome, opts.modelRef))
	if err != nil {
		return nil, fmt.Errorf("failed to read local storage: %s\n", err)
	}

	if _, err := localStore.Resolve(ctx, opts.modelRef.Reference); err == nil {
		// Reference is present in local storage
		return localStore, nil
	}

	// Not in local storage, check remote
	remoteRegistry, err := remote.NewRegistry(opts.modelRef.Registry)
	if err != nil {
		return nil, fmt.Errorf("could not resolve registry %s: %w", opts.modelRef.Registry, err)
	}
	if opts.usehttp {
		remoteRegistry.PlainHTTP = true
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
	output.Debugf("Using storage path: %s", opts.storageHome)
	output.Debugf("Overwrite: %t", opts.overwrite)
	output.Debugf("Exporting %s", opts.modelRef.String())
}
