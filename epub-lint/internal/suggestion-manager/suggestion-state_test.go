//go:build unit

package suggestionmanager

import (
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type getStringDiffAsDisplayTestCase struct {
	original          string
	currentSuggestion string

	expectedDisplay string
}

type replaceBrokenDisplayCharactersTestCase struct {
	inputText string

	expectedText string
	expectedFlag bool
}

type undoReplaceBrokenDisplayCharactersTestCase struct {
	currentSuggestion string
	originalFlag      bool

	expectedSuggestion string
	expectedFlag       bool
}

var getStringDiffAsDisplayTestCases = map[string]getStringDiffAsDisplayTestCase{
	"latin1 characters are repaired correctly": {
		original:          "Original…",
		currentSuggestion: "Original….",
		expectedDisplay:   "Original….",
	},
}

var replaceBrokenDisplayCharactersTestCases = map[string]replaceBrokenDisplayCharactersTestCase{
	"no halfwidth katakana leaves text unchanged": {
		inputText:    "Temperature is 20°.",
		expectedText: "Temperature is 20°.",
		expectedFlag: false,
	},
	"halfwidth katakana is replaced with degree symbol": {
		inputText:    "ﾊﾟ",
		expectedText: "ﾊ°",
		expectedFlag: true,
	},
}

var undoReplaceBrokenDisplayCharactersTestCases = map[string]undoReplaceBrokenDisplayCharactersTestCase{
	"does nothing when original did not contain halfwidth katakana": {
		currentSuggestion:  "Temperature is 20°.",
		originalFlag:       false,
		expectedSuggestion: "Temperature is 20°.",
		expectedFlag:       false,
	},
	"replaces degree symbols back to halfwidth katakana and clears flag": {
		currentSuggestion:  "ﾊ°",
		originalFlag:       true,
		expectedSuggestion: "ﾊﾟ",
		expectedFlag:       false,
	},
}

func TestSuggestionState(t *testing.T) {
	t.Parallel()

	t.Run("GetStringDiffAsDisplay", func(t *testing.T) {
		t.Parallel()

		for name, tc := range getStringDiffAsDisplayTestCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				state := SuggestionState{
					Original:          tc.original,
					CurrentSuggestion: tc.currentSuggestion,
				}

				err := state.GetStringDiffAsDisplay()

				require.NoError(t, err)
				assert.Equal(t, tc.expectedDisplay, ansi.Strip(state.Display))
			})
		}
	})

	t.Run("ReplaceBrokenDisplayCharacters", func(t *testing.T) {
		t.Parallel()

		for name, tc := range replaceBrokenDisplayCharactersTestCases {
			tc := tc

			t.Run(name, func(t *testing.T) {
				t.Parallel()

				var state SuggestionState

				actual := state.ReplaceBrokenDisplayCharacters(tc.inputText)

				assert.Equal(t, tc.expectedText, actual)
				assert.Equal(t, tc.expectedFlag, state.OriginallyHadHalfwidthCircleKatakana)
			})
		}
	})

	t.Run("undoReplaceBrokenDisplayCharacters", func(t *testing.T) {
		t.Parallel()

		for name, tc := range undoReplaceBrokenDisplayCharactersTestCases {
			tc := tc

			t.Run(name, func(t *testing.T) {
				t.Parallel()

				state := SuggestionState{
					CurrentSuggestion:                    tc.currentSuggestion,
					OriginallyHadHalfwidthCircleKatakana: tc.originalFlag,
				}

				state.undoReplaceBrokenDisplayCharacters()

				assert.Equal(t, tc.expectedSuggestion, state.CurrentSuggestion)
				assert.Equal(t, tc.expectedFlag, state.OriginallyHadHalfwidthCircleKatakana)
			})
		}
	})
}
