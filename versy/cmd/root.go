package cmd

import (
	"os"

	"github.com/pjkaufman/go-go-gadgets/pkg/cli/flags"
	"github.com/pjkaufman/go-go-gadgets/pkg/crawler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	verbose            bool
	version1, version2 string
	rootFlags          = flags.Flags{
		Flags: []flags.Flag{
			flags.NewBoolFlag(false, true, &verbose, "verbose", "v", false, "show more info about what is going on"),
			flags.NewStringFlag(false, true, &version1, "translation-a", "", "ESV", "gets the verse reference specified in this translation first (default is ESV)"),
			flags.NewStringFlag(false, true, &version2, "translation-b", "", "NVI", "gets the verse reference specified in this translation second (default is NVI)"),
		},
	}
)

var rootCmd = &cobra.Command{
	Use:           "versy",
	Short:         "A verse of the day retriever for two translations",
	SilenceErrors: true, // avoids double printing of errors when thrown
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return rootFlags.Validate()
	},
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
	err := rootFlags.AddToCmd(rootCmd)
	if err != nil {
		logger.WriteError(err.Error())
	}

	rootCmd.SetOut(os.Stdout)
}
