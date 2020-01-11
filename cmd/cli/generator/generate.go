package generator

import (
	"github.com/spf13/cobra"
)

// GenerateCmd represents the generate command when called without any subcommands
var GenerateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generate files",
	Aliases: []string{"g"},
}
