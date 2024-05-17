//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/linter"
	"github.com/stretchr/testify/assert"
)

type GetPotentialSquareBracketConversationInstancesTestCase struct {
	InputText           string
	ExpectedSuggestions map[string]string
}

var GetPotentialSquareBracketConversationInstancesTestCases = map[string]GetPotentialSquareBracketConversationInstancesTestCase{
	"make sure that a file with no square brackets gives no suggestions": {
		InputText: `<p>Here is some content.</p>
<p>Here is some more content</p>`,
		ExpectedSuggestions: map[string]string{},
	},
	"make sure that a file with a paragraph with square brackets in it, but not for all of the content except whitespace has no matches": {
		InputText:           `<p>Meanwhile, he casts [Gram Dispersion] on the knife acting as a shield. The remote control magic was erased, and the knife fell into the weeds. </p>`,
		ExpectedSuggestions: map[string]string{},
	},
	"make sure that a file with a paragraph with all of its non-space contents contained in square brackets gives a suggestion": {
		InputText: `<p>[I'll take you to him.] </p>
<p>Here is some more content</p>`,
		ExpectedSuggestions: map[string]string{
			"<p>[I'll take you to him.] </p>": "<p>\"I'll take you to him.\" </p>",
		},
	},
	"make sure that a file with a paragraph with all of its non-space contents contained in square brackets with another set of square bracketed text in it gives a suggestion": {
		InputText: `<p>[I'll take you [to] him.] </p>
<p>Here is some more content</p>`,
		ExpectedSuggestions: map[string]string{
			"<p>[I'll take you [to] him.] </p>": "<p>\"I'll take you [to] him.\" </p>",
		},
	},
}

func TestGetPotentialSquareBracketConversationInstances(t *testing.T) {
	for name, args := range GetPotentialSquareBracketConversationInstancesTestCases {
		t.Run(name, func(t *testing.T) {
			actual := linter.GetPotentialSquareBracketConversationInstances(args.InputText)

			assert.Equal(t, args.ExpectedSuggestions, actual)
		})
	}
}
