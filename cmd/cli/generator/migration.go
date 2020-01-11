package generator

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/nrfta/go-log"
	"github.com/nrfta/go-tiger/helpers"
	"github.com/spf13/cobra"
)

var migrationCmd = &cobra.Command{
	Use:   "migration [name]",
	Short: "Generates a migration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		t := time.Now()
		timestamp := t.Format("20060102150405")
		baseFileName := path.Join(helpers.FindRootPath(), "./db/migrations/"+timestamp+"_"+args[0])

		err := writeFile(baseFileName+".up.sql", "")
		if err != nil {
			log.Error(err)
			return
		}
		err = writeFile(baseFileName+".down.sql", "")
		if err != nil {
			log.Error(err)
			return
		}

		fmt.Println(baseFileName + ".up.sql")
		fmt.Println(baseFileName + ".down.sql\n")
		fmt.Println("Created!")
	},
}

func init() {
	GenerateCmd.AddCommand(migrationCmd)
}

func writeFile(path string, content string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = f.WriteString(content)
	if err != nil {
		f.Close()
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}
