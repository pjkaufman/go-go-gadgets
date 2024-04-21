package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/pjkaufman/go-go-gadgets/cat-ascii/internal/ascii"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

// listCmd represents the list names command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists the names of all cat ascii options",
	Example: heredoc.Doc(`To list all of the names for the cat ascii art:
	cat-ascii list
	`),
	Run: func(cmd *cobra.Command, args []string) {
		listNames()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listNames() {
	for _, catAscii := range ascii.CAT_ASCII {
		logger.WriteInfo(catAscii.Name)
	}
}
