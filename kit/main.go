package main

import (
	"log"

	"github.com/comeonjy/go-kit/kit/cmd"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{}

func main() {
	root.AddCommand(cmd.NewProject)
	if err := root.Execute(); err != nil {
		log.Fatalln(err)
	}
}
