//go:build unit

package filehandler

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/pjkaufman/go-go-gadgets/pkg/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ZipsAreEqual compares two files with the same name in two separate directories to see if they are logically equivalent.
// This is a test helper function and is not meant for use in non-testing code.
func ZipsAreEqual(t *testing.T, filename, originalFileDir, expectedFileDir string) {
	var originalPath = JoinPath(originalFileDir, filename)
	actualZip, err := zip.OpenReader(originalPath)
	require.NoErrorf(t, err, "Failed to open zip file %q", originalPath)

	defer tests.MustClose(t, actualZip)

	var expectedFilePath = JoinPath(expectedFileDir, filename)
	expectedZip, err := zip.OpenReader(expectedFilePath)
	require.NoErrorf(t, err, "Failed to open zip file %q", expectedFilePath)

	defer tests.MustClose(t, expectedZip)

	assert.Lenf(t, actualZip.File, len(expectedZip.File), "expected %d files in zip, but got %d files", len(expectedZip.File), len(actualZip.File))
	assert.Equalf(t, "mimetype", actualZip.File[0].Name, "actual zip should have the mimetype as the first file")
	assert.Equalf(t, "mimetype", expectedZip.File[0].Name, "expected zip should have the mimetype as the first file")

	for _, zipFile := range actualZip.File {
		var found bool
		for _, expectedZipFile := range expectedZip.File {
			if zipFile.Name == expectedZipFile.Name {
				if filesAreTheSame, issue := zipFilesAreEqual(t, zipFile, expectedZipFile); !filesAreTheSame {
					assert.FailNow(t, issue)
				}

				found = true
				break
			}
		}

		if found {
			continue
		}

		assert.FailNowf(t, "did not find file %q in the actual zip", zipFile.Name)
	}
}

func zipFilesAreEqual(t *testing.T, actual, expected *zip.File) (bool, string) {
	if actual.Method != expected.Method || actual.CompressedSize64 != expected.CompressedSize64 || actual.UncompressedSize64 != expected.UncompressedSize64 {
		return false, fmt.Sprintf("%q has file metadata that does not match what is expected.\nMethod is %d and expected %d\nCompressedSize64 is %d and expected %d\nUncompressedSize64 is %d and expected %d", actual.Name, actual.Method, expected.Method, actual.CompressedSize64, expected.CompressedSize64, actual.UncompressedSize64, expected.UncompressedSize64)
	}

	actualReader, err := actual.Open()
	require.NoErrorf(t, err, "failed to open actual zip contents for %q", actual.Name)

	defer tests.MustClose(t, actualReader)

	var actualContents = &bytes.Buffer{}
	_, err = io.Copy(actualContents, actualReader)
	require.NoErrorf(t, err, "failed to read in actual zip contents for %q", actual.Name)

	expectedReader, err := actual.Open()
	require.NoErrorf(t, err, "failed to open expected zip contents for %q", expected.Name)

	defer tests.MustClose(t, expectedReader)

	var expectedContents = &bytes.Buffer{}
	_, err = io.Copy(expectedContents, expectedReader)
	require.NoErrorf(t, err, "failed to read in expected zip contents for %q", expected.Name)

	return bytes.Equal(expectedContents.Bytes(), actualContents.Bytes()), fmt.Sprintf("%q does not have the expected bytes", actual.Name) // the message here will only be used when the bytes are not equal
}
