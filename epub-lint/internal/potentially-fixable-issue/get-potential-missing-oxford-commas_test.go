//go:build unit

package potentiallyfixableissue_test

import (
	"testing"

	potentiallyfixableissue "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/potentially-fixable-issue"
)

var getPotentialMissingOxfordCommasTestCases = map[string]suggesterTestCase{
	"make sure that a file with no missing and's or or's without a comma proceeding it gets no suggestions": {
		inputText: `<p>Here is some content.</p>
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that a file with a missing comma before an and gets a suggestion": {
		inputText: `<p>Here is some content.</p>
		<p>Here is a situation where I run, skip and jump for a long time.</p>`,
		expectedSuggestions: map[string]string{
			`		<p>Here is a situation where I run, skip and jump for a long time.</p>`: `		<p>Here is a situation where I run, skip, and jump for a long time.</p>`,
		},
	},
	"make sure that a file with a missing comma before an or gets a suggestion": {
		inputText: `<p>Here is some content.</p>
		<p>Here is a situation where I run, skip or jump for a long time.</p>`,
		expectedSuggestions: map[string]string{
			`		<p>Here is a situation where I run, skip or jump for a long time.</p>`: `		<p>Here is a situation where I run, skip, or jump for a long time.</p>`,
		},
	},
}

func TestGetPotentialMissingOxfordCommas(t *testing.T) {
	testSuggesterNoError(t, getPotentialMissingOxfordCommasTestCases, potentiallyfixableissue.GetPotentialMissingOxfordCommas)
}
