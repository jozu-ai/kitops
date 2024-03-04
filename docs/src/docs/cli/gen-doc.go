package main

import (
	"log"

	"github.com/spf13/cobra/doc"
	"kitops/cmd"
)

func main() {

	root := cmd.RunCommand()
	root.DisableAutoGenTag = true

	err := doc.GenMarkdownTree(root, "../cli")
	if err != nil {
		log.Fatal(err)
	}
}
