//go:build unit

package epubcheck_test

import (
	"testing"

	epubcheck "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check"
	"github.com/stretchr/testify/assert"
)

type parseEPUBCheckTestCase struct {
	input    string
	expected epubcheck.ValidationErrors
}

var parseEPUBCheckTestCases = map[string]parseEPUBCheckTestCase{
	"No validation issues returns no values": {
		input:    "Validating using EPUB version 2.0.1 rules.\nCheck finished with errors\nMessages: 0 fatals / 0 errors / 0 warnings / 0 infos\nEPUBCheck completed",
		expected: epubcheck.ValidationErrors{},
	},
	"A single validation issue returns correct model": {
		input: `ERROR(RSC-005): /home/user/Documents/Book.epub/chapter1.html(5,10): Error while parsing file: element "img" missing required attribute "alt"`,
		expected: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-005",
					FilePath: "chapter1.html",
					Location: &epubcheck.Position{Line: 5, Column: 10},
					Message:  `Error while parsing file: element "img" missing required attribute "alt"`,
				},
			},
		},
	},
	"Multiple validation issues returns multiple models": {
		input: `ERROR(RSC-005): /home/user/Documents/Book.epub/chapter1.html(5,10): Error while parsing file: element "img" missing required attribute "alt"
ERROR(RSC-007): /home/user/Documents/Book.epub/chapter2.html(15,20): Referenced resource "chapter3.html" could not be found in the EPUB.`,
		expected: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-005",
					FilePath: "chapter1.html",
					Location: &epubcheck.Position{Line: 5, Column: 10},
					Message:  `Error while parsing file: element "img" missing required attribute "alt"`,
				},
				{
					Code:     "RSC-007",
					FilePath: "chapter2.html",
					Location: &epubcheck.Position{Line: 15, Column: 20},
					Message:  `Referenced resource "chapter3.html" could not be found in the EPUB.`,
				},
			},
		},
	},
	"A validation issue with -1,-1 results in nil Position": {
		input: `ERROR(RSC-999): /home/user/Documents/Book.epub/chapter4.html(-1,-1): Some general error with no position`,
		expected: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-999",
					FilePath: "chapter4.html",
					Location: nil,
					Message:  "Some general error with no position",
				},
			},
		},
	},
	"Validation issues for duplicate id references should be cut down to a single instance per file per id": {
		input: `ERROR(RSC-005): /home/user/Documents/Book.epub/OPS/section-0009.html(15,54): Error while parsing file: Duplicate ID "auto_bookmark_toc_9"
ERROR(RSC-005): /home/user/Documents/Book.epub/OPS/section-0009.html(14,54): Error while parsing file: Duplicate ID "auto_bookmark_toc_9"`,
		expected: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-005",
					FilePath: "OPS/section-0009.html",
					Location: &epubcheck.Position{
						Line:   14,
						Column: 54,
					},
					Message: `Error while parsing file: Duplicate ID "auto_bookmark_toc_9"`,
				},
			},
		},
	},
}

func TestParseEPUBCheckOutput(t *testing.T) {
	for name, args := range parseEPUBCheckTestCases {
		t.Run(name, func(t *testing.T) {
			actual, err := epubcheck.ParseEPUBCheckOutput(args.input)
			assert.NoError(t, err)
			assert.Equal(t, args.expected, actual)
		})
	}
}
