//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

type CommonStringReplaceTestCase struct {
	Input    string
	Expected string
}

var commonStringReplaceTestCases = map[string]CommonStringReplaceTestCase{
	"make sure that html comments are left alone": {
		Input:    "<!--this is a comment. comments are not displayed in the browser-->",
		Expected: "<!--this is a comment. comments are not displayed in the browser-->",
	},
	"make sure that two en dashes are replaced with an em dash": {
		Input:    "-- test --",
		Expected: "— test —",
	},
	"make sure that three periods with 0 spaces between them get cut down to proper ellipsis": {
		Input: `
		  ...
		`,
		Expected: `
		  …
		`,
	},
	"make sure that an uppercase 'Sneaked' results in an uppercase 'Snuck'": {
		Input:    "Sneaked",
		Expected: "Snuck",
	},
	"make sure that a lowercase 'snuck' results in a lowercase 'snuck'": {
		Input:    "On his way he sneaked out the door",
		Expected: "On his way he snuck out the door",
	},
	"make sure that words with 2 or more spaces between them have the multiple spaces cut down to 1": {
		Input:    "This  is an    interestingly spaced   sentence.  See the multiple    blanks?",
		Expected: "This is an interestingly spaced sentence. See the multiple blanks?",
	},
	"make sure that spacing before a paragraph tag is not removed": {
		Input:    "  <p>This  is an    interestingly spaced   sentence.  See the multiple    blanks?</p>",
		Expected: "  <p>This is an interestingly spaced sentence. See the multiple blanks?</p>",
	},
	"make sure that smart double quotes are replaced with straight quotes</p>": {
		Input: `“Hey. How are you?”
		“I am doing great!”`,
		Expected: `"Hey. How are you?"
		"I am doing great!"`,
	},
	"make sure that smart single quotes are replaced with straight quotes": {
		Input: `‘Hey. How are you?’
		‘I am doing great!’`,
		Expected: `'Hey. How are you?'
		'I am doing great!'`,
	},
}

func TestCommonStringReplace(t *testing.T) {
	for name, args := range commonStringReplaceTestCases {
		t.Run(name, func(t *testing.T) {
			actual := linter.CommonStringReplace(args.Input)

			assert.Equal(t, args.Expected, actual)
		})
	}

}
