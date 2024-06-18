package images

import (
	"bytes"
	"strings"

	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/image"
)

var (
	width   int = 800
	quality int = 40
)

func CompressImage(filePath string) error {
	if !isCompressableImage(filePath) {
		return nil
	}

	fileSize, err := filehandler.MustGetFileSize(filePath)
	if err != nil {
		return err
	}

	if fileSize < 150 {
		return nil
	}

	var newData []byte
	data, err := filehandler.ReadInBinaryFileContents(filePath)
	if err != nil {
		return err
	}

	var isPng = strings.HasSuffix(filePath, ".png")
	if isPng {
		newData, err = image.PngRemoveExifData(data)
	} else {
		newData, err = image.JpegRemoveExifData(data)
	}
	if err != nil {
		return err
	}

	if isPng {
		newData, err = image.PngResize(newData, width)
	} else {
		newData, err = image.JpegResize(newData, width, &quality)
	}
	if err != nil {
		return err
	}

	if !bytes.Equal(data, newData) {
		err = filehandler.WriteBinaryFileContents(filePath, newData)

		if err != nil {
			return err
		}
	}

	return nil
}

func isCompressableImage(imagePath string) bool {
	for _, ext := range image.CompressableImageExts {
		if strings.HasSuffix(strings.ToLower(imagePath), ext) {
			return true
		}
	}

	return false
}
