package cmd

import (
	"github.com/spf13/cobra"

	"github.com/safe-k/gonnect/internal/app/server"
)

func init() {
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

			switch args[0] {
			case "match":
				user := cmd.Flag("username").Value.String()
				pass := cmd.Flag("password").Value.String()

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
