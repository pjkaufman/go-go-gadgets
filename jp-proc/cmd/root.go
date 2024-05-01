package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/image/jpeg"
	"github.com/pjkaufman/go-go-gadgets/pkg/image/png"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

const defaultQuality = 75

var (
	quiet, removeExif, overwrite bool
	quality, height, width       int
	file                         string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jpproc",
	Short: "JPEG and PNG image processor",
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

		var newData = data
		if isPng {
			newData, err = png.PngRemoveExifData(data)
		} else {
			newData, err = jpeg.JpegRemoveExifData(data)
		}
		if err != nil {
			logger.WriteError(err.Error())
		}

		var resizeImage = height != 0 || width != 0
		if !quiet && resizeImage {
			logger.WriteInfo(fmt.Sprintf("resizing image to %dx%d for %s", width, height, file))
		}

		if isPng {
			newData, err = png.PngResize(newData, width, height)
		} else {
			newData, err = jpeg.JpegResize(newData, width, height, &quality)
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&file, "file", "f", "", "the image file to operate on")
	err := rootCmd.MarkFlagRequired("file")
	if err != nil {
		logger.WriteError(fmt.Sprintf(`failed to mark flag "file" as required on root command: %v`, err))
	}

	err = rootCmd.MarkFlagFilename("file", "png", "jpg", "jpeg")
	if err != nil {
		logger.WriteError(fmt.Sprintf(`failed to mark flag "file" as looking for specific file types on root command: %v`, err))
	}

	rootCmd.Flags().BoolVarP(&removeExif, "remove-exif", "e", false, "whether or not to remove exif data from the image")
	rootCmd.Flags().BoolVarP(&quiet, "mute", "m", false, "whether or not to keep from printing out values to standard out")
	rootCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "whether or not to overwrite the original file when done")

	rootCmd.Flags().IntVarP(&quality, "quality", "q", defaultQuality, "the quality of the jpeg to use when encoding the image (default is 75)")
	rootCmd.Flags().IntVarP(&height, "height", "H", 0, "the height of the image to use when the image is resized (leave blank to keep original)")
	rootCmd.Flags().IntVarP(&width, "width", "W", 0, "the width of the image to use when the image is resized (leave blank to keep original)")
}

func ValidateJpProcFlags(file string) error {
	if strings.TrimSpace(file) == "" {
		return errors.New(`file cannot be empty`)
	}

	if !strings.HasSuffix(file, ".jpg") && !strings.HasSuffix(file, ".jpeg") && !strings.HasSuffix(file, ",png") {
		return errors.New(`file must be a jpeg or png file`)
	}

	return nil
}
