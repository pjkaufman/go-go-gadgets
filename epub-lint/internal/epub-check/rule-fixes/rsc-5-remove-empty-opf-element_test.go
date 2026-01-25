package rulefixes_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
	"github.com/stretchr/testify/assert"
)

type removeEmptyOpfElementsTestCase struct {
	elementName    string
	lineNum        int
	opfContents    string
	expectedChange positions.TextEdit
	expectedDelete bool
}

var removeEmptyOpfElementsTestCases = map[string]removeEmptyOpfElementsTestCase{
	"Remove dc:identifier element's line with ending tag": {
		elementName: "dc:identifier",
		lineNum:     2,
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:identifier></dc:identifier>
    <dc:title>Example Book</dc:title>
</metadata>`,
		expectedChange: positions.TextEdit{
			Range: positions.Range{
				Start: positions.Position{
					Line:   2,
					Column: 1,
				},
				End: positions.Position{
					Line:   3,
					Column: 1,
				},
			},
		},
		expectedDelete: true,
	},
	"Remove dc:identifier element's line with self-closing tag": {
		elementName: "dc:identifier",
		lineNum:     2,
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:identifier />
    <dc:title>Example Book</dc:title>
</metadata>`,
		expectedChange: positions.TextEdit{
			Range: positions.Range{
				Start: positions.Position{
					Line:   2,
					Column: 1,
				},
				End: positions.Position{
					Line:   3,
					Column: 1,
				},
			},
		},
		expectedDelete: true,
	},
	"Remove dc:description element's line with ending tag": {
		elementName: "dc:description",
		lineNum:     3,
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:title>Example Book</dc:title>
    <dc:description></dc:description>
</metadata>`,
		expectedChange: positions.TextEdit{
			Range: positions.Range{
				Start: positions.Position{
					Line:   3,
					Column: 1,
				},
				End: positions.Position{
					Line:   4,
					Column: 1,
				},
			},
		},
		expectedDelete: true,
	},
	"Remove dc:description element's line with self-closing tag": {
		elementName: "dc:description",
		lineNum:     3,
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:title>Example Book</dc:title>
    <dc:description />
</metadata>`,
		expectedChange: positions.TextEdit{
			Range: positions.Range{
				Start: positions.Position{
					Line:   3,
					Column: 1,
				},
				End: positions.Position{
					Line:   4,
					Column: 1,
				},
			},
		},
		expectedDelete: true,
	},
	"Remove dc:description element, but not the line with ending tag": {
		elementName: "dc:description",
		lineNum:     3,
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:title>Example Book</dc:title>
    <dc:description></dc:description></metadata>`,
		expectedChange: positions.TextEdit{
			Range: positions.Range{
				Start: positions.Position{
					Line:   3,
					Column: 5,
				},
				End: positions.Position{
					Line:   3,
					Column: 38,
				},
			},
		},
		expectedDelete: false,
	},
	"Remove dc:description element, but not the line with self-closing tag": {
		elementName: "dc:description",
		lineNum:     3,
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:title>Example Book</dc:title>
    <dc:description /></metadata>`,
		expectedChange: positions.TextEdit{
			Range: positions.Range{
				Start: positions.Position{
					Line:   3,
					Column: 5,
				},
				End: positions.Position{
					Line:   3,
					Column: 23,
				},
			},
		},
		expectedDelete: false,
	},
}

func TestRemoveEmptyOpfElements(t *testing.T) {
	for name, tc := range removeEmptyOpfElementsTestCases {
		t.Run(name, func(t *testing.T) {
			actualOutput, actualDelete, err := rulefixes.RemoveEmptyOpfElements(tc.elementName, tc.lineNum, tc.opfContents)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedChange, actualOutput)
			assert.Equal(t, tc.expectedDelete, actualDelete)
		})
	}
}
