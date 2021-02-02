package am

import (
	am "github.com/nrfta/go-access-management"
	"github.com/nrfta/go-log"
	"github.com/nrfta/go-tiger/helpers"
	"github.com/spf13/cobra"
)

// AMCmd represents the am command when called without any subcommands
var AMCmd = &cobra.Command{
	Use:   "am",
	Short: "Manage Access Management",
}

func initAM() {
	amConfig := helpers.LoadConfig().AccessManagement
	if !am.IsInitialized() {
		am.Initialize(amConfig, log.GetLogger())
	}
}

func closeAM() {
	am.Close()
}
