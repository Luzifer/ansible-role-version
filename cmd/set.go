package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set <role> <version>",
	Short: "Set the version of a role statically",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("You must specify role and version")
		}

		return patchRoleFile(cfg.RolesFile, map[string]string{args[0]: args[1]})
	},
}

func init() {
	RootCmd.AddCommand(setCmd)
}
