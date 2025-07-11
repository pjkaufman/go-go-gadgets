//go:build unit

package cmd_test

import (
	"log"
	"os"
	"testing"

	epub "github.com/pjkaufman/go-go-gadgets/epub-lint/cmd"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/stretchr/testify/assert"
)

type lintEpubTestCase struct {
	filename                string
	compressImages, verbose bool
	removableFileExts       []string
}

const (
	originalFileDir = "testdata/original"
	lintedFileDir   = "testdata/linted"
)

var LintEpubTestCases = map[string]lintEpubTestCase{
	"Linting a file with image compression should work consistently": {
		filename:       "jules-verne_from-the-earth-to-the-moon_ward-lock-co.epub",
		compressImages: true,
	},
	"Linting a file without image compression should work consistently and not affect the images": {
		filename:       "jules-verne_in-search-of-the-castaways_j-b-lippincott-co.epub",
		compressImages: false,
	},
	"Linting a file without a mimetype should have the mimetype added": {
		filename:       "jules-verne_in-search-of-the-castaways_j-b-lippincott-co-missing_mimetype.epub",
		compressImages: false,
	},
	"Linting a file with an extra text file should have it removed": {
		filename:          "jules-verne_from-the-earth-to-the-moon_ward-lock-co-extra_txt.epub",
		compressImages:    true,
		verbose:           true,
		removableFileExts: []string{".txt"},
	},
	"Linting a file with an extra text file with no removable file exts should not have it removed": {
		filename:          "jules-verne_from-the-earth-to-the-moon_ward-lock-co-extra_txt_no_change.epub",
		compressImages:    true,
		verbose:           true,
		removableFileExts: []string{},
	},
}

func TestLintEpub(t *testing.T) {
	for name, test := range LintEpubTestCases {
		t.Run(name, func(t *testing.T) {
			err := epub.LintEpub(originalFileDir, test.filename, test.compressImages, test.verbose, test.removableFileExts)
			assert.Nil(t, err)

			// This runs after the operation of LintEpub which leads to the linted file taking the place of the original.
			// This means that we are able to assume that the original file's location should have data comparable to that found
			// in the linted file.
			equalityStatus, issue := filehandler.ZipsAreEqual(test.filename, originalFileDir, lintedFileDir, true)
			assert.True(t, equalityStatus, issue)

			var originalEpubPath = originalFileDir + string(os.PathSeparator) + test.filename
			err = os.RemoveAll(originalEpubPath)
			if err != nil {
				log.Fatalf("failed to remove the result of lint epub %q: %s", originalEpubPath, err)
			}

			err = os.Rename(originalEpubPath+".original", originalEpubPath)
			if err != nil {
				log.Fatalf("failed move original file back to its starting location for %q: %s", test.filename, err)
			}
		})
	}
}

func BenchmarkLintEpub(b *testing.B) {
	var (
		filename                = "jules-verne_from-the-earth-to-the-moon_ward-lock-co.epub"
		compressImages, verbose = true, false
	)

	for n := 0; n < b.N; n++ {
		var originalEpubPath = originalFileDir + string(os.PathSeparator) + filename
		err := epub.LintEpub(originalFileDir, filename, compressImages, verbose, []string{})
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
