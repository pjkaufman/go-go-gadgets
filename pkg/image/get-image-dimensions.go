package image

import (
	"bytes"
	"image"
	"log"

	_ "image/jpeg"
	_ "image/png"
)

func GetImageDimensions(data []byte) (int, int) {
	im, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		log.Fatalf("failed to decode image to get dimensions: %s\n", err)
	}
	return im.Height, im.Width
}
