package cmd

import (
	"github.com/spf13/cobra"

	"github.com/safe-k/gonnect/internal/app/server"
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
			port, err := cmd.Flags().GetString("port")
			if err != nil {
				cmd.PrintErr(err)
				return
			}

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

				server.ServeMatch(port, server.BasicAuthenticator(user, pass))
			case "player":
				server.ServePlayer(port)
			}
		},
	}

	serveCmd.Flags().String("port", ":5000", "Port address to listen to")
	serveCmd.Flags().String("username", "admin", "Basic authentication username")
	serveCmd.Flags().String("password", "admin", "Basic authentication password")
	rootCmd.AddCommand(serveCmd)
}
