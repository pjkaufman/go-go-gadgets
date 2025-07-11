//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

type commonStringReplaceTestCase struct {
	input    string
	expected string
}

var commonStringReplaceTestCases = map[string]commonStringReplaceTestCase{
	"make sure that html comments are left alone": {
		input:    "<!--this is a comment. comments are not displayed in the browser-->",
		expected: "<!--this is a comment. comments are not displayed in the browser-->",
	},
	"make sure that two en dashes are replaced with an em dash": {
		input:    "-- test --",
		expected: "— test —",
	},
	"make sure that three periods with 0 spaces between them get cut down to proper ellipsis": {
		input: `
		  ...
		`,
		expected: `
		  …
		`,
	},
	"make sure that an uppercase 'Sneaked' results in an uppercase 'Snuck'": {
		input:    "Sneaked",
		expected: "Snuck",
	},
	"make sure that a lowercase 'snuck' results in a lowercase 'snuck'": {
		input:    "On his way he sneaked out the door",
		expected: "On his way he snuck out the door",
	},
	"make sure that words with 2 or more spaces between them have the multiple spaces cut down to 1": {
		input:    "This  is an    interestingly spaced   sentence.  See the multiple    blanks?",
		expected: "This is an interestingly spaced sentence. See the multiple blanks?",
	},
	"make sure that spacing before a paragraph tag is not removed": {
		input:    "  <p>This  is an    interestingly spaced   sentence.  See the multiple    blanks?</p>",
		expected: "  <p>This is an interestingly spaced sentence. See the multiple blanks?</p>",
	},
	"make sure that smart double quotes are replaced with straight quotes</p>": {
		input: `“Hey. How are you?”
		“I am doing great!”`,
		expected: `"Hey. How are you?"
		"I am doing great!"`,
	},
	"make sure that smart single quotes are replaced with straight quotes": {
		input: `‘Hey. How are you?’
		‘I am doing great!’`,
		expected: `'Hey. How are you?'
		'I am doing great!'`,
	},
}

func TestCommonStringReplace(t *testing.T) {
	for name, args := range commonStringReplaceTestCases {
		t.Run(name, func(t *testing.T) {
			actual := linter.CommonStringReplace(args.input)

			assert.Equal(t, args.expected, actual)
		})
	}

}
