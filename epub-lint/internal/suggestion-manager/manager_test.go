//go:build unit

package suggestionmanager_test

import (
	"testing"

	_ "embed"

	potentiallyfixableissue "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/potentially-fixable-issue"
	sm "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/suggestion-manager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/**
I am working on creating a suggestion manager. The latest should be on the branch for the issue. It has a file in there called epub-lint/internal/suggestion-manager/manager.go on the suggestion-manager branch. I would like you to go ahead and create the following test cases:
- Initialization:
  - Files are properly sorted
  - Files already sorted stay in sorted order
  - Passed in properties are as passed in on the model
*/

var (
	//go:embed testdata/APotentialConversationInstance.html
	simplePotentialConversationInstanceHtml string
	//go:embed testdata/empty.html
	noPotentialIssuesHtml string
	//go:embed testdata/APotentialWordOmmisionInstance.html
	simplePotentialWordOmissionInstanceHtml string
)

// Test case structures
type suggestionManagerTestCase struct {
	suggestions []potentiallyfixableissue.PotentiallyFixableIssue
	files       []sm.FileSuggestionInfo
	runAll      bool
	skipCss     bool
	setupFunc   func(*sm.SuggestionManager)
	assertions  func(*testing.T, *sm.SuggestionManager)
}

func createPotentiallyFixableIssues(t *testing.T, contextBreak string) []potentiallyfixableissue.PotentiallyFixableIssue {
	t.Helper()

	return []potentiallyfixableissue.PotentiallyFixableIssue{
		{
			Name:           "Potential Conversation Instances",
			GetSuggestions: potentiallyfixableissue.GetPotentialSquareBracketConversationInstances,
			IsEnabled:      pointerToBool(false),
		},
		{
			Name:           "Potential Necessary Word Omission Instances",
			GetSuggestions: potentiallyfixableissue.GetPotentialSquareBracketNecessaryWords,
			IsEnabled:      pointerToBool(false),
		},
		{
			Name:           "Potential Broken Lines",
			GetSuggestions: potentiallyfixableissue.GetPotentiallyBrokenLines,
			IsEnabled:      pointerToBool(false),
		},
		{
			Name:           "Potential Incorrect Single Quotes",
			GetSuggestions: potentiallyfixableissue.GetPotentialIncorrectSingleQuotes,
			IsEnabled:      pointerToBool(false),
		},
		{
			Name: "Potential Section Breaks",
			// wrapper here allows calling the get potential section breaks logic without needing to change the function definition
			GetSuggestions: func(text string) (map[string]string, error) {
				return potentiallyfixableissue.GetPotentialSectionBreaks(text, contextBreak)
			},
			IsEnabled:                   pointerToBool(false),
			UpdateAllInstances:          true,
			AddCssSectionBreakIfMissing: true,
		},
		{
			Name:                     "Potential Page Breaks",
			GetSuggestions:           potentiallyfixableissue.GetPotentialPageBreaks,
			IsEnabled:                pointerToBool(false),
			UpdateAllInstances:       true,
			AddCssPageBreakIfMissing: true,
		},
		{
			Name:           "Potential Missing Oxford Commas",
			GetSuggestions: potentiallyfixableissue.GetPotentialMissingOxfordCommas,
			IsEnabled:      pointerToBool(false),
		},
		{
			Name:           "Potentially Lacking Subordinate Clause Instances",
			GetSuggestions: potentiallyfixableissue.GetPotentiallyLackingSubordinateClauseInstances,
			IsEnabled:      pointerToBool(false),
		},
		{
			Name:           "Potential Thought Instances",
			GetSuggestions: potentiallyfixableissue.GetPotentialThoughtInstances,
			IsEnabled:      pointerToBool(false),
		},
	}
}

func pointerToBool(value bool) *bool {
	return &value
}

type suggestionManagerSetupForNextSuggestionsTestCase struct {
	filePathsToText           map[string]string
	runAll                    bool
	skipCss                   bool
	expectedCurrentFileName   string
	expectedCurrentIssueIndex int
	expectedFoundSuggestion   bool
}

var suggestionManagerSetupForNextSuggestionsTestCases = map[string]suggestionManagerSetupForNextSuggestionsTestCase{
	"When the first file and first potentially fixable issue has a potential issue, it should be found and suggested": {
		filePathsToText: map[string]string{
			"simple.html": simplePotentialConversationInstanceHtml,
		},
		runAll:                  true,
		expectedFoundSuggestion: true,
		expectedCurrentFileName: "simple.html",
	},
	"When the first file and second potentially fixable issue has a potential issue, it should be found and suggested": {
		filePathsToText: map[string]string{
			"simple.html": simplePotentialWordOmissionInstanceHtml,
		},
		runAll:                    true,
		expectedCurrentIssueIndex: 1,
		expectedFoundSuggestion:   true,
		expectedCurrentFileName:   "simple.html",
	},
	"When there is a potentially fixable issue, but the logic for it is disabled, it should not return that suggestion": {
		filePathsToText: map[string]string{
			"simple.html": simplePotentialConversationInstanceHtml,
		},
		runAll:                  false,
		expectedFoundSuggestion: false,
		expectedCurrentFileName: "simple.html",
	},
	"When the first file has no potential issue, but the second does, it should be found and suggested": {
		filePathsToText: map[string]string{
			"empty.html":  noPotentialIssuesHtml,
			"simple.html": simplePotentialWordOmissionInstanceHtml,
		},
		runAll:                    true,
		expectedCurrentIssueIndex: 1,
		expectedFoundSuggestion:   true,
		expectedCurrentFileName:   "simple.html",
	},
	"When there is a single file and it has no potentially fixable issues, it should return that no suggestions were found": {
		filePathsToText: map[string]string{
			"empty.html": noPotentialIssuesHtml,
		},
		runAll:                  true,
		expectedFoundSuggestion: false,
		expectedCurrentFileName: "empty.html",
	},
	"When there are multiple files and they have no potentially fixable issues, it should return that no suggestions were found": {
		filePathsToText: map[string]string{
			"dont.html":  noPotentialIssuesHtml,
			"empty.html": noPotentialIssuesHtml,
			"full.html":  noPotentialIssuesHtml,
		},
		runAll:                  true,
		expectedFoundSuggestion: false,
		expectedCurrentFileName: "full.html",
	},
}

func TestSuggestionManager(t *testing.T) {
	t.Parallel()

	t.Run("SetupForNextSuggestions", func(t *testing.T) {
		t.Parallel()

		for name, tc := range suggestionManagerSetupForNextSuggestionsTestCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				manager := sm.NewSuggestionManager(createPotentiallyFixableIssues(t, "----"), tc.filePathsToText, tc.runAll, tc.skipCss, nil)

				foundSuggestion, err := manager.SetupForNextSuggestions()

				require.NoError(t, err)
				assert.Equal(t, tc.expectedFoundSuggestion, foundSuggestion)
				assert.Equal(t, tc.expectedCurrentIssueIndex, manager.CurrentIssueIndex, "Current issue index did not match")
				assert.Equal(t, 0, manager.CurrentSuggestionIndex, "Current suggestion index was not 0 for some reason...")
				assert.Equal(t, tc.expectedCurrentFileName, manager.CurrentFileName, "Current file name was not the expected file name")
			})
		}
	})
}
