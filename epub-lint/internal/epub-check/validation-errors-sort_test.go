//go:build unit

package epubcheck_test

import (
	"testing"

	epubcheck "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check"
	"github.com/stretchr/testify/assert"
)

type sortValidationErrorsTestCase struct {
	name     string
	input    epubcheck.ValidationErrors
	expected epubcheck.ValidationErrors
}

// TODO: use a map instead
var sortValidationErrorsTestCases = []sortValidationErrorsTestCase{
	{
		name: "When an empty metadata property is present in the list of validation errors, it should be the first entry",
		input: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-001",
					FilePath: "chapter1.html",
					Location: &epubcheck.Position{Line: 5, Column: 10},
					Message:  "Some error",
				},
				{
					Code:     "RSC-002",
					FilePath: "chapter1.html",
					Location: &epubcheck.Position{Line: 3, Column: 5},
					Message:  epubcheck.EmptyMetadataProperty + "publisher\" invalid",
				},
				{
					Code:     "RSC-005",
					FilePath: "apendix.html",
					Location: &epubcheck.Position{Line: 1, Column: 1},
					Message:  "Some error",
				},
			},
		},
		expected: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-002",
					FilePath: "chapter1.html",
					Location: &epubcheck.Position{Line: 3, Column: 5},
					Message:  epubcheck.EmptyMetadataProperty + "publisher\" invalid",
				},
				{
					Code:     "RSC-005",
					FilePath: "apendix.html",
					Location: &epubcheck.Position{Line: 1, Column: 1},
					Message:  "Some error",
				},
				{
					Code:     "RSC-001",
					FilePath: "chapter1.html",
					Location: &epubcheck.Position{Line: 5, Column: 10},
					Message:  "Some error",
				},
			},
		},
	},
	{
		name: "When there are multiple different file paths present, they should be sorted in ascending order",
		input: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-005",
					FilePath: "b-roll.html",
					Location: &epubcheck.Position{Line: 7, Column: 25},
					Message:  "Some error",
				},
				{
					Code:     "RSC-005",
					FilePath: "zing.html",
					Location: &epubcheck.Position{Line: 529, Column: 7},
					Message:  "Some error",
				},
				{
					Code:     "RSC-005",
					FilePath: "apendix.html",
					Location: &epubcheck.Position{Line: 1, Column: 1},
					Message:  "Some error",
				},
			},
		},
		expected: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-005",
					FilePath: "apendix.html",
					Location: &epubcheck.Position{Line: 1, Column: 1},
					Message:  "Some error",
				},
				{
					Code:     "RSC-005",
					FilePath: "b-roll.html",
					Location: &epubcheck.Position{Line: 7, Column: 25},
					Message:  "Some error",
				},
				{
					Code:     "RSC-005",
					FilePath: "zing.html",
					Location: &epubcheck.Position{Line: 529, Column: 7},
					Message:  "Some error",
				},
			},
		},
	},
	{
		name: "When there are multiple validation errors in the same file, any with no location present should be sorted to be after ones with a location and those with locations should be sorted in descending order",
		input: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "OPF-014",
					FilePath: "content.opf",
					Location: &epubcheck.Position{Line: 529, Column: 7},
					Message:  "Some error",
				},
				{
					Code:     "OPF-015",
					FilePath: "content.opf",
					Location: nil,
					Message:  "Some error",
				},
				{
					Code:     "OPF-014",
					FilePath: "content.opf",
					Location: &epubcheck.Position{Line: 16, Column: 15},
					Message:  "Some error",
				},
			},
		},
		expected: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "OPF-014",
					FilePath: "content.opf",
					Location: &epubcheck.Position{Line: 529, Column: 7},
					Message:  "Some error",
				},
				{
					Code:     "OPF-014",
					FilePath: "content.opf",
					Location: &epubcheck.Position{Line: 16, Column: 15},
					Message:  "Some error",
				},
				{
					Code:     "OPF-015",
					FilePath: "content.opf",
					Location: nil,
					Message:  "Some error",
				},
			},
		},
	},
	{
		name: "When there are multiple validation errors in the same file on the same line, then errors are sorted by column in a descending order",
		input: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-005",
					FilePath: "alpha.html",
					Location: &epubcheck.Position{Line: 7, Column: 25},
					Message:  "Some error",
				},
				{
					Code:     "RSC-005",
					FilePath: "alpha.html",
					Location: &epubcheck.Position{Line: 7, Column: 578},
					Message:  "Some error",
				},
				{
					Code:     "RSC-005",
					FilePath: "alpha.html",
					Location: &epubcheck.Position{Line: 7, Column: 7},
					Message:  "Some error",
				},
			},
		},
		expected: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-005",
					FilePath: "alpha.html",
					Location: &epubcheck.Position{Line: 7, Column: 578},
					Message:  "Some error",
				},
				{
					Code:     "RSC-005",
					FilePath: "alpha.html",
					Location: &epubcheck.Position{Line: 7, Column: 25},
					Message:  "Some error",
				},
				{
					Code:     "RSC-005",
					FilePath: "alpha.html",
					Location: &epubcheck.Position{Line: 7, Column: 7},
					Message:  "Some error",
				},
			},
		},
	},
	{
		name: "When dealing with multiple validation errors, then the sort order should be empty property errors first, then file paths sorted ascending, then file line locations in descending order with nil being after those, file column locations in descending order",
		input: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-001",
					FilePath: "chapter1.html",
					Location: &epubcheck.Position{Line: 5, Column: 10},
					Message:  "Some error",
				},
				{
					Code:     "RSC-002",
					FilePath: "chapter1.html",
					Location: &epubcheck.Position{Line: 3, Column: 5},
					Message:  epubcheck.EmptyMetadataProperty + "publisher\" invalid",
				},
				{
					Code:     "RSC-005",
					FilePath: "apendix.html",
					Location: &epubcheck.Position{Line: 1, Column: 1},
					Message:  "Some error",
				},
				{
					Code:     "RSC-005",
					FilePath: "new.html",
					Location: nil,
					Message:  "Some error",
				},
				{
					Code:     "RSC-005",
					FilePath: "new.html",
					Location: &epubcheck.Position{Line: 18, Column: 10},
					Message:  "Some error",
				},
				{
					Code:     "RSC-005",
					FilePath: "new.html",
					Location: &epubcheck.Position{Line: 18, Column: 20},
					Message:  "Some error",
				},
				{
					Code:     "RSC-005",
					FilePath: "new.html",
					Location: &epubcheck.Position{Line: 5, Column: 2},
					Message:  "Some error",
				},
			},
		},
		expected: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-002",
					FilePath: "chapter1.html",
					Location: &epubcheck.Position{Line: 3, Column: 5},
					Message:  epubcheck.EmptyMetadataProperty + "publisher\" invalid",
				},
				{
					Code:     "RSC-005",
					FilePath: "apendix.html",
					Location: &epubcheck.Position{Line: 1, Column: 1},
					Message:  "Some error",
				},
				{
					Code:     "RSC-001",
					FilePath: "chapter1.html",
					Location: &epubcheck.Position{Line: 5, Column: 10},
					Message:  "Some error",
				},
				{
					Code:     "RSC-005",
					FilePath: "new.html",
					Location: &epubcheck.Position{Line: 18, Column: 20},
					Message:  "Some error",
				},
				{
					Code:     "RSC-005",
					FilePath: "new.html",
					Location: &epubcheck.Position{Line: 18, Column: 10},
					Message:  "Some error",
				},
				{
					Code:     "RSC-005",
					FilePath: "new.html",
					Location: &epubcheck.Position{Line: 5, Column: 2},
					Message:  "Some error",
				},
				{
					Code:     "RSC-005",
					FilePath: "new.html",
					Location: nil,
					Message:  "Some error",
				},
			},
		},
	},
}

func TestSortValidationErrors(t *testing.T) {
	for _, tc := range sortValidationErrorsTestCases {
		t.Run(tc.name, func(t *testing.T) {
			input := tc.input
			input.Sort()
			assert.Equal(t, tc.expected, input)
		})
	}
}
