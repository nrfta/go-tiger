package generator

import (
	"path"

	"github.com/nrfta/go-tiger/helpers"
	"github.com/nrfta/go-tiger/pkg/generators/scaffolding"
	"github.com/spf13/cobra"
)

var scaffoldCmd = &cobra.Command{
	Use:   "scaffold [pathToModel.go]",
	Short: "Generates a scaffold file",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		rootPath := helpers.FindRootPath()

		scaffolding.Process(rootPath, path.Join(rootPath, args[0]))
	},
}

func init() {
	GenerateCmd.AddCommand(scaffoldCmd)
}
