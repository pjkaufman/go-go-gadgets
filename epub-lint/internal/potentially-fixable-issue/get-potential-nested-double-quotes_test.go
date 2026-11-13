//go:build unit

package potentiallyfixableissue_test

import (
	"testing"

	potentiallyfixableissue "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/potentially-fixable-issue"
)

var getPotentialNestedDoubleQuotesTestCases = map[string]suggesterTestCase{
	"nested double quotes inside an outer double quote should be suggested": {
		inputText: `<p>"Obiwan said, "Hello Kenobi" to General Grievous."</p>`,
		expectedSuggestions: map[string]string{
			`<p>"Obiwan said, "Hello Kenobi" to General Grievous."</p>`: `<p>"Obiwan said, 'Hello Kenobi' to General Grievous."</p>`,
		},
	},
	"multiple nested quote phrases within one paragraph should be converted": {
		inputText: `<p>"Here lies the "old man" of Georgetown. May he  "rest well"."</p>`,
		expectedSuggestions: map[string]string{
			`<p>"Here lies the "old man" of Georgetown. May he  "rest well"."</p>`: `<p>"Here lies the 'old man' of Georgetown. May he  'rest well'."</p>`,
		},
	},
	"nested quotes with inline HTML should be converted without changing the markup": {
		inputText: `<p>"Obiwan said, <em>"Hello Kenobi"</em> to <strong>General Grievous</strong>."</p>`,
		expectedSuggestions: map[string]string{
			`<p>"Obiwan said, <em>"Hello Kenobi"</em> to <strong>General Grievous</strong>."</p>`: `<p>"Obiwan said, <em>'Hello Kenobi'</em> to <strong>General Grievous</strong>."</p>`,
		},
	},
	"quote characters in nested HTML attributes should not affect quote detection": {
		inputText: `<p>"He said, <span title="greeting">"Hello Kenobi"</span> to me."</p>`,
		expectedSuggestions: map[string]string{
			`<p>"He said, <span title="greeting">"Hello Kenobi"</span> to me."</p>`: `<p>"He said, <span title="greeting">'Hello Kenobi'</span> to me."</p>`,
		},
	},
	"separate quoted sentences with only spaces between them should not be suggested": {
		inputText:           `<p>"What is going on?" "Why are you here?"</p>`,
		expectedSuggestions: map[string]string{},
	},
}

func TestGetPotentialNestedDoubleQuotes(t *testing.T) {
	testSuggesterNoError(t, getPotentialNestedDoubleQuotesTestCases, potentiallyfixableissue.GetPotentialIncorrectNestedDoubleQuotes)
}
