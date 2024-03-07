package inspect

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"kitops/pkg/cmd/options"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/repo"
	"kitops/pkg/output"
	"strings"

	"github.com/spf13/cobra"
	"oras.land/oras-go/v2/errdef"
	"oras.land/oras-go/v2/registry"
)

const (
	shortDesc = `Inspect a ModelKit's manifest`
	longDesc  = `Print the contents of a ModelKit manifest to the screen.

If the ModelKit is present locally, this ModelKit is inspected by default. If
it is not present and the ModelKit reference includes a registry, kit will
download the manifest from the remote registry if possible.`
	example = `# Inspect a local ModelKit:
kit inspect mymodel:mytag

# Inspect a local ModelKit by digest:
kit inspect mymodel@sha256:44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a

# Inspect a remote ModelKit if not present locally:
kit inspect registry.example.com/my-model:1.0.0`
)

type inspectOptions struct {
	options.NetworkOptions
	configHome  string
	checkRemote bool
	modelRef    *registry.Reference
}

func InspectCommand() *cobra.Command {
	opts := &inspectOptions{}

	cmd := &cobra.Command{
		Use:     "inspect [flags] MODELKIT",
		Short:   shortDesc,
		Long:    longDesc,
		Example: example,
		Run:     runCommand(opts),
		Args:    cobra.ExactArgs(1),
	}

	opts.AddNetworkFlags(cmd)
	cmd.Flags().BoolVarP(&opts.checkRemote, "remote", "r", false, "Check remote registry even if ModelKit is present locally")
	return cmd
}

func runCommand(opts *inspectOptions) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := opts.complete(cmd.Context(), args); err != nil {
			output.Fatalf("Failed to parse arguments: %s", err)
		}
		manifest, err := inspectReference(cmd.Context(), opts)
		if err != nil {
			if errors.Is(err, errdef.ErrNotFound) {
				output.Fatalf("Could not find ModelKit %s", repo.FormatRepositoryForDisplay(opts.modelRef.String()))
			}
			output.Fatalf("Error resolving ModelKit: %s", err)
		}
		jsonBytes, err := json.MarshalIndent(manifest, "", "  ")
		if err != nil {
			output.Fatalf("Error formatting manifest: %w", err)
		}
		fmt.Println(string(jsonBytes))
	}
}

func (opts *inspectOptions) complete(ctx context.Context, args []string) error {
	configHome, ok := ctx.Value(constants.ConfigKey{}).(string)
	if !ok {
		return fmt.Errorf("default config path not set on command context")
	}
	opts.configHome = configHome

	ref, extraTags, err := repo.ParseReference(args[0])
	if err != nil {
		return err
	}
	if len(extraTags) > 0 {
		return fmt.Errorf("invalid reference format: extra tags are not supported: %s", strings.Join(extraTags, ", "))
	}
	opts.modelRef = ref

	if opts.modelRef.Registry == repo.DefaultRegistry && opts.checkRemote {
		return fmt.Errorf("can not check remote: %s does not contain registry", repo.FormatRepositoryForDisplay(opts.modelRef.String()))
	}

	return nil
}
