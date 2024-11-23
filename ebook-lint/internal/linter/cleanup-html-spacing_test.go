//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

type cleanupHtmlSpacingTestCase struct {
	inputText    string
	expectedText string
}

const (
	simpleParagraphTagExpectedText = `<p>Text here</p>
		<p>More text here</p>
`
)

var cleanupHtmlSpacingTestCases = map[string]cleanupHtmlSpacingTestCase{
	"make sure that empty lines are removed from text in between tags": {
		inputText: `<p>Text here</p>

		<p>More text here</p>`,
		expectedText: simpleParagraphTagExpectedText,
	},
	"make sure that starting whitespace in the text is removed": {
		inputText: `




		<p>Text here</p>

		<p>More text here</p>`,
		expectedText: simpleParagraphTagExpectedText,
	},
	"make sure that trailing whitespace is converted to a single blank line at the end of the file": {
		inputText: `<p>Text here</p>

		<p>More text here</p>
		`,
		expectedText: simpleParagraphTagExpectedText,
	},
	"make sure that there is a blank line at the end of a file": {
		inputText: `<p>Text here</p>
		<p>More text here</p>`,
		expectedText: simpleParagraphTagExpectedText,
	},
	"make sure that a paragraph tag with a new line and whitespace at the start gets it removed": {
		inputText: `<p class="test">     
		Text here</p>
		<p>More text here</p>`,
		expectedText: `<p class="test">Text here</p>
		<p>More text here</p>
`,
	},
	"make sure that a paragraph tag with a new line and whitespace at the end gets it removed": {
		inputText: `<p class="test">Text here      
		
		
		
		</p>
		<p>More text here</p>`,
		expectedText: `<p class="test">Text here</p>
		<p>More text here</p>
`,
	},
	"make sure that a paragraph tag with a new line and whitespace at the end after a nested element that ends the content also gets it removed": {
		inputText: `<p class="test">Text <i>here</i>      
	
	
	
	</p>
	<p>More text here</p>`,
		expectedText: `<p class="test">Text <i>here</i></p>
	<p>More text here</p>
`,
	},
	"make sure that blank lines with whitespace get removed": {
		inputText: `<p class="test">Text <i>here</i>      
	
    
		
	</p>
	<p>More text here</p>`,
		expectedText: `<p class="test">Text <i>here</i></p>
	<p>More text here</p>
`,
	},
}

func TestCleanupHtmlSpacing(t *testing.T) {
	for name, args := range cleanupHtmlSpacingTestCases {
		t.Run(name, func(t *testing.T) {
			actual := linter.CleanupHtmlSpacing(args.inputText)

			assert.Equal(t, args.expectedText, actual, "output text doesn't match")
		})
	}
}
