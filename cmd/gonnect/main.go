package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	mainCmd := &cobra.Command{}
	mainCmd.AddCommand(matchCommand())
	mainCmd.AddCommand(serveCommand())

	if err := mainCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
