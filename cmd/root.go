/*
Copyright © 2024 Jozu.com
*/
package cmd

import (
	"kitops/pkg/cmd/build"
	"kitops/pkg/cmd/export"
	"kitops/pkg/cmd/list"
	"kitops/pkg/cmd/login"
	"kitops/pkg/cmd/pull"
	"kitops/pkg/cmd/push"
	"kitops/pkg/cmd/version"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type (
	RootOptions struct {
		ConfigHome string
	}

	RootFlags struct {
		ConfigHome string
	}
)

var (
	shortDesc = `KitOps model manager`
	longDesc  = `KitOps is a tool to manage AI and ML models`
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = newRootCmd()

func init() {
	rootCmd.AddCommand(build.NewCmdBuild())
	rootCmd.AddCommand(login.NewCmdLogin())
	rootCmd.AddCommand(pull.PullCommand())
	rootCmd.AddCommand(push.PushCommand())
	rootCmd.AddCommand(list.ListCommand())
	rootCmd.AddCommand(export.ExportCommand())
	rootCmd.AddCommand(version.NewCmdVersion())
}

func newRootCmd() *cobra.Command {
	flags := &RootFlags{}
	cmd := &cobra.Command{
		Use:   "kit",
		Short: shortDesc,
		Long:  longDesc,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			options, err := flags.ToOptions()
			if err != nil {
				panic(err)
			}
			err = options.Complete()
			if err != nil {
				panic(err)
			}
		},
	}
	flags.addFlags(cmd)
	return cmd
}

func (f *RootFlags) addFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&f.ConfigHome, "config", "", "config file (default is $HOME/.kitops)")
	viper.BindPFlag("config", cmd.PersistentFlags().Lookup("config"))
}

func (f *RootFlags) ToOptions() (*RootOptions, error) {
	return &RootOptions{
		ConfigHome: f.ConfigHome,
	}, nil
}

func (o *RootOptions) Complete() error {
	if o.ConfigHome == "" {
		currentUser, err := user.Current()
		if err != nil {
			return err
		}
		configpath := filepath.Join(currentUser.HomeDir, ".kitops")
		viper.Set("config", configpath)
		o.ConfigHome = configpath
	}
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
