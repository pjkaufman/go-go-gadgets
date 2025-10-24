package cmd

import (
	"github.com/spf13/cobra"
)

// fixCmd represents the fix command
var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Deals with fixing things with an epub file",
}

func init() {
	rootCmd.AddCommand(fixCmd)
}
