//go:build unit

package potentiallyfixableissue_test

import (
	"fmt"
	"testing"

	potentiallyfixableissue "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/potentially-fixable-issue"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			fmt.Sprintf("<p>%s</p>", contextBreak):                    potentiallyfixableissue.SectionBreakEl,
			fmt.Sprintf("<p><a id=\"pg10\"></a>%s</p>", contextBreak): "<p><a id=\"pg10\"></a>" + potentiallyfixableissue.SectionBreakEl + "</p>",
		},
	},
}

func TestGetPotentialSectionBreaks(t *testing.T) {
	t.Parallel()

	for name, args := range getPotentialSectionBreaksTestCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actual, err := potentiallyfixableissue.GetPotentialSectionBreaks(args.inputText, args.inputContextBreak)

			require.NoError(t, err)
			assert.Equal(t, args.expectedSuggestions, actual)
		})
	}
}
