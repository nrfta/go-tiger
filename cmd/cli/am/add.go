package am

import (
	"log"

	am "github.com/nrfta/go-access-management"
	"github.com/spf13/cobra"
)

func init() {
	AMCmd.AddCommand(addCmd)

}

var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add User to Domain with Role",
	Long:    "Add User to Domain with Role in Casbin AM. Domains will still need to be populated in JWTs issued to user by other tools.",
	Example: "tiger am add {{user id}} {{domain}} {{role}}",
	Args:    cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		initAM()
		defer closeAM()
		if err := am.AddUserToRoles(args[0], args[1], am.Role(args[2])); err != nil {
			log.Fatalf("error adding user to roles %s, %s", args, err.Error())
		}
	},
}

