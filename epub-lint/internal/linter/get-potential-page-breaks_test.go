//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

type getPotentialPageBreaksTestCase struct {
	inputText           string
	expectedSuggestions map[string]string
}

var getPotentialPageBreaksTestCases = map[string]getPotentialPageBreaksTestCase{
	"make sure that a file with no empty div or paragraph elements gives no suggestions": {
		inputText: `<p>Here is some content.</p>
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that a file with an empty paragraph gets that value as a suggestion": {
		inputText: `<p>Here is some content.</p>
		<p>    	</p>
		<p class="calibre1"><a id="p169"></a>If there are objects with a simple structure and the same properties, then they can be recognized as a single "set" allowing decomposition of the </p>
		<p class="calibre1">"set" rather than each object separately. </p>`,
		expectedSuggestions: map[string]string{
			`		<p>    	</p>`: linter.PageBrakeEl,
		},
	},
	"make sure that a file with an empty div gets that value as a suggestion": {
		inputText: `<p>Here is some content.</p>
		<div>    	</div>
		<p class="calibre1"><a id="p169"></a>If there are objects with a simple structure and the same properties, then they can be recognized as a single "set" allowing decomposition of the </p>
		<p class="calibre1">"set" rather than each object separately. </p>`,
		expectedSuggestions: map[string]string{
			`		<div>    	</div>`: linter.PageBrakeEl,
		},
	},
}

func TestGetPotentialPageBreaks(t *testing.T) {
	for name, args := range getPotentialPageBreaksTestCases {
		t.Run(name, func(t *testing.T) {
			actual := linter.GetPotentialPageBreaks(args.inputText)

			assert.Equal(t, args.expectedSuggestions, actual)
		})
	}
}
