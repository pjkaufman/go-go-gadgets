//go:build unit

package suggestionmanager

import (
	"bytes"
	"io"
	"testing"

	_ "embed"

	potentiallyfixableissue "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/potentially-fixable-issue"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/APotentialConversationInstance.html
	simplePotentialConversationInstanceHtml string
	//go:embed testdata/empty.html
	noPotentialIssuesHtml string
	//go:embed testdata/APotentialWordOmmisionInstance.html
	simplePotentialWordOmissionInstanceHtml string
)

// Test case structures
type newSuggestionManagerTestCase struct {
	filePathsToText map[string]string
	suggestions     []potentiallyfixableissue.PotentiallyFixableIssue
	runAll          bool
	skipCss         bool
	logFile         io.Writer

	expectedFileNames        []string
	expectedFileTexts        []string
	expectedSuggestionLength int

	expectedCurrentFileIndex       int
	expectedCurrentIssueIndex      int
	expectedCurrentSuggestionIndex int

	expectedRunAll  bool
	expectedSkipCss bool
	expectedLogFile io.Writer
}

var newSuggestionManagerTestCases = map[string]newSuggestionManagerTestCase{
	"Files are properly sorted": {
		filePathsToText: map[string]string{
			"c.html": "c",
			"a.html": "a",
			"b.html": "b",
		},
		suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "Issue 1"},
			{Name: "Issue 2"},
		},
		expectedFileNames:         []string{"a.html", "b.html", "c.html"},
		expectedFileTexts:         []string{"a", "b", "c"},
		expectedSuggestionLength:  2,
		expectedCurrentIssueIndex: -1,
	},
	"Files already sorted stay in sorted order": {
		filePathsToText: map[string]string{
			"a.html": "a",
			"b.html": "b",
			"c.html": "c",
		},
		suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "Issue 1"},
			{Name: "Issue 2"},
		},
		expectedFileNames:         []string{"a.html", "b.html", "c.html"},
		expectedFileTexts:         []string{"a", "b", "c"},
		expectedSuggestionLength:  2,
		expectedCurrentIssueIndex: -1,
	},
	"Empty initialization": {
		filePathsToText:           map[string]string{},
		suggestions:               nil,
		expectedFileNames:         []string{},
		expectedFileTexts:         []string{},
		expectedSuggestionLength:  0,
		expectedCurrentIssueIndex: -1,
	},
	"State is initialized correctly": {
		filePathsToText: map[string]string{
			"test.html": "test",
		},
		suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "Issue 1"},
		},
		expectedFileNames:              []string{"test.html"},
		expectedFileTexts:              []string{"test"},
		expectedSuggestionLength:       1,
		expectedCurrentFileIndex:       0,
		expectedCurrentIssueIndex:      -1,
		expectedCurrentSuggestionIndex: 0,
	},
	"Constructor parameters are preserved": {
		filePathsToText: map[string]string{
			"test.html": "text",
		},
		suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "Issue 1"},
			{Name: "Issue 2"},
		},
		runAll:  true,
		skipCss: true,
		logFile: &bytes.Buffer{},

		expectedFileNames:         []string{"test.html"},
		expectedFileTexts:         []string{"text"},
		expectedSuggestionLength:  2,
		expectedCurrentIssueIndex: -1,
		expectedRunAll:            true,
		expectedSkipCss:           true,
	},
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
	setup func(*SuggestionManager)

	expectedFound             bool
	expectedCurrentFileIndex  int
	expectedCurrentIssueIndex int
	expectedSuggestionIndex   int
	expectedFileName          string
	expectedSuggestionName    string
	expectedSuggestionState   *SuggestionState
}

