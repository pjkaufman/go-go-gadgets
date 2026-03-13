package image

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"sync"

	"golang.org/x/image/draw"
)

var bufPool = sync.Pool{
	New: func() any { return new(bytes.Buffer) },
}

var rgbaPool = sync.Pool{
	New: func() any { return image.NewRGBA(image.Rect(0, 0, 0, 0)) },
}

// originally based on https://roeber.dev/posts/resize-an-image-in-go/
func JpegResize(data []byte, width int, quality *int) ([]byte, error) {
	src, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode jpeg: %w", err)
	}

	bounds := src.Bounds()
	height := bounds.Max.Y * width / bounds.Max.X

	dst := rgbaPool.Get().(*image.RGBA)
	if dst.Bounds().Dx() != width || dst.Bounds().Dy() != height {
		*dst = *image.NewRGBA(image.Rect(0, 0, width, height))
	}
	defer rgbaPool.Put(dst)

	draw.NearestNeighbor.Scale(dst, dst.Rect, src, bounds, draw.Over, nil)

	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)

	var jpegOptions *jpeg.Options
	if quality != nil {
		jpegOptions = &jpeg.Options{
			Quality: *quality,
		}
	}

	if err := jpeg.Encode(buf, dst, jpegOptions); err != nil {
		return nil, fmt.Errorf("failed to jpeg encode image: %w", err)
	}

	out := make([]byte, buf.Len())
	copy(out, buf.Bytes())
	return out, nil
}
