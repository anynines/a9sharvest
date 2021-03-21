package cmd

import (
	"github.com/anynines/a9sharvest/pkg/harvest"

	"github.com/spf13/cobra"
)

func newCmdDownload() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Downloads entries for further processing",
		Long:  `Downloads entries for further processing with the group commands. This saves the time to download the entirs if you need to apply several filters on the same data.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFlag, _ := cmd.Flags().GetString("output")
			verboseFlag, _ := cmd.Flags().GetBool("verbose")

			return harvest.DownloadAndOutputEntries(verboseFlag, outputFlag)
		},
	}

	cmd.Flags().StringP("output", "o", "", "If empty, outputs to STDOUT, otherwise to <filename>.")
	cmd.Flags().BoolP("verbose", "v", false, "Provide additional debug details.")
	return cmd
}
