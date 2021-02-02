package cli

import (
	"github.com/spf13/cobra"

	"github.com/nrfta/go-log"
	"github.com/nrfta/go-tiger/cmd/cli/am"
	"github.com/nrfta/go-tiger/cmd/cli/db"
	"github.com/nrfta/go-tiger/cmd/cli/generator"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "tiger",
	Short: "tiger",
	Long:  `Command line utility tiger`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	RootCmd.AddCommand(db.DBCmd)
	RootCmd.AddCommand(generator.GenerateCmd)
	RootCmd.AddCommand(am.AMCmd)
}
