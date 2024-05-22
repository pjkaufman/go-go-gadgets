package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/image"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

const defaultQuality = 75

var (
	quiet, removeExif, overwrite bool
	quality, width               int
	file                         string
)

// procCmd represents the process command for processing an image
var procCmd = &cobra.Command{
	Use:   "proc",
	Short: "Processes the provided image in the specified ways",
	Run: func(cmd *cobra.Command, args []string) {
		err := ValidateJpProcFlags(file)
		if err != nil {
			logger.WriteError(err.Error())
		}

		filehandler.FileMustExist(file, "file")

		data := filehandler.ReadInBinaryFileContents(file)

		var isPng = strings.HasSuffix(file, ".png")

		// remove exif data if specified
		if !quiet && removeExif {
			logger.WriteInfo(fmt.Sprintf("removing exif data for %s", file))
		}

		var newData []byte
		if isPng {
			newData, err = image.PngRemoveExifData(data)
		} else {
			newData, err = image.JpegRemoveExifData(data)
		}
		if err != nil {
			logger.WriteError(err.Error())
		}

		var resizeImage = width != 0
		if !quiet && resizeImage {
			logger.WriteInfo(fmt.Sprintf("resizing image to width %d for %s", width, file))
		}

		if isPng {
			newData, err = image.PngResize(newData, width)
		} else {
			newData, err = image.JpegResize(newData, width, &quality)
		}
		if err != nil {
			logger.WriteError(err.Error())
		}

		if !bytes.Equal(data, newData) {
			if overwrite {
				filehandler.WriteBinaryFileContents(file, newData)
			} else {
				var newFile = strings.Split(file, ".")
				var ext = newFile[len(newFile)-1]
				newFile[len(newFile)-1] = "test"
				newFile = append(newFile, ext)

				filehandler.WriteBinaryFileContents(strings.Join(newFile, "."), newData)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(procCmd)

	procCmd.Flags().StringVarP(&file, "file", "f", "", "the image file to operate on")
	err := procCmd.MarkFlagRequired("file")
	if err != nil {
		logger.WriteError(fmt.Sprintf(`failed to mark flag "file" as required on root command: %v`, err))
	}

	err = procCmd.MarkFlagFilename("file", "png", "jpg", "jpeg")
	if err != nil {
		logger.WriteError(fmt.Sprintf(`failed to mark flag "file" as looking for specific file types on root command: %v`, err))
	}

	procCmd.Flags().BoolVarP(&removeExif, "remove-exif", "e", false, "whether or not to remove exif data from the image")
	procCmd.Flags().BoolVarP(&quiet, "mute", "m", false, "whether or not to keep from printing out values to standard out")
	procCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "whether or not to overwrite the original file when done")

	procCmd.Flags().IntVarP(&quality, "quality", "q", defaultQuality, "the quality of the jpeg to use when encoding the image (default is 75)")
	procCmd.Flags().IntVarP(&width, "width", "w", 0, "the width of the image to use when the image is resized (leave blank to keep original)")
}

func ValidateJpProcFlags(file string) error {
	if strings.TrimSpace(file) == "" {
		return errors.New(`file cannot be empty`)
	}

	if !strings.HasSuffix(file, ".jpg") && !strings.HasSuffix(file, ".jpeg") && !strings.HasSuffix(file, ".png") {
		return errors.New(`file must be a jpeg or png file`)
	}

	return nil
}
