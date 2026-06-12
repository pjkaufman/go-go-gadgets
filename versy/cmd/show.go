package cmd

import (
	"github.com/pjkaufman/go-go-gadgets/pkg/cli/flags"
	"github.com/pjkaufman/go-go-gadgets/pkg/crawler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	verseReference string
	showFlags      = flags.Flags{
		Flags: []flags.Flag{
			flags.NewStringFlag(true, false, &verseReference, "verse", "", "", "the Bible verse to get the two versions of"),
		},
	}
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Displays the specified verse reference in the two specified Bible versions",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := rootFlags.Validate()
		if err != nil {
			return err
		}

		return showFlags.Validate()
	},
	Run: func(cmd *cobra.Command, args []string) {
		var scrapper = crawler.CreateNewCollyCrawler(userAgent, verbose, allowedDomains)

		getAndDisplayBothVerses(verseReference, version1, version2, scrapper)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	err := showFlags.AddToCmd(showCmd)
	if err != nil {
		logger.WriteError(err.Error())
	}
}
