//go:build unit

package converter_test

import (
	"fmt"
	"testing"

	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/converter"
	"github.com/stretchr/testify/assert"
)

type buildHtmlSongsTestCase struct {
	inputMdInfo   []converter.MdFileInfo
	expectedSongs []string
	expectedHtml  string
	expectError   bool
}

var multipleFileMdInfo = []converter.MdFileInfo{
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
}

var buildHtmlSongsTestCases = map[string]buildHtmlSongsTestCase{
	"no files provided should just result in an empty string an no headers": {
		expectedSongs: []string{},
	},
	"multiple files should be a new line character followed by each song with a new line character after it": {
		inputMdInfo:   multipleFileMdInfo,
		expectedSongs: []string{"above-it-all-there-stands-jesus", "be-thou-exalted", "behold-the-heavens", "he-is"},
		expectedHtml:  fmt.Sprintf("%s\n%s\n%s\n%s\n", AboveItAllFileHtml, BeThouExaltedFileHtml, BeholdTheHeavensFileHtml, HeIsFileHtml),
	},
	"multiple headers with the same heading get broken into unique header ids": {
		inputMdInfo: []converter.MdFileInfo{
			{
				FilePath:     "Be Thou Exalted.md",
				FileContents: BeThouExaltedFileMd,
			},
			{
				FilePath:     "Be Thou Exalted.md",
				FileContents: BeThouExaltedFileMd,
			},
		},
		expectedSongs: []string{"be-thou-exalted", "be-thou-exalted-2"},
		expectedHtml:  fmt.Sprintf("%s\n%s\n", BeThouExaltedFileHtml, BeThouExalted2FileHtml),
	},
}

func TestBuildHtmlSongs(t *testing.T) {
	for name, args := range buildHtmlSongsTestCases {
		t.Run(name, func(t *testing.T) {
			actual, actualSongIds, err := converter.BuildHtmlSongs(args.inputMdInfo, converter.Digital)
			if args.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, args.expectedHtml, actual)
			assert.Equal(t, args.expectedSongs, actualSongIds)
		})
	}
}

func BenchmarkBuildHtmlSongs(b *testing.B) {
	for n := 0; n < b.N; n++ {
		converter.BuildHtmlSongs(multipleFileMdInfo, converter.Digital)
	}
}
