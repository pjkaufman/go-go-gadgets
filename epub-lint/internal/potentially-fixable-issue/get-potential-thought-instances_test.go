//go:build unit

package potentiallyfixableissue_test

import (
	"testing"

	potentiallyfixableissue "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/potentially-fixable-issue"
)

var getPotentialThoughtInstancesTestCases = map[string]suggesterTestCase{
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
	"make sure that a file with a paragraph with parentheses inside an italics tag gives a suggestion without italics added": {
		inputText: `<p><i>(I wonder why this is happening like this.)</i></p>
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{
			"<p><i>(I wonder why this is happening like this.)</i></p>": "<p><i>I wonder why this is happening like this.</i></p>",
		},
	},
	"make sure that a file with a paragraph with parentheses inside an emphasis tag gives a suggestion without italics added": {
		inputText: `<p><em>(I wonder why this is happening like this.)</em></p>	
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{
			"<p><em>(I wonder why this is happening like this.)</em></p>": "<p><em>I wonder why this is happening like this.</em></p>",
		},
	},
	"make sure that a file with a paragraph with parentheses enclosing an italics tag gives a suggestion without italics added": {
		inputText: `<p>(<i>I wonder why this is happening like this.</i>)</p>
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{
			"<p>(<i>I wonder why this is happening like this.</i>)</p>": "<p><i>I wonder why this is happening like this.</i></p>",
		},
	},
	"make sure that a file with a paragraph with parentheses enclosing an emphasis tag gives a suggestion without italics added": {
		inputText: `<p>(<em>I wonder why this is happening like this.</em>)</p>	
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{
			"<p>(<em>I wonder why this is happening like this.</em>)</p>": "<p><em>I wonder why this is happening like this.</em></p>",
		},
	},
	"make sure that a file with a paragraph with parentheses with an italics tag before both parentheses gives a suggestion without italics added": {
		inputText: `<p><i>(I wonder why this is happening like this.</i>)</p>
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{
			"<p><i>(I wonder why this is happening like this.</i>)</p>": "<p><i>I wonder why this is happening like this.</i></p>",
		},
	},
	"make sure that a file with a paragraph with parentheses with an emphasis tag before both parentheses gives a suggestion without italics added": {
		inputText: `<p><em>(I wonder why this is happening like this.</em>)</p>	
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{
			"<p><em>(I wonder why this is happening like this.</em>)</p>": "<p><em>I wonder why this is happening like this.</em></p>",
		},
	},
	"make sure that a file with a paragraph with parentheses with an italics tag after both parentheses gives a suggestion without italics added": {
		inputText: `<p>(<i>I wonder why this is happening like this.)</i></p>
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{
			"<p>(<i>I wonder why this is happening like this.)</i></p>": "<p><i>I wonder why this is happening like this.</i></p>",
		},
	},
	"make sure that a file with a paragraph with parentheses with an emphasis tag after both parentheses gives a suggestion without italics added": {
		inputText: `<p>(<em>I wonder why this is happening like this.)</em></p>	
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{
			"<p>(<em>I wonder why this is happening like this.)</em></p>": "<p><em>I wonder why this is happening like this.</em></p>",
		},
	},
	"make sure that a file with a paragraph with parentheses with an opening italics tag and closing emphasis tag adds the opening and closing italics as part of the suggestion": {
		inputText: `<p><em>(I wonder why this is happening like this.)</i></p>	
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{
			"<p><em>(I wonder why this is happening like this.)</i></p>": "<p><em><i>I wonder why this is happening like this.</i></i></p>",
		},
	},
}

func TestGetPotentialThoughtInstances(t *testing.T) {
	testSuggesterNoError(t, getPotentialThoughtInstancesTestCases, potentiallyfixableissue.GetPotentialThoughtInstances)
}
