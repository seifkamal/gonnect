package cmd

import (
	"github.com/spf13/cobra"

	"github.com/safe-k/gonnect/internal"
	"github.com/safe-k/gonnect/internal/app/server"
	"github.com/safe-k/gonnect/internal/app/server/match"
	"github.com/safe-k/gonnect/internal/app/server/player"
)

func init() {
	validArgs := []string{"match", "player"}
	serveCmd := &cobra.Command{
		Use:   "serve [ROUTER]",
		Short: "Starts an API server",
		Long: `Starts an API server with the specified router.
Currently the available routers are:
- match [for use of the game server to retrieve match information]
- player [for use by the player clients to search for a match]
`,
		Example:   "gonnect serve player",
		Version:   "1.0.0",
		Args:      cobra.ExactValidArgs(1),
		ValidArgs: validArgs,
		Run: func(cmd *cobra.Command, args []string) {
			storage := internal.Storage()
			defer storage.Close()

			var h server.Handler
			switch args[0] {
			case "match":
				user, err := cmd.Flags().GetString("user")
				if err != nil {
					cmd.PrintErr(err)
					return
				}

				pass, err := cmd.Flags().GetString("pass")
				if err != nil {
					cmd.PrintErr(err)
					return
				}

				auth := server.NewBasicAuthenticator(user, pass)
				h = &match.Handler{
					BasicAuthenticator: auth,
					Storage:            storage,
				}
			case "player":
				h = &player.Handler{Storage: storage}
			}

			server.Serve(h)
		},
	}

	serveCmd.Flags().String("username", "admin", "Basic authentication username")
	serveCmd.Flags().String("password", "admin", "Basic authentication password")
	rootCmd.AddCommand(serveCmd)
}