var moveToPreviousFileTestCases = map[string]moveToPreviousFileTestCase{
	"When on first file, false is returned": {
		setup: func(manager *SuggestionManager) {
			manager.CurrentFileIndex = 0
		},
		expectedFound:             false,
		expectedCurrentFileIndex:  0,
		expectedCurrentIssueIndex: 0,
		expectedSuggestionIndex:   0,
	},
	"When no prior file has suggestions, false is returned (2 prior files exist)": {
		setup: func(manager *SuggestionManager) {
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
		setup: func(manager *SuggestionManager) {
			manager.CurrentFileIndex = 1
			manager.CurrentIssueIndex = 1
		},
		expectedFound:             false,
		expectedCurrentFileIndex:  1,
		expectedCurrentIssueIndex: 1,
		expectedSuggestionIndex:   0,
	},
	"When the previous file has a suggestion, the current file index is updated, and the suggestion is updated to be the first suggestion, and the issue is the corresponding issue for that suggestion": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[1].Suggestions[0] = []SuggestionState{
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
		expectedSuggestionState: &SuggestionState{
			Original:          "old",
			CurrentSuggestion: "new",
		},
	},
	"When two files prior there is a suggestion, but on the prior file there is not a suggestion, then moving to the previous file will go to the file two files back and the suggestion will be the first and the issue will correspond to that one": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[1] = []SuggestionState{
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
		expectedSuggestionState: &SuggestionState{
			Original:          "first",
			CurrentSuggestion: "replacement",
		},
	},
	"When on the third file and the prior two files both have suggestions, moving to the previous file only goes back a single file and suggestion and issue data matches that first suggestion even when there are suggestions for multiple issues": {
		setup: func(manager *SuggestionManager) {
			// File 1 has suggestions for both issues.
			manager.FileSuggestionData[0].Suggestions[0] = []SuggestionState{
				{
					Original:          "file1-issue1",
					CurrentSuggestion: "replacement1",
				},
			}
			manager.FileSuggestionData[0].Suggestions[1] = []SuggestionState{
				{
					Original:          "file1-issue2",
					CurrentSuggestion: "replacement2",
				},
			}

			// File 2 also has suggestions for both issues.
			manager.FileSuggestionData[1].Suggestions[0] = []SuggestionState{
				{
					Original:          "file2-issue1",
					CurrentSuggestion: "replacement3",
				},
			}
			manager.FileSuggestionData[1].Suggestions[1] = []SuggestionState{
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
		expectedSuggestionState: &SuggestionState{
			Original:          "file2-issue1",
			CurrentSuggestion: "replacement3",
		},
	},
}

type moveToPreviousSuggestionTestCase struct {
	setup func(*SuggestionManager)

	expectedFound             bool
	expectedCurrentFileIndex  int
	expectedCurrentIssueIndex int
	expectedSuggestionIndex   int
	expectedFileName          string
	expectedSuggestionName    string
	expectedSuggestionState   *SuggestionState
}

var moveToPreviousSuggestionTestCases = map[string]moveToPreviousSuggestionTestCase{
	"When on first file and no prior suggestion exists, false is returned": {
		setup: func(manager *SuggestionManager) {
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
		setup: func(manager *SuggestionManager) {
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
		setup: func(manager *SuggestionManager) {
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
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[0] = []SuggestionState{
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
		expectedSuggestionState: &SuggestionState{
			Original: "two",
		},
	},
	"When on the second file, and the first file has multiple suggestions for multiple issues, moving to the previous suggestion moves to the prior suggestion and its corresponding suggestion and not any prior suggestions": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[0] = []SuggestionState{
				{Original: "one"},
				{Original: "two"},
			}

			manager.FileSuggestionData[0].Suggestions[1] = []SuggestionState{
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
		expectedSuggestionState: &SuggestionState{
			Original: "three",
		},
	},
	"When on third file, on the first suggestion and the second file has no suggestions, but the first does, moving to the prior suggestion moves to the last suggestion in first file": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[0] = []SuggestionState{
				{Original: "one"},
				{Original: "two"},
			}

			manager.FileSuggestionData[0].Suggestions[1] = []SuggestionState{
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
		expectedSuggestionState: &SuggestionState{
			Original: "three",
		},
	},
}

type moveToNextSuggestionTestCase struct {
	setup func(*SuggestionManager)

	expectedFound             bool
	expectedCurrentFileIndex  int
	expectedCurrentIssueIndex int
	expectedSuggestionIndex   int
	expectedSuggestionState   *SuggestionState
}

var moveToNextSuggestionTestCases = map[string]moveToNextSuggestionTestCase{
	"When on the last suggestion for the issue, moving forward returns false": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[0] = []SuggestionState{
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
		expectedSuggestionState: &SuggestionState{
			Original: "two",
		},
	},
	"When on the first suggestion and there are multiple next suggestions, moving forward only moves forward a single suggestion": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[0] = []SuggestionState{
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
		expectedSuggestionState: &SuggestionState{
			Original: "two",
		},
	},
}

type acceptSuggestionTestCase struct {
	setup func(*SuggestionManager)

	expectedError error

	expectedText     string
	expectedAccepted bool
}

var acceptSuggestionTestCases = map[string]acceptSuggestionTestCase{
	"When there is no current suggestion, ErrNoCurrentSuggestion is returned": {
		setup: func(manager *SuggestionManager) {
			manager.CurrentSuggestionState = nil
		},
		expectedError: ErrNoCurrentSuggestion,
	},
	"When the current suggestion has already been accepted, ErrSuggestionAlreadyAccepted is returned": {
		setup: func(manager *SuggestionManager) {
			manager.CurrentSuggestionState = &manager.FileSuggestionData[0].Suggestions[0][0]
			manager.CurrentSuggestionState.IsAccepted = true
		},
		expectedError: ErrSuggestionAlreadyAccepted,
	},
	"When there is no current issue, ErrNoCurrentIssueAvailable is returned": {
		setup: func(manager *SuggestionManager) {
			manager.CurrentSuggestionState = &manager.FileSuggestionData[0].Suggestions[0][0]
			manager.CurrentSuggestion = nil
		},
		expectedError: ErrNoCurrentIssueAvailable,
	},
	"When UpdateAllInstances is false, only the first instance is replaced": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[0].Text = "old old"

			manager.Suggestions[0].UpdateAllInstances = false

			manager.FileSuggestionData[0].Suggestions[0] = []SuggestionState{
				{
					Original:          "old",
					CurrentSuggestion: "new",
				},
			}

			manager.CurrentFileIndex = 0
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 0
			manager.CurrentSuggestionState = &manager.FileSuggestionData[0].Suggestions[0][0]
			manager.CurrentSuggestion = &manager.Suggestions[0]
		},
		expectedText:     "new old",
		expectedAccepted: true,
	},
	"When UpdateAllInstances is true, all instances are replaced": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[0].Text = "old old"

			manager.Suggestions[0].UpdateAllInstances = true

			manager.FileSuggestionData[0].Suggestions[0] = []SuggestionState{
				{
					Original:          "old",
					CurrentSuggestion: "new",
				},
			}

			manager.CurrentFileIndex = 0
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 0
			manager.CurrentSuggestionState = &manager.FileSuggestionData[0].Suggestions[0][0]
			manager.CurrentSuggestion = &manager.Suggestions[0]
		},
		expectedText:     "new new",
		expectedAccepted: true,
	},
}

type moveToPreviousIssueTestCase struct {
	setup func(*SuggestionManager)

	expectedFound             bool
	expectedCurrentFileIndex  int
	expectedCurrentIssueIndex int
	expectedSuggestionIndex   int
	expectedFileName          string
	expectedSuggestionName    string
	expectedSuggestionState   *SuggestionState
}

func newManagerForPreviousIssueTests() *SuggestionManager {
	return &SuggestionManager{
		Suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "Issue 1"},
			{Name: "Issue 2"},
			{Name: "Issue 3"},
		},
		FileSuggestionData: []FileSuggestionInfo{
			{
				Name: "file1.html",
				Suggestions: [][]SuggestionState{
					nil,
					nil,
					nil,
				},
			},
			{
				Name: "file2.html",
				Suggestions: [][]SuggestionState{
					nil,
					nil,
					nil,
				},
			},
			{
				Name: "file3.html",
				Suggestions: [][]SuggestionState{
					nil,
					nil,
					nil,
				},
			},
		},
	}
}

var moveToPreviousIssueTestCases = map[string]moveToPreviousIssueTestCase{
	"When on the first file and first issue, no previous issue exists, false is returned": {
		setup: func(manager *SuggestionManager) {
			manager.CurrentFileIndex = 0
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             false,
		expectedCurrentFileIndex:  0,
		expectedCurrentIssueIndex: 0,
		expectedSuggestionIndex:   0,
	},
	"When on the first file and second issue, moving back goes to the previous issue": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[0] = []SuggestionState{
				{
					Original: "issue1",
				},
			}

			manager.FileSuggestionData[0].Suggestions[1] = []SuggestionState{
				{
					Original: "issue2",
				},
			}

			manager.CurrentFileIndex = 0
			manager.CurrentIssueIndex = 1
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             true,
		expectedCurrentFileIndex:  0,
		expectedCurrentIssueIndex: 0,
		expectedSuggestionIndex:   0,
		expectedFileName:          "file1.html",
		expectedSuggestionName:    "Issue 1",
		expectedSuggestionState: &SuggestionState{
			Original: "issue1",
		},
	},
	"When the immediate previous issue has no suggestions, moving back skips to the previous issue with suggestions": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[0] = []SuggestionState{
				{
					Original: "issue1",
				},
			}

			manager.CurrentFileIndex = 0
			manager.CurrentIssueIndex = 2
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             true,
		expectedCurrentFileIndex:  0,
		expectedCurrentIssueIndex: 0,
		expectedSuggestionIndex:   0,
		expectedFileName:          "file1.html",
		expectedSuggestionName:    "Issue 1",
		expectedSuggestionState: &SuggestionState{
			Original: "issue1",
		},
	},
	"When multiple previous issues have suggestions, only the closest previous issue is selected": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[0] = []SuggestionState{
				{
					Original: "issue1",
				},
			}

			manager.FileSuggestionData[0].Suggestions[1] = []SuggestionState{
				{
					Original: "issue2",
				},
			}

			manager.CurrentFileIndex = 0
			manager.CurrentIssueIndex = 2
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             true,
		expectedCurrentFileIndex:  0,
		expectedCurrentIssueIndex: 1,
		expectedSuggestionIndex:   0,
		expectedFileName:          "file1.html",
		expectedSuggestionName:    "Issue 2",
		expectedSuggestionState: &SuggestionState{
			Original: "issue2",
		},
	},
	"When on the first issue of a file, moving back goes to the previous file's issue with suggestions": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[2] = []SuggestionState{
				{
					Original: "previous-file",
				},
			}

			manager.CurrentFileIndex = 1
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             true,
		expectedCurrentFileIndex:  0,
		expectedCurrentIssueIndex: 2,
		expectedSuggestionIndex:   0,
		expectedFileName:          "file1.html",
		expectedSuggestionName:    "Issue 3",
		expectedSuggestionState: &SuggestionState{
			Original: "previous-file",
		},
	},
	"When previous file has empty issues after its last suggestion issue, empty issues are skipped": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[0] = []SuggestionState{
				{
					Original: "issue1",
				},
			}

			manager.CurrentFileIndex = 1
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             true,
		expectedCurrentFileIndex:  0,
		expectedCurrentIssueIndex: 0,
		expectedSuggestionIndex:   0,
		expectedFileName:          "file1.html",
		expectedSuggestionName:    "Issue 1",
		expectedSuggestionState: &SuggestionState{
			Original: "issue1",
		},
	},
	"When previous file has no suggestions, searching continues to earlier files": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[0] = []SuggestionState{
				{
					Original: "file1",
				},
			}

			manager.CurrentFileIndex = 2
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             true,
		expectedCurrentFileIndex:  0,
		expectedCurrentIssueIndex: 0,
		expectedSuggestionIndex:   0,
		expectedFileName:          "file1.html",
		expectedSuggestionName:    "Issue 1",
		expectedSuggestionState: &SuggestionState{
			Original: "file1",
		},
	},
	"When multiple previous files have suggestions, only the closest previous file is selected": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[0] = []SuggestionState{
				{
					Original: "file1",
				},
			}

			manager.FileSuggestionData[1].Suggestions[0] = []SuggestionState{
				{
					Original: "file2",
				},
			}

			manager.CurrentFileIndex = 2
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             true,
		expectedCurrentFileIndex:  1,
		expectedCurrentIssueIndex: 0,
		expectedSuggestionIndex:   0,
		expectedFileName:          "file2.html",
		expectedSuggestionName:    "Issue 1",
		expectedSuggestionState: &SuggestionState{
			Original: "file2",
		},
	},
	"When the previous file has multiple issues, the last issue with suggestions is selected": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[0] = []SuggestionState{
				{
					Original: "issue1",
				},
			}

			manager.FileSuggestionData[0].Suggestions[2] = []SuggestionState{
				{
					Original: "issue3",
				},
			}

			manager.CurrentFileIndex = 1
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             true,
		expectedCurrentFileIndex:  0,
		expectedCurrentIssueIndex: 2,
		expectedSuggestionIndex:   0,
		expectedFileName:          "file1.html",
		expectedSuggestionName:    "Issue 3",
		expectedSuggestionState: &SuggestionState{
			Original: "issue3",
		},
	},
	"When the selected issue has multiple suggestions, the first suggestion is selected": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[0] = []SuggestionState{
				{
					Original: "first",
				},
				{
					Original: "second",
				},
			}

			manager.CurrentFileIndex = 1
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             true,
		expectedCurrentFileIndex:  0,
		expectedCurrentIssueIndex: 0,
		expectedSuggestionIndex:   0,
		expectedFileName:          "file1.html",
		expectedSuggestionName:    "Issue 1",
		expectedSuggestionState: &SuggestionState{
			Original: "first",
		},
	},
}

