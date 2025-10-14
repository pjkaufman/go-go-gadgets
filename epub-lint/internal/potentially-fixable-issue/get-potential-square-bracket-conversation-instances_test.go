//go:build unit

package potentiallyfixableissue_test

import (
	"testing"

	potentiallyfixableissue "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/potentially-fixable-issue"
	"github.com/stretchr/testify/assert"
)

type getPotentialSquareBracketConversationInstancesTestCase struct {
	inputText           string
	expectedSuggestions map[string]string
}

var getPotentialSquareBracketConversationInstancesTestCases = map[string]getPotentialSquareBracketConversationInstancesTestCase{
	"make sure that a file with no square brackets gives no suggestions": {
		inputText: `<p>Here is some content.</p>
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that a file with a paragraph with square brackets in it, but not for all of the content except whitespace has no matches": {
		inputText:           `<p>Meanwhile, he casts [Gram Dispersion] on the knife acting as a shield. The remote control magic was erased, and the knife fell into the weeds. </p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that a file with a paragraph with all of its non-space contents contained in square brackets gives a suggestion": {
		inputText: `<p>[I'll take you to him.] </p>
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{
			"<p>[I'll take you to him.] </p>": "<p>\"I'll take you to him.\" </p>",
		},
	},
	"make sure that a file with a paragraph with all of its non-space contents contained in square brackets with another set of square bracketed text in it gives a suggestion": {
		inputText: `<p>[I'll take you [to] him.] </p>
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{
			"<p>[I'll take you [to] him.] </p>": "<p>\"I'll take you [to] him.\" </p>",
		},
	},
	"make sure that a file with a paragraph with an empty anchor tag at the start and all of its non-whitespace content in the square bracket gives a suggestion": {
		inputText: `<p><a id="something" href=""></a>[I'll take you [to] him.] </p>`,
		expectedSuggestions: map[string]string{
			"<p><a id=\"something\" href=\"\"></a>[I'll take you [to] him.] </p>": "<p><a id=\"something\" href=\"\"></a>\"I'll take you [to] him.\" </p>",
		},
	},
}

func TestGetPotentialSquareBracketConversationInstances(t *testing.T) {
	for name, args := range getPotentialSquareBracketConversationInstancesTestCases {
		t.Run(name, func(t *testing.T) {
			actual, err := potentiallyfixableissue.GetPotentialSquareBracketConversationInstances(args.inputText)

			assert.Nil(t, err)
			assert.Equal(t, args.expectedSuggestions, actual)
		})
	}
}
