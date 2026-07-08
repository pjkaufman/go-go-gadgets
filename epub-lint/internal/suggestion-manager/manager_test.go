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

type moveToPreviousFileTestCase struct {
	setup func(*sm.SuggestionManager)

	expectedFound             bool
	expectedCurrentFileIndex  int
	expectedCurrentIssueIndex int
	expectedSuggestionIndex   int
	expectedFileName          string
	expectedSuggestionName    string
	expectedSuggestionState   *sm.SuggestionState
}

var moveToPreviousFileTestCases = map[string]moveToPreviousFileTestCase{
	"When on first file, false is returned": {
		setup: func(manager *sm.SuggestionManager) {
			manager.CurrentFileIndex = 0
		},
		expectedFound:             false,
		expectedCurrentFileIndex:  0,
		expectedCurrentIssueIndex: 0,
		expectedSuggestionIndex:   0,
	},
	"When no prior file has suggestions, false is returned (2 prior files exist)": {
		setup: func(manager *sm.SuggestionManager) {
			manager.CurrentFileIndex = 2
			manager.CurrentIssueIndex = 1
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             false,
		expectedCurrentFileIndex:  2,
		expectedCurrentIssueIndex: 1,
		expectedSuggestionIndex:   0,
	},
	"When no prior file has suggestions, false is returned (1 prior file exists)": {
		setup: func(manager *sm.SuggestionManager) {
			manager.CurrentFileIndex = 1
			manager.CurrentIssueIndex = 1
		},
		expectedFound:             false,
		expectedCurrentFileIndex:  1,
		expectedCurrentIssueIndex: 1,
		expectedSuggestionIndex:   0,
	},
	"When the previous file has a suggestion, the current file index is updated, and the suggestion is updated to be the first suggestion, and the issue is the corresponding issue for that suggestion": {
		setup: func(manager *sm.SuggestionManager) {
			manager.FileSuggestionData[1].Suggestions[0] = []sm.SuggestionState{
				{
					Original:          "old",
					CurrentSuggestion: "new",
				},
			}

			manager.CurrentFileIndex = 2
			manager.CurrentIssueIndex = 1
			manager.CurrentSuggestionIndex = 1
		},
		expectedFound:             true,
		expectedCurrentFileIndex:  1,
		expectedCurrentIssueIndex: 0,
		expectedSuggestionIndex:   0,
		expectedFileName:          "file2.html",
		expectedSuggestionName:    "Issue 1",
		expectedSuggestionState: &sm.SuggestionState{
			Original:          "old",
			CurrentSuggestion: "new",
		},
	},
	"When two files prior there is a suggestion, but on the prior file there is not a suggestion, then moving to the previous file will go to the file two files back and the suggestion will be the first and the issue will correspond to that one": {
		setup: func(manager *sm.SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[1] = []sm.SuggestionState{
				{
					Original:          "first",
					CurrentSuggestion: "replacement",
				},
			}

			manager.CurrentFileIndex = 2
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             true,
		expectedCurrentFileIndex:  0,
		expectedCurrentIssueIndex: 1,
		expectedSuggestionIndex:   0,
		expectedFileName:          "file1.html",
		expectedSuggestionName:    "Issue 2",
		expectedSuggestionState: &sm.SuggestionState{
			Original:          "first",
			CurrentSuggestion: "replacement",
		},
	},
	"When on the third file and the prior two files both have suggestions, moving to the previous file only goes back a single file and suggestion and issue data matches that first suggestion even when there are suggestions for multiple issues": {
		setup: func(manager *sm.SuggestionManager) {
			// File 1 has suggestions for both issues.
			manager.FileSuggestionData[0].Suggestions[0] = []sm.SuggestionState{
				{
					Original:          "file1-issue1",
					CurrentSuggestion: "replacement1",
				},
			}
			manager.FileSuggestionData[0].Suggestions[1] = []sm.SuggestionState{
				{
					Original:          "file1-issue2",
					CurrentSuggestion: "replacement2",
				},
			}

			// File 2 also has suggestions for both issues.
			manager.FileSuggestionData[1].Suggestions[0] = []sm.SuggestionState{
				{
					Original:          "file2-issue1",
					CurrentSuggestion: "replacement3",
				},
			}
			manager.FileSuggestionData[1].Suggestions[1] = []sm.SuggestionState{
				{
					Original:          "file2-issue2",
					CurrentSuggestion: "replacement4",
				},
			}

			// Start on the third file.
			manager.CurrentFileIndex = 2
			manager.CurrentIssueIndex = 1
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             true,
		expectedCurrentFileIndex:  1,
		expectedCurrentIssueIndex: 0,
		expectedSuggestionIndex:   0,
		expectedFileName:          "file2.html",
		expectedSuggestionName:    "Issue 1",
		expectedSuggestionState: &sm.SuggestionState{
			Original:          "file2-issue1",
			CurrentSuggestion: "replacement3",
		},
	},
}

type moveToPreviousSuggestionTestCase struct {
	setup func(*sm.SuggestionManager)

	expectedFound             bool
	expectedCurrentFileIndex  int
	expectedCurrentIssueIndex int
	expectedSuggestionIndex   int
	expectedFileName          string
	expectedSuggestionName    string
	expectedSuggestionState   *sm.SuggestionState
}

var moveToPreviousSuggestionTestCases = map[string]moveToPreviousSuggestionTestCase{
	"When on first file and no prior suggestion exists, false is returned": {
		setup: func(manager *sm.SuggestionManager) {
			manager.CurrentFileIndex = 0
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             false,
		expectedCurrentFileIndex:  0,
		expectedCurrentIssueIndex: 0,
		expectedSuggestionIndex:   0,
	},
	"When on second file, on the first suggestion, and no prior suggestion exists, false is returned": {
		setup: func(manager *sm.SuggestionManager) {
			manager.CurrentFileIndex = 1
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             false,
		expectedCurrentFileIndex:  1,
		expectedCurrentIssueIndex: 0,
		expectedSuggestionIndex:   0,
	},
	"When on the third file, on the first suggestion, and no prior suggestion exists, false is returned": {
		setup: func(manager *sm.SuggestionManager) {
			manager.CurrentFileIndex = 2
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             false,
		expectedCurrentFileIndex:  2,
		expectedCurrentIssueIndex: 0,
		expectedSuggestionIndex:   0,
	},
	"When on the first file and there are multiple suggestions prior to the current one, moving to the prior suggestion only moves back one suggestion": {
		setup: func(manager *sm.SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[0] = []sm.SuggestionState{
				{Original: "one"},
				{Original: "two"},
				{Original: "three"},
			}

			manager.CurrentFileIndex = 0
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 2
			manager.CurrentSuggestionState = &manager.FileSuggestionData[0].Suggestions[0][2]
		},
		expectedFound:             true,
		expectedCurrentFileIndex:  0,
		expectedCurrentIssueIndex: 0,
		expectedSuggestionIndex:   1,
		expectedSuggestionState: &sm.SuggestionState{
			Original: "two",
		},
	},
	"When on the second file, and the first file has multiple suggestions for multiple issues, moving to the previous suggestion moves to the prior suggestion and its corresponding suggestion and not any prior suggestions": {
		setup: func(manager *sm.SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[0] = []sm.SuggestionState{
				{Original: "one"},
				{Original: "two"},
			}

			manager.FileSuggestionData[0].Suggestions[1] = []sm.SuggestionState{
				{Original: "three"},
			}

			manager.CurrentFileIndex = 1
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             true,
		expectedCurrentFileIndex:  0,
		expectedCurrentIssueIndex: 1,
		expectedSuggestionIndex:   0,
		expectedFileName:          "file1.html",
		expectedSuggestionName:    "Issue 2",
		expectedSuggestionState: &sm.SuggestionState{
			Original: "three",
		},
	},
	"When on third file, on the first suggestion and the second file has no suggestions, but the first does, moving to the prior suggestion moves to the last suggestion in first file": {
		setup: func(manager *sm.SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[0] = []sm.SuggestionState{
				{Original: "one"},
				{Original: "two"},
			}

			manager.FileSuggestionData[0].Suggestions[1] = []sm.SuggestionState{
				{Original: "three"},
			}

			manager.CurrentFileIndex = 2
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             true,
		expectedCurrentFileIndex:  0,
		expectedCurrentIssueIndex: 1,
		expectedSuggestionIndex:   0,
		expectedFileName:          "file1.html",
		expectedSuggestionName:    "Issue 2",
		expectedSuggestionState: &sm.SuggestionState{
			Original: "three",
		},
	},
}

type moveToNextSuggestionTestCase struct {
	setup func(*sm.SuggestionManager)

	expectedFound             bool
	expectedCurrentFileIndex  int
	expectedCurrentIssueIndex int
	expectedSuggestionIndex   int
	expectedSuggestionState   *sm.SuggestionState
}

var moveToNextSuggestionTestCases = map[string]moveToNextSuggestionTestCase{
	"When on the last suggestion for the issue, moving forward returns false": {
		setup: func(manager *sm.SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[0] = []sm.SuggestionState{
				{Original: "one"},
				{Original: "two"},
			}

			manager.CurrentFileIndex = 0
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 1
			manager.CurrentSuggestionState = &manager.FileSuggestionData[0].Suggestions[0][1]
		},
		expectedFound:             false,
		expectedCurrentFileIndex:  0,
		expectedCurrentIssueIndex: 0,
		expectedSuggestionIndex:   1,
		expectedSuggestionState: &sm.SuggestionState{
			Original: "two",
		},
	},
	"When on the first suggestion and there are multiple next suggestions, moving forward only moves forward a single suggestion": {
		setup: func(manager *sm.SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[0] = []sm.SuggestionState{
				{Original: "one"},
				{Original: "two"},
				{Original: "three"},
			}

			manager.CurrentFileIndex = 0
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 0
			manager.CurrentSuggestionState = &manager.FileSuggestionData[0].Suggestions[0][0]
		},
		expectedFound:             true,
		expectedCurrentFileIndex:  0,
		expectedCurrentIssueIndex: 0,
		expectedSuggestionIndex:   1,
		expectedSuggestionState: &sm.SuggestionState{
			Original: "two",
		},
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

	t.Run("MoveToPreviousSuggestion", func(t *testing.T) {
		t.Parallel()

		for name, tc := range moveToPreviousSuggestionTestCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				manager := newManagerForPreviousSuggestionTests()

				tc.setup(manager)

				found := manager.MoveToPreviousSuggestion()

				assert.Equal(t, tc.expectedFound, found)
				assert.Equal(t, tc.expectedCurrentFileIndex, manager.CurrentFileIndex)
				assert.Equal(t, tc.expectedCurrentIssueIndex, manager.CurrentIssueIndex)
				assert.Equal(t, tc.expectedSuggestionIndex, manager.CurrentSuggestionIndex)
				assert.Equal(t, tc.expectedFileName, manager.CurrentFileName)
				assert.Equal(t, tc.expectedSuggestionName, manager.CurrentSuggestionName)

				if tc.expectedSuggestionState == nil {
					assert.Nil(t, manager.CurrentSuggestionState)
				} else {
					require.NotNil(t, manager.CurrentSuggestionState)
					assert.Equal(t, *tc.expectedSuggestionState, *manager.CurrentSuggestionState)
				}
			})
		}
	})

	t.Run("MoveToPreviousFile", func(t *testing.T) {
		t.Parallel()

		for name, tc := range moveToPreviousFileTestCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				manager := newManagerForPreviousFileTests()

				tc.setup(manager)

				found := manager.MoveToPreviousFile()

				assert.Equal(t, tc.expectedFound, found)
				assert.Equal(t, tc.expectedCurrentFileIndex, manager.CurrentFileIndex)
				assert.Equal(t, tc.expectedCurrentIssueIndex, manager.CurrentIssueIndex)
				assert.Equal(t, tc.expectedSuggestionIndex, manager.CurrentSuggestionIndex)
				assert.Equal(t, tc.expectedFileName, manager.CurrentFileName)
				assert.Equal(t, tc.expectedSuggestionName, manager.CurrentSuggestionName)

				if tc.expectedSuggestionState == nil {
					assert.Nil(t, manager.CurrentSuggestionState)
				} else {
					require.NotNil(t, manager.CurrentSuggestionState)
					assert.Equal(t, *tc.expectedSuggestionState, *manager.CurrentSuggestionState)
				}
			})
		}
	})

	t.Run("MoveToNextSuggestion", func(t *testing.T) {
		t.Parallel()

		for name, tc := range moveToNextSuggestionTestCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				manager := newManagerForNextSuggestionTests()

				tc.setup(manager)

				found := manager.MoveToNextSuggestion()

				assert.Equal(t, tc.expectedFound, found)
				assert.Equal(t, tc.expectedCurrentFileIndex, manager.CurrentFileIndex)
				assert.Equal(t, tc.expectedCurrentIssueIndex, manager.CurrentIssueIndex)
				assert.Equal(t, tc.expectedSuggestionIndex, manager.CurrentSuggestionIndex)

				if tc.expectedSuggestionState == nil {
					assert.Nil(t, manager.CurrentSuggestionState)
				} else {
					require.NotNil(t, manager.CurrentSuggestionState)
					assert.Equal(t, *tc.expectedSuggestionState, *manager.CurrentSuggestionState)
				}
			})
		}
	})
}

