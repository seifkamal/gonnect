package cmd

import (
	"github.com/spf13/cobra"

	"github.com/safe-k/gonnect/internal/app"
	"github.com/safe-k/gonnect/internal/app/server"
	"github.com/safe-k/gonnect/internal/app/server/match"
	"github.com/safe-k/gonnect/internal/app/server/player"
)

func init() {
	serveCmd := &cobra.Command{
		Use:       "serve [ROUTER]",
		Short:     "Starts an HTTP API server",
		Args:      cobra.ExactValidArgs(1),
		ValidArgs: []string{"match", "player"},
		Run: func(cmd *cobra.Command, args []string) {
			DB := app.DB()
			defer DB.Close()

			var h server.Handler
			switch args[0] {
			case "match":
				h = &match.Handler{DB: DB}
			case "player":
				h = &player.Handler{DB: DB}
			}

			server.Serve(h)
		},
	}

	rootCmd.AddCommand(serveCmd)
}
