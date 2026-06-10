//go:build unit

package potentiallyfixableissue_test

import (
	"testing"

	potentiallyfixableissue "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/potentially-fixable-issue"
)

var getPotentialSquareBracketNecessaryWordsTestCases = map[string]suggesterTestCase{
	"make sure that a file with no square brackets gives no suggestions": {
		inputText: `<p>Here is some content.</p>
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that a file with a paragraph with its contents entirely contained in square brackets gives no suggestion": {
		inputText: `<p>[I wonder why this is happening like this.]</p>
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that a file with a paragraph with its contents entirely contained in square brackets and whitespace after the ending bracket gives no suggestion": {
		inputText: `<p>[I wonder why this is happening like this.]   </p>
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that a file with a paragraph with its contents entirely contained in square brackets and whitespace before the opening bracket gives no suggestion": {
		inputText: `<p>  [I wonder why this is happening like this.]</p>
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that a file with a paragraph with square brackets for some of the words in a paragraph gives a suggestion": {
		inputText: `<p> I [wonder why] this is happening like this.</p>
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{
			"<p> I [wonder why] this is happening like this.</p>": "<p> I wonder why this is happening like this.</p>",
		},
	},
	"make sure that a file with a paragraph with multiple square brackets for some of the words in a paragraph gives a suggestion": {
		inputText: `<p> I [wonder why] this is [happening] like this.</p>
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{
			"<p> I [wonder why] this is [happening] like this.</p>": "<p> I wonder why this is happening like this.</p>",
		},
	},
}

func TestGetPotentialSquareBracketNecessaryWords(t *testing.T) {
	testSuggesterNoError(t, getPotentialSquareBracketNecessaryWordsTestCases, potentiallyfixableissue.GetPotentialSquareBracketNecessaryWords)
}
