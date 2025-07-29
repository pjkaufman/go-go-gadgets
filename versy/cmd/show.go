package cmd

import (
	"github.com/pjkaufman/go-go-gadgets/pkg/crawler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var verseReference string

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Displays the specified verse reference in the two specified Bible versions",
	Run: func(cmd *cobra.Command, args []string) {
		var scrapper = crawler.CreateNewCollyCrawler(userAgent, verbose, allowedDomains)

		getAndDisplayBothVerses(verseReference, version1, version2, scrapper)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	showCmd.Flags().StringVarP(&verseReference, "verse", "", "", "the Bible verse to get the two versions of")
	err := showCmd.MarkFlagRequired("verse")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"verse\" as required on show command: %v\n", err)
	}
}
