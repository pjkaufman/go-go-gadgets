package epub

import (
	"archive/zip"
	"errors"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/pjkaufman/go-go-gadgets/epub-lint/cmd"
	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
	filesize "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/file-size"
	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/images"
	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/linter"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	lintDir            string
	lang               string
	removableFileTypes string
	runCompressImages  bool
	verbose            bool
	ErrLintDirArgEmpty = errors.New("directory must have a non-whitespace value")
	ErrLangArgEmpty    = errors.New("lang must have a non-whitespace value")
)

// compressAndLintCmd represents the compressAndLint command
var compressAndLintCmd = &cobra.Command{
	Use:   "compress-and-lint",
	Short: "Compresses and lints all of the epub files in the specified directory even compressing images using imgp if that option is specified.",
	Example: heredoc.Doc(`To compress images and make general modifications to all epubs in a folder:
	epub-lint epub compress-and-lint -d folder -i
	
	To compress images and make general modifications to all epubs in the current directory:
	epub-lint epub compress-and-lint -i

	To just make general modifications to all epubs in the current directory:
	epub-lint epub compress-and-lint
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

		var removableFileExts []string
		if len(removableFileTypes) != 0 {
			removableFileExts = strings.Split(removableFileTypes, ",")
		}

		var totalBeforeFileSize, totalAfterFileSize float64
		for _, epub := range epubs {
			logger.WriteInfof("starting epub compressing for %s...\n", epub)

			err = LintEpub(lintDir, epub, runCompressImages, verbose, removableFileExts)
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
	cmd.EpubCmd.AddCommand(compressAndLintCmd)

	compressAndLintCmd.Flags().StringVarP(&lintDir, "directory", "d", ".", "the location to run the epub lint logic")
	compressAndLintCmd.Flags().StringVarP(&lang, "lang", "l", "en", "the language to add to the xhtml, htm, or html files if the lang is not already specified")
	compressAndLintCmd.Flags().StringVarP(&removableFileTypes, "removable-file-types", "", ".jpg,.jpeg,.png,.gif,.bmp,.js,.html,.htm,.xhtml,.txt,.css", "A comma separated list of file extensions of files to remove if they are not in the manifest (i.e. '.jpeg,.jpg')")
	compressAndLintCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "whether or not to show extra logs like what files were removed from the epub")
	compressAndLintCmd.Flags().BoolVarP(&runCompressImages, "compress-images", "i", false, "whether or not to also compress images which requires imgp to be installed")
}

func LintEpub(lintDir, epub string, runCompressImages, verbose bool, removableFileExts []string) error {
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
		var manifestFiles = make(map[string]struct{}, len(epubInfo.HtmlFiles)+len(epubInfo.ImagesFiles)+len(epubInfo.CssFiles)+len(epubInfo.OtherFiles))

		// fix up all xhtml files first
		for file := range epubInfo.HtmlFiles {
			var filePath = getFilePath(opfFolder, file)
			manifestFiles[filePath] = struct{}{}

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

		if runCompressImages {
			for imagePath := range epubInfo.ImagesFiles {
				var filePath = filehandler.JoinPath(opfFolder, imagePath)
				manifestFiles[filePath] = struct{}{}

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

		// handle the files that are present in the epub, but not present in the actual manifest
		if len(removableFileExts) == 0 {
			return handledFiles, nil
		}

		for otherPath := range epubInfo.OtherFiles {
			manifestFiles[filehandler.JoinPath(opfFolder, otherPath)] = struct{}{}
		}

		for otherPath := range epubInfo.CssFiles {
			manifestFiles[filehandler.JoinPath(opfFolder, otherPath)] = struct{}{}
		}

		for filePath := range zipFiles {
			if _, exists := manifestFiles[filePath]; exists {
				continue
			}

			if hasExt(removableFileExts, filePath) {
				// label file as handled despite not saving it to the destination
				handledFiles = append(handledFiles, filePath)

				if verbose {
					logger.WriteInfof("Removed file %q from the epub since it is not in the manifest.\n", filePath)
				}
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

func hasExt(slice []string, file string) bool {
	for _, item := range slice {
		if strings.HasSuffix(file, item) {
			return true
		}
	}

	return false
}
