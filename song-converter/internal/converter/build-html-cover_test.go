//go:build unit

package converter_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/converter"
	"github.com/stretchr/testify/assert"
)

type buildHtmlCoverTestCase struct {
	inputCoverMd  string
	expectedHtml  string
	coverType     string
	extraStyleCss string
	dateCreated   time.Time
}

var buildHtmlCoverTestCases = map[string]buildHtmlCoverTestCase{
	"a valid file should properly get turned into html with no style content": {
		inputCoverMd: coverFileMd,
		expectedHtml: fmt.Sprintf(coverFileHtmlFormat, "", "abridged", "Abridged", "Jul 2024"),
		coverType:    "Abridged",
		dateCreated:  time.Date(2024, time.July, 1, 0, 0, 0, 0, time.Local),
	},
	"a valid file should properly get turned into html and respect the type and created date info": {
		inputCoverMd: coverFileMd,
		expectedHtml: fmt.Sprintf(coverFileHtmlFormat, "", "unabridged", "Unabridged", "Jan 2021"),
		coverType:    "Unabridged",
		dateCreated:  time.Date(2021, time.January, 1, 0, 0, 0, 0, time.Local),
	},
	"a valid file with extra css provided should properly get turned into html and respect the type and created date info": {
		inputCoverMd:  coverFileMd,
		expectedHtml:  fmt.Sprintf(coverFileHtmlFormat, "; font-size: 52pt;", "unabridged", "Unabridged", "Jan 2021"),
		coverType:     "Unabridged",
		extraStyleCss: "font-size: 52pt;",
		dateCreated:   time.Date(2021, time.January, 1, 0, 0, 0, 0, time.Local),
	},
}

func TestBuildHtmlCover(t *testing.T) {
	for name, args := range buildHtmlCoverTestCases {
		t.Run(name, func(t *testing.T) {
			actual := converter.BuildHtmlCover(args.inputCoverMd, args.coverType, args.extraStyleCss, args.dateCreated)

			assert.Equal(t, args.expectedHtml, actual)
		})
	}
}
