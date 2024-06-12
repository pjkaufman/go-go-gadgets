package cbz

import (
	"errors"
	"strings"

	"github.com/MakeNowJust/heredoc"
	filesize "github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/file-size"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
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

		cbzs, err := filehandler.MustGetAllFilesWithExtInASpecificFolder(dir, ".cbz")
		if err != nil {
			logger.WriteError(err.Error())
		}

		var totalBeforeFileSize, totalAfterFileSize float64
		for _, cbz := range cbzs {
			logger.WriteInfof("starting cbz compression for %s...\n", cbz)

			compressCbz(dir, cbz)

			var originalFile = cbz + ".original"

			newKbSize, err := filehandler.MustGetFileSize(filehandler.JoinPath(dir, cbz))
			if err != nil {
				logger.WriteError(err.Error())
			}

			oldKbSize, err := filehandler.MustGetFileSize(filehandler.JoinPath(dir, originalFile))
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

func compressCbz(lintDir, cbz string) {
	// var src = filehandler.JoinPath(lintDir, cbz)
	// var dest = filehandler.JoinPath(lintDir, "cbz")

	// filehandler.UnzipRunOperationAndRezip(src, dest, func() {
	// 	var imageFiles = filehandler.MustGetAllFilesWithExtsInASpecificFolderAndSubFolders(dest, image.CompressableImageExts...)

	// 	for i, imageFile := range imageFiles {
	// 		if verbose {
	// 			logger.WriteInfo(fmt.Sprintf(`%d of %d: compressing %q`, i, len(imageFiles), imageFile))
	// 		}

	// 		images.CompressImage(imageFile)
	// 	}
	// })
}

func ValidateCompressFlags(dir string) error {
	if strings.TrimSpace(dir) == "" {
		return errors.New(DirArgEmpty)
	}

	return nil
}
