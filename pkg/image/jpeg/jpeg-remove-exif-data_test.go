package jpeg_test

import (
	_ "embed"
	"testing"

	"github.com/pjkaufman/go-go-gadgets/pkg/image/jpeg"

	"github.com/stretchr/testify/assert"
)

// go:embed add back by removing space after //
var imageJpeg []byte

func TestGetNextTableAndItsEndPosition(t *testing.T) {
	t.Run("exif data gets removed", func(t *testing.T) {
		newData, err := jpeg.JpegRemoveExifData(imageJpeg)
		assert.Nil(t, err)

		assert.NotEqual(t, imageJpeg, newData)
		// assert.Fail(t, "fail...")
	})
}
