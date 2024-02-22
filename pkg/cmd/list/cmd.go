package list

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
	shortDesc = `List model kits in a repository`
	longDesc  = `List model kits TODO`
)

var (
	flags *ListFlags
	opts  *ListOptions
)

type ListFlags struct {
	UseHTTP bool
}

type ListOptions struct {
	configHome  string
	storageHome string
	remoteRef   *registry.Reference
	usehttp     bool
}

func (opts *ListOptions) complete(flags *ListFlags, args []string) error {
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

func (opts *ListOptions) validate() error {
	return nil
}

// ListCommand represents the models command
func ListCommand() *cobra.Command {
	flags = &ListFlags{}
	opts = &ListOptions{}

	cmd := &cobra.Command{
		Use:   "list [repository]",
		Short: shortDesc,
		Long:  longDesc,
		Run:   RunCommand(opts),
	}

	cmd.Args = cobra.MaximumNArgs(1)
	cmd.Flags().BoolVar(&flags.UseHTTP, "http", false, "Use plain HTTP when connecting to remote registries")
	return cmd
}

func RunCommand(options *ListOptions) func(*cobra.Command, []string) {
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
			lines, err := listLocalKits(opts.storageHome)
			if err != nil {
				fmt.Println(err)
				return
			}
			allInfoLines = lines
		} else {
			lines, err := listRemoteKits(cmd.Context(), opts.remoteRef, opts.usehttp)
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
	fmt.Fprintln(tw, listTableHeader)
	for _, line := range lines {
		fmt.Fprintln(tw, line)
	}
	tw.Flush()
}
