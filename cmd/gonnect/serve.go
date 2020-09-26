package main

import (
	"github.com/spf13/cobra"

	"github.com/seifkamal/gonnect/internal"
	"github.com/seifkamal/gonnect/internal/server"
	"github.com/seifkamal/gonnect/internal/server/websocket"
)

func serveCommand() *cobra.Command {
	validArgs := []string{"match", "player"}
	serveCmd := &cobra.Command{
		Use:   "serve [ROUTER]",
		Short: "Starts an API api",
		Long: `Starts an API api with the specified router.
Currently the available routers are:
- match [for use of the game api to retrieve match information]
- player [for use by the player clients to search for a match]
`,
		Example:   "gonnect serve player",
		Version:   "1.0.0",
		Args:      cobra.ExactValidArgs(1),
		ValidArgs: validArgs,
		Run: func(cmd *cobra.Command, args []string) {
			port := cmd.Flag("port").Value.String()

			storage := internal.Storage()
			defer storage.Close()

			switch args[0] {
			case "match":
				user := cmd.Flag("username").Value.String()
				pass := cmd.Flag("password").Value.String()

				server.MatchmakingServer(server.BasicAuthenticator(user, pass), storage).Serve(port)
			case "player":
				server.PlayerServer(websocket.ConnectionUpgrader(), storage).Serve(port)
			}
		},
	}

	serveCmd.Flags().String("port", ":5000", "Port address to listen to")
	serveCmd.Flags().StringP("username", "u", "admin", "Basic authentication username")
	serveCmd.Flags().StringP("password", "p", "admin", "Basic authentication password")

	return serveCmd
}
