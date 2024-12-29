package cli

import (
	"github.com/kenean-50/container-cast/internal/domain/deploy"
	"github.com/spf13/cobra"
)

func DeployCommand(deploy deploy.DeployService) *cobra.Command {

	command := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy manifest to the server",
		// Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			deploy.Apply()
		},
	}

	return command
}
