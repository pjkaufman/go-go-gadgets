//go:build unit

package image_test

import (
	"embed"
	"log"
	"testing"

	"github.com/dsoprea/go-exif/v2"
	pngstructure "github.com/dsoprea/go-png-image-structure"
	"github.com/pjkaufman/go-go-gadgets/pkg/image"
)

//go:embed testdata/png/*.png
var pngs embed.FS

func TestPngExifDataRemoval(t *testing.T) {
	iterateAndTestImageExifDataRemoval(t, pngs, getPngExifData, image.PngRemoveExifData)
}

func getPngExifData(data []byte) []byte {
	pmp := pngstructure.NewPngMediaParser()

	intfc, parseErr := pmp.ParseBytes(data)
	if parseErr != nil && intfc == nil {
		log.Fatalf("failed to parse out png bytes: %s\n", parseErr)
	}

	_, et, err := intfc.Exif()
	if err != nil {
		if err.Error() == pngstructure.ErrNoExif.Error() || err.Error() == exif.ErrNoExif.Error() {
			return nil
		}

		log.Fatalf("failed to dump exif info: %s\n", err)
	}

	return et
}
