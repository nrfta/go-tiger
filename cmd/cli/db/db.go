package db

import (
	_ "github.com/lib/pq" // Postgres
	"github.com/spf13/cobra"
)

// DBCmd represents the DB command when called without any subcommands
var DBCmd = &cobra.Command{
	Use:   "db",
	Short: "Database related operations",
	Long:  "Execute operations againt the database",
}

func init() {
}
