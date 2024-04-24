/*
Copyright © 2024 Jozu.com
*/
package cmd

import (
	"context"
	"os"

	"kitops/pkg/cmd/info"
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
	configHome   string
	verbose      bool
	progressBars string
}

func RunCommand() *cobra.Command {
	opts := &rootOptions{}

	cmd := &cobra.Command{
		Use:   `kit`,
		Short: shortDesc,
		Long:  longDesc,
		// Commands do all their printing directly and return an error only to signal that we should exit with
		// nonzero status. We don't want to print usage or the error message in this case.
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			output.SetOut(cmd.OutOrStdout())
			output.SetErr(cmd.ErrOrStderr())
			output.SetDebug(opts.verbose)
			output.SetProgressBars(opts.progressBars)

			configHome, err := getConfigHome(opts)
			if err != nil {
				output.Errorf("Failed to read base config directory")
				output.Infof("Use the --config flag or set the $KITOPS_HOME environment variable to provide a default")
				output.Debugf("Error: %s", err)
			}
			ctx := context.WithValue(cmd.Context(), constants.ConfigKey{}, configHome)
			cmd.SetContext(ctx)
		},
	}
	addSubcommands(cmd)
	cmd.PersistentFlags().StringVar(&opts.configHome, "config", "", "Alternate path to root storage directory for CLI")
	cmd.PersistentFlags().BoolVarP(&opts.verbose, "verbose", "v", false, "Include additional information in output (default false)")
	cmd.PersistentFlags().StringVar(&opts.progressBars, "progress", "plain", "Configure progress bars for longer operations (options: none, plain, fancy)")

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
	rootCmd.AddCommand(info.InfoCommand())
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

func getConfigHome(opts *rootOptions) (string, error) {
	if opts.configHome != "" {
		output.Debugf("Using config directory from flag: %s", opts.configHome)
		return opts.configHome, nil
	}

	envHome := os.Getenv("KITOPS_HOME")
	if envHome != "" {
		output.Debugf("Using config directory from environment variable: %s", envHome)
		return envHome, nil
	}

	defaultHome, err := constants.DefaultConfigPath()
	if err != nil {
		return "", err
	}
	output.Debugf("Using default config directory: %s", defaultHome)
	return defaultHome, nil
}
