package db

import (
	"fmt"

	_ "github.com/lib/pq" // Postgres

	"github.com/neighborly/go-pghelpers"
	"github.com/nrfta/go-log"
	"github.com/nrfta/go-tiger/helpers"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create database",
	Long:  "Create database defined in the config",
	Run: func(cmd *cobra.Command, args []string) {
		pgConfig := helpers.LoadConfig().PostgresDatabase

		query := fmt.Sprintf("CREATE DATABASE \"%s\";", pgConfig.Database)

		pgConfigCopy := pgConfig
		pgConfigCopy.Database = "postgres"
		db, err := pghelpers.ConnectPostgres(pgConfigCopy)

		if err != nil {
			log.Fatal("Unable to connect to postgres database. Error: ", err)
		} else {
			log.Info(query)
		}
		_, err = db.Exec(query)

		if err != nil {
			log.Fatal("Unable to create database, Error: ", err)
		} else {
			log.Info(query)
		}
	},
}

func init() {
	DBCmd.AddCommand(createCmd)
}
