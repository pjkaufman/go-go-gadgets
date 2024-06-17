//go:build unit

package epub_test

import (
	"archive/zip"
	"bytes"
	"io"
	"log"
	"os"
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/cmd/epub"
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
		Filename:       "jules-verne_from-the-earth-to-the-moon_ward-lock-co.epub",
		CompressImages: false,
	},
}

func TestLintEpub(t *testing.T) {
	for name, test := range LintEpubTestCases {
		t.Run(name, func(t *testing.T) {
			epub.LintEpub(originalFileDir, test.Filename, test.CompressImages)

			assert.True(t, epubsAreEqual(test.Filename))

			var originalEpubPath = originalFileDir + string(os.PathSeparator) + test.Filename
			err := os.RemoveAll(originalEpubPath)
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

// epubsAreEqual runs after the operation of LintEpub which leads to the linted file taking the place of the original.
// This means that we are able to assume that the original file's location should have data comparable to that found
// in the linted file.
func epubsAreEqual(filename string) bool {
	var originalEpubPath = originalFileDir + string(os.PathSeparator) + filename
	lintedEpub, err := zip.OpenReader(originalEpubPath)
	if err != nil {
		log.Fatalf("Failed to open zip file %q: %s", originalEpubPath, err)
	}
	defer lintedEpub.Close()

	var lintedEpubPath = lintedFileDir + string(os.PathSeparator) + filename
	expectedEpub, err := zip.OpenReader(lintedEpubPath)
	if err != nil {
		log.Fatalf("Failed to open zip file %q: %s", lintedEpubPath, err)
	}
	defer expectedEpub.Close()

	if len(lintedEpub.File) != len(expectedEpub.File) {
		return false
	}

	for _, zipFile := range lintedEpub.File {
		var found bool
		for _, expectedZipFile := range expectedEpub.File {
			if zipFile.Name == expectedZipFile.Name {
				if !zipFilesAreEqual(zipFile, expectedZipFile) {
					return false
				}

				found = true
				break
			}
		}

		if found {
			continue
		}

		return false
	}

	return true
}

func zipFilesAreEqual(actual, expected *zip.File) bool {
	if actual.Method != expected.Method || actual.CompressedSize64 != expected.CompressedSize64 || actual.UncompressedSize64 != expected.UncompressedSize64 {
		return false
	}

	actualReader, err := actual.Open()
	if err != nil {
		log.Fatalf("failed to open actual zip contents for %q: %s", actual.Name, err)
	}

	defer actualReader.Close()

	var actualContents = &bytes.Buffer{}
	_, err = io.Copy(actualContents, actualReader)
	if err != nil {
		log.Fatalf("failed to read in actual zip contents for %q: %s", actual.Name, err)
	}

	expectedReader, err := actual.Open()
	if err != nil {
		log.Fatalf("failed to open expected zip contents for %q: %s", expected.Name, err)
	}

	defer expectedReader.Close()

	var expectedContents = &bytes.Buffer{}
	_, err = io.Copy(expectedContents, expectedReader)
	if err != nil {
		log.Fatalf("failed to read in expected zip contents for %q: %s", expected.Name, err)
	}

	return bytes.Equal(expectedContents.Bytes(), actualContents.Bytes())
}
