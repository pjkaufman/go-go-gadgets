package cmd

import (
	"archive/zip"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/linter"
	"github.com/pjkaufman/go-go-gadgets/pkg/cli/flags"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	extraReplacesFilePath string
	replaceFlags          = flags.Flags{
		Flags: []flags.Flag{
			flags.NewFileFlag(true, false, &extraReplacesFilePath, "replacements", "e", "", "the path to the file with extra strings to replace", []string{"md"}, true),
			flags.NewFileFlag(true, false, &epubFile, "file", "f", "", "the epub file to replace strings in in", []string{"epub"}, true),
		},
	}
)

// replaceCmd represents the replace string command
var replaceCmd = &cobra.Command{
	Use:   "replace",
	Short: "Replaces a list of common strings and the extra strings for all content/xhtml files in the provided epub",
	Long: heredoc.Doc(`Uses the provided epub and extra replace Markdown file to replace a common set of strings and any extra instances specified in the extra file replace. After all replacements are made, the original epub will be moved to a .original file and the new file will take the place of the old file. It will also print out the successful extra replacements with the number of replacements made followed by warnings for any extra strings that it tried to find and replace values for, but did not find any instances to replace.
		Note: it only replaces strings in content/xhtml files listed in the opf file.`),
	Example: heredoc.Doc(`
		epub-lint replace -f test.epub -e replacements.md
		will replace the common strings and extra strings parsed out of replacements.md in content/xhtml files located in test.epub.
		The original test.epub will be moved to test.epub.original and test.epub will have the updated files.

		replacements.md is expected to be in the following format:
		| Text to replace | Text to replace with |
		| --------------- | -------------------- |
		| I am typo | I the correct value |
		...
		| I am another issue to correct | the correction |
	`),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return replaceFlags.Validate()
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger.WriteInfo("Starting epub string replacement...\n")

		var numHits = make(map[string]int)
		extraReplaceContents, err := filehandler.ReadInFileContents(extraReplacesFilePath)
		if err != nil {
			logger.WriteFatal(err.Error())
		}

		extraTextReplacements, err := linter.ParseTextReplacements(extraReplaceContents)
		if err != nil {
			logger.WriteFatal(err.Error())
		}

		err = epubhandler.UpdateEpub(epubFile, func(zipFiles map[string]*zip.File, w *zip.Writer, epubInfo epubhandler.EpubInfo, opfFolder string) ([]string, error) {
			err = validateFilesExist(opfFolder, epubInfo.HtmlFiles, zipFiles)
			if err != nil {
				return nil, err
			}

			var handledFiles []string

			for file := range epubInfo.HtmlFiles {
				var filePath = getFilePath(opfFolder, file)
				zipFile := zipFiles[filePath]

				fileText, err := filehandler.ReadInZipFileContents(zipFile)
				if err != nil {
					return nil, err
				}

				var newText = linter.CommonStringReplace(fileText)
				newText = linter.ExtraStringReplace(newText, extraTextReplacements, numHits)

				err = filehandler.WriteZipCompressedString(w, filePath, newText)
				if err != nil {
					return nil, err
				}

				handledFiles = append(handledFiles, filePath)
			}

			var successfulReplaces []string
			var failedReplaces []string
			for searchText, hits := range numHits {
				if hits == 0 {
					failedReplaces = append(failedReplaces, searchText)
				} else {
					var timeText = "time"
					if hits > 1 {
						timeText += "s"
					}

					successfulReplaces = append(successfulReplaces, fmt.Sprintf("`%s` was replaced %d %s", searchText, hits, timeText))
				}
			}

			logger.WriteInfo("Successful Replaces:")
			for _, successfulReplace := range successfulReplaces {
				logger.WriteInfo(successfulReplace)
			}

			if len(failedReplaces) == 0 {
				return handledFiles, nil
			}

			logger.WriteWarn("\nFailed Replaces:")
			for i, failedReplace := range failedReplaces {
				logger.WriteWarnf("%d. %s\n", i+1, failedReplace)
			}

			return handledFiles, nil
		})
		if err != nil {
			logger.WriteFatalf("failed to replace strings in %q: %s", epubFile, err)
		}

		logger.WriteInfo("\nFinished epub string replacement...")
	},
}

func init() {
	rootCmd.AddCommand(replaceCmd)

	err := replaceFlags.AddToCmd(replaceCmd)
	if err != nil {
		logger.WriteFatal(err.Error())
	}
}
