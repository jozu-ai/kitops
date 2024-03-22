/*
Copyright © 2024 Jozu.com
*/
package logout

import (
	"context"
	"fmt"
	"kitops/pkg/lib/constants"
	"kitops/pkg/output"

	"github.com/spf13/cobra"
)

const (
	shortDesc = `Log out from an OCI registry`
	longDesc  = `Log out from a specified OCI-compatible registry. Any saved credentials are
removed from storage.`

	example = `# Log out from ghcr.io
kit logout ghcr.io`
)

type logoutOptions struct {
	credentialStoreHome string
	registry            string
}

func LogoutCommand() *cobra.Command {
	opts := &logoutOptions{}

	cmd := &cobra.Command{
		Use:     "logout [flags] REGISTRY",
		Short:   shortDesc,
		Long:    longDesc,
		Example: example,
		Run:     runLogout(opts),
	}
	cmd.Args = cobra.ExactArgs(1)
	return cmd
}

func runLogout(opts *logoutOptions) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := opts.complete(cmd.Context(), args); err != nil {
			output.Fatalf("Invalid arguments: %s", err)
		}
		err := logout(cmd.Context(), opts.registry, opts.credentialStoreHome)
		if err != nil {
			output.Fatalln(err)
		}
	}
}

func (opts *logoutOptions) complete(ctx context.Context, args []string) error {
	opts.registry = args[0]
	configHome, ok := ctx.Value(constants.ConfigKey{}).(string)
	if !ok {
		return fmt.Errorf("default config path not set on command context")
	}
	opts.credentialStoreHome = constants.CredentialsPath(configHome)
	return nil
}
