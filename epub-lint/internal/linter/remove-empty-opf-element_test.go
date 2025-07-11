package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

type removeEmptyOpfElementsTestCase struct {
	elementName    string
	lineNum        int
	opfContents    string
	expectedOutput string
	expectedDelete bool
}

var removeEmptyOpfElementsTestCases = map[string]removeEmptyOpfElementsTestCase{
	"Remove dc:identifier element's line with ending tag": {
		elementName: "dc:identifier",
		lineNum:     1,
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:identifier></dc:identifier>
    <dc:title>Example Book</dc:title>
</metadata>`,
		expectedOutput: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:title>Example Book</dc:title>
</metadata>`,
		expectedDelete: true,
	},
	"Remove dc:identifier element's line with self-closing tag": {
		elementName: "dc:identifier",
		lineNum:     1,
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:identifier />
    <dc:title>Example Book</dc:title>
</metadata>`,
		expectedOutput: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:title>Example Book</dc:title>
</metadata>`,
		expectedDelete: true,
	},
	"Remove dc:description element's line with ending tag": {
		elementName: "dc:description",
		lineNum:     2,
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:title>Example Book</dc:title>
    <dc:description></dc:description>
</metadata>`,
		expectedOutput: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:title>Example Book</dc:title>
</metadata>`,
		expectedDelete: true,
	},
	"Remove dc:description element's line with self-closing tag": {
		elementName: "dc:description",
		lineNum:     2,
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:title>Example Book</dc:title>
    <dc:description />
</metadata>`,
		expectedOutput: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:title>Example Book</dc:title>
</metadata>`,
		expectedDelete: true,
	},
	"Remove dc:description element, but not the line with ending tag": {
		elementName: "dc:description",
		lineNum:     2,
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:title>Example Book</dc:title>
    <dc:description></dc:description></metadata>`,
		expectedOutput: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:title>Example Book</dc:title>
    </metadata>`,
		expectedDelete: false,
	},
	"Remove dc:description element, but not the line with self-closing tag": {
		elementName: "dc:description",
		lineNum:     2,
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:title>Example Book</dc:title>
    <dc:description /></metadata>`,
		expectedOutput: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:title>Example Book</dc:title>
    </metadata>`,
		expectedDelete: false,
	},
}

func TestRemoveEmptyOpfElements(t *testing.T) {
	for name, tc := range removeEmptyOpfElementsTestCases {
		t.Run(name, func(t *testing.T) {
			actualOutput, actualDelete, err := linter.RemoveEmptyOpfElements(tc.elementName, tc.lineNum, tc.opfContents)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedOutput, actualOutput)
			assert.Equal(t, tc.expectedDelete, actualDelete)
		})
	}
}
