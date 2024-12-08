package cli

import (
	// "fmt"

	"github.com/spf13/cobra"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type cobraCli struct {
	name          string
	Command       *cobra.Command
	logger        zerolog.Logger
}

func NewCobraCli(name string) *cobraCli {

	var command = &cobra.Command{Use: name}

  command.AddCommand(ValidateCommand())
	command.AddCommand(RunCommand())

	return &cobraCli{
		name: name,
		Command: command,
		logger: log.
			With().
			Str("actor", "cobra").
			Logger(),
	}
}
