package cmd

import (
	"time"

	"github.com/Luzifer/ansible-role-version/tags"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Seek for updates in git repositories and update the roles file",
	RunE: func(cmd *cobra.Command, args []string) error {
		roles, err := getRoleDefinitions(cfg.RolesFile)
		if err != nil {
			return err
		}

		updates := map[string]string{}

		for _, role := range roles {
			logger := log.WithFields(log.Fields{
				"role": role.Name,
			})

			tag, err := tags.GetLatestTag(role.Src, true)
			if err != nil {
				logger.WithError(err).Error("Failed to fetch latest tag")
				continue
			}

			if tag.Name != role.Version {
				updates[role.Name] = tag.Name
				logger.WithFields(log.Fields{
					"from":     role.Version,
					"to":       tag.Name,
					"released": tag.When.Format(time.RFC1123),
				}).Info("Update queued")
			}
		}

		return patchRoleFile(cfg.RolesFile, updates)
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
}
