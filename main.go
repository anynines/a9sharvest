package main

import (
	"os"

	"github.com/anynines/a9sharvest/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
