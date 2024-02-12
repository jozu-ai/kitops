package models

import (
	"fmt"
	"io/fs"
	"jmm/pkg/lib/storage"
	"os"
	"path"
	"path/filepath"
	"text/tabwriter"

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
	configHome  string
	storageHome string
}

func (opts *ModelsOptions) complete() {
	opts.configHome = viper.GetString("config")
	opts.storageHome = path.Join(opts.configHome, "storage")
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

		storeDirs, err := findRepos(opts.storageHome)
		if err != nil {
			fmt.Println(err)
		}

		var allInfoLines []string
		for _, storeDir := range storeDirs {
			store := storage.NewLocalStore(opts.storageHome, storeDir)

			infolines, err := listModels(store)
			if err != nil {
				fmt.Println(err)
				return
			}
			allInfoLines = append(allInfoLines, infolines...)
		}

		printSummary(allInfoLines)

	}
}

func findRepos(storePath string) ([]string, error) {
	var indexPaths []string
	err := filepath.WalkDir(storePath, func(file string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "index.json" && !info.IsDir() {
			dir := filepath.Dir(file)
			relDir, err := filepath.Rel(storePath, dir)
			if err != nil {
				return err
			}
			if relDir == "." {
				relDir = ""
			}
			indexPaths = append(indexPaths, relDir)
		}
		return nil
	})
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read local storage: %w", err)
	}
	return indexPaths, nil
}

func printSummary(lines []string) {
	tw := tabwriter.NewWriter(os.Stdout, 0, 2, 3, ' ', 0)
	fmt.Fprintln(tw, ModelsTableHeader)
	for _, line := range lines {
		fmt.Fprintln(tw, line)
	}
	tw.Flush()
}
