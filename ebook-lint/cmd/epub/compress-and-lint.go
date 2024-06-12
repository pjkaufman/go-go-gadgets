package epub

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	epubhandler "github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/epub-handler"
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

		epubs, err := filehandler.MustGetAllFilesWithExtInASpecificFolder(lintDir, ".epub")
		if err != nil {
			logger.WriteError(err.Error())
		}

		var totalBeforeFileSize, totalAfterFileSize float64
		for _, epub := range epubs {
			logger.WriteInfof("starting epub compressing for %s...\n", epub)

			err = LintEpub(lintDir, epub, runCompressImages)
			if err != nil {
				logger.WriteError(err.Error())
			}

			var originalFile = epub + ".original"
			newKbSize, err := filehandler.MustGetFileSize(epub)
			if err != nil {
				logger.WriteError(err.Error())
			}

			oldKbSize, err := filehandler.MustGetFileSize(originalFile)
			if err != nil {
				logger.WriteError(err.Error())
			}

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

func LintEpub(lintDir, epub string, runCompressImages bool) error {
	var src = filehandler.JoinPath(lintDir, epub)
	// var dest = filehandler.JoinPath(lintDir, "epub")

	err := epubhandler.UpdateEpub(src, func(zipFiles map[string]*zip.File, w *zip.Writer, epubInfo epubhandler.EpubInfo, opfFolder string) []string {
		validateFilesExist(opfFolder, epubInfo.HtmlFiles, zipFiles)
		validateFilesExist(opfFolder, epubInfo.ImagesFiles, zipFiles)
		validateFilesExist(opfFolder, epubInfo.OtherFiles, zipFiles)

		var handledFiles []string

		// fix up all xhtml files first
		for file := range epubInfo.HtmlFiles {
			var filePath = getFilePath(opfFolder, file)

			zipFile := zipFiles[filePath]

			fileText, err := filehandler.ReadInZipFileContents(zipFile)
			if err != nil {
				logger.WriteError(err.Error())
			}

			var newText = linter.EnsureEncodingIsPresent(fileText)
			newText = linter.CommonStringReplace(newText)

			newText = linter.EnsureLanguageIsSet(newText, lang)

			if fileText == newText {
				continue
			}

			err = filehandler.WriteZipCompressedString(w, filePath, newText)
			if err != nil {
				logger.WriteError(err.Error())
			}

			handledFiles = append(handledFiles, filePath)
		}

		//TODO: get all files in the repo and prompt the user whether they want to delete them if they are not in the manifest

		if runCompressImages {
			for imagePath := range epubInfo.ImagesFiles {
				var filePath = filehandler.JoinPath(opfFolder, imagePath)

				imageFile := zipFiles[filePath]

				data, err := filehandler.ReadInZipFileBytes(imageFile)
				if err != nil {
					logger.WriteError(err.Error())
				}

				newData, err := images.CompressImage(filePath, data)
				if err != nil {
					logger.WriteError(err.Error())
				}

				if bytes.Equal(data, newData) {
					continue
				}

				err = filehandler.WriteZipCompressedBytes(w, filePath, newData)
				if err != nil {
					logger.WriteError(err.Error())
				}

				handledFiles = append(handledFiles, filePath)
			}
		}

		return handledFiles
	})
	if err != nil {
		logger.WriteError(fmt.Sprintf("failed to update epub %q: %s", src, err))
	}

	// filehandler.UnzipRunOperationAndRezip(src, dest, func() {
	// 	opfFolder, epubInfo := getEpubInfo(dest, epub)

	// validateFilesExist(opfFolder, epubInfo.HtmlFiles)
	// validateFilesExist(opfFolder, epubInfo.ImagesFiles)
	// validateFilesExist(opfFolder, epubInfo.OtherFiles)

	// // fix up all xhtml files first
	// for file := range epubInfo.HtmlFiles {
	// 	var filePath = getFilePath(opfFolder, file)
	// 	fileText := filehandler.ReadInFileContents(filePath)
	// 	var newText = linter.EnsureEncodingIsPresent(fileText)
	// 	newText = linter.CommonStringReplace(newText)

	// 	newText = linter.EnsureLanguageIsSet(newText, lang)

	// 	if fileText == newText {
	// 		continue
	// 	}

	// 	filehandler.WriteFileContents(filePath, newText)
	// }

	// //TODO: get all files in the repo and prompt the user whether they want to delete them if they are not in the manifest

	// if runCompressImages {
	// 	images.CompressRelativeImages(opfFolder, epubInfo.ImagesFiles)
	// }
	// })
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
