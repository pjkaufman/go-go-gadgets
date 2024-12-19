package cmd

import (
	"os"

	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jpproc",
	Short: "JPEG and PNG image processor",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.WriteErrorf("Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {}
