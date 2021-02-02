package am

import (
	"log"

	am "github.com/nrfta/go-access-management"
	"github.com/spf13/cobra"
)

func init() {
	AMCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use: "remove",

	Short:   "Remove user from am",
	Long:    "Add User to Domain with Role in Casbin AM. Domains will still need to be populated in JWTs issued to user by other tools.",
	Example: "tiger am remove {{user id}}",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		initAM()
		defer closeAM()
		if err := am.DeleteUser(args[0]); err != nil {
			log.Fatalf("error removing user to roles %s, %s", args, err.Error())
		}
	},
}
