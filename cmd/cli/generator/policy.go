package generator

import (
	"github.com/nrfta/go-tiger/helpers"
	"github.com/nrfta/go-tiger/pkg/generators/policy"
	"github.com/spf13/cobra"
)

var policyCmd = &cobra.Command{
	Use:   "policy [SingularResourceName]",
	Short: "Generates a policy file",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		rootPath := helpers.FindRootPath()

		policy.Process(rootPath, args[0])
	},
}

func init() {
	GenerateCmd.AddCommand(policyCmd)
}
