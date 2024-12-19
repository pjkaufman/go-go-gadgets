package cmd

import (
	"os"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/cmd/cbr"
	"github.com/pjkaufman/go-go-gadgets/ebook-lint/cmd/cbz"
	"github.com/pjkaufman/go-go-gadgets/ebook-lint/cmd/epub"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ebook-lint",
	Short: "A set of functions that are helpful for linting ebooks",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.WriteErrorf("Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(epub.EpubCmd)
	rootCmd.AddCommand(cbz.CbzCmd)
	rootCmd.AddCommand(cbr.CbrCmd)
}
