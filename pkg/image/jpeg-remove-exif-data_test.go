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
	"github.com/stretchr/testify/require"
)

//go:embed testdata/jpeg/*.jpg
var jpegs embed.FS

func TestJpegExifDataRemoval(t *testing.T) {
	iterateAndTestImageExifDataRemoval(t, jpegs, getJpegExifData, image.JpegRemoveExifData)
}

func iterateAndTestImageExifDataRemoval(t *testing.T, files embed.FS, getExifData func([]byte) []byte, removeExifData func([]byte) ([]byte, error)) {
	_ = fs.WalkDir(files, ".", func(path string, d fs.DirEntry, err error) error {
		require.NoError(t, err, "failed to walk path")

		if !d.IsDir() {
			imageFile, err := files.ReadFile(path)

			require.NoErrorf(t, err, "failed to read file %q", path)

			t.Run(fmt.Sprintf(`%q: exif data gets removed`, path), func(t *testing.T) {
				existingTags := getExifData(imageFile)
				newData, err := removeExifData(imageFile)
				require.NoError(t, err)

				if len(existingTags) != 0 {
					assert.NotEqual(t, imageFile, newData)
				} else {
					assert.Equal(t, imageFile, newData)
				}

				// validate that exif data was removed
				newExifData := getExifData(newData)

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
