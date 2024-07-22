//go:build unit

package converter_test

import (
	"testing"
	"time"

	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/converter"
	"github.com/stretchr/testify/assert"
)

type BuildHtmlCoverTestCase struct {
	InputCoverMd string
	ExpectedHtml string
}

var BuildHtmlCoverTestCases = map[string]BuildHtmlCoverTestCase{
	"a valid file should properly get turned into html with no style content": {
		InputCoverMd: coverFileMd,
		ExpectedHtml: coverFileHtml,
	},
}

func TestBuildHtmlCover(t *testing.T) {
	for name, args := range BuildHtmlCoverTestCases {
		t.Run(name, func(t *testing.T) {
			actual := converter.BuildHtmlCover(args.InputCoverMd, "Abridged", time.Now())

			assert.Equal(t, args.ExpectedHtml, actual)
		})
	}
}
