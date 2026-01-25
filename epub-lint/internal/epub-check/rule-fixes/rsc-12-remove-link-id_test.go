//go:build unit

package rulefixes_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
	"github.com/stretchr/testify/assert"
)

func TestRemoveLinkId(t *testing.T) {
	tests := []struct {
		name            string
		fileContents    string
		lineToUpdate    int
		startOfFragment int
		change          positions.TextEdit
	}{
		{
			name:            "Line number does not exist",
			fileContents:    "line1\nline2\nline3",
			lineToUpdate:    5,
			startOfFragment: 1,
			change:          positions.TextEdit{},
		},
		{
			name:            "Line has less characters than start of fragment",
			fileContents:    "line1\nline2\nline3",
			lineToUpdate:    2,
			startOfFragment: 11,
			change:          positions.TextEdit{},
		},
		{
			name: "Line has href with # in the link",
			fileContents: `line1
<a href="link#id">link</a>
line3`,
			lineToUpdate:    2,
			startOfFragment: 15,
			change: positions.TextEdit{
				Range: positions.Range{
					Start: positions.Position{
						Line:   2,
						Column: 14,
					},
					End: positions.Position{
						Line:   2,
						Column: 17,
					},
				},
			},
		},
		{
			name: "Line has href without # in the link",
			fileContents: `line1
<a href="link">link</a>
line3`,
			lineToUpdate:    2,
			startOfFragment: 16,
			change:          positions.TextEdit{},
		},
		{
			name: "Line has src without # in the link",
			fileContents: `line1
<content src="link"/>
line3`,
			lineToUpdate:    2,
			startOfFragment: 21,
			change:          positions.TextEdit{},
		},
		{
			name: "Line has src with # in the link",
			fileContents: `line1
<content src="link#id"/>
line3`,
			lineToUpdate:    2,
			startOfFragment: 20,
			change: positions.TextEdit{
				Range: positions.Range{
					Start: positions.Position{
						Line:   2,
						Column: 19,
					},
					End: positions.Position{
						Line:   2,
						Column: 22,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rulefixes.RemoveLinkId(tt.fileContents, tt.lineToUpdate, tt.startOfFragment)
			assert.Equal(t, tt.change, result)
		})
	}
}
