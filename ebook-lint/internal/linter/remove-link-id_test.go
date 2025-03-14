//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

func TestRemoveLinkId(t *testing.T) {
	tests := []struct {
		name            string
		fileContents    string
		lineToUpdate    int
		startOfFragment int
		expected        string
	}{
		{
			name:            "Line number does not exist",
			fileContents:    "line1\nline2\nline3",
			lineToUpdate:    5,
			startOfFragment: 0,
			expected:        "line1\nline2\nline3",
		},
		{
			name:            "Line has less characters than start of fragment",
			fileContents:    "line1\nline2\nline3",
			lineToUpdate:    1,
			startOfFragment: 10,
			expected:        "line1\nline2\nline3",
		},
		{
			name: "Line has href with # in the link",
			fileContents: `line1
<a href="link#id">link</a>
line3`,
			lineToUpdate:    1,
			startOfFragment: 15,
			expected: `line1
<a href="link">link</a>
line3`,
		},
		{
			name: "Line has href without # in the link",
			fileContents: `line1
<a href="link">link</a>
line3`,
			lineToUpdate:    1,
			startOfFragment: 15,
			expected: `line1
<a href="link">link</a>
line3`,
		},
		{
			name: "Line has src without # in the link",
			fileContents: `line1
<content src="link"/>
line3`,
			lineToUpdate:    1,
			startOfFragment: 20,
			expected: `line1
<content src="link"/>
line3`,
		},
		{
			name: "Line has src with # in the link",
			fileContents: `line1
<content src="link#id"/>
line3`,
			lineToUpdate:    1,
			startOfFragment: 23,
			expected: `line1
<content src="link"/>
line3`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := linter.RemoveLinkId(tt.fileContents, tt.lineToUpdate, tt.startOfFragment)
			assert.Equal(t, tt.expected, result)
		})
	}
}