type moveToNextIssueTestCase struct {
	setup func(*SuggestionManager)

	expectedFound             bool
	expectedError             error
	expectedCurrentFileIndex  int
	expectedCurrentIssueIndex int
	expectedSuggestionIndex   int
	expectedFileName          string
	expectedSuggestionName    string
	expectedSuggestionState   *SuggestionState
}

var moveToNextIssueTestCases = map[string]moveToNextIssueTestCase{
	"When on the last issue and no next issue exists, false is returned": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[2] = []SuggestionState{
				{Original: "current"},
			}

			manager.CurrentFileIndex = 0
			manager.CurrentIssueIndex = 2
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             false,
		expectedCurrentFileIndex:  2,
		expectedCurrentIssueIndex: 0,
		expectedSuggestionIndex:   1, // gets incremented by 1 when move to next issue is called
		expectedFileName:          "file2.html",
	},
	"When another issue exists in the same file, moving forward selects the next issue": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[0] = []SuggestionState{
				{Original: "current"},
			}

			manager.FileSuggestionData[0].Suggestions[1] = []SuggestionState{
				{Original: "next"},
			}

			manager.CurrentFileIndex = 0
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             true,
		expectedCurrentFileIndex:  0,
		expectedCurrentIssueIndex: 1,
		expectedSuggestionIndex:   0,
		expectedFileName:          "file1.html",
		expectedSuggestionName:    "Issue 2",
		expectedSuggestionState: &SuggestionState{
			Original: "next",
		},
	},
	"When the next issue has no suggestions, moving forward skips it": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[0] = []SuggestionState{
				{Original: "current"},
			}

			manager.FileSuggestionData[0].Suggestions[2] = []SuggestionState{
				{Original: "next"},
			}

			manager.CurrentFileIndex = 0
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             true,
		expectedCurrentFileIndex:  0,
		expectedCurrentIssueIndex: 2,
		expectedSuggestionIndex:   0,
		expectedFileName:          "file1.html",
		expectedSuggestionName:    "Issue 3",
		expectedSuggestionState: &SuggestionState{
			Original: "next",
		},
	},
	"When the current file has no more issues, moving forward goes to the next file": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[2] = []SuggestionState{
				{Original: "current"},
			}

			manager.FileSuggestionData[1].Suggestions[0] = []SuggestionState{
				{Original: "next-file"},
			}

			manager.CurrentFileIndex = 0
			manager.CurrentIssueIndex = 2
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             true,
		expectedCurrentFileIndex:  1,
		expectedCurrentIssueIndex: 0,
		expectedSuggestionIndex:   0,
		expectedFileName:          "file2.html",
		expectedSuggestionName:    "Issue 1",
		expectedSuggestionState: &SuggestionState{
			Original: "next-file",
		},
	},
	"When multiple future issues have suggestions, moving forward only moves to the first issue with suggestions": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[0] = []SuggestionState{
				{
					Original: "current",
				},
			}

			manager.FileSuggestionData[0].Suggestions[1] = []SuggestionState{
				{
					Original: "issue2-first",
				},
				{
					Original: "issue2-second",
				},
			}

			manager.FileSuggestionData[0].Suggestions[2] = []SuggestionState{
				{
					Original: "issue3",
				},
			}

			manager.CurrentFileIndex = 0
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 0
			manager.CurrentSuggestionState = &manager.FileSuggestionData[0].Suggestions[0][0]
		},
		expectedFound:             true,
		expectedCurrentFileIndex:  0,
		expectedCurrentIssueIndex: 1,
		expectedSuggestionIndex:   0,
		expectedFileName:          "file1.html",
		expectedSuggestionName:    "Issue 2",
		expectedSuggestionState: &SuggestionState{
			Original: "issue2-first",
		},
	},
}

