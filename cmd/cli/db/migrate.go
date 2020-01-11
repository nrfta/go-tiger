package db

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Add source file for migrate
	_ "github.com/lib/pq"                                // Postgres
	"github.com/spf13/cobra"

	"github.com/neighborly/go-errors"
	"github.com/neighborly/go-pghelpers"
	"github.com/nrfta/go-log"
	"github.com/nrfta/go-tiger/helpers"
)

func setupMigrate() *migrate.Migrate {
	pgConfig := helpers.LoadConfig().PostgresDatabase

	db, err := pghelpers.ConnectPostgres(pgConfig)
	if err != nil {
		log.Fatal(errors.Wrap(err, "migrate unable to connect to postgres instance"))
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(errors.Wrap(err, "migrate unable to connect to postgres instance"))
	}

	migrateInstance, err := migrate.NewWithDatabaseInstance(
		"file:"+helpers.FindRootPath()+"/db/migrations",
		"postgres",
		driver,
	)

	if err != nil {
		log.Fatal(errors.Wrap(err, "unable to instantiate migrate tool"))
	}

	return migrateInstance
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Execute database migrations",
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run migrations in the up direction",
	Run: func(cmd *cobra.Command, args []string) {
		mi := setupMigrate()
		err := mi.Up()

		if err != nil && err != migrate.ErrNoChange {
			log.Fatal(errors.Wrap(err, "unable to run migrations"))
		} else {

			if err == migrate.ErrNoChange {
				log.Info("No migrations to run.")
			} else {
				log.Info("Database successfuly migrated!")
			}
		}
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Run migrations in the down direction",
	Run: func(cmd *cobra.Command, args []string) {
		mi := setupMigrate()
		err := mi.Steps(-1)

		if err != nil {
			log.Fatal(errors.Wrap(err, "unable to run migrations"))
		}
	},
}

func init() {
	migrateCmd.AddCommand(upCmd)
	migrateCmd.AddCommand(downCmd)
	DBCmd.AddCommand(migrateCmd)
}
