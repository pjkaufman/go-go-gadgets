//go:build unit

package converter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/converter"
	"github.com/stretchr/testify/assert"
)

type convertMdToHtmlSongTestCase struct {
	inputFilePath string
	inputContent  string
	expectedHtml  string
}

var convertMdToHtmlSongTestCases = map[string]convertMdToHtmlSongTestCase{
	"a valid file should properly get turned into html": {
		inputFilePath: "He Is.md",
		inputContent:  HeIsFileMd,
		expectedHtml:  HeIsFileHtml,
	},
	"a valid file with another title should properly get converted into an html file": {
		inputFilePath: "Above It All (There Stands Jesus).md",
		inputContent:  AboveItAllFileMd,
		expectedHtml:  AboveItAllFileHtml,
	},
	"a valid file with just the melody in the second row should properly get converted into an html file": {
		inputFilePath: "Behold The Heavens.md",
		inputContent:  BeholdTheHeavensFileMd,
		expectedHtml:  BeholdTheHeavensFileHtml,
	},
	"a valid file with a verse in the second row should properly get converted into an html file": {
		inputFilePath: "Be Thou Exalted.md",
		inputContent:  BeThouExaltedFileMd,
		expectedHtml:  BeThouExaltedFileHtml,
	},
}

func TestConvertMdToHtmlSong(t *testing.T) {
	for name, args := range convertMdToHtmlSongTestCases {
		t.Run(name, func(t *testing.T) {
			actual, err := converter.ConvertMdToHtmlSong(args.inputFilePath, args.inputContent, converter.Digital, false)
			assert.Nil(t, err, "there should be no errors when parsing the song contents for the UTs")

			assert.Equal(t, args.expectedHtml, actual)
		})
	}
}
