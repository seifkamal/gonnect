package cmd

import (
	"github.com/spf13/cobra"

	"github.com/safe-k/gonnect/internal/app/server"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts an HTTP API server",
	Run: func(cmd *cobra.Command, args []string) {
		server.Serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
