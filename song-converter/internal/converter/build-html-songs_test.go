//go:build unit

package converter_test

import (
	"fmt"
	"testing"

	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/converter"
	"github.com/stretchr/testify/assert"
)

type BuildHtmlSongsTestCase struct {
	InputMdInfo   []converter.MdFileInfo
	ExpectedSongs []string
	ExpectedHtml  string
	ExpectError   bool
}

var BuildHtmlSongsTestCases = map[string]BuildHtmlSongsTestCase{
	"no files provided should just result in an empty string an no headers": {
		ExpectedSongs: []string{},
	},
	"multiple files should be a new line character followed by each song with a new line character after it": {
		InputMdInfo: []converter.MdFileInfo{
			{
				FilePath:     "Above It All (There Stands Jesus).md",
				FileContents: AboveItAllFileMd,
			},
			{
				FilePath:     "Be Thou Exalted.md",
				FileContents: BeThouExaltedFileMd,
			},
			{
				FilePath:     "Behold The Heavens.md",
				FileContents: BeholdTheHeavensFileMd,
			},
			{
				FilePath:     "He Is.md",
				FileContents: HeIsFileMd,
			},
		},
		ExpectedSongs: []string{"above-it-all-there-stands-jesus", "be-thou-exalted", "behold-the-heavens", "he-is"},
		ExpectedHtml:  fmt.Sprintf("%s\n%s\n%s\n%s\n", AboveItAllFileHtml, BeThouExaltedFileHtml, BeholdTheHeavensFileHtml, HeIsFileHtml),
	},
	"multiple headers with teh same heading get broken into unique header ids": {
		InputMdInfo: []converter.MdFileInfo{
			{
				FilePath:     "Be Thou Exalted.md",
				FileContents: BeThouExaltedFileMd,
			},
			{
				FilePath:     "Be Thou Exalted.md",
				FileContents: BeThouExaltedFileMd,
			},
		},
		ExpectedSongs: []string{"be-thou-exalted", "be-thou-exalted-2"},
		ExpectedHtml:  fmt.Sprintf("%s\n%s\n", BeThouExaltedFileHtml, BeThouExalted2FileHtml),
	},
}

func TestBuildHtmlSongs(t *testing.T) {
	for name, args := range BuildHtmlSongsTestCases {
		t.Run(name, func(t *testing.T) {
			actual, actualSongIds, err := converter.BuildHtmlSongs(args.InputMdInfo)
			if args.ExpectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, args.ExpectedHtml, actual)
			assert.Equal(t, args.ExpectedSongs, actualSongIds)
		})
	}
}
