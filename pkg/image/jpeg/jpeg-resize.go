package jpeg

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"math"

	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"golang.org/x/image/draw"
)

// originally based on https://roeber.dev/posts/resize-an-image-in-go/
func JpegResize(data []byte, width int, quality *int) ([]byte, error) {
	src, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf(`failed to decode jpeg: %w`, err)
	}

	ratio := (float64)(src.Bounds().Max.Y) / (float64)(src.Bounds().Max.X)
	height := int(math.Round(float64(width) * ratio))

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	logger.WriteInfo(fmt.Sprintf("Height ratio is %d", height))
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

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
