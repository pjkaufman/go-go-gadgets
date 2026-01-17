package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Deals with creating files from the song Markdown files",
}

func init() {
	rootCmd.AddCommand(createCmd)
}