type moveToNextFileTestCase struct {
	setup func(*SuggestionManager)

	expectedFound             bool
	expectedError             error
	expectedCurrentFileIndex  int
	expectedCurrentIssueIndex int
	expectedSuggestionIndex   int
	expectedFileName          string
	expectedSuggestionName    string
	expectedSuggestionState   *SuggestionState
}

var moveToNextFileTestCases = map[string]moveToNextFileTestCase{
	"When on the last file and no next file exists, false is returned": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[2].Suggestions[0] = []SuggestionState{
				{
					Original: "current",
				},
			}

			manager.CurrentFileIndex = 2
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             false,
		expectedCurrentFileIndex:  3,
		expectedCurrentIssueIndex: 0,
		expectedSuggestionIndex:   0,
		expectedFileName:          "file3.html",
	},
	"When the next file has suggestions, moving forward selects the first suggestion": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[0] = []SuggestionState{
				{
					Original: "current",
				},
			}

			manager.FileSuggestionData[1].Suggestions[0] = []SuggestionState{
				{
					Original: "next",
				},
			}

			manager.CurrentFileIndex = 0
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             true,
		expectedCurrentFileIndex:  1,
		expectedCurrentIssueIndex: 0,
		expectedSuggestionIndex:   0,
		expectedFileName:          "file2.html",
		expectedSuggestionName:    "Issue 1",
		expectedSuggestionState: &SuggestionState{
			Original: "next",
		},
	},
	"When the next file has no suggestions, moving forward skips to the next file with suggestions": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[0].Suggestions[0] = []SuggestionState{
				{
					Original: "current",
				},
			}

			manager.FileSuggestionData[2].Suggestions[0] = []SuggestionState{
				{
					Original: "third-file",
				},
			}

			manager.CurrentFileIndex = 0
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             true,
		expectedCurrentFileIndex:  2,
		expectedCurrentIssueIndex: 0,
		expectedSuggestionIndex:   0,
		expectedFileName:          "file3.html",
		expectedSuggestionName:    "Issue 1",
		expectedSuggestionState: &SuggestionState{
			Original: "third-file",
		},
	},
	"When multiple next files have suggestions, only the closest next file is selected": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[1].Suggestions[0] = []SuggestionState{
				{
					Original: "second-file",
				},
			}

			manager.FileSuggestionData[2].Suggestions[0] = []SuggestionState{
				{
					Original: "third-file",
				},
			}

			manager.CurrentFileIndex = 0
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             true,
		expectedCurrentFileIndex:  1,
		expectedCurrentIssueIndex: 0,
		expectedSuggestionIndex:   0,
		expectedFileName:          "file2.html",
		expectedSuggestionName:    "Issue 1",
		expectedSuggestionState: &SuggestionState{
			Original: "second-file",
		},
	},
	"When the next file has multiple issues with suggestions, the first issue with suggestions is selected": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[1].Suggestions[1] = []SuggestionState{
				{
					Original: "issue2",
				},
			}

			manager.FileSuggestionData[1].Suggestions[2] = []SuggestionState{
				{
					Original: "issue3",
				},
			}

			manager.CurrentFileIndex = 0
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             true,
		expectedCurrentFileIndex:  1,
		expectedCurrentIssueIndex: 1,
		expectedSuggestionIndex:   0,
		expectedFileName:          "file2.html",
		expectedSuggestionName:    "Issue 2",
		expectedSuggestionState: &SuggestionState{
			Original: "issue2",
		},
	},
	"When the selected issue has multiple suggestions, the first suggestion is selected": {
		setup: func(manager *SuggestionManager) {
			manager.FileSuggestionData[1].Suggestions[0] = []SuggestionState{
				{
					Original: "first",
				},
				{
					Original: "second",
				},
			}

			manager.CurrentFileIndex = 0
			manager.CurrentIssueIndex = 0
			manager.CurrentSuggestionIndex = 0
		},
		expectedFound:             true,
		expectedCurrentFileIndex:  1,
		expectedCurrentIssueIndex: 0,
		expectedSuggestionIndex:   0,
		expectedFileName:          "file2.html",
		expectedSuggestionName:    "Issue 1",
		expectedSuggestionState: &SuggestionState{
			Original: "first",
		},
	},
}

