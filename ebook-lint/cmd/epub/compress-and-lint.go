package epub

import (
	"archive/zip"
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
	lintDir            string
	lang               string
	runCompressImages  bool
	ErrLintDirArgEmpty = errors.New("directory must have a non-whitespace value")
	ErrLangArgEmpty    = errors.New("lang must have a non-whitespace value")
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

		epubs, err := filehandler.GetAllFilesWithExtInASpecificFolder(lintDir, ".epub")
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
			newKbSize, err := filehandler.GetFileSize(epub)
			if err != nil {
				logger.WriteError(err.Error())
			}

			oldKbSize, err := filehandler.GetFileSize(originalFile)
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
	err := epubhandler.UpdateEpub(src, func(zipFiles map[string]*zip.File, w *zip.Writer, epubInfo epubhandler.EpubInfo, opfFolder string) ([]string, error) {
		err := validateFilesExist(opfFolder, epubInfo.HtmlFiles, zipFiles)
		if err != nil {
			return nil, err
		}

		err = validateFilesExist(opfFolder, epubInfo.ImagesFiles, zipFiles)
		if err != nil {
			return nil, err
		}

		err = validateFilesExist(opfFolder, epubInfo.OtherFiles, zipFiles)
		if err != nil {
			return nil, err
		}

		var handledFiles []string

		// fix up all xhtml files first
		for file := range epubInfo.HtmlFiles {
			var filePath = getFilePath(opfFolder, file)

			zipFile := zipFiles[filePath]

			fileText, err := filehandler.ReadInZipFileContents(zipFile)
			if err != nil {
				return nil, err
			}

			var newText = linter.EnsureEncodingIsPresent(fileText)
			newText = linter.CommonStringReplace(newText)

			newText = linter.EnsureLanguageIsSet(newText, lang)

			err = filehandler.WriteZipCompressedString(w, filePath, newText)
			if err != nil {
				return nil, err
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
					return nil, err
				}

				newData, err := images.CompressImage(filePath, data)
				if err != nil {
					return nil, err
				}

				err = filehandler.WriteZipCompressedBytes(w, filePath, newData)
				if err != nil {
					return nil, err
				}

				handledFiles = append(handledFiles, filePath)
			}
		}

		return handledFiles, nil
	})
	if err != nil {
		return fmt.Errorf("failed to update epub %q: %s", src, err)
	}

	return nil
}

func ValidateCompressAndLintFlags(lintDir, lang string) error {
	if strings.TrimSpace(lintDir) == "" {
		return ErrLintDirArgEmpty
	}

	if strings.TrimSpace(lang) == "" {
		return ErrLangArgEmpty
	}

	return nil
}

func getFilePath(opfFolder, file string) string {
	return filehandler.JoinPath(opfFolder, file)
}
