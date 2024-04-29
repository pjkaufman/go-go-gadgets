//go:build unit

package jpeg_test

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/jpeg"
	"log"
	"testing"

	"github.com/pjkaufman/go-go-gadgets/pkg/image/jpeg"
	"github.com/stretchr/testify/assert"
)

//go:embed test-data/22-canon_tags.jpg
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
		NewHeight:      800,
		NewWidth:       800,
	},
	"Resizing an image to be smaller should work with quality specified": {
		InputFileData:  canonTagsJpeg,
		OriginalHeight: 1200,
		OriginalWidth:  1600,
		NewHeight:      800,
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

			newData, err := jpeg.JpegResize(test.InputFileData, test.NewHeight, test.NewWidth, test.DesiredQuality)
			assert.Nil(t, err)

			height, width = getHeightAndWidth(newData)
			assert.Equal(t, test.NewHeight, height, "height was not the expected value")
			assert.Equal(t, test.NewWidth, width, "width was not the expected value")
		})
	}
	// fs.WalkDir(jpegs, ".", func(path string, d fs.DirEntry, err error) error {
	// 	if err != nil {
	// 		log.Fatalf("failed to walk path: %s\n", err)
	// 	}

	// 	if !d.IsDir() {
	// 		jpegFile, err := jpegs.ReadFile(path)
	// 		if err != nil {
	// 			log.Fatalf("failed to read file \"%s\": %s\n", path, err)
	// 		}

	// t.Run(fmt.Sprintf(`"%s": resize`, path), func(t *testing.T) {
	// 	newData, err := jpeg.JpegResize(jpegFile, 800, 800)
	// 	assert.Nil(t, err)

	// 	height, width := getHeightAndWidth(newData)
	// 	assert.Equal(t, 800, height, "height was not the expected value")
	// 	assert.Equal(t, 800, width, "width was not the expected value")
	// })
	// 	}

	// 	return nil
	// })
}

func getHeightAndWidth(data []byte) (int, int) {
	im, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		log.Fatalf("failed to decode image to jpeg to get dimensions: %s\n", err)
	}
	return im.Height, im.Width
}
