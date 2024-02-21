package export

import (
	"context"
	"errors"
	"fmt"
	"jmm/pkg/lib/storage"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

var (
	flags *ExportFlags
	opts  *ExportOptions
)

type ExportFlags struct {
	Overwrite      bool
	UseHTTP        bool
	ExportConfig   bool
	ExportModels   bool
	ExportDatasets bool
	ExportCode     bool
	ExportDir      string
}

type ExportOptions struct {
	configHome  string
	storageHome string
	exportDir   string
	overwrite   bool
	exportConf  ExportConf
	modelRef    *registry.Reference
	usehttp     bool
}

type ExportConf struct {
	ExportConfig   bool
	ExportModels   bool
	ExportCode     bool
	ExportDatasets bool
}

func (opts *ExportOptions) complete(args []string) error {
	opts.configHome = viper.GetString("config")
	opts.storageHome = storage.StorageHome(opts.configHome)
	modelRef, extraTags, err := storage.ParseReference(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse reference %s: %w", args[0], err)
	}
	if len(extraTags) > 0 {
		return fmt.Errorf("can not export multiple tags")
	}
	opts.modelRef = modelRef
	opts.overwrite = flags.Overwrite
	opts.usehttp = flags.UseHTTP
	opts.exportDir = flags.ExportDir

	if !flags.ExportConfig && !flags.ExportModels && !flags.ExportCode && !flags.ExportDatasets {
		opts.exportConf.ExportConfig = true
		opts.exportConf.ExportModels = true
		opts.exportConf.ExportCode = true
		opts.exportConf.ExportDatasets = true
	} else {
		opts.exportConf.ExportConfig = flags.ExportConfig
		opts.exportConf.ExportModels = flags.ExportModels
		opts.exportConf.ExportCode = flags.ExportCode
		opts.exportConf.ExportDatasets = flags.ExportDatasets
	}

	return nil
}

func (opts *ExportOptions) validate() error {
	return nil
}

func ExportCommand() *cobra.Command {
	opts = &ExportOptions{}
	flags = &ExportFlags{}

	cmd := &cobra.Command{
		Use:   "export",
		Short: shortDesc,
		Long:  longDesc,
		Run:   runCommand(opts),
	}

	cmd.Args = cobra.ExactArgs(1)
	cmd.Flags().StringVarP(&flags.ExportDir, "dir", "d", "", "Directory to export into. Will be created if it does not exist")
	cmd.Flags().BoolVarP(&flags.Overwrite, "overwrite", "o", false, "Overwrite existing files and directories in the export dir")
	cmd.Flags().BoolVar(&flags.ExportConfig, "config", false, "Export only config file")
	cmd.Flags().BoolVar(&flags.ExportModels, "models", false, "Export only models")
	cmd.Flags().BoolVar(&flags.ExportCode, "code", false, "Export only code")
	cmd.Flags().BoolVar(&flags.ExportDatasets, "datasets", false, "Export only datasets")
	cmd.Flags().BoolVar(&flags.UseHTTP, "http", false, "Use plain HTTP when connecting to remote registries")

	return cmd
}

func runCommand(opts *ExportOptions) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := opts.complete(args); err != nil {
			fmt.Printf("Failed to process arguments: %s", err)
			return
		}
		err := opts.validate()
		if err != nil {
			fmt.Println(err)
			return
		}

		store, err := getStoreForRef(cmd.Context(), opts)
		if err != nil {
			fmt.Println(err)
			return
		}

		exportTo := opts.exportDir
		if exportTo == "" {
			exportTo = "current directory"
		}
		fmt.Printf("Exporting to %s\n", exportTo)
		err = ExportModel(cmd.Context(), store, opts.modelRef, opts)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func getStoreForRef(ctx context.Context, opts *ExportOptions) (oras.Target, error) {
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
