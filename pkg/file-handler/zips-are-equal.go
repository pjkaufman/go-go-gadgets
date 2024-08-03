//go:build unit

package filehandler

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
)

// ZipsAreEqual compares two files with the same name in two separate directories to see if they are logically equivalent.
// This is a test helper function and is not meant for use in non-testing code.
func ZipsAreEqual(filename, originalFileDir, expectedFileDir string, firstFileMustBeMimetype bool) (bool, string) {
	var originalCbzPath = JoinPath(originalFileDir, filename)
	actualZip, err := zip.OpenReader(originalCbzPath)
	if err != nil {
		log.Fatalf("Failed to open zip file %q: %s", originalCbzPath, err)
	}
	defer actualZip.Close()

	var expectedFilePath = JoinPath(expectedFileDir, filename)
	expectedZip, err := zip.OpenReader(expectedFilePath)
	if err != nil {
		log.Fatalf("Failed to open zip file %q: %s", expectedFilePath, err)
	}
	defer expectedZip.Close()

	if len(actualZip.File) != len(expectedZip.File) {
		return false, fmt.Sprintf("expected %d files in cbz, but got %d files", len(expectedZip.File), len(actualZip.File))
	}

	if firstFileMustBeMimetype {
		if actualZip.File[0].Name != "mimetype" {
			return false, "actual zip should have the mimetype as the first file"
		} else if expectedZip.File[0].Name != "mimetype" {
			return false, "expected zip should have the mimetype as the first file"
		}
	}

	for _, zipFile := range actualZip.File {
		var found bool
		for _, expectedZipFile := range expectedZip.File {
			if zipFile.Name == expectedZipFile.Name {
				if filesAreTheSame, issue := zipFilesAreEqual(zipFile, expectedZipFile); !filesAreTheSame {
					return false, issue
				}

				found = true
				break
			}
		}

		if found {
			continue
		}

		return false, fmt.Sprintf("did not find file %q in the actual zip", zipFile.Name)
	}

	return true, ""
}

func zipFilesAreEqual(actual, expected *zip.File) (bool, string) {
	if actual.Method != expected.Method || actual.CompressedSize64 != expected.CompressedSize64 || actual.UncompressedSize64 != expected.UncompressedSize64 {
		return false, fmt.Sprintf("%q has file metadata that does not match what is expected.\nMethod is %d and expected %d\nCompressedSize64 is %d and expected%d\nUncompressedSize64 is %d and expected %d", actual.Name, actual.Method, expected.Method, actual.CompressedSize64, expected.CompressedSize64, actual.UncompressedSize64, expected.UncompressedSize64)
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

	return bytes.Equal(expectedContents.Bytes(), actualContents.Bytes()), fmt.Sprintf("%q does not have the expected bytes", actual.Name) // the message here will only be used when the bytes are not equal
}
