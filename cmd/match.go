package cmd

import (
	"github.com/spf13/cobra"

	"github.com/safe-k/gonnect/internal/app/matcher"
)

var matchCmd = &cobra.Command{
	Use:   "match",
	Short: "Runs a matchmaking worker",
	Run: func(cmd *cobra.Command, args []string) {
		bch, err := cmd.Flags().GetInt("batch")
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		matcher.Work(bch)
	},
}

func init() {
	rootCmd.AddCommand(matchCmd)

	matchCmd.Flags().Int("batch", 10, "The player count per match")
}
