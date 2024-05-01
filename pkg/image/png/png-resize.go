package png

import (
	"bytes"
	"fmt"
	"image"
	"image/png"

	"golang.org/x/image/draw"
)

// originally based on https://roeber.dev/posts/resize-an-image-in-go/
func PngResize(data []byte, width, height int) ([]byte, error) {
	src, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf(`failed to decode png: %w`, err)
	}

	// ratio := (float64)(src.Bounds().Max.Y) / (float64)(src.Bounds().Max.X)
	// height := int(math.Round(float64(width) * ratio))

	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	buf := new(bytes.Buffer)
	err = png.Encode(buf, dst)
	if err != nil {
		return nil, fmt.Errorf(`failed to png encode image: %w`, err)
	}

	return buf.Bytes(), nil
}
