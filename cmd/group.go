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
			inputFlag, _ := cmd.Flags().GetString("input")
			outputFlag, _ := cmd.Flags().GetString("output")
			verboseFlag, _ := cmd.Flags().GetBool("verbose")
			skipUnknown, _ := cmd.Flags().GetBool("skip-unknown")

			return harvest.Group(verboseFlag, inputFlag, outputFlag, skipUnknown)
		},
	}

	cmd.Flags().StringP("input", "i", "", "Input file. If not present, entries get downloaded.")
	cmd.Flags().StringP("output", "o", "text", "Output format: text, csv.")
	cmd.Flags().BoolP("verbose", "v", false, "Provide additional debug details.")
	cmd.Flags().BoolP("skip-unknown", "s", false, "Skip unknown entries and do not print those.")
	return cmd
}