func TestSuggestionManager(t *testing.T) {
	t.Parallel()

	t.Run("NewSuggestionManager", func(t *testing.T) {
		t.Parallel()

		for name, tc := range newSuggestionManagerTestCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				manager := NewSuggestionManager(
					tc.suggestions,
					tc.filePathsToText,
					tc.runAll,
					tc.skipCss,
					tc.logFile,
				)

				// File initialization
				require.Len(t, manager.FileSuggestionData, len(tc.expectedFileNames))

				for i := range tc.expectedFileNames {
					assert.Equal(t, tc.expectedFileNames[i], manager.FileSuggestionData[i].Name)
					assert.Equal(t, tc.expectedFileTexts[i], manager.FileSuggestionData[i].Text)

					require.Len(t, manager.FileSuggestionData[i].Suggestions, tc.expectedSuggestionLength)

					for _, suggestions := range manager.FileSuggestionData[i].Suggestions {
						assert.Empty(t, suggestions)
					}
				}

				// Initial state
				assert.Equal(t, tc.expectedCurrentFileIndex, manager.CurrentFileIndex)
				assert.Equal(t, tc.expectedCurrentIssueIndex, manager.CurrentIssueIndex)
				assert.Equal(t, tc.expectedCurrentSuggestionIndex, manager.CurrentSuggestionIndex)

				assert.Nil(t, manager.CurrentSuggestion)
				assert.Nil(t, manager.CurrentSuggestionState)
				assert.Empty(t, manager.CurrentFileName)
				assert.Empty(t, manager.CurrentSuggestionName)

				// Constructor parameters
				assert.Equal(t, tc.suggestions, manager.Suggestions)
				assert.Equal(t, tc.expectedRunAll, manager.runAll)
				assert.Equal(t, tc.expectedSkipCss, manager.skipCss)
				assert.Equal(t, tc.logFile, manager.logFile)
			})
		}
	})

	t.Run("SetupForNextSuggestions", func(t *testing.T) {
		t.Parallel()

		for name, tc := range suggestionManagerSetupForNextSuggestionsTestCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				manager := NewSuggestionManager(createPotentiallyFixableIssues(t, "----"), tc.filePathsToText, tc.runAll, tc.skipCss, nil)

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

	t.Run("MoveToPreviousIssue", func(t *testing.T) {
		t.Parallel()

		for name, tc := range moveToPreviousIssueTestCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				manager := newManagerForPreviousIssueTests()

				tc.setup(manager)

				found := manager.MoveToPreviousIssue()

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

	t.Run("MoveToNextIssue", func(t *testing.T) {
		t.Parallel()

		for name, tc := range moveToNextIssueTestCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				manager := newManagerForNextIssueTests()

				tc.setup(manager)

				found, err := manager.MoveToNextIssue()

				if tc.expectedError != nil {
					require.ErrorIs(t, err, tc.expectedError)
					return
				}

				require.NoError(t, err)

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

	t.Run("MoveToNextFile", func(t *testing.T) {
		t.Parallel()

		for name, tc := range moveToNextFileTestCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				manager := newManagerForNextFileTests()

				tc.setup(manager)

				found, err := manager.MoveToNextFile()

				if tc.expectedError != nil {
					require.ErrorIs(t, err, tc.expectedError)
					return
				}

				require.NoError(t, err)

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

	t.Run("AcceptSuggestion", func(t *testing.T) {
		t.Parallel()

		for name, tc := range acceptSuggestionTestCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				manager := newManagerForAcceptSuggestionTests()

				tc.setup(manager)

				err := manager.AcceptSuggestion()

				if tc.expectedError != nil {
					require.ErrorIs(t, err, tc.expectedError)
					return
				}

				require.NoError(t, err)
				assert.Equal(t, tc.expectedText, manager.FileSuggestionData[0].Text)
				assert.Equal(t, tc.expectedAccepted, manager.CurrentSuggestionState.IsAccepted)
			})
		}
	})
}

// Helpers

func newManagerForAcceptSuggestionTests() *SuggestionManager {
	return &SuggestionManager{
		Suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{
				Name: "Issue 1",
			},
		},
		FileSuggestionData: []FileSuggestionInfo{
			{
				Name: "file1.html",
				Text: "",
				Suggestions: [][]SuggestionState{
					{
						{
							Original:          "old",
							CurrentSuggestion: "new",
						},
					},
				},
			},
		},
	}
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

func newManagerForNextSuggestionTests() *SuggestionManager {
	return &SuggestionManager{
		Suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "Issue 1"},
		},
		FileSuggestionData: []FileSuggestionInfo{
			{
				Name: "file1.html",
				Suggestions: [][]SuggestionState{
					nil,
				},
			},
		},
	}
}

func newManagerForNextIssueTests() *SuggestionManager {
	return &SuggestionManager{
		Suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "Issue 1", GetSuggestions: func(s string) (map[string]string, error) { return make(map[string]string), nil }},
			{Name: "Issue 2", GetSuggestions: func(s string) (map[string]string, error) { return make(map[string]string), nil }},
			{Name: "Issue 3", GetSuggestions: func(s string) (map[string]string, error) { return make(map[string]string), nil }},
		},
		FileSuggestionData: []FileSuggestionInfo{
			{
				Name: "file1.html",
				Suggestions: [][]SuggestionState{
					{},
					{},
					{},
				},
			},
			{
				Name: "file2.html",
				Suggestions: [][]SuggestionState{
					{},
					{},
					{},
				},
			},
		},
		runAll: true,
	}
}

