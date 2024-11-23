//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

type getPotentialSquareBracketNecessaryWordsTestCase struct {
	inputText           string
	expectedSuggestions map[string]string
}

var getPotentialSquareBracketNecessaryWordsTestCases = map[string]getPotentialSquareBracketNecessaryWordsTestCase{
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
	for name, args := range getPotentialSquareBracketNecessaryWordsTestCases {
		t.Run(name, func(t *testing.T) {
			actual := linter.GetPotentialSquareBracketNecessaryWords(args.inputText)

			assert.Equal(t, args.expectedSuggestions, actual)
		})
	}
}
