/*
Copyright © 2024 Jozu.com
*/
package models

import (
	"fmt"

	"github.com/spf13/cobra"
)

type ModelsFlags struct {
}
type ModelsOptions struct {
}

// modelsCmd represents the models command
func NewCmdModels() *cobra.Command {
	modelsFlags := NewModelsFlags()

	cmd := &cobra.Command{
		Use:   "models",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			options, err := modelsFlags.ToOptions()
			if err != nil {
				fmt.Println(err)
				return
			}
			err = options.Validate()
			if err != nil {
				fmt.Println(err)
				return
			}
			options.RunModels()
			if err != nil {
				fmt.Println(err)
				return
			}
		},
	}
	modelsFlags.AddFlags(cmd)
	return cmd
}

func NewModelsFlags() *ModelsFlags {
	return &ModelsFlags{}
}

func (f *ModelsFlags) AddFlags(cmd *cobra.Command) {

}

func (f *ModelsFlags) ToOptions() (*ModelsOptions, error) {
	return &ModelsOptions{}, nil
}

func (o *ModelsOptions) Validate() error {
	return nil
}

func (o *ModelsOptions) RunModels() error {
	return nil
}
