//go:build unit

package image_test

import (
	_ "embed"
	"testing"

	image_pkg "github.com/pjkaufman/go-go-gadgets/pkg/image"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/png/exif.png
var exifPng []byte

//go:embed testdata/png/ml.png
var mlPng []byte

type PngResizeTestCase struct {
	InputFileData  []byte
	NewHeight      int
	NewWidth       int
	OriginalHeight int
	OriginalWidth  int
}

var PngResizeTestCases = map[string]PngResizeTestCase{
	"Resizing a PNG to be larger when it has a larger width should work": {
		InputFileData:  exifPng,
		OriginalHeight: 69,
		OriginalWidth:  91,
		NewHeight:      138,
		NewWidth:       182,
	},
	"Resizing a PNG to be smaller when it has a larger width should work should work": {
		InputFileData:  exifPng,
		OriginalHeight: 69,
		OriginalWidth:  91,
		NewHeight:      46,
		NewWidth:       61,
	},
	"Resizing a PNG to be larger when it has a larger height should work": {
		InputFileData:  mlPng,
		OriginalHeight: 380,
		OriginalWidth:  308,
		NewHeight:      760,
		NewWidth:       616,
	},
	"Resizing a PNG to be smaller when it has a larger height should work should work": {
		InputFileData:  mlPng,
		OriginalHeight: 380,
		OriginalWidth:  308,
		NewHeight:      190,
		NewWidth:       154,
	},
}

func TestPngResize(t *testing.T) {
	for name, test := range PngResizeTestCases {
		t.Run(name, func(t *testing.T) {
			height, width := image_pkg.GetImageDimensions(test.InputFileData)
			assert.Equal(t, test.OriginalHeight, height, "original height was not the expected value")
			assert.Equal(t, test.OriginalWidth, width, "original width was not the expected value")

			newData, err := image_pkg.PngResize(test.InputFileData, test.NewWidth)
			assert.Nil(t, err)

			height, width = image_pkg.GetImageDimensions(newData)
			assert.Equal(t, test.NewHeight, height, "height was not the expected value")
			assert.Equal(t, test.NewWidth, width, "width was not the expected value")
		})
	}
}
