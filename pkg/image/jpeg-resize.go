package image

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"math"

	"golang.org/x/image/draw"
)

// originally based on https://roeber.dev/posts/resize-an-image-in-go/
func JpegResize(data []byte, width int, quality *int) ([]byte, error) {
	src, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf(`failed to decode jpeg: %w`, err)
	}

	var bounds = src.Bounds()
	ratio := (float64)(bounds.Max.Y) / (float64)(bounds.Max.X)
	height := int(math.Round(float64(width) * ratio))

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, bounds, draw.Over, nil)

	var jpegOptions *jpeg.Options
	if quality != nil {
		jpegOptions = &jpeg.Options{
			Quality: *quality,
		}
	}
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, dst, jpegOptions)
	if err != nil {
		return nil, fmt.Errorf(`failed to jpeg encode image: %w`, err)
	}

	return buf.Bytes(), nil
}