func newManagerForNextSuggestionTests() *sm.SuggestionManager {
	return &sm.SuggestionManager{
		Suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "Issue 1"},
		},
		FileSuggestionData: []sm.FileSuggestionInfo{
			{
				Name: "file1.html",
				Suggestions: [][]sm.SuggestionState{
					nil,
				},
			},
		},
	}
}

func newManagerForPreviousSuggestionTests() *sm.SuggestionManager {
	return &sm.SuggestionManager{
		Suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "Issue 1"},
			{Name: "Issue 2"},
		},
		FileSuggestionData: []sm.FileSuggestionInfo{
			{
				Name:        "file1.html",
				Suggestions: make([][]sm.SuggestionState, 2),
			},
			{
				Name:        "file2.html",
				Suggestions: make([][]sm.SuggestionState, 2),
			},
			{
				Name:        "file3.html",
				Suggestions: make([][]sm.SuggestionState, 2),
			},
		},
	}
}

func newManagerForPreviousFileTests() *sm.SuggestionManager {
	return &sm.SuggestionManager{
		Suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "Issue 1"},
			{Name: "Issue 2"},
		},
		FileSuggestionData: []sm.FileSuggestionInfo{
			{
				Name:        "file1.html",
				Suggestions: make([][]sm.SuggestionState, 2),
			},
			{
				Name:        "file2.html",
				Suggestions: make([][]sm.SuggestionState, 2),
			},
			{
				Name:        "file3.html",
				Suggestions: make([][]sm.SuggestionState, 2),
			},
		},
	}
}

func pointerToBool(value bool) *bool {
	return &value
}
