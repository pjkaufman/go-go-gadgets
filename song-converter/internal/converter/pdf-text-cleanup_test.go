//go:build unit

package converter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/converter"
	"github.com/stretchr/testify/assert"
)

type pdfTextCleanupTestCase struct {
	InputText     string
	CombineNLines int
	ExpectedLines []string
}

var pdfTextCleanupTestCases = map[string]pdfTextCleanupTestCase{
	"text with blank lines should have blank lines removed from the output": {
		InputText: `Here is some text

Here is some more text`,
		CombineNLines: 0,
		ExpectedLines: []string{"Here is some text", "Here is some more text"},
	},
	"text with just whitespace and a number should be skipped (often just considered to be line numbers in PDFs)": {
		InputText: `Here is some text
674 
Here is some more text`,
		CombineNLines: 0,
		ExpectedLines: []string{"Here is some text", "Here is some more text"},
	},
	"text with a line that starts with 4 or more spaces should have whitespace collapsed from multiple down to a single space and any starting whitespace on a line is removed": {
		InputText: `Here is some text
    Hello  It Is I
Here is some more text`,
		CombineNLines: 0,
		ExpectedLines: []string{"Here is some text", "Hello It Is I", "Here is some more text"},
	},
	"text with a line where it has two or more spaces between text and a number should have the spaces removed from them": {
		InputText: `Here is some text
Some more text      78
Here is some more text`,
		CombineNLines: 0,
		ExpectedLines: []string{"Here is some text", "Some more text78", "Here is some more text"},
	},
}

func TestPdfTextCleanup(t *testing.T) {
	for name, args := range pdfTextCleanupTestCases {
		t.Run(name, func(t *testing.T) {
			actual := converter.PdfTextCleanup(args.InputText, args.CombineNLines)

			assert.Equal(t, args.ExpectedLines, actual)
		})
	}
}
