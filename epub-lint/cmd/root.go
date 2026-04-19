package cmd

import (
	"os"

	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "epub-lint",
	Short:         "A set of functions that are helpful for linting epubs",
	SilenceErrors: true, // avoids double printing of errors when thrown
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.WriteErrorf("Error: %v\n", err)
	}
}

func init() {
	rootCmd.SetOut(os.Stdout)
}
