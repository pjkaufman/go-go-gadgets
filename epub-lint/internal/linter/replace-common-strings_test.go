//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/linter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		input:    "<p>-- test --</p>",
		expected: "<p>— test —</p>",
	},
	"make sure that three periods with 0 spaces between them get cut down to proper ellipsis": {
		input: `
		  <p>...</p>
		`,
		expected: `
		  <p>…</p>
		`,
	},
	"make sure that an uppercase 'Sneaked' results in an uppercase 'Snuck'": {
		input:    "<p>Sneaked</p>",
		expected: "<p>Snuck</p>",
	},
	"make sure that a lowercase 'snuck' results in a lowercase 'snuck'": {
		input:    "<p>On his way he sneaked out the door</p>",
		expected: "<p>On his way he snuck out the door</p>",
	},
	"make sure that words with 2 or more spaces between them have the multiple spaces cut down to 1": {
		input:    "<p>This  is an    interestingly spaced   sentence.  See the multiple    blanks?</p>",
		expected: "<p>This is an interestingly spaced sentence. See the multiple blanks?</p>",
	},
	"make sure that spacing before a paragraph tag is not removed": {
		input:    "  <p>This  is an    interestingly spaced   sentence.  See the multiple    blanks?</p>",
		expected: "  <p>This is an interestingly spaced sentence. See the multiple blanks?</p>",
	},
	"make sure that smart double quotes are replaced with straight quotes</p>": {
		input: `<p>“Hey. How are you?”</p>
		<p>“I am doing great!”</p>`,
		expected: `<p>"Hey. How are you?"</p>
		<p>"I am doing great!"</p>`,
	},
	"make sure that smart single quotes are replaced with straight quotes": {
		input: `<p>‘Hey. How are you?’</p>
		<p>‘I am doing great!’</p>`,
		expected: `<p>'Hey. How are you?'</p>
		<p>'I am doing great!'</p>`,
	},
	"make sure that markup is left untouched": {
		input:    `  <p class='foo' data-test="bar">Sneaked</p>`,
		expected: `  <p class='foo' data-test="bar">Snuck</p>`,
	},
	"make sure that self-closing tags do not affect replacements": {
		input:    `<p>Before<br/>Sneaked</p>`,
		expected: `<p>Before<br/>Snuck</p>`,
	},
	"make sure that an href in the code is not affected by the replacements": {
		input:    `<p>Here is some <a href="hello--worldSneaked.txt">Text</a></p>`,
		expected: `<p>Here is some <a href="hello--worldSneaked.txt">Text</a></p>`,
	},
}

func TestCommonStringReplace(t *testing.T) {
	t.Parallel()

	for name, args := range commonStringReplaceTestCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actual, err := linter.CommonStringReplace(args.input)

			require.NoError(t, err)

			assert.Equal(t, args.expected, actual)
		})
	}
}
