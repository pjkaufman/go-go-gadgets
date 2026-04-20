package cmd

import (
	"os"

	"github.com/pjkaufman/go-go-gadgets/pkg/crawler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	verbose            bool
	version1, version2 string
)

var rootCmd = &cobra.Command{
	Use:           "versy",
	Short:         "A verse of the day retriever for two translations",
	SilenceErrors: true, // avoids double printing of errors when thrown
	Run: func(cmd *cobra.Command, args []string) {
		var scrapper = crawler.CreateNewCollyCrawler(userAgent, verbose, allowedDomains)

		reference, err := getVerseOfTheDayReference(scrapper)
		if err != nil {
			logger.WriteError(err.Error())
		}

		getAndDisplayBothVerses(reference, version1, version2, scrapper.Clone())
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.WriteErrorf("Error: %v\n", err)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "show more info about what is going on")
	rootCmd.PersistentFlags().StringVarP(&version1, "translation-a", "", "ESV", "gets the verse reference specified in this translation first (default is ESV)")
	rootCmd.PersistentFlags().StringVarP(&version2, "translation-b", "", "NVI", "gets the verse reference specified in this translation second (default is NVI)")

	rootCmd.SetOut(os.Stdout)
}
