package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd represents the root Cobra command
var RootCmd = &cobra.Command{
	Use:   "a9sharvest",
	Short: "a9sharvest is a CLI to aha.io.",
	Long:  `a9sharvest is a CLI to aha.io.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	RootCmd.AddCommand(newCmdGroup())
	RootCmd.AddCommand(newCmdVersion())
}
