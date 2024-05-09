//go:build unit

package image_test

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/jpeg"
	"log"
	"testing"

	image_pkg "github.com/pjkaufman/go-go-gadgets/pkg/image"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/jpeg/22-canon_tags.jpg
var canonTagsJpeg []byte

type JpegResizeTestCase struct {
	InputFileData  []byte
	NewHeight      int
	NewWidth       int
	OriginalHeight int
	OriginalWidth  int
	DesiredQuality *int
}

var quality40 = 40

var JpegResizeTestCases = map[string]JpegResizeTestCase{
	"Resizing an image to be smaller should work": {
		InputFileData:  canonTagsJpeg,
		OriginalHeight: 1200,
		OriginalWidth:  1600,
		NewHeight:      600,
		NewWidth:       800,
	},
	"Resizing an image to be smaller should work with quality specified": {
		InputFileData:  canonTagsJpeg,
		OriginalHeight: 1200,
		OriginalWidth:  1600,
		NewHeight:      600,
		NewWidth:       800,
		DesiredQuality: &quality40,
	},
}

func TestJpegResize(t *testing.T) {
	for name, test := range JpegResizeTestCases {
		t.Run(name, func(t *testing.T) {
			height, width := getHeightAndWidth(test.InputFileData)
			assert.Equal(t, test.OriginalHeight, height, "original height was not the expected value")
			assert.Equal(t, test.OriginalWidth, width, "original width was not the expected value")

			newData, err := image_pkg.JpegResize(test.InputFileData, test.NewWidth, test.DesiredQuality)
			assert.Nil(t, err)

			height, width = getHeightAndWidth(newData)
			assert.Equal(t, test.NewHeight, height, "height was not the expected value")
			assert.Equal(t, test.NewWidth, width, "width was not the expected value")
		})
	}
}

func getHeightAndWidth(data []byte) (int, int) {
	im, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		log.Fatalf("failed to decode image to jpeg to get dimensions: %s\n", err)
	}
	return im.Height, im.Width
}
