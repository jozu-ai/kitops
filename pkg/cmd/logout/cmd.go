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
	shortDesc = `Log out from a registry`
	longDesc  = `Log out from a registry TODO`
)

type logoutOptions struct {
	credentialStoreHome string
	registry            string
}

func LogoutCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout [flags] registry",
		Short: shortDesc,
		Long:  longDesc,
		Run:   runLogout(),
	}
	cmd.Args = cobra.ExactArgs(1)
	return cmd
}

func runLogout() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		opts := &logoutOptions{}
		if err := opts.complete(cmd.Context(), args); err != nil {
			output.Fatalf("Failed to process flags: %s", err)
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
