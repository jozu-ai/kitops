/*
Copyright © 2024 Jozu.com
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"kitops/pkg/cmd/inspect"
	"kitops/pkg/cmd/list"
	"kitops/pkg/cmd/login"
	"kitops/pkg/cmd/logout"
	"kitops/pkg/cmd/pack"
	"kitops/pkg/cmd/pull"
	"kitops/pkg/cmd/push"
	"kitops/pkg/cmd/remove"
	"kitops/pkg/cmd/tag"
	"kitops/pkg/cmd/unpack"
	"kitops/pkg/cmd/version"
	"kitops/pkg/lib/constants"
	"kitops/pkg/output"

	"github.com/spf13/cobra"
)

var (
	shortDesc = `Streamline the lifecycle of AI/ML models`
	longDesc  = `Kit is a tool for efficient AI/ML model lifecycle management.

Find more information at: http://kitops.ml`
)

type rootOptions struct {
	configHome string
	verbose    bool
}

func RunCommand() *cobra.Command {
	opts := &rootOptions{}

	cmd := &cobra.Command{
		Use:   `kit`,
		Short: shortDesc,
		Long:  longDesc,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			configHome := opts.configHome
			if configHome == "" {
				currentUser, err := user.Current()
				if err != nil {
					output.Fatalf("Failed to resolve default storage path '$HOME/%s: could not get current user", constants.DefaultConfigSubdir)
				}
				configHome = filepath.Join(currentUser.HomeDir, constants.DefaultConfigSubdir)
			}
			if opts.verbose {
				output.SetDebug(true)
			}
			ctx := context.WithValue(cmd.Context(), constants.ConfigKey{}, configHome)
			cmd.SetContext(ctx)
		},
	}
	addSubcommands(cmd)
	cmd.PersistentFlags().StringVar(&opts.configHome, "config", "", fmt.Sprintf("Config file (default $HOME/%s)", constants.DefaultConfigSubdir))
	cmd.PersistentFlags().BoolVarP(&opts.verbose, "verbose", "v", false, "Include additional information in output (default false)")

	cmd.SetHelpTemplate(helpTemplate)
	cmd.SetUsageTemplate(usageTemplate)
	cobra.AddTemplateFunc("indent", indentBlock)
	cobra.AddTemplateFunc("sectionHead", sectionHead)
	cobra.AddTemplateFunc("ensureTrailingNewline", ensureTrailingNewline)

	return cmd
}

func addSubcommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(pack.PackCommand())
	rootCmd.AddCommand(unpack.UnpackCommand())
	rootCmd.AddCommand(push.PushCommand())
	rootCmd.AddCommand(pull.PullCommand())
	rootCmd.AddCommand(tag.TagCommand())
	rootCmd.AddCommand(list.ListCommand())
	rootCmd.AddCommand(inspect.InspectCommand())
	rootCmd.AddCommand(remove.RemoveCommand())
	rootCmd.AddCommand(login.LoginCommand())
	rootCmd.AddCommand(logout.LogoutCommand())
	rootCmd.AddCommand(version.VersionCommand())
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RunCommand().Execute()
	if err != nil {
		os.Exit(1)
	}
}
