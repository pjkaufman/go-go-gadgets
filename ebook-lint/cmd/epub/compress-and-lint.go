package epub

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	filesize "github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/file-size"
	"github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/images"
	"github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/linter"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	lintDir           string
	lang              string
	runCompressImages bool
)

const (
	LintDirArgEmpty = "directory must have a non-whitespace value"
	LangArgEmpty    = "lang must have a non-whitespace value"
)

// compressAndLintCmd represents the compressAndLint command
var compressAndLintCmd = &cobra.Command{
	Use:   "compress-and-lint",
	Short: "Compresses and lints all of the epub files in the specified directory even compressing images using imgp if that option is specified.",
	Example: heredoc.Doc(`To compress images and make general modifications to all epubs in a folder:
	ebook-lint epub compress-and-lint -d folder -i
	
	To compress images and make general modifications to all epubs in the current directory:
	ebook-lint epub compress-and-lint -i

	To just make general modifications to all epubs in the current directory:
	ebook-lint epub compress-and-lint
	`),
	Long: heredoc.Doc(`Gets all of the .epub files in the specified directory.
	Then it lints each epub separately making sure to compress the images if specified.
	Some of the things that the linting includes:
	- Replacing a list of common strings
	- Adds language encoding specified if it is not present already (default is "en")
	- Sets encoding on content files to utf-8 to prevent errors in some readers
	`),
	Run: func(cmd *cobra.Command, args []string) {
		err := ValidateCompressAndLintFlags(lintDir, lang)
		if err != nil {
			logger.WriteError(err.Error())
		}

		logger.WriteInfo("Starting compression and linting for each epub\n")

		epubs := filehandler.MustGetAllFilesWithExtInASpecificFolder(lintDir, ".epub")

		var totalBeforeFileSize, totalAfterFileSize float64
		for _, epub := range epubs {
			logger.WriteInfo(fmt.Sprintf("starting epub compressing for %s...", epub))

			LintEpub(lintDir, epub, runCompressImages)

			var originalFile = epub + ".original"
			var newKbSize = filehandler.MustGetFileSize(epub)
			var oldKbSize = filehandler.MustGetFileSize(originalFile)

			logger.WriteInfo(filesize.FileSizeSummary(originalFile, epub, oldKbSize, newKbSize))

			totalBeforeFileSize += oldKbSize
			totalAfterFileSize += newKbSize
		}

		logger.WriteInfo(filesize.FilesSizeSummary(totalBeforeFileSize, totalAfterFileSize))
		logger.WriteInfo("Finished compression and linting")
	},
}

func init() {
	EpubCmd.AddCommand(compressAndLintCmd)

	compressAndLintCmd.Flags().StringVarP(&lintDir, "directory", "d", ".", "the location to run the epub lint logic")
	compressAndLintCmd.Flags().StringVarP(&lang, "lang", "l", "en", "the language to add to the xhtml, htm, or html files if the lang is not already specified")
	compressAndLintCmd.Flags().BoolVarP(&runCompressImages, "compress-images", "i", false, "whether or not to also compress images which requires imgp to be installed")
}

// TODO: make this function return an error
func LintEpub(lintDir, epub string, runCompressImages bool) {
	var src = filehandler.JoinPath(lintDir, epub)
	var dest = filehandler.JoinPath(lintDir, "epub")

	err := filehandler.UnzipRunOperationAndRezip(src, dest, func() {
		opfFolder, epubInfo := getEpubInfo(dest, epub)

		validateFilesExist(opfFolder, epubInfo.HtmlFiles)
		validateFilesExist(opfFolder, epubInfo.ImagesFiles)
		validateFilesExist(opfFolder, epubInfo.OtherFiles)

		// fix up all xhtml files first
		for file := range epubInfo.HtmlFiles {
			var filePath = getFilePath(opfFolder, file)
			fileText, err := filehandler.ReadInFileContents(filePath)
			if err != nil {
				logger.WriteError(err.Error())
			}

			var newText = linter.EnsureEncodingIsPresent(fileText)
			newText = linter.CommonStringReplace(newText)

			newText = linter.EnsureLanguageIsSet(newText, lang)

			if fileText == newText {
				continue
			}

			filehandler.WriteFileContents(filePath, newText)
		}

		//TODO: get all files in the repo and prompt the user whether they want to delete them if they are not in the manifest

		if runCompressImages {
			images.CompressRelativeImages(opfFolder, epubInfo.ImagesFiles)
		}
	})

	if err != nil {
		logger.WriteError(err.Error())
	}
}

func ValidateCompressAndLintFlags(lintDir, lang string) error {
	if strings.TrimSpace(lintDir) == "" {
		return errors.New(LintDirArgEmpty)
	}

	if strings.TrimSpace(lang) == "" {
		return errors.New(LangArgEmpty)
	}

	return nil
}

func getFilePath(opfFolder, file string) string {
	return filehandler.JoinPath(opfFolder, file)
}
