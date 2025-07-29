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
	Use:   "versy",
	Short: "A verse of the day retriever for two translations",
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
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "show more info about what is going on")
	rootCmd.PersistentFlags().StringVarP(&version1, "version-one", "", "ESV", "gets the first instance of the verse in the specified version (default is ESV)")
	rootCmd.PersistentFlags().StringVarP(&version2, "version-two", "", "NVI", "gets the second instance of the verse in the specified version (default is NVI)")
}
