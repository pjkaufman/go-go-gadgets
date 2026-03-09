//go:build unit

package compare_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/compare"
	"github.com/stretchr/testify/assert"
)

type compareLinesTestCase struct {
	pdfLines, htmlLines []string
	differences         []compare.Difference
}

var compareLinesTestCases = map[string]compareLinesTestCase{
	"When there is a difference in lines between the PDF and HTML content, there should be a difference mentioning that difference": {
		pdfLines:  []string{"Line 1", "Line 2"},
		htmlLines: []string{"Line 1"},
		differences: []compare.Difference{
			{
				Message:  "Line count mismatch for HTML and PDF file: expected 1 but was 2",
				DiffType: compare.LikelyMismatch,
			},
			{
				Message:  "Ran out of lines in the HTML to compare to the PDF: had 1 line to go",
				DiffType: compare.DefiniteMismatch,
			},
		},
		/**
		Scenarios to test yet:
		- difference in line numbers where the HTML has more lines than the pdf (lest make it 5 more than the PDF)
		- difference where the one line in the HTML is broken into 2 or more in the PDF
		- difference where the one line in the HTML is a partial wrap onto 2 or more lines in the PDF
		- difference where the one line does not match the other when it comes to a single whitespace character
		- difference where the last line of the HTML and the PDF do not match at all
		- difference where there are at least 2 line wraps, 3 partial line wraps, a whitespace character issue, and the last line not actually matching
		*/
	},
}

func TestCompareLines(t *testing.T) {
	for name, args := range compareLinesTestCases {
		t.Run(name, func(t *testing.T) {
			actual := compare.CompareLines(args.pdfLines, args.htmlLines)

			assert.Equal(t, args.differences, actual)
		})
	}
}
