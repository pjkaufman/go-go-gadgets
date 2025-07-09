//go:build unit

package linter_test

import (
	"fmt"
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

type getPotentialSectionBreaksTestCase struct {
	inputText           string
	inputContextBreak   string
	expectedSuggestions map[string]string
}

var getPotentialSectionBreaksTestCases = map[string]getPotentialSectionBreaksTestCase{
	"make sure that a file with no section breaks gives no suggestions": {
		inputText: `<p>Here is some content.</p>
<p>Here is some more content</p>`,
		inputContextBreak:   contextBreak,
		expectedSuggestions: map[string]string{},
	},
	"make sure that a file a couple section breaks gives suggestions": {
		inputText: fmt.Sprintf(`<p>Here is some content.</p>
<p>%[1]s</p>
<p><a id="pg10"></a>%[1]s</p>
<p>Here is some more content</p>`, contextBreak),
		inputContextBreak: contextBreak,
		expectedSuggestions: map[string]string{
			fmt.Sprintf("<p>%s</p>", contextBreak):                    linter.SectionBreakEl,
			fmt.Sprintf("<p><a id=\"pg10\"></a>%s</p>", contextBreak): "<p><a id=\"pg10\"></a>" + linter.SectionBreakEl + "</p>",
		},
	},
}

func TestGetPotentialSectionBreaks(t *testing.T) {
	for name, args := range getPotentialSectionBreaksTestCases {
		t.Run(name, func(t *testing.T) {
			actual := linter.GetPotentialSectionBreaks(args.inputText, args.inputContextBreak)

			assert.Equal(t, args.expectedSuggestions, actual)
		})
	}
}
