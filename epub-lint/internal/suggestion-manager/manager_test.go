//go:build unit

package suggestionmanager_test

import (
	"testing"

	potentiallyfixableissue "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/potentially-fixable-issue"
	sm "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/suggestion-manager"
)

// Test case structures
type suggestionManagerTestCase struct {
	name        string
	suggestions []potentiallyfixableissue.PotentiallyFixableIssue
	files       []sm.FileSuggestionInfo
	runAll      bool
	skipCss     bool
	setupFunc   func(*sm.SuggestionManager)
	assertions  func(*testing.T, *sm.SuggestionManager)
}
