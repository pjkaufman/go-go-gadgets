package cbz

import (
	"archive/zip"
	"errors"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	filesize "github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/file-size"
	"github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/images"
	ziphandler "github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/zip-handler"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/image"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	dir     string
	verbose bool
)

const (
	DirArgEmpty = "directory must have a non-whitespace value"
)

// compressCmd represents the compress command
var compressCmd = &cobra.Command{
	Use:   "compress",
	Short: "Compresses all of the png and jpeg files in the cbz files in the specified directory",
	Example: heredoc.Doc(`To compress images in all cbzs in a folder:
	ebook-lint cbz compress -d folder
	
	To compress images in all cbzs in the current directory:
	ebook-lint cbz compress
	`),
	Run: func(cmd *cobra.Command, args []string) {
		err := ValidateCompressFlags(dir)
		if err != nil {
			logger.WriteError(err.Error())
		}

		logger.WriteInfo("Started compressing all cbzs\n")

		cbzs, err := filehandler.GetAllFilesWithExtInASpecificFolder(dir, ".cbz")
		if err != nil {
			logger.WriteError(err.Error())
		}

		var totalBeforeFileSize, totalAfterFileSize float64
		for _, cbz := range cbzs {
			logger.WriteInfof("starting cbz compression for %s...\n", cbz)

			err = compressCbz(dir, cbz)
			if err != nil {
				logger.WriteError(err.Error())
			}

			var originalFile = cbz + ".original"

			newKbSize, err := filehandler.GetFileSize(filehandler.JoinPath(dir, cbz))
			if err != nil {
				logger.WriteError(err.Error())
			}

			oldKbSize, err := filehandler.GetFileSize(filehandler.JoinPath(dir, originalFile))
			if err != nil {
				logger.WriteError(err.Error())
			}

			logger.WriteInfo(filesize.FileSizeSummary(originalFile, cbz, oldKbSize, newKbSize))

			totalBeforeFileSize += oldKbSize
			totalAfterFileSize += newKbSize
		}

		logger.WriteInfo(filesize.FilesSizeSummary(totalBeforeFileSize, totalAfterFileSize))

		logger.WriteInfo("Finished compressing all cbzs")
	},
}

func init() {
	CbzCmd.AddCommand(compressCmd)

	compressCmd.Flags().StringVarP(&dir, "directory", "d", ".", "the location to run the cbz image compression in")
	compressCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "whether or not to show extra information about the image compression")
}

func compressCbz(lintDir, cbz string) error {
	var src = filehandler.JoinPath(lintDir, cbz)

	err := ziphandler.UpdateZip(src, func(zipFiles map[string]*zip.File, w *zip.Writer) ([]string, error) {
		var (
			handledFiles, imagePaths []string
		)
		for filePath := range zipFiles {
			if !fileHasOneOfExts(filePath, image.CompressableImageExts) {
				continue
			}

			imagePaths = append(imagePaths, filePath)
		}

		var numFiles = len(imagePaths)
		for i, imagePath := range imagePaths {
			if verbose {
				logger.WriteInfof(`%d of %d: compressing %q`, i, numFiles, imagePath)
			}

			imageFile := zipFiles[imagePath]

			data, err := filehandler.ReadInZipFileBytes(imageFile)
			if err != nil {
				return nil, err
			}

			newData, err := images.CompressImage(imagePath, data)
			if err != nil {
				return nil, err
			}

			err = filehandler.WriteZipCompressedBytes(w, imagePath, newData)
			if err != nil {
				return nil, err
			}

			handledFiles = append(handledFiles, imagePath)
		}

		return handledFiles, nil
	})
	if err != nil {
		return fmt.Errorf("failed to compress cbz images for %q: %s", cbz, err)
	}

	return nil
}

func ValidateCompressFlags(dir string) error {
	if strings.TrimSpace(dir) == "" {
		return errors.New(DirArgEmpty)
	}

	return nil
}

func fileHasOneOfExts(fileName string, exts []string) bool {
	for _, ext := range exts {
		if strings.HasSuffix(fileName, ext) {
			return true
		}
	}

	return false
}
