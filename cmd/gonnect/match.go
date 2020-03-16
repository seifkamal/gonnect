package main

import (
	"github.com/spf13/cobra"

	"github.com/safe-k/gonnect/internal"
	"github.com/safe-k/gonnect/internal/matchmaking"
)

func matchCommand() *cobra.Command {
	matchCmd := &cobra.Command{
		Use:     "match",
		Short:   "Runs a matchmaking worker",
		Long:    "Runs a matchmaking worker that will create a new match whenever enough players are searching",
		Example: "gonnect match --batch 5",
		Version: "1.0.0",
		Run: func(cmd *cobra.Command, args []string) {
			batch, err := cmd.Flags().GetInt("batch")
			if err != nil {
				cmd.PrintErr(err)
				return
			}

			retryInterval, err := cmd.Flags().GetInt("retry-interval")
			if err != nil {
				cmd.PrintErr(err)
				return
			}

			storage := internal.Storage()
			defer storage.Close()

			matchmaking.Worker(storage).WorkIndefinitely(batch, retryInterval)
		},
	}

	matchCmd.Flags().IntP("batch", "b", 10, "Number of players per match")
	matchCmd.Flags().IntP("retry-interval", "r", 2, "Amount of time (in seconds) to wait before retrying when not enough players are available")

	return matchCmd
}
