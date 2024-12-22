package cli

import (
	"github.com/kenean-50/vm-container-manager/internal/domain/deploy"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type cobraCli struct {
	name    string
	command *cobra.Command
	logger  zerolog.Logger
}

func NewCobraCli(name string, deploy deploy.DeployService) *cobraCli {

	var command = &cobra.Command{Use: name}

	command.AddCommand(ValidateCommand())
	command.AddCommand(DeployCommand(deploy))

	return &cobraCli{
		name:    name,
		command: command,
		logger: log.
			With().
			Str("actor", "cobra").
			Logger(),
	}
}

func (r *cobraCli) Execute() error {
	return r.command.Execute()
}
