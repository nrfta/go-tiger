package cli

// Source: https://github.com/gobuffalo/buffalo/blob/master/buffalo/cmd/task.go

import (
	grifts "github.com/markbates/grift/cmd"
	"github.com/spf13/cobra"
)

// task command is a forward to grift tasks
var taskCommand = &cobra.Command{
	Use:                "task",
	Aliases:            []string{"t", "tasks"},
	Short:              "Run grift tasks",
	DisableFlagParsing: true,
	RunE: func(c *cobra.Command, args []string) error {
		return grifts.Run("tiger task", args)
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	RootCmd.AddCommand(taskCommand)
}
