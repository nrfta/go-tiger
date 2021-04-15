package generator

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/nrfta/go-log"
	"github.com/nrfta/go-tiger/helpers"
	"github.com/spf13/cobra"
)

const mainPolicyFileContent = `package dbPolicy

import (
  policyMigration "github.com/nrfta/go-platform-security-policy-migration/pkg/policy_migration"
)

var (
  migrations = make(policyMigration.Migrations)
)

func MigrateUp() error {
	op := policyMigration.NewOperation(migrations, config.Config.Security)
	return op.Up()
}

func MigrateDown() error {
	op := policyMigration.NewOperation(migrations, config.Config.Security)
	return op.Down()
}
`

var policyMigrationCmd = &cobra.Command{
	Use:   "policyMigration [platform security policy version] [policyID,policyID,...] [description]",
	Short: "Generates a security policy migration file",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		t := time.Now()
		timestamp := t.Format("20060102150405")
		policyIDs := strings.Split(args[1], ",")
		policyIDsStr := `"` + strings.Join(policyIDs, `", "`) + `"`
		description := strings.ReplaceAll(args[2], " ", "_")
		name := args[1] + "_" + description
		baseFileName := path.Join(helpers.FindRootPath(), "./db/policy/"+timestamp+"_"+name)

		if err := ensureMigrationsFileExists(path.Join(helpers.FindRootPath(), "./db/policy/(migrations).go")); err != nil {
			log.Error(err)
			return
		}

		content := fmt.Sprintf(`package dbPolicy

import (
	"errors"

  policyMigration "github.com/nrfta/go-platform-security-policy-migration/pkg/policy_migration"
)

func init() {
	migrations[%s] = &policyMigration.Migration{
		Version:       %s,
		PolicyVersion: "%s",
		PolicyIDs:     []string{%s},
		Description:   "%s",
		Up:            migration%sUp,
		Down:          migration%sDown,
	}
}

func migration%sUp(policyVersion string) error {
	// migrate up implementation
	return errors.New("not implemented")
}

func migration%sDown(policyVersion string) error {
	// migrate up implementation
	return errors.New("not implemented")
}
`, timestamp, timestamp, args[0], policyIDsStr, args[2], timestamp, timestamp, timestamp, timestamp)

		err := writeFile(baseFileName+".go", content)
		if err != nil {
			log.Error(err)
			return
		}

		fmt.Println(baseFileName + ".go created")
	},
}

func init() {
	GenerateCmd.AddCommand(policyMigrationCmd)
}

func ensureMigrationsFileExists(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return writeFile(path, mainPolicyFileContent)
	}
	return nil
}
