package cmd

import (
	"github.com/anynines/a9sharvest/pkg/harvest"

	"github.com/spf13/cobra"
)

func newCmdGroup() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group",
		Short: "Groups time entries by given tags",
		Long:  `Groups time entries by given tags.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			verboseFlag, _ := cmd.Flags().GetBool("verbose")

			return harvest.Group(verboseFlag)
		},
	}

	cmd.Flags().BoolP("verbose", "v", false, "Provide additional debug details.")
	return cmd
}
