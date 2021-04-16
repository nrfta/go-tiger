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

const griftFileContent = `package grifts

import (
	"fmt"
	"os"
	"strings"

	dbPolicy "github.com/nrfta/%s/db/policy"
	. "github.com/markbates/grift/grift"
)

const (
	up   = "up"
	down = "down"
)

var _ = Namespace("policy", func() {
	Desc("migrate", "migrate accounts-api access management")
	Add("migrate", func(c *Context) error {
		if len(c.Args) != 1 || (strings.ToLower(c.Args[0]) != up && strings.ToLower(c.Args[0]) != down) {
			fmt.Println("Usage: policy:migrate [up|down]")
			os.Exit(1)
		}

		isUp := strings.ToLower(c.Args[0]) == up
		if isUp {
			return dbPolicy.MigrateUp()
		}
		return dbPolicy.MigrateDown()
	})
})
`

const mainPolicyFileContent = `package dbPolicy

import (
	"github.com/nrfta/%s/config"
  policyMigration "github.com/nrfta/go-platform-security-policy-migration/pkg/policy_migration"
)

var (
  migrations = make(policyMigration.Migrations)
)

// MigrateUp runs the security policy up migrations from the next version to the last version.
func MigrateUp() error {
	op := policyMigration.NewOperation(migrations, config.Config.Security)
	return op.Up()
}

// MigrationDown runs the security policy down migration for the current version.
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
		rootPath := helpers.FindRootPath()
		baseFileName := path.Join(rootPath, "./db/policy/"+timestamp+"_"+name)

		cfg := helpers.LoadConfig()
		serviceName := cfg.Meta.ServiceName

		if err := ensureMigrationsFileExists(path.Join(rootPath, "./db/policy/0_migrations.go"), serviceName); err != nil {
			log.Error(err)
			return
		}
		if err := ensureGriftFileExists(path.Join(rootPath, "./grifts/policy.go"), serviceName); err != nil {
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

func ensureMigrationsFileExists(fileName, serviceName string) error {
	dir := path.Dir(fileName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return writeFile(fileName, fmt.Sprintf(mainPolicyFileContent, serviceName))
	}
	return nil
}

func ensureGriftFileExists(fileName, serviceName string) error {
	dir := path.Dir(fileName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return writeFile(fileName, fmt.Sprintf(griftFileContent, serviceName))
	}
	return nil
}
