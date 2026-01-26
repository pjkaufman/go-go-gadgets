//go:build unit

package rulefixes_test

import (
	"testing"

	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
)

type removeLinkIDTestCase struct {
	fileContents    string
	lineToUpdate    int
	startOfFragment int
	expected        string
}

var removeLinkIDTestCases = map[string]removeLinkIDTestCase{
	"Line number does not exist": {
		fileContents:    "line1\nline2\nline3",
		lineToUpdate:    5,
		startOfFragment: 1,
		expected:        "line1\nline2\nline3",
	},
	"Line has less characters than start of fragment": {
		fileContents:    "line1\nline2\nline3",
		lineToUpdate:    2,
		startOfFragment: 11,
		expected:        "line1\nline2\nline3",
	},
	"Line has href with # in the link": {
		fileContents: `line1
<a href="link#id">link</a>
line3`,
		lineToUpdate:    2,
		startOfFragment: 16,
		expected: `line1
<a href="link">link</a>
line3`,
	},
	"Line has href without # in the link": {
		fileContents: `line1
<a href="link">link</a>
line3`,
		lineToUpdate:    2,
		startOfFragment: 16,
		expected: `line1
<a href="link">link</a>
line3`,
	},
	"Line has src without # in the link": {
		fileContents: `line1
<content src="link"/>
line3`,
		lineToUpdate:    2,
		startOfFragment: 21,
		expected: `line1
<content src="link"/>
line3`,
	},
	"Line has src with # in the link": {
		fileContents: `line1
<content src="link#id"/>
line3`,
		lineToUpdate:    2,
		startOfFragment: 23,
		expected: `line1
<content src="link"/>
line3`,
	},
}

func TestRemoveLinkId(t *testing.T) {
	for name, args := range removeLinkIDTestCases {
		t.Run(name, func(t *testing.T) {
			edit := rulefixes.RemoveLinkId(args.fileContents, args.lineToUpdate, args.startOfFragment)

			checkFinalOutputMatches(t, args.fileContents, args.expected, edit)
		})
	}
}
