//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

type getPotentialThoughtInstancesTestCase struct {
	inputText           string
	expectedSuggestions map[string]string
}

var getPotentialThoughtInstancesTestCases = map[string]getPotentialThoughtInstancesTestCase{
	"make sure that a file with no parentheses gives no suggestions": {
		inputText: `<p>Here is some content.</p>
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that a file with a paragraph with parentheses gives a suggestion": {
		inputText: `<p>(I wonder why this is happening like this.)</p>
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{
			"<p>(I wonder why this is happening like this.)</p>": "<p><i>I wonder why this is happening like this.</i></p>",
		},
	},
	"make sure that a file with a paragraph with multiple instances of parentheses gives a suggestion": {
		inputText: `<p>(I wonder why this is happening like this.) (How come is this happening to me?)</p>
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{
			"<p>(I wonder why this is happening like this.) (How come is this happening to me?)</p>": "<p><i>I wonder why this is happening like this.</i> <i>How come is this happening to me?</i></p>",
		},
	},
	"make sure that a file with a paragraph with an instance of parentheses with other content on the same line gives a suggestion": {
		inputText: `<p>(I wonder why this is happening like this.) This is what John was thinking about.</p>
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{
			"<p>(I wonder why this is happening like this.) This is what John was thinking about.</p>": "<p><i>I wonder why this is happening like this.</i> This is what John was thinking about.</p>",
		},
	},
}

func TestGetPotentialThoughtInstances(t *testing.T) {
	for name, args := range getPotentialThoughtInstancesTestCases {
		t.Run(name, func(t *testing.T) {
			actual := linter.GetPotentialThoughtInstances(args.inputText)

			assert.Equal(t, args.expectedSuggestions, actual)
		})
	}
}
