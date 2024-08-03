//go:build unit

package cbz_test

import (
	"log"
	"os"
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/cmd/cbz"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/stretchr/testify/assert"
)

type CompressCbzTestCase struct {
	Filename string
}

const (
	originalFileDir   = "testdata/original"
	compressedFileDir = "testdata/compressed"
)

var compressCbzTestCases = map[string]CompressCbzTestCase{
	"Compressing a cbz file should work consistently": {
		Filename: "Whiz_Comics_Volume_1-1948-Fawcett_Publications.cbz",
	},
}

func TestCompressCbz(t *testing.T) {
	for name, test := range compressCbzTestCases {
		t.Run(name, func(t *testing.T) {
			err := cbz.CompressCbz(originalFileDir, test.Filename)
			assert.Nil(t, err)

			equalityStatus, issue := filehandler.ZipsAreEqual(test.Filename, originalFileDir, compressedFileDir, false)
			assert.True(t, equalityStatus, issue)

			var originalCbzPath = originalFileDir + string(os.PathSeparator) + test.Filename
			err = os.RemoveAll(originalCbzPath)
			if err != nil {
				log.Fatalf("failed to remove the result of compressing the cbz %q: %s", originalCbzPath, err)
			}

			err = os.Rename(originalCbzPath+".original", originalCbzPath)
			if err != nil {
				log.Fatalf("failed move original file back to its starting location for %q: %s", test.Filename, err)
			}
		})
	}
}

func BenchmarkCompressCbz(b *testing.B) {
	var filename = "Whiz-Comics_Shazam_Golden_Arow_Volume_1-Fawcett_Publications.cbr"

	for n := 0; n < b.N; n++ {
		err := cbz.CompressCbz(originalFileDir, filename)
		if err != nil {
			log.Fatalf("failed to compress cbz %q: %s", filename, err)
		}

		var originalCbzPath = originalFileDir + string(os.PathSeparator) + filename
		err = os.RemoveAll(originalCbzPath)
		if err != nil {
			log.Fatalf("failed to remove the result of compressing the cbz %q: %s", originalCbzPath, err)
		}

		err = os.Rename(originalCbzPath+".original", originalCbzPath)
		if err != nil {
			log.Fatalf("failed move original file back to its starting location for %q: %s", filename, err)
		}
	}
}
