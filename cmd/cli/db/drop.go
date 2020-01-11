package db

import (
	"fmt"

	_ "github.com/lib/pq" // Postgres
	"github.com/spf13/cobra"

	"github.com/neighborly/go-errors"
	"github.com/neighborly/go-pghelpers"
	"github.com/nrfta/go-log"
	"github.com/nrfta/go-tiger/helpers"
)

var dropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop database",
	Long:  "Drop database defined in the config",
	Run: func(cmd *cobra.Command, args []string) {
		config := helpers.LoadConfig()

		if config.Meta.Environment == "prod" {
			log.Fatal(errors.New("you cannot drop database in production environment"))
		}

		pgConfig := config.PostgresDatabase

		query := fmt.Sprintf("DROP DATABASE \"%s\";", pgConfig.Database)

		pgConfigCopy := pgConfig
		pgConfigCopy.Database = "postgres"
		db, _ := pghelpers.ConnectPostgres(pgConfigCopy)
		_, err := db.Exec(query)

		if err != nil {
			log.Info(errors.Wrap(err, "unable to drop database"))
		} else {
			log.Info(query)
		}
	},
}

func init() {
	DBCmd.AddCommand(dropCmd)
}
