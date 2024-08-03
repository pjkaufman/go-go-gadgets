//go:build unit

package filehandler_test

import (
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

			equalityStatus, issue := filehandler.ZipsAreEqual(strings.Replace(test.Filename, ".cbr", ".cbz", 1), originalFileDir, convertedFileDir, false)
			assert.True(t, equalityStatus, issue)

			var originalEpubPath = originalFileDir + string(os.PathSeparator) + test.Filename
			err = os.RemoveAll(strings.Replace(originalEpubPath, ".cbr", ".cbz", 1))
			if err != nil {
				log.Fatalf("failed to remove the result of rar to cbz %q: %s", originalEpubPath, err)
			}
		})
	}
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
