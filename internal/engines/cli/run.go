package cli

import (
	"github.com/spf13/cobra"
)

// Run запускает переданную команду и возвращает код завершения
func RunCommand(cmd *cobra.Command) int {
	if err := cmd.Execute(); err != nil {
		return 1
	}
	return 0
}
