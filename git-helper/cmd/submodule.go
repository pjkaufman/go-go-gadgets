package cmd

import (
	"github.com/spf13/cobra"
)

// submoduleCmd represents the submodule command
var submoduleCmd = &cobra.Command{
	Use:   "submodule",
	Short: "Deals with submodules in git",
	Long:  `Handles operations on git submodules to help simplify interacting with them`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(submoduleCmd)
}
