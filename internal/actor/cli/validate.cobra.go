package cli

import (
  "fmt"
  "strings"

  "github.com/spf13/cobra"
)

func ValidateCommand() *cobra.Command {

	command := &cobra.Command{
    Use:   "validate",
    Short: "Validate the repository configuration",
    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
      fmt.Println("Print: validating config " + strings.Join(args, " "))
    },
  }

	return command
}