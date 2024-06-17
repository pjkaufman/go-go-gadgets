package images

import (
	"bytes"
	"strings"

	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/image"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

var (
	width   int = 800
	quality int = 40
)

func CompressImage(filePath string) {
	if !isCompressableImage(filePath) || filehandler.MustGetFileSize(filePath) <= 150 {
		return
	}

	var newData []byte
	data, err := filehandler.ReadInBinaryFileContents(filePath)
	if err != nil {
		logger.WriteError(err.Error())
	}

	var isPng = strings.HasSuffix(filePath, ".png")
	if isPng {
		newData, err = image.PngRemoveExifData(data)
	} else {
		newData, err = image.JpegRemoveExifData(data)
	}
	if err != nil {
		logger.WriteError(err.Error())
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
		filehandler.WriteBinaryFileContents(filePath, newData)
	}
}

func isCompressableImage(imagePath string) bool {
	for _, ext := range image.CompressableImageExts {
		if strings.HasSuffix(strings.ToLower(imagePath), ext) {
			return true
		}
	}

	return false
}
