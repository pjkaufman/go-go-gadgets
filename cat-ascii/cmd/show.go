package cmd

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/pjkaufman/go-go-gadgets/cat-ascii/internal/ascii"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

// showCmd represents the show art command
var showCmd = &cobra.Command{
	Use:   "show [flags] [ASCII_NAME]",
	Short: "Shows the cat ascii and info about the cat ascii that matches the provided name",
	Args:  cobra.ExactArgs(1),
	Example: heredoc.Doc(`To list information about the cat ascii and the cat ascii itself:
	cat-ascii show stalking-cat
	`),
	Run: func(cmd *cobra.Command, args []string) {
		var name = args[0]
		logger.WriteInfo(fmt.Sprintf("Name %s:", name))
		for _, catAscii := range ascii.CAT_ASCII {
			if strings.EqualFold(name, catAscii.Name) {
				logger.WriteInfo(catAscii.Ascii)
				logger.WriteInfo(catAscii.From)

				return
			}
		}

		logger.WriteWarn("\nNot found. Please use one of the following names:")
		listNames()
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
