package remove

import (
	"context"
	"fmt"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/repo"
	"kitops/pkg/output"
	"strings"

	"github.com/spf13/cobra"
	"oras.land/oras-go/v2/registry"
)

const (
	shortDesc = `Remove a modelkit from local storage`
	longDesc  = `Removes a modelkit from storage on the local disk.

The model to be removed may be specifed either by a tag or by a digest. If
specified by digest, that modelkit will be removed along with any tags that
might refer to it. If specified by tag (and the --force flag is not used),
the modelkit will only be removed if no other tags refer to it; otherwise
it is only untagged.`

	examples = `kit remove my-registry.com/my-org/my-repo:my-tag
kit remove my-registry.com/my-org/my-repo@sha256:<digest>
kit remove my-registry.com/my-org/my-repo:tag1,tag2,tag3`
)

type removeOptions struct {
	configHome  string
	forceDelete bool
	modelRef    *registry.Reference
	extraTags   []string
}

func (opts *removeOptions) complete(ctx context.Context, args []string) error {
	configHome, ok := ctx.Value(constants.ConfigKey{}).(string)
	if !ok {
		return fmt.Errorf("default config path not set on command context")
	}
	opts.configHome = configHome

	modelRef, extraTags, err := repo.ParseReference(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse reference %s: %w", modelRef, err)
	}
	opts.modelRef = modelRef
	opts.extraTags = extraTags

	printConfig(opts)
	return nil
}

func RemoveCommand() *cobra.Command {
	opts := &removeOptions{}
	cmd := &cobra.Command{
		Use:     "remove [flags] registry/repository[:tag|@digest]",
		Short:   shortDesc,
		Long:    longDesc,
		Example: examples,
		Run:     runCommand(opts),
	}
	cmd.Args = cobra.ExactArgs(1)
	cmd.Flags().BoolVarP(&opts.forceDelete, "force", "f", false, "remove manifest even if other tags refer to it")
	return cmd
}

func runCommand(opts *removeOptions) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := opts.complete(cmd.Context(), args); err != nil {
			output.Fatalf("Failed to process arguments: %s", err)
		}
		storageRoot := constants.StoragePath(opts.configHome)
		localStore, err := repo.NewLocalStore(storageRoot, opts.modelRef)
		if err != nil {
			output.Fatalf("Failed to read local storage: %s", storageRoot)
		}
		desc, err := removeModel(cmd.Context(), localStore, opts.modelRef, opts.forceDelete)
		if err != nil {
			output.Fatalf("Failed to remove: %s", err)
		}
		output.Infof("Removed %s (digest %s)", opts.modelRef.String(), desc.Digest)

		for _, tag := range opts.extraTags {
			ref := *opts.modelRef
			ref.Reference = tag
			desc, err := removeModel(cmd.Context(), localStore, &ref, opts.forceDelete)
			if err != nil {
				output.Errorf("Failed to remove: %s", err)
			} else {
				output.Infof("Removed %s (digest %s)", ref.String(), desc.Digest)
			}
		}
	}
}

func printConfig(opts *removeOptions) {
	output.Debugf("Removing %s and additional tags: [%s]", opts.modelRef.String(), strings.Join(opts.extraTags, ", "))
}
