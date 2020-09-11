package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/anynines/a9sharvest/pkg/version"

	"github.com/spf13/cobra"
)

type versionOptions struct {
	shortVersion bool
}

func newVersionOptions() *versionOptions {
	return &versionOptions{
		shortVersion: false,
	}
}

func newCmdVersion() *cobra.Command {
	options := newVersionOptions()

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the client version information",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			runVersion(options, os.Stdout)
		},
	}

	cmd.PersistentFlags().BoolVar(&options.shortVersion, "short", options.shortVersion, "Print the version number(s) only, with no additional output")

	return cmd
}

func runVersion(options *versionOptions, stdout io.Writer) {
	clientVersion := version.Version

	if options.shortVersion {
		fmt.Fprintln(stdout, clientVersion)
	} else {
		fmt.Fprintf(stdout, "Client version: %s\n", clientVersion)
	}
}
