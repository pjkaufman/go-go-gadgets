//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/linter"
	"github.com/stretchr/testify/assert"
)

type GetPotentialPageBreaksTestCase struct {
	InputText           string
	ExpectedSuggestions map[string]string
}

var GetPotentialPageBreaksTestCases = map[string]GetPotentialPageBreaksTestCase{
	"make sure that a file with no empty div or paragraph elements gives no suggestions": {
		InputText: `<p>Here is some content.</p>
<p>Here is some more content</p>`,
		ExpectedSuggestions: map[string]string{},
	},
	"make sure that a file with an empty paragraph gets that value as a suggestion": {
		InputText: `<p>Here is some content.</p>
		<p>    	</p>
		<p class="calibre1"><a id="p169"></a>If there are objects with a simple structure and the same properties, then they can be recognized as a single "set" allowing decomposition of the </p>
		<p class="calibre1">"set" rather than each object separately. </p>`,
		ExpectedSuggestions: map[string]string{
			`
		<p>    	</p>`: "\n" + linter.PageBrakeEl,
		},
	},
	"make sure that a file with an empty div gets that value as a suggestion": {
		InputText: `<p>Here is some content.</p>
		<div>    	</div>
		<p class="calibre1"><a id="p169"></a>If there are objects with a simple structure and the same properties, then they can be recognized as a single "set" allowing decomposition of the </p>
		<p class="calibre1">"set" rather than each object separately. </p>`,
		ExpectedSuggestions: map[string]string{
			`
		<div>    	</div>`: "\n" + linter.PageBrakeEl,
		},
	},
}

func TestGetPotentialPageBreaks(t *testing.T) {
	for name, args := range GetPotentialPageBreaksTestCases {
		t.Run(name, func(t *testing.T) {
			actual := linter.GetPotentialPageBreaks(args.InputText)

			assert.Equal(t, args.ExpectedSuggestions, actual)
		})
	}
}
