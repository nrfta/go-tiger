package generator

import (
	"path"

	"github.com/nrfta/go-tiger/helpers"
	"github.com/nrfta/go-tiger/pkg/scaffolding"
	"github.com/spf13/cobra"
)

var scaffoldCmd = &cobra.Command{
	Use:   "scaffold [pathToModel.go]",
	Short: "Generates a scaffold file",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		file := path.Join(helpers.FindRootPath(), args[0])

		scaffolding.Process(file)
	},
}

func init() {
	GenerateCmd.AddCommand(scaffoldCmd)
}