func newManagerForNextFileTests() *SuggestionManager {
	return &SuggestionManager{
		Suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "Issue 1", GetSuggestions: func(s string) (map[string]string, error) { return make(map[string]string), nil }},
			{Name: "Issue 2", GetSuggestions: func(s string) (map[string]string, error) { return make(map[string]string), nil }},
			{Name: "Issue 3", GetSuggestions: func(s string) (map[string]string, error) { return make(map[string]string), nil }},
		},
		FileSuggestionData: []FileSuggestionInfo{
			{
				Name: "file1.html",
				Suggestions: [][]SuggestionState{
					nil,
					nil,
					nil,
				},
			},
			{
				Name: "file2.html",
				Suggestions: [][]SuggestionState{
					nil,
					nil,
					nil,
				},
			},
			{
				Name: "file3.html",
				Suggestions: [][]SuggestionState{
					nil,
					nil,
					nil,
				},
			},
		},
		runAll: true,
	}
}

func newManagerForPreviousSuggestionTests() *SuggestionManager {
	return &SuggestionManager{
		Suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "Issue 1"},
			{Name: "Issue 2"},
		},
		FileSuggestionData: []FileSuggestionInfo{
			{
				Name:        "file1.html",
				Suggestions: make([][]SuggestionState, 2),
			},
			{
				Name:        "file2.html",
				Suggestions: make([][]SuggestionState, 2),
			},
			{
				Name:        "file3.html",
				Suggestions: make([][]SuggestionState, 2),
			},
		},
	}
}

func newManagerForPreviousFileTests() *SuggestionManager {
	return &SuggestionManager{
		Suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "Issue 1"},
			{Name: "Issue 2"},
		},
		FileSuggestionData: []FileSuggestionInfo{
			{
				Name:        "file1.html",
				Suggestions: make([][]SuggestionState, 2),
			},
			{
				Name:        "file2.html",
				Suggestions: make([][]SuggestionState, 2),
			},
			{
				Name:        "file3.html",
				Suggestions: make([][]SuggestionState, 2),
			},
		},
	}
}

func pointerToBool(value bool) *bool {
	return &value
}
