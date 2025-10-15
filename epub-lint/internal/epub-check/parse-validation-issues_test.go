//go:build unit

package epubcheck

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type parseEPUBCheckTestCase struct {
	name     string
	input    string
	expected ValidationErrors
}

var parseEPUBCheckTestCases = []parseEPUBCheckTestCase{
	{
		name:     "no validation issues returns no values",
		input:    "Validating using EPUB version 2.0.1 rules.\nCheck finished with errors\nMessages: 0 fatals / 0 errors / 0 warnings / 0 infos\nEPUBCheck completed",
		expected: ValidationErrors{},
	},
	{
		name:  "single validation issue returns correct model",
		input: `ERROR(RSC-005): /home/user/Documents/Book.epub/chapter1.html(5,10): Error while parsing file: element "img" missing required attribute "alt"`,
		expected: ValidationErrors{
			ValidationIssues: []ValidationError{
				{
					Code:     "RSC-005",
					FilePath: "chapter1.html",
					Location: &Position{Line: 5, Column: 10},
					Message:  `Error while parsing file: element "img" missing required attribute "alt"`,
				},
			},
		},
	},
	{
		name: "multiple validation issues returns multiple models",
		input: `ERROR(RSC-005): /home/user/Documents/Book.epub/chapter1.html(5,10): Error while parsing file: element "img" missing required attribute "alt"
ERROR(RSC-007): /home/user/Documents/Book.epub/chapter2.html(15,20): Referenced resource "chapter3.html" could not be found in the EPUB.`,
		expected: ValidationErrors{
			ValidationIssues: []ValidationError{
				{
					Code:     "RSC-005",
					FilePath: "chapter1.html",
					Location: &Position{Line: 5, Column: 10},
					Message:  `Error while parsing file: element "img" missing required attribute "alt"`,
				},
				{
					Code:     "RSC-007",
					FilePath: "chapter2.html",
					Location: &Position{Line: 15, Column: 20},
					Message:  `Referenced resource "chapter3.html" could not be found in the EPUB.`,
				},
			},
		},
	},
	{
		name:  "validation issue with -1,-1 results in nil Position",
		input: `ERROR(RSC-999): /home/user/Documents/Book.epub/chapter4.html(-1,-1): Some general error with no position`,
		expected: ValidationErrors{
			ValidationIssues: []ValidationError{
				{
					Code:     "RSC-999",
					FilePath: "chapter4.html",
					Location: nil,
					Message:  "Some general error with no position",
				},
			},
		},
	},
}

func TestParseEPUBCheckOutput(t *testing.T) {
	for _, tc := range parseEPUBCheckTestCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ParseEPUBCheckOutput(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
