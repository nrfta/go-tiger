package db

import (
	_ "github.com/lib/pq" // Postgres
	grifts "github.com/markbates/grift/cmd"
	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Runs grift task db:seed",
	Long:  "It runs the grift task called db:seed",
	Run: func(cmd *cobra.Command, args []string) {
		grifts.Run("tiger task", []string{"db:seed"})
	},
}

func init() {
	DBCmd.AddCommand(seedCmd)
}
