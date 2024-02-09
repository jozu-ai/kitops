/*
Copyright © 2024 Jozu.com
*/
package cmd

import (
	"jmm/pkg/cmd/build"
	"jmm/pkg/cmd/login"
	"jmm/pkg/cmd/models"
	"jmm/pkg/cmd/pull"
	"jmm/pkg/cmd/push"
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
	shortDesc = `Jozu Model Manager`
	longDesc  = `Jozu Model Manager is a tool to manage AI and ML models`
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = newRootCmd()

func init() {
	rootCmd.AddCommand(build.NewCmdBuild())
	rootCmd.AddCommand(login.NewCmdLogin())
	rootCmd.AddCommand(pull.NewCmdPull())
	rootCmd.AddCommand(push.NewCmdPush())
	rootCmd.AddCommand(models.ModelsCommand())
}

func newRootCmd() *cobra.Command {
	flags := &RootFlags{}
	cmd := &cobra.Command{
		Use:   "jmm",
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
	cmd.PersistentFlags().StringVar(&f.ConfigHome, "config", "", "config file (default is $HOME/.jozu)")
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
		configpath := filepath.Join(currentUser.HomeDir, ".jozu")
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
