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

	isPng := strings.HasSuffix(filePath, ".png")
	_, width := image.GetImageDimensions(data)

	// Skip resize if already smaller than desired for PNGs
	widthToUse := desiredWidth
	if width <= desiredWidth {
		if isPng {
			return data, nil
		}

		widthToUse = width
	}
	if isPng {
		return image.PngResize(data, widthToUse)
	} else {
		return image.JpegResize(data, widthToUse, &quality)
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
