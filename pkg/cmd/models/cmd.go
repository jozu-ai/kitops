package models

import (
	"fmt"
	"kitops/pkg/lib/storage"
	"os"
	"path"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"oras.land/oras-go/v2/registry"
)

const (
	shortDesc = `List models`
	longDesc  = `List models TODO`
)

var (
	flags *ModelsFlags
	opts  *ModelsOptions
)

type ModelsFlags struct {
	UseHTTP bool
}

type ModelsOptions struct {
	configHome  string
	storageHome string
	remoteRef   *registry.Reference
	usehttp     bool
}

func (opts *ModelsOptions) complete(flags *ModelsFlags, args []string) error {
	opts.configHome = viper.GetString("config")
	opts.storageHome = path.Join(opts.configHome, "storage")
	if len(args) > 0 {
		remoteRef, extraTags, err := storage.ParseReference(args[0])
		if err != nil {
			return fmt.Errorf("invalid reference: %w", err)
		}
		if len(extraTags) > 0 {
			return fmt.Errorf("repository cannot reference multiple tags")
		}
		opts.remoteRef = remoteRef
	}
	opts.usehttp = flags.UseHTTP
	return nil
}

func (opts *ModelsOptions) validate() error {
	return nil
}

// ModelsCommand represents the models command
func ModelsCommand() *cobra.Command {
	flags = &ModelsFlags{}
	opts = &ModelsOptions{}

	cmd := &cobra.Command{
		Use:   "models [repository]",
		Short: shortDesc,
		Long:  longDesc,
		Run:   RunCommand(opts),
	}

	cmd.Args = cobra.MaximumNArgs(1)
	cmd.Flags().BoolVar(&flags.UseHTTP, "http", false, "Use plain HTTP when connecting to remote registries")
	return cmd
}

func RunCommand(options *ModelsOptions) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := options.complete(flags, args); err != nil {
			fmt.Printf("Failed to parse argument: %s", err)
			return
		}
		if err := options.validate(); err != nil {
			fmt.Println(err)
			return
		}

		var allInfoLines []string
		if opts.remoteRef == nil {
			lines, err := listLocalModels(opts.storageHome)
			if err != nil {
				fmt.Println(err)
				return
			}
			allInfoLines = lines
		} else {
			lines, err := listRemoteModels(cmd.Context(), opts.remoteRef, opts.usehttp)
			if err != nil {
				fmt.Println(err)
				return
			}
			allInfoLines = lines
		}

		printSummary(allInfoLines)

	}
}

func printSummary(lines []string) {
	tw := tabwriter.NewWriter(os.Stdout, 0, 2, 3, ' ', 0)
	fmt.Fprintln(tw, ModelsTableHeader)
	for _, line := range lines {
		fmt.Fprintln(tw, line)
	}
	tw.Flush()
}
