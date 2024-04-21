//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/linter"
	"github.com/stretchr/testify/assert"
)

type GetPotentialThoughtInstancesTestCase struct {
	InputText           string
	ExpectedSuggestions map[string]string
}

var GetPotentialThoughtInstancesTestCases = map[string]GetPotentialThoughtInstancesTestCase{
	"make sure that a file with no parentheses gives no suggestions": {
		InputText: `<p>Here is some content.</p>
<p>Here is some more content</p>`,
		ExpectedSuggestions: map[string]string{},
	},
	"make sure that a file with a paragraph with parentheses gives a suggestion": {
		InputText: `<p>(I wonder why this is happening like this.)</p>
<p>Here is some more content</p>`,
		ExpectedSuggestions: map[string]string{
			"<p>(I wonder why this is happening like this.)</p>": "<p><i>I wonder why this is happening like this.</i></p>",
		},
	},
	"make sure that a file with a paragraph with multiple instances of parentheses gives a suggestion": {
		InputText: `<p>(I wonder why this is happening like this.) (How come is this happening to me?)</p>
<p>Here is some more content</p>`,
		ExpectedSuggestions: map[string]string{
			"<p>(I wonder why this is happening like this.) (How come is this happening to me?)</p>": "<p><i>I wonder why this is happening like this.</i> <i>How come is this happening to me?</i></p>",
		},
	},
	"make sure that a file with a paragraph with an instance of parentheses with other content on the same line gives a suggestion": {
		InputText: `<p>(I wonder why this is happening like this.) This is what John was thinking about.</p>
<p>Here is some more content</p>`,
		ExpectedSuggestions: map[string]string{
			"<p>(I wonder why this is happening like this.) This is what John was thinking about.</p>": "<p><i>I wonder why this is happening like this.</i> This is what John was thinking about.</p>",
		},
	},
}

func TestGetPotentialThoughtInstances(t *testing.T) {
	for name, args := range GetPotentialThoughtInstancesTestCases {
		t.Run(name, func(t *testing.T) {
			actual := linter.GetPotentialThoughtInstances(args.InputText)

			assert.Equal(t, args.ExpectedSuggestions, actual)
		})
	}
}
