//go:build unit

package epub_test

import (
	"log"
	"os"
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/cmd/epub"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/stretchr/testify/assert"
)

type LintEpubTestCase struct {
	Filename       string
	CompressImages bool
}

const (
	originalFileDir = "testdata/original"
	lintedFileDir   = "testdata/linted"
)

var LintEpubTestCases = map[string]LintEpubTestCase{
	"Linting a file with image compression should work consistently": {
		Filename:       "jules-verne_from-the-earth-to-the-moon_ward-lock-co.epub",
		CompressImages: true,
	},
	"Linting a file without image compression should work consistently and not affect the images": {
		Filename:       "jules-verne_in-search-of-the-castaways_j-b-lippincott-co.epub",
		CompressImages: false,
	},
	"Linting a file without a mimetype should have the mimetype added": {
		Filename:       "jules-verne_in-search-of-the-castaways_j-b-lippincott-co-missing_mimetype.epub",
		CompressImages: false,
	},
}

func TestLintEpub(t *testing.T) {
	for name, test := range LintEpubTestCases {
		t.Run(name, func(t *testing.T) {
			err := epub.LintEpub(originalFileDir, test.Filename, test.CompressImages)
			assert.Nil(t, err)

			// This runs after the operation of LintEpub which leads to the linted file taking the place of the original.
			// This means that we are able to assume that the original file's location should have data comparable to that found
			// in the linted file.
			equalityStatus, issue := filehandler.ZipsAreEqual(test.Filename, originalFileDir, lintedFileDir, true)
			assert.True(t, equalityStatus, issue)

			var originalEpubPath = originalFileDir + string(os.PathSeparator) + test.Filename
			err = os.RemoveAll(originalEpubPath)
			if err != nil {
				log.Fatalf("failed to remove the result of lint epub %q: %s", originalEpubPath, err)
			}

			err = os.Rename(originalEpubPath+".original", originalEpubPath)
			if err != nil {
				log.Fatalf("failed move original file back to its starting location for %q: %s", test.Filename, err)
			}
		})
	}
}

func BenchmarkLintEpub(b *testing.B) {
	var (
		filename       = "jules-verne_from-the-earth-to-the-moon_ward-lock-co.epub"
		compressImages = true
	)

	for n := 0; n < b.N; n++ {
		var originalEpubPath = originalFileDir + string(os.PathSeparator) + filename
		err := epub.LintEpub(originalFileDir, filename, compressImages)
		if err != nil {
			log.Fatalf("failed to lint epub %q: %s", originalEpubPath, err)
		}

		err = os.RemoveAll(originalEpubPath)
		if err != nil {
			log.Fatalf("failed to remove the result of lint epub %q: %s", originalEpubPath, err)
		}

		err = os.Rename(originalEpubPath+".original", originalEpubPath)
		if err != nil {
			log.Fatalf("failed move original file back to its starting location for %q: %s", filename, err)
		}
	}
}
