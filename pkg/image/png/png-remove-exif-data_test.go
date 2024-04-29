//go:build unit

package png_test

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"testing"

	"github.com/dsoprea/go-exif/v2"
	pngstructure "github.com/dsoprea/go-png-image-structure"
	"github.com/pjkaufman/go-go-gadgets/pkg/image/png"
	"github.com/stretchr/testify/assert"
)

//go:embed test-data/*.png
var pngs embed.FS

func TestPngExifDataRemoval(t *testing.T) {
	fs.WalkDir(pngs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatalf("failed to walk path: %s\n", err)
		}

		if !d.IsDir() {
			pngFile, err := pngs.ReadFile(path)
			if err != nil {
				log.Fatalf("failed to read file \"%s\": %s\n", path, err)
			}

			t.Run(fmt.Sprintf(`"%s": exif data gets removed`, path), func(t *testing.T) {
				existingTags := getExifData(pngFile)
				newData, err := png.PngRemoveExifData(pngFile)
				assert.Nil(t, err)

				if len(existingTags) != 0 {
					assert.NotEqual(t, pngFile, newData)
				} else {
					assert.Equal(t, pngFile, newData)
				}

				// validate that exif data was removed
				newExifData := getExifData(newData)

				assert.Nil(t, newExifData)
			})
		}

		return nil
	})
}

func getExifData(data []byte) []byte {
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
