package cmd

import (
	"bytes"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/pkg/cli/flags"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/image"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

const defaultQuality = 75

var (
	quiet, removeExif, updateExisting bool
	quality, width                    int
	file                              string
	procFlags                         = flags.Flags{
		Flags: []flags.Flag{
			flags.NewFileFlag(true, false, &file, "file", "f", "", "the image file to operate on", []string{"png", "jpg", "jpeg"}, true),
			flags.NewBoolFlag(false, false, &removeExif, "remove-exif", "e", false, "whether or not to remove exif data from the image"),
			flags.NewBoolFlag(false, false, &quiet, "quiet", "q", false, "whether or not to keep from printing out values to standard out"),
			flags.NewBoolFlag(false, false, &updateExisting, "update", "u", false, "whether or not to update the original file when done"),
			flags.NewIntFlag(false, false, &quality, "quality", "", defaultQuality, "the quality of the jpeg to use when encoding the image (default is 75)"),
			flags.NewIntFlag(false, false, &width, "width", "w", 0, "the width of the image to use when the image is resized (leave blank to keep original)"),
		},
	}
)

// procCmd represents the process command for processing an image
var procCmd = &cobra.Command{
	Use:   "proc",
	Short: "Processes the provided image in the specified ways",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return procFlags.Validate()
	},
	Run: func(cmd *cobra.Command, args []string) {
		data, err := filehandler.ReadInBinaryFileContents(file)
		if err != nil {
			logger.WriteFatal(err.Error())
		}

		var isPng = strings.HasSuffix(file, ".png")

		// remove exif data if specified
		if !quiet && removeExif {
			logger.WriteInfof("removing exif data for %s\n", file)
		}

		var newData []byte
		if isPng {
			newData, err = image.PngRemoveExifData(data)
		} else {
			newData, err = image.JpegRemoveExifData(data)
		}
		if err != nil {
			logger.WriteFatal(err.Error())
		}

		var resizeImage = width != 0
		if !quiet && resizeImage {
			logger.WriteInfof("resizing image to width %d for %s\n", width, file)
		}

		if isPng {
			newData, err = image.PngResize(newData, width)
		} else {
			newData, err = image.JpegResize(newData, width, &quality)
		}
		if err != nil {
			logger.WriteFatal(err.Error())
		}

		if !bytes.Equal(data, newData) {
			if updateExisting {
				err = filehandler.WriteBinaryFileContents(file, newData)
			} else {
				var newFile = strings.Split(file, ".")
				var ext = newFile[len(newFile)-1]
				newFile[len(newFile)-1] = "test"
				newFile = append(newFile, ext)

				err = filehandler.WriteBinaryFileContents(strings.Join(newFile, "."), newData)
			}

			if err != nil {
				logger.WriteFatal(err.Error())
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(procCmd)

	err := procFlags.AddToCmd(procCmd)
	if err != nil {
		logger.WriteFatal(err.Error())
	}
}
