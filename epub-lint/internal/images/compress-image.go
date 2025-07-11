package images

import (
	"strings"

	"github.com/pjkaufman/go-go-gadgets/pkg/image"
)

var (
	desiredWidth  int = 800
	quality       int = 40
	minimumKbSize     = 150
)

func CompressImage(filePath string, data []byte) ([]byte, error) {
	if !isCompressableImage(filePath) || len(data)/1024 <= minimumKbSize {
		return data, nil
	}

	var (
		newData []byte
		isPng   = strings.HasSuffix(filePath, ".png")
		err     error
	)
	if isPng {
		newData, err = image.PngRemoveExifData(data)
	} else {
		newData, err = image.JpegRemoveExifData(data)
	}
	if err != nil {
		return nil, err
	}

	var widthToUse = desiredWidth
	_, width := image.GetImageDimensions(newData)
	if width < desiredWidth {
		desiredWidth = width
	}

	if isPng && widthToUse == desiredWidth {
		newData, err = image.PngResize(newData, widthToUse)
	} else if !isPng {
		newData, err = image.JpegResize(newData, widthToUse, &quality)
	}

	return newData, err
}

func isCompressableImage(imagePath string) bool {
	for _, ext := range image.CompressableImageExts {
		if strings.HasSuffix(strings.ToLower(imagePath), ext) {
			return true
		}
	}

	return false
}
