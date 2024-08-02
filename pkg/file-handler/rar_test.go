//go:build unit

package filehandler_test

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/stretchr/testify/assert"
)

type ConvertRarToCbzTestCase struct {
	Filename string
}

const (
	originalFileDir  = "testdata/original"
	convertedFileDir = "testdata/converted"
)

var convertRarToCbzTestCases = map[string]ConvertRarToCbzTestCase{
	"Converting a cbr to a cbz file should work consistently": {
		Filename: "Whiz-Comics_Shazam_Golden_Arow_Volume_1-Fawcett_Publications.cbr",
	},
}

func TestConvertRarToCbz(t *testing.T) {
	for name, test := range convertRarToCbzTestCases {
		t.Run(name, func(t *testing.T) {
			err := filehandler.ConvertRarToCbz(filehandler.JoinPath(originalFileDir, test.Filename))
			assert.Nil(t, err)

			equalityStatus, issue := cbzsAreEqual(strings.Replace(test.Filename, ".cbr", ".cbz", 1))
			assert.True(t, equalityStatus, issue)

			var originalEpubPath = originalFileDir + string(os.PathSeparator) + test.Filename
			err = os.RemoveAll(strings.Replace(originalEpubPath, ".cbr", ".cbz", 1))
			if err != nil {
				log.Fatalf("failed to remove the result of rar to cbz %q: %s", originalEpubPath, err)
			}
		})
	}
}

func cbzsAreEqual(filename string) (bool, string) {
	var originalCbzPath = filehandler.JoinPath(originalFileDir, filename)
	convertedCbz, err := zip.OpenReader(originalCbzPath)
	if err != nil {
		log.Fatalf("Failed to open zip file %q: %s", originalCbzPath, err)
	}
	defer convertedCbz.Close()

	var convertedCbzPath = filehandler.JoinPath(convertedFileDir, filename)
	expectedCbz, err := zip.OpenReader(convertedCbzPath)
	if err != nil {
		log.Fatalf("Failed to open zip file %q: %s", convertedCbzPath, err)
	}
	defer expectedCbz.Close()

	fmt.Println(originalCbzPath, convertedCbzPath)
	if len(convertedCbz.File) != len(expectedCbz.File) {
		return false, fmt.Sprintf("expected %d files in cbz, but got %d files", len(expectedCbz.File), len(convertedCbz.File))
	}

	for _, zipFile := range convertedCbz.File {
		var found bool
		for _, expectedZipFile := range expectedCbz.File {
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

		return false, fmt.Sprintf("did not find file %q in the actual cbz", zipFile.Name)
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

func BenchmarkConvertRarToCbz(b *testing.B) {
	var filename = "Whiz-Comics_Shazam_Golden_Arow_Volume_1-Fawcett_Publications.cbr"

	for n := 0; n < b.N; n++ {
		var originalCbrPath = filehandler.JoinPath(originalFileDir, filename)
		err := filehandler.ConvertRarToCbz(originalCbrPath)
		if err != nil {
			log.Fatalf("failed to convert cbr to cbz %q: %s", originalCbrPath, err)
		}

		err = os.RemoveAll(strings.Replace(originalCbrPath, ".cbr", ".cbz", 1))
		if err != nil {
			log.Fatalf("failed to remove the result of cbr to cbz %q: %s", originalCbrPath, err)
		}
	}
}
