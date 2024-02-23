package list

import (
	"context"
	"fmt"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/storage"
	"kitops/pkg/output"
	"os"
	"path"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"oras.land/oras-go/v2/registry"
)

const (
	shortDesc = `List model kits in a repository`
	longDesc  = `List model kits TODO`
)

type listFlags struct {
	useHTTP bool
}

type listOptions struct {
	configHome  string
	storageHome string
	remoteRef   *registry.Reference
	usehttp     bool
}

func (opts *listOptions) complete(ctx context.Context, flags *listFlags, args []string) error {
	configHome, ok := ctx.Value(constants.ConfigKey{}).(string)
	if !ok {
		return fmt.Errorf("default config path not set on command context")
	}
	opts.configHome = configHome
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
	opts.usehttp = flags.useHTTP
	return nil
}

// ListCommand represents the models command
func ListCommand() *cobra.Command {
	flags := &listFlags{}

	cmd := &cobra.Command{
		Use:   "list [repository]",
		Short: shortDesc,
		Long:  longDesc,
		Run:   runCommand(flags),
	}

	cmd.Args = cobra.MaximumNArgs(1)
	cmd.Flags().BoolVar(&flags.useHTTP, "http", false, "Use plain HTTP when connecting to remote registries")
	return cmd
}

func runCommand(flags *listFlags) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		opts := &listOptions{}
		if err := opts.complete(cmd.Context(), flags, args); err != nil {
			output.Fatalf("Failed to parse argument: %s", err)
		}

		var allInfoLines []string
		if opts.remoteRef == nil {
			lines, err := listLocalKits(cmd.Context(), opts.storageHome)
			if err != nil {
				output.Fatalln(err)
			}
			allInfoLines = lines
		} else {
			lines, err := listRemoteKits(cmd.Context(), opts.remoteRef, opts.usehttp)
			if err != nil {
				output.Fatalln(err)
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

func printConfig(opts *listOptions) {
	output.Debugf("Using storage path: %s", opts.storageHome)
	if opts.remoteRef != nil {
		output.Debugf("Listing remote model kits in %s", opts.remoteRef.String())
	}
}
