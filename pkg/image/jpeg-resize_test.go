//go:build unit

package image_test

import (
	_ "embed"
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
	"Resizing a JPEG to be smaller when its width is larger should work": {
		InputFileData:  canonTagsJpeg,
		OriginalHeight: 1200,
		OriginalWidth:  1600,
		NewHeight:      600,
		NewWidth:       800,
	},
	"Resizing a JPEG to be larger when its width is larger should work": {
		InputFileData:  canonTagsJpeg,
		OriginalHeight: 1200,
		OriginalWidth:  1600,
		NewHeight:      2400,
		NewWidth:       3200,
	},
	"Resizing a JPEG to be smaller when its width is larger should work with quality specified": {
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
			height, width := image_pkg.GetImageDimensions(test.InputFileData)
			assert.Equal(t, test.OriginalHeight, height, "original height was not the expected value")
			assert.Equal(t, test.OriginalWidth, width, "original width was not the expected value")

			newData, err := image_pkg.JpegResize(test.InputFileData, test.NewWidth, test.DesiredQuality)
			assert.Nil(t, err)

			height, width = image_pkg.GetImageDimensions(newData)
			assert.Equal(t, test.NewHeight, height, "height was not the expected value")
			assert.Equal(t, test.NewWidth, width, "width was not the expected value")
		})
	}
}
