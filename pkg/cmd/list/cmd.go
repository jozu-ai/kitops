package list

import (
	"context"
	"fmt"
	"kitops/pkg/cmd/options"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/repo"
	"kitops/pkg/output"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"oras.land/oras-go/v2/registry"
)

const (
	shortDesc = `List modelkits in a repository`
	longDesc  = `Displays a list of modelkits available in a repository.

This command provides an overview of modelkits stored either in the local
repository or a specified remote repository. It displays each modelkit along
with its associated tags and the cumulative size of its contents. Modelkits
comprise multiple artifacts, including models, datasets, code, and
configuration, designed to enhance reusability and modularity. However, this
command focuses on the aggregate rather than listing individual artifacts.

Each modelkit entry includes its DIGEST, a unique identifier that ensures
distinct versions of a modelkit are easily recognizable, even if they share
the same name or tags. Modelkits with multiple tags or repository names will
appear multiple times in the list, distinguished by their DIGEST.

The SIZE displayed for each modelkit represents the total storage space
occupied by all its components.`

	example = `# List local modelkits
kit list

# List modelkits from a remote repository
kit list registry.example.com/my-namespace/my-model`
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
		remoteRef, extraTags, err := repo.ParseReference(args[0])
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
		Use:     "list [flags] [REPOSITORY]",
		Short:   shortDesc,
		Long:    longDesc,
		Example: example,
		Run:     runCommand(opts),
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
	if opts.remoteRef != nil {
		output.Debugf("Listing remote model kits in %s", opts.remoteRef.String())
	}
}
