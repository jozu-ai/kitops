package tag

import (
	"context"
	"fmt"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/repo"
	"kitops/pkg/output"

	"github.com/spf13/cobra"
	"oras.land/oras-go/v2/registry"
)

var (
	shortDesc = "Tag a modelkit"
	longDesc  = `Tag a modelkit with a new tag.`
)

type tagOptions struct {
	configHome string
	sourceRef  *registry.Reference
	targetRef  *registry.Reference
}

func (opts *tagOptions) complete(ctx context.Context, args []string) error {

	configHome, ok := ctx.Value(constants.ConfigKey{}).(string)
	if !ok {
		return fmt.Errorf("default config path not set on command context")
	}
	opts.configHome = configHome
	modelRef, _, err := repo.ParseReference(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse reference %s: %w", opts.sourceRef, err)
	}
	opts.sourceRef = modelRef
	modelRef, _, err = repo.ParseReference(args[1])
	if err != nil {
		return fmt.Errorf("failed to parse reference %s: %w", opts.targetRef, err)
	}
	opts.targetRef = modelRef
	return nil
}

func TagCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "tag SOURCE_MODELKIT[:TAG] TARGET_MODELKIT[:TAG]",
		Short:   shortDesc,
		Long:    longDesc,
		Example: `kit tag myregistry.com/myrepo/mykit:latest myregistry.com/myrepo/mykit:v1.0.0`,
		Run:     runCommand(&tagOptions{}),
	}

	cmd.Args = cobra.ExactArgs(2)
	return cmd
}

func runCommand(opts *tagOptions) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := opts.complete(cmd.Context(), args); err != nil {
			output.Fatalf("Failed to parse argument: %s", err)
		}

		err := RunTag(cmd.Context(), opts)
		if err != nil {
			output.Fatalf("Failed to tag modelkit: %s", err)
		}
		output.Infof("Modelkit %s tagged as %s", opts.sourceRef, opts.targetRef)
	}
}
