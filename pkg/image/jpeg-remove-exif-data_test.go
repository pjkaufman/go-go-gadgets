//go:build unit

package image_test

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"testing"

	"github.com/dsoprea/go-exif/v2"
	jpegstructure "github.com/dsoprea/go-jpeg-image-structure"
	"github.com/pjkaufman/go-go-gadgets/pkg/image"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/jpeg/*.jpg
var jpegs embed.FS

func TestJpegExifDataRemoval(t *testing.T) {
	fs.WalkDir(jpegs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatalf("failed to walk path: %s\n", err)
		}

		if !d.IsDir() {
			jpegFile, err := jpegs.ReadFile(path)
			if err != nil {
				log.Fatalf("failed to read file \"%s\": %s\n", path, err)
			}

			t.Run(fmt.Sprintf(`"%s": exif data gets removed`, path), func(t *testing.T) {
				existingTags := getJpegExifData(jpegFile)
				newData, err := image.JpegRemoveExifData(jpegFile)
				assert.Nil(t, err)

				if len(existingTags) != 0 {
					assert.NotEqual(t, jpegFile, newData)
				} else {
					assert.Equal(t, jpegFile, newData)
				}

				// validate that exif data was removed
				newExifData := getJpegExifData(newData)

				assert.Nil(t, newExifData)
			})
		}

		return nil
	})
}

func getJpegExifData(data []byte) []byte {
	jmp := jpegstructure.NewJpegMediaParser()

	intfc, parseErr := jmp.ParseBytes(data)
	if parseErr != nil && intfc == nil {
		log.Fatalf("failed to parse out jpeg bytes: %s\n", parseErr)
	}

	_, et, err := intfc.Exif()
	if err != nil {
		// for some reasons errors.Is did not work here so we will compare the error text instead
		if err.Error() == exif.ErrNoExif.Error() {
			return nil
		}

		log.Fatalf("failed to dump exif info: %s\n", err)
	}

	return et
}
