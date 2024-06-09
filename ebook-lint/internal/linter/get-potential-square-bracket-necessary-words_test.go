//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

type GetPotentialSquareBracketNecessaryWordsTestCase struct {
	InputText           string
	ExpectedSuggestions map[string]string
}

var GetPotentialSquareBracketNecessaryWordsTestCases = map[string]GetPotentialSquareBracketNecessaryWordsTestCase{
	"make sure that a file with no square brackets gives no suggestions": {
		InputText: `<p>Here is some content.</p>
<p>Here is some more content</p>`,
		ExpectedSuggestions: map[string]string{},
	},
	"make sure that a file with a paragraph with its contents entirely contained in square brackets gives no suggestion": {
		InputText: `<p>[I wonder why this is happening like this.]</p>
<p>Here is some more content</p>`,
		ExpectedSuggestions: map[string]string{},
	},
	"make sure that a file with a paragraph with its contents entirely contained in square brackets and whitespace after the ending bracket gives no suggestion": {
		InputText: `<p>[I wonder why this is happening like this.]   </p>
<p>Here is some more content</p>`,
		ExpectedSuggestions: map[string]string{},
	},
	"make sure that a file with a paragraph with its contents entirely contained in square brackets and whitespace before the opening bracket gives no suggestion": {
		InputText: `<p>  [I wonder why this is happening like this.]</p>
<p>Here is some more content</p>`,
		ExpectedSuggestions: map[string]string{},
	},
	"make sure that a file with a paragraph with square brackets for some of the words in a paragraph gives a suggestion": {
		InputText: `<p> I [wonder why] this is happening like this.</p>
<p>Here is some more content</p>`,
		ExpectedSuggestions: map[string]string{
			"<p> I [wonder why] this is happening like this.</p>": "<p> I wonder why this is happening like this.</p>",
		},
	},
	"make sure that a file with a paragraph with multiple square brackets for some of the words in a paragraph gives a suggestion": {
		InputText: `<p> I [wonder why] this is [happening] like this.</p>
<p>Here is some more content</p>`,
		ExpectedSuggestions: map[string]string{
			"<p> I [wonder why] this is [happening] like this.</p>": "<p> I wonder why this is happening like this.</p>",
		},
	},
}

func TestGetPotentialSquareBracketNecessaryWords(t *testing.T) {
	for name, args := range GetPotentialSquareBracketNecessaryWordsTestCases {
		t.Run(name, func(t *testing.T) {
			actual := linter.GetPotentialSquareBracketNecessaryWords(args.InputText)

			assert.Equal(t, args.ExpectedSuggestions, actual)
		})
	}
}
