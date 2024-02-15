package export

import (
	"fmt"
	"jmm/pkg/lib/storage"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"oras.land/oras-go/v2/content/oci"
	"oras.land/oras-go/v2/registry"
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
	UseHTTP   bool
	ExportDir string
}

type ExportOptions struct {
	usehttp     bool
	configHome  string
	storageHome string
	exportDir   string
	modelRef    *registry.Reference
}

func (opts *ExportOptions) complete(args []string) error {
	opts.configHome = viper.GetString("config")
	opts.storageHome = storage.StorageHome(opts.configHome)
	modelRef, extraTags, err := storage.ParseReference(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse reference %s: %w", modelRef, err)
	}
	if len(extraTags) > 0 {
		return fmt.Errorf("can not export multiple tags")
	}
	opts.modelRef = modelRef
	opts.usehttp = flags.UseHTTP
	opts.exportDir = flags.ExportDir

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
	cmd.Flags().BoolVar(&flags.UseHTTP, "http", false, "Use plain HTTP when connecting to remote registries")
	cmd.Flags().StringVarP(&flags.ExportDir, "dir", "d", "", "Directory to export into. Will be created if it does not exist")

	return cmd
}

func runCommand(opts *ExportOptions) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := opts.complete(args); err != nil {
			fmt.Printf("Failed to process arguments: %s", err)
		}
		err := opts.validate()
		if err != nil {
			fmt.Println(err)
			return
		}
		store, err := oci.New(storage.LocalStorePath(opts.storageHome, opts.modelRef))
		if err != nil {
			fmt.Printf("failed to read local storage: %s\n", err)
			return
		}

		if _, err := store.Resolve(cmd.Context(), opts.modelRef.Reference); err != nil {
			fmt.Printf("reference %s not found in local store\n", opts.modelRef.String())
			return
		}

		fmt.Printf("Exporting to %s\n", opts.exportDir)
		err = ExportModel(cmd.Context(), store, opts.modelRef, opts.exportDir)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
