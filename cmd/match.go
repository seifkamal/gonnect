package cmd

import (
	"github.com/spf13/cobra"

	"github.com/safe-k/gonnect/internal/app"
	"github.com/safe-k/gonnect/internal/app/matchmaker"
)

func init() {
	matchCmd := &cobra.Command{
		Use:   "match",
		Short: "Runs a matchmaking worker",
		Run: func(cmd *cobra.Command, args []string) {
			bch, err := cmd.Flags().GetInt("batch")
			if err != nil {
				cmd.PrintErr(err)
				return
			}

			DB := app.DB()
			defer DB.Close()

			w := &matchmaker.Worker{DB: DB}
			w.Match(bch)
		},
	}

	matchCmd.Flags().Int("batch", 10, "The player count per match")
	rootCmd.AddCommand(matchCmd)
}
