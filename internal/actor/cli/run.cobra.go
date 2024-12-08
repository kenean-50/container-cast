package cli

import (
  "fmt"
  "strings"

  "github.com/spf13/cobra"
)

func RunCommand() *cobra.Command {

	command := &cobra.Command{
    Use:   "run",
    Short: "run service on the server",
    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
      fmt.Println("Print: running service " + strings.Join(args, " "))
    },
  }

	return command
}