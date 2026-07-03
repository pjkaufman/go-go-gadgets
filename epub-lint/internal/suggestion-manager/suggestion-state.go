package suggestionmanager

import (
	"strings"

	stringdiff "github.com/pjkaufman/go-go-gadgets/pkg/string-diff"
)

// SuggestionState represents the state of a single suggestion.
type SuggestionState struct {
	IsAccepted                           bool
	OriginallyHadHalfwidthCircleKatakana bool
	Original                             string
	OriginalSuggestion                   string
	CurrentSuggestion                    string
	Display                              string
}

func (s *SuggestionState) GetStringDiffAsDisplay() error {
	var err error
	s.Display, err = stringdiff.GetPrettyDiffString(strings.TrimLeft(s.Original, "\n"), strings.TrimLeft(s.CurrentSuggestion, "\n"))

	return err
}

// replacing in reverse is not guaranteed to work correctly, but for now this works.
// This can be changed if necessary down the road.
func (s *SuggestionState) undoReplaceBrokenDisplayCharacters() {
	if !s.OriginallyHadHalfwidthCircleKatakana {
		return
	}

	s.CurrentSuggestion = strings.ReplaceAll(s.CurrentSuggestion, "°", "ﾟ")
	s.OriginallyHadHalfwidthCircleKatakana = false
}

// replaces the provided text's broken display characters and returns a string that should display fine in the terminal.
// This should note be used with text that is not directly related to the suggestion in question as it also sets
// the flag for whether or not there was halfwidth circle katakana present.
func (s *SuggestionState) ReplaceBrokenDisplayCharacters(text string) string {
	// text with handakuten in them are not having their width calculated correctly, so I will just remove them
	// and we can display a warning if need bee
	if strings.Contains(text, "ﾟ") {
		s.OriginallyHadHalfwidthCircleKatakana = true

		return strings.ReplaceAll(text, "ﾟ", "°")
	}

	s.OriginallyHadHalfwidthCircleKatakana = false

	return text
}
