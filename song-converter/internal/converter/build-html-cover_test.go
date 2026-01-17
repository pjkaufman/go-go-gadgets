//go:build unit

package converter_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/converter"
	"github.com/stretchr/testify/assert"
)

type BuildHtmlCoverTestCase struct {
	InputCoverMd  string
	ExpectedHtml  string
	Type          string
	ExtraStyleCss string
	DateCreated   time.Time
}

var BuildHtmlCoverTestCases = map[string]BuildHtmlCoverTestCase{
	"a valid file should properly get turned into html with no style content": {
		InputCoverMd: coverFileMd,
		ExpectedHtml: fmt.Sprintf(coverFileHtmlFormat, "", "abridged", "Abridged", "Jul 2024"),
		Type:         "Abridged",
		DateCreated:  time.Date(2024, time.July, 1, 0, 0, 0, 0, time.Local),
	},
	"a valid file should properly get turned into html and respect the type and created date info": {
		InputCoverMd: coverFileMd,
		ExpectedHtml: fmt.Sprintf(coverFileHtmlFormat, "", "unabridged", "Unabridged", "Jan 2021"),
		Type:         "Unabridged",
		DateCreated:  time.Date(2021, time.January, 1, 0, 0, 0, 0, time.Local),
	},
	"a valid file with extra css provided should properly get turned into html and respect the type and created date info": {
		InputCoverMd:  coverFileMd,
		ExpectedHtml:  fmt.Sprintf(coverFileHtmlFormat, "; font-size: 52pt;", "unabridged", "Unabridged", "Jan 2021"),
		Type:          "Unabridged",
		ExtraStyleCss: "font-size: 52pt;",
		DateCreated:   time.Date(2021, time.January, 1, 0, 0, 0, 0, time.Local),
	},
}

func TestBuildHtmlCover(t *testing.T) {
	for name, args := range BuildHtmlCoverTestCases {
		t.Run(name, func(t *testing.T) {
			actual := converter.BuildHtmlCover(args.InputCoverMd, args.Type, args.ExtraStyleCss, args.DateCreated)

			assert.Equal(t, args.ExpectedHtml, actual)
		})
	}
}
