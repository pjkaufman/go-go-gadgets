//go:build unit

package converter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/converter"
	"github.com/stretchr/testify/assert"
)

type ConvertMdToHtmlSongTestCase struct {
	InputFilePath string
	InputContent  string
	ExpectedHtml  string
}

var ConvertMdToHtmlSongTestCases = map[string]ConvertMdToHtmlSongTestCase{
	"a valid file should properly get turned into html": {
		InputFilePath: "He Is.md",
		InputContent:  HeIsFileMd,
		ExpectedHtml:  HeIsFileHtml,
	},
	"a valid file with another title should properly get converted into an html file": {
		InputFilePath: "Above It All (There Stands Jesus).md",
		InputContent:  AboveItAllFileMd,
		ExpectedHtml:  AboveItAllFileHtml,
	},
	"a valid file with just the melody in the second row should properly get converted into an html file": {
		InputFilePath: "Behold The Heavens.md",
		InputContent:  BeholdTheHeavensFileMd,
		ExpectedHtml:  BeholdTheHeavensFileHtml,
	},
	"a valid file with a verse in the second row should properly get converted into an html file": {
		InputFilePath: "Be Thou Exalted.md",
		InputContent:  BeThouExaltedFileMd,
		ExpectedHtml:  BeThouExaltedFileHtml,
	},
}

func TestConvertMdToHtmlSong(t *testing.T) {
	for name, args := range ConvertMdToHtmlSongTestCases {
		t.Run(name, func(t *testing.T) {
			actual, err := converter.ConvertMdToHtmlSong(args.InputFilePath, args.InputContent)
			assert.Nil(t, err, "there should be no errors when parsing the song contents for the UTs")

			assert.Equal(t, args.ExpectedHtml, actual)
		})
	}
}
