package cmd

import (
	"os"

	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

const (
	gitProgramName = "git"
	upADirectory   = ".."
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-helper",
	Short: "Some basic commands to help with common git actions I encounter",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.WriteErrorf("Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
}
