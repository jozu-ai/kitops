package list

import (
	"context"
	"fmt"
	"kitops/pkg/cmd/options"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/storage"
	"kitops/pkg/output"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"oras.land/oras-go/v2/registry"
)

const (
	shortDesc = `List model kits in a repository`
	longDesc  = `List model kits TODO`
)

type listOptions struct {
	options.NetworkOptions
	configHome string
	remoteRef  *registry.Reference
}

func (opts *listOptions) complete(ctx context.Context, args []string) error {
	configHome, ok := ctx.Value(constants.ConfigKey{}).(string)
	if !ok {
		return fmt.Errorf("default config path not set on command context")
	}
	opts.configHome = configHome
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

	printConfig(opts)
	return nil
}

// ListCommand represents the models command
func ListCommand() *cobra.Command {
	opts := &listOptions{}

	cmd := &cobra.Command{
		Use:   "list [repository]",
		Short: shortDesc,
		Long:  longDesc,
		Run:   runCommand(opts),
	}

	cmd.Args = cobra.MaximumNArgs(1)
	opts.AddNetworkFlags(cmd)

	return cmd
}

func runCommand(opts *listOptions) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := opts.complete(cmd.Context(), args); err != nil {
			output.Fatalf("Failed to parse argument: %s", err)
		}

		var allInfoLines []string
		if opts.remoteRef == nil {
			lines, err := listLocalKits(cmd.Context(), opts)
			if err != nil {
				output.Fatalln(err)
			}
			allInfoLines = lines
		} else {
			lines, err := listRemoteKits(cmd.Context(), opts)
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
	output.Debugf("Using config path: %s", opts.configHome)
	if opts.remoteRef != nil {
		output.Debugf("Listing remote model kits in %s", opts.remoteRef.String())
	}
}
