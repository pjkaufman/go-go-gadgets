package images

import (
	"strings"

	"github.com/pjkaufman/go-go-gadgets/pkg/image"
)

var (
	width         int = 800
	quality       int = 40
	minimumKbSize     = 150
)

func CompressImage(filePath string, data []byte) ([]byte, error) {
	if !isCompressableImage(filePath) || len(data)/1024 <= minimumKbSize {
		return data, nil
	}

	var newData []byte

	var isPng = strings.HasSuffix(filePath, ".png")
	if isPng {
		newData, err = image.PngRemoveExifData(data)
	} else {
		newData, err = image.JpegRemoveExifData(data)
	}
	if err != nil {
		return nil, err
	}

	if isPng {
		newData, err = image.PngResize(newData, width)
	} else {
		newData, err = image.JpegResize(newData, width, &quality)
	}
	if err != nil {
		return nil, err
	}

	return newData, nil
}

func isCompressableImage(imagePath string) bool {
	for _, ext := range image.CompressableImageExts {
		if strings.HasSuffix(strings.ToLower(imagePath), ext) {
			return true
		}
	}

	return false
}
