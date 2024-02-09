package models

import (
	"fmt"
	"jmm/pkg/lib/storage"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	shortDesc = `List models`
	longDesc  = `List models TODO`
)

var (
	opts *ModelsOptions
)

type ModelsOptions struct {
	configHome string
}

func (opts *ModelsOptions) complete() {
	opts.configHome = viper.GetString("config")
}

func (opts *ModelsOptions) validate() error {
	return nil
}

// ModelsCommand represents the models command
func ModelsCommand() *cobra.Command {
	opts = &ModelsOptions{}

	cmd := &cobra.Command{
		Use:   "models",
		Short: shortDesc,
		Long:  longDesc,
		Run:   RunCommand(opts),
	}

	return cmd
}

func RunCommand(options *ModelsOptions) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		options.complete()
		err := options.validate()
		if err != nil {
			fmt.Println(err)
			return
		}
		store := storage.NewLocalStore(opts.configHome)

		summary, err := listModels(store)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(summary)
	}
}
