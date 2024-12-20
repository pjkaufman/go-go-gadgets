//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

type getPotentialAlthoughButInstancesTestCase struct {
	inputText           string
	expectedSuggestions map[string]string
}

var getPotentialAlthoughButInstancesTestCases = map[string]getPotentialAlthoughButInstancesTestCase{
	"make sure that a file with no missing and's or or's without a comma proceeding it gets no suggestions": {
		inputText: `<p>Here is some content.</p>
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that a file with a missing comma before an and gets a suggestion": {
		inputText: `<p class="calibre1">This is exactly what Tatsuya was thinking now. Only he was released from school, and Miyuki couldn't miss school with him. </p>
		<p class="calibre1">If it will be a short period, less than a week, this has happened more than once. </p>
		<p class="calibre1"><a id="p57"></a>However, this time there was a chance that this could drag on for a month or more. Although he could "watch" her from afar, but he couldn't stop worrying that Miyuki and Minami would be left alone in this house. </p>
		<p class="calibre1">Maya on the screen looked at Miyuki standing next to Tatsuya, then diagonally behind her from Minami, then looked back at Tatsuya. </p>`,
		expectedSuggestions: map[string]string{
			`
		<p class="calibre1"><a id="p57"></a>However, this time there was a chance that this could drag on for a month or more. Although he could "watch" her from afar, but he couldn't stop worrying that Miyuki and Minami would be left alone in this house. </p>`: `
		<p class="calibre1"><a id="p57"></a>However, this time there was a chance that this could drag on for a month or more. Although he could "watch" her from afar, he couldn't stop worrying that Miyuki and Minami would be left alone in this house. </p>`,
		},
	},
}

func TestGetPotentialAlthoughButInstances(t *testing.T) {
	for name, args := range getPotentialAlthoughButInstancesTestCases {
		t.Run(name, func(t *testing.T) {
			actual := linter.GetPotentialAlthoughButInstances(args.inputText)

			assert.Equal(t, args.expectedSuggestions, actual)
		})
	}
}
