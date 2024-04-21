package cmd

import (
	"os"

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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
