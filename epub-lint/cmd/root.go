package cmd

import (
	"os"

	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

// EpubCmd represents the base command when called without any subcommands
var EpubCmd = &cobra.Command{
	Use:   "epub-lint",
	Short: "A set of functions that are helpful for linting ebooks",
}

func Execute() {
	if err := EpubCmd.Execute(); err != nil {
		logger.WriteErrorf("Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
}
