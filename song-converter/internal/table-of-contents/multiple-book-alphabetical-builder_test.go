//go:build unit

package tableofcontents_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/converter"
	tableofcontents "github.com/pjkaufman/go-go-gadgets/song-converter/internal/table-of-contents"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type buildMultipleBookAlphabeticalListItemsTestCase struct {
	expected    string
	expectError bool
	headerInfo  []converter.MdFileInfo
}

var buildMultipleBookAlphabeticalListItemsTestCases = map[string]buildMultipleBookAlphabeticalListItemsTestCase{
	"When a song has only secondary book page numbers, it has exactly 2 book page numbers, and an alternate title, the first page number should get the original title and the second one should get the alternate title": {
		headerInfo: []converter.MdFileInfo{
			{
				FilePath:             "Song.md",
				FileName:             "Song",
				Header:               "Zephaniah Song",
				AlternateTitle:       "Appleton",
				SecondaryPageNumbers: []int{50, 60},
			},
		},
		expected: `<div class="divider">&#9836;</div>
<li><span class="name">Appleton</span><span class="page">(50)</span></li>
<div class="divider">&#9836;</div>
<li><span class="name">Zephaniah Song</span><span class="page">(60)</span></li>
`,
	},
	"When a song has only secondary book page numbers, it has 1 book page number, and an alternate title, the first page number should get two entries: the original title and the alternate title": {
		headerInfo: []converter.MdFileInfo{
			{
				FilePath:             "Song.md",
				FileName:             "Song",
				Header:               "Zephaniah Song",
				AlternateTitle:       "Appleton",
				SecondaryPageNumbers: []int{50},
			},
		},
		expected: `<div class="divider">&#9836;</div>
<li><span class="name">Appleton</span><span class="page">(50)</span></li>
<div class="divider">&#9836;</div>
<li><span class="name">Zephaniah Song</span><span class="page">(50)</span></li>
`,
	},
	"When a song has only secondary book page numbers, it has exactly 2 book page numbers, and no alternate title, the title gets printed back to back for the two page numbers": {
		headerInfo: []converter.MdFileInfo{
			{
				FilePath:             "Song.md",
				FileName:             "Song",
				Header:               "Zephaniah Song",
				SecondaryPageNumbers: []int{50, 60},
			},
		},
		expected: `<div class="divider">&#9836;</div>
<li><span class="name">Zephaniah Song</span><span class="page">(50)</span></li>
<li><span class="name">Zephaniah Song</span><span class="page">(60)</span></li>
`,
	},
	"When a song has a single primary and secondary page number and an alternate title, the primary page number gets the header and the secondary gets the alternate title": {
		headerInfo: []converter.MdFileInfo{
			{
				FilePath:             "Song.md",
				FileName:             "Song",
				Header:               "Zephaniah Song",
				AlternateTitle:       "Appleton",
				PrimaryPageNumbers:   []int{13},
				SecondaryPageNumbers: []int{50},
			},
		},
		expected: `<div class="divider">&#9836;</div>
<li><span class="name">Appleton</span><span class="page">(50)</span></li>
<div class="divider">&#9836;</div>
<li><span class="name">Zephaniah Song</span><span class="page">13</span></li>
`,
	},
}

func TestBuildMultipleBookAlphabeticalListItems(t *testing.T) {
	t.Parallel()

	for name, args := range buildMultipleBookAlphabeticalListItemsTestCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actual, err := tableofcontents.BuildMultipleBookAlphabeticalListItems(nil, args.headerInfo)

			if args.expectError {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}

			assert.Equal(t, args.expected, actual)
		})
	}
}
