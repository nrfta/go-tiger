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
	dbPkg "github.com/nrfta/%s/db"
	policyMigration "github.com/nrfta/go-platform-security-policy-migration/pkg/policy_migration"
)

var (
	DB dbPkg.DB
  Migrations = make(policyMigration.Migrations)
)

// MigrateUp runs the security policy up migrations from the next version to the last version.
func MigrateUp() error {
	openDB()
	defer closeDB()
	op := policyMigration.NewOperation(Migrations, config.Config.Security)
	return op.Up()
}

// MigrationDown runs the security policy down migration for the current version.
func MigrateDown() error {
	openDB()
	defer closeDB()
	op := policyMigration.NewOperation(Migrations, config.Config.Security)
	return op.Down()
}

func openDB() {
	DB = dbPkg.NewDBConnection(config.Config.PostgresDatabase)
}

func closeDB() {
	if DB != nil {
		DB.Close()
	}
}

`

const migrationFileContent = `package dbPolicy

import (
	"errors"

  policyMigration "github.com/nrfta/go-platform-security-policy-migration/pkg/policy_migration"
)

func init() {
	Migrations[%s] = &policyMigration.Migration{
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
`

const testSuiteFileContent = `package dbPolicy_test

import (
	"strconv"
	"testing"

	dbPkg "github.com/nrfta/%s/db"
	"github.com/nrfta/%s/tests"
	amTests "github.com/nrfta/go-access-management/v3/tests"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPolicyMigrations(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Policy Migrations Suite")
}

var (
	DB dbPkg.DB
)

func openDB() {
	// Gets current test to create a unique identifier
	test := CurrentGinkgoTestDescription()
	id := test.FileName + ":" + strconv.Itoa(test.LineNumber)

	// Creates a new DB connection with a empty database for each test
	DB = tests.NewDBConnection(id)
	dbPolicy.DB = DB
	tests.ResetAccessManagementPolicies()
}

func closeDB() {
	if DB != nil {
		DB.Close()
	}
	amTests.Finalize()
}

var _ = BeforeSuite(openDB)
var _ = AfterSuite(closeDB)
var _ = BeforeEach(openDB)
var _ = AfterEach(closeDB)
`

const testFileContent = `package dbPolicy_test

import (
	dbPolicy "github.com/nrfta/%s/db/policy"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("%s", func() {
	migrationVersion := int64(%s)
	migration, ok := dbPolicy.Migrations[migrationVersion]
	Expect(ok).To(BeTrue())
	Expect(migration.Version).To(Equal(migrationVersion))

	Context("Up", func() {
		It("should migrate up and apply correct policies", func() {
			// err := migration.Up(migration.PolicyVersion)
			// Expect(err).To(BeNil())
			Expect(false).To(BeTrue(), "policy migration tests are required")
		})
	})

	Context("Down", func() {
		It("should migrate down and revoke correct policies", func() {
			// err := migration.Down(migration.PolicyVersion)
			// Expect(err).To(BeNil())
			Expect(false).To(BeTrue(), "policy migration tests are required")
		})
	})
})
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
		baseDir := path.Join(rootPath, "db/policy")
		baseFileName := path.Join(baseDir, timestamp+"_"+name)

		cfg := helpers.LoadConfig()
		serviceName := cfg.Meta.ServiceName

		if err := ensureMigrationsFileExists(path.Join(rootPath, "./db/policy/0_migrations.go"), serviceName); err != nil {
			log.Error(err)
			return
		}
		if err := ensureTestSuiteFileExists(path.Join(rootPath, "./db/policy/0_migrations_suite_test.go"), serviceName); err != nil {
			log.Error(err)
			return
		}
		if err := ensureGriftFileExists(path.Join(rootPath, "./grifts/policy.go"), serviceName); err != nil {
			log.Error(err)
			return
		}

		content := fmt.Sprintf(migrationFileContent, timestamp, timestamp, args[0], policyIDsStr, args[2], timestamp,
			timestamp, timestamp, timestamp)
		err := writeFile(baseFileName+".go", content)
		if err != nil {
			log.Error(err)
			return
		}
		fmt.Println(baseFileName + ".go created")

		testContent := fmt.Sprintf(testFileContent, serviceName, timestamp, timestamp)
		testFilePath := path.Join(baseDir, timestamp+"_test.go")
		err = writeFile(testFilePath, testContent)
		if err != nil {
			log.Error(err)
			return
		}

		fmt.Println(testFilePath + " created")
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
		return writeFile(fileName, fmt.Sprintf(mainPolicyFileContent, serviceName, serviceName))
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

func ensureTestSuiteFileExists(fileName, serviceName string) error {
	dir := path.Dir(fileName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return writeFile(fileName, fmt.Sprintf(testSuiteFileContent, serviceName, serviceName))
	}
	return nil
}
