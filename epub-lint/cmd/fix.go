package cmd

import (
	"github.com/spf13/cobra"
)

// fixCmd represents the fix command
var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Deals with fixing things with an epub file",
	// Example: heredoc.Doc(`
	// 	epub-lint fix-validation -f test.epub --issue-file epubCheckOutput.txt
	// 	will read in the contents of the file and try to fix any of the fixable
	// 	validation issues

	// 	epub-lint fix-validation -f test.epub --issue-file epubCheckOutput.txt --cleanup-jnovels
	// 	will read in the contents of the file and try to fix any of the fixable
	// 	validation issues as well as remove any jnovels specific files
	// `),
}

func init() {
	rootCmd.AddCommand(fixCmd)
}
