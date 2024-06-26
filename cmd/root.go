/*
Copyright © 2024 Jozu.com
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"kitops/pkg/cmd/dev"
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
	loglevel     string
	progressBars string
}

func RunCommand() *cobra.Command {
	opts := &rootOptions{}

	cmd := &cobra.Command{
		Use:   `kit`,
		Short: shortDesc,
		Long:  longDesc,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			output.SetOut(cmd.OutOrStdout())
			output.SetErr(cmd.ErrOrStderr())
			if err := output.SetLogLevelFromString(opts.loglevel); err != nil {
				output.Fatalln(err)
			}
			if opts.verbose {
				output.SetLogLevel(output.LogLevelDebug)
			}

			output.SetProgressBars(opts.progressBars)

			configHome, err := getConfigHome(opts)
			if err != nil {
				output.Errorf("Failed to read base config directory")
				output.Infof("Use the --config flag or set the $KITOPS_HOME environment variable to provide a default")
				output.Debugf("Error: %s", err)
			}
			ctx := context.WithValue(cmd.Context(), constants.ConfigKey{}, configHome)
			cmd.SetContext(ctx)
			// At this point, we've parsed the command tree and args; the CLI is being correctly
			// so we don't want to print usage. Each subcommand should print its error message before
			// returning
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true
		},
	}
	addSubcommands(cmd)
	cmd.PersistentFlags().StringVar(&opts.configHome, "config", "", "Alternate path to root storage directory for CLI")
	cmd.PersistentFlags().BoolVarP(&opts.verbose, "verbose", "v", false, "Include additional information in output. Alias for --log-level=debug")
	cmd.PersistentFlags().StringVar(&opts.loglevel, "log-level", "info", "Log messages above specified level ('trace', 'debug', 'info', 'warn', 'error') (default 'info')")
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
	rootCmd.AddCommand(dev.DevCommand())
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
		absHome, err := filepath.Abs(opts.configHome)
		if err != nil {
			return "", fmt.Errorf("failed to get absolute path for %s: %w", opts.configHome, err)
		}
		return absHome, nil
	}

	envHome := os.Getenv("KITOPS_HOME")
	if envHome != "" {
		output.Debugf("Using config directory from environment variable: %s", envHome)
		absHome, err := filepath.Abs(envHome)
		if err != nil {
			return "", fmt.Errorf("failed to get absolute path for $KITOPS_HOME: %w", err)
		}
		return absHome, nil
	}

	defaultHome, err := constants.DefaultConfigPath()
	if err != nil {
		return "", err
	}
	output.Debugf("Using default config directory: %s", defaultHome)
	return defaultHome, nil
}
