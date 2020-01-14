package db

import (
	"os"
	"os/exec"
	"strings"

	_ "github.com/lib/pq" // Postgres
	"github.com/nrfta/go-log"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Runs drop, create and migrate up",
	Long:  "It runs drop, create and migrate up",
	Run: func(cmd *cobra.Command, args []string) {

		err := run(exec.Command("tiger", "db", "drop"))
		if err != nil {
			return
		}

		err = run(exec.Command("tiger", "db", "create"))
		if err != nil {
			return
		}

		err = run(exec.Command("tiger", "db", "migrate", "up"))
		if err != nil {
			return
		}
	},
}

func init() {
	DBCmd.AddCommand(resetCmd)
}

func run(cmd *exec.Cmd) error {
	log.Infof("--> %s", strings.Join(cmd.Args, " "))
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
