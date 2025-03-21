package cli

import "github.com/spf13/cobra"

// Run will initialize and execute the given Cobra command
func Run(command *cobra.Command) int {
	if err := command.Execute(); err != nil {
		return 1
	}
	return 0
}
