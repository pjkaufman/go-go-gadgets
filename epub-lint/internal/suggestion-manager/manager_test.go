//go:build unit

package suggestionmanager

import (
	"bytes"
	"fmt"
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
	"When files are unsorted, then they are properly sorted": {
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
	"When files are already sorted, then they stay in sorted order": {
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
	"When initializing with no files, then initialization is empty": {
		filePathsToText:           map[string]string{},
		suggestions:               nil,
		expectedFileNames:         []string{},
		expectedFileTexts:         []string{},
		expectedSuggestionLength:  0,
		expectedCurrentIssueIndex: -1,
	},
	"When initializing with files and suggestions, then state is initialized correctly": {
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
	"When providing constructor parameters, then they are preserved": {
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
	expectedFileIndex         int
	expectedFoundSuggestion   bool
}

var suggestionManagerSetupForNextSuggestionsTestCases = map[string]suggestionManagerSetupForNextSuggestionsTestCase{
	"When the first file and first issue has a potential issue, then it is found and suggested": {
		filePathsToText: map[string]string{
			"simple.html": simplePotentialConversationInstanceHtml,
		},
		runAll:                  true,
		expectedFoundSuggestion: true,
		expectedCurrentFileName: "simple.html",
	},
	"When the first file and second issue has a potential issue, then it is found and suggested": {
		filePathsToText: map[string]string{
			"simple.html": simplePotentialWordOmissionInstanceHtml,
		},
		runAll:                    true,
		expectedCurrentIssueIndex: 1,
		expectedFoundSuggestion:   true,
		expectedCurrentFileName:   "simple.html",
	},
	"When a potentially fixable issue exists but its logic is disabled, then no suggestion is returned": {
		filePathsToText: map[string]string{
			"simple.html": simplePotentialConversationInstanceHtml,
		},
		runAll:                  false,
		expectedFoundSuggestion: false,
		expectedFileIndex:       1,
		expectedCurrentFileName: "simple.html",
	},
	"When the first file has no issue but the second does, then the suggestion from the second file is found": {
		filePathsToText: map[string]string{
			"empty.html":  noPotentialIssuesHtml,
			"simple.html": simplePotentialWordOmissionInstanceHtml,
		},
		runAll:                    true,
		expectedCurrentIssueIndex: 1,
		expectedFileIndex:         1,
		expectedFoundSuggestion:   true,
		expectedCurrentFileName:   "simple.html",
	},
	"When there is a single file with no issues, then no suggestions are found": {
		filePathsToText: map[string]string{
			"empty.html": noPotentialIssuesHtml,
		},
		runAll:                  true,
		expectedFoundSuggestion: false,
		expectedFileIndex:       1,
		expectedCurrentFileName: "empty.html",
	},
	"When there are multiple files with no issues, then no suggestions are found": {
		filePathsToText: map[string]string{
			"dont.html":  noPotentialIssuesHtml,
			"empty.html": noPotentialIssuesHtml,
			"full.html":  noPotentialIssuesHtml,
		},
		runAll:                  true,
		expectedFoundSuggestion: false,
		expectedFileIndex:       3,
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
	"When on the first file, then false is returned": {
		setup: func(manager *SuggestionManager) {
			manager.CurrentFileIndex = 0
		},
		expectedFound:             false,
		expectedCurrentFileIndex:  0,
		expectedCurrentIssueIndex: 0,
		expectedSuggestionIndex:   0,
	},
	"When no prior file has suggestions with 2 prior files existing, then false is returned": {
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
	"When no prior file has suggestions with 1 prior file existing, then false is returned": {
		setup: func(manager *SuggestionManager) {
			manager.CurrentFileIndex = 1
			manager.CurrentIssueIndex = 1
		},
		expectedFound:             false,
		expectedCurrentFileIndex:  1,
		expectedCurrentIssueIndex: 1,
		expectedSuggestionIndex:   0,
	},
	"When the previous file has a suggestion, then indices are updated and suggestion is found": {
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
	"When two files back has a suggestion but prior file does not, then the file two back is selected": {
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
	"When on the third file with suggestions in two prior files, then only the closest file is selected": {
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
	"When on the first file with no prior suggestion, then false is returned": {
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
	"When on the second file on the first suggestion with no prior suggestion, then false is returned": {
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
	"When on the third file on the first suggestion with no prior suggestion, then false is returned": {
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
	"When on the first file with multiple prior suggestions, then only one suggestion is moved back": {
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
	"When on the second file with multiple suggestions in first file across multiple issues, then prior suggestion is found": {
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
	"When on the third file on first suggestion with second file having no suggestions, then last suggestion in first file is found": {
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
	"When on the last suggestion for the issue, then false is returned": {
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
	"When on the first suggestion with multiple next suggestions, then only one suggestion is moved forward": {
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
	"When there is no current suggestion, then ErrNoCurrentSuggestion is returned": {
		setup: func(manager *SuggestionManager) {
			manager.CurrentSuggestionState = nil
		},
		expectedError: ErrNoCurrentSuggestion,
	},
	"When the current suggestion has already been accepted, then ErrSuggestionAlreadyAccepted is returned": {
		setup: func(manager *SuggestionManager) {
			manager.CurrentSuggestionState = &manager.FileSuggestionData[0].Suggestions[0][0]
			manager.CurrentSuggestionState.IsAccepted = true
		},
		expectedError: ErrSuggestionAlreadyAccepted,
	},
	"When there is no current issue available, then ErrNoCurrentIssueAvailable is returned": {
		setup: func(manager *SuggestionManager) {
			manager.CurrentSuggestionState = &manager.FileSuggestionData[0].Suggestions[0][0]
			manager.CurrentSuggestion = nil
		},
		expectedError: ErrNoCurrentIssueAvailable,
	},
	"When UpdateAllInstances is false, then only the first instance is replaced": {
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
	"When UpdateAllInstances is true, then all instances are replaced": {
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

var moveToPreviousIssueTestCases = map[string]moveToPreviousIssueTestCase{
	"When on the first file and first issue with no previous issue, then false is returned": {
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
	"When on the first file and second issue, then moving back goes to the previous issue": {
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
	"When the immediate previous issue has no suggestions, then moving back skips to the previous issue with suggestions": {
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
	"When multiple previous issues have suggestions, then only the closest previous issue is selected": {
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
	"When on the first issue of a file, then moving back goes to the previous file's issue with suggestions": {
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
	"When previous file has empty issues after its last suggestion issue, then empty issues are skipped": {
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
	"When previous file has no suggestions, then searching continues to earlier files": {
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
	"When multiple previous files have suggestions, then only the closest previous file is selected": {
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
	"When the previous file has multiple issues, then the last issue with suggestions is selected": {
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
	"When the selected issue has multiple suggestions, then the first suggestion is selected": {
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
	"When on the last issue and no next issue exists, then false is returned": {
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
	"When another issue exists in the same file, then moving forward selects the next issue": {
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
	"When the next issue has no suggestions, then moving forward skips it": {
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
	"When the current file has no more issues, then moving forward goes to the next file": {
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
	"When multiple future issues have suggestions, then moving forward only moves to the first issue with suggestions": {
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
	"When on the last file and no next file exists, then false is returned": {
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
	"When the next file has suggestions, then moving forward selects the first suggestion": {
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
	"When the next file has no suggestions, then moving forward skips to the next file with suggestions": {
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
	"When multiple next files have suggestions, then only the closest next file is selected": {
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
	"When the next file has multiple issues with suggestions, then the first issue with suggestions is selected": {
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
	"When the selected issue has multiple suggestions, then the first suggestion is selected": {
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

type updateCurrentSuggestionValueTestCase struct {
	setup func(*SuggestionManager)

	newValue string

	expectedError error

	expectedCurrentSuggestion string
	expectedOtherSuggestion   string
}

var updateCurrentSuggestionValueTestCases = map[string]updateCurrentSuggestionValueTestCase{
	"When there is no current suggestion, then ErrNoCurrentSuggestion is returned": {
		setup: func(manager *SuggestionManager) {
			manager.CurrentSuggestionState = nil
		},
		newValue:      "new value",
		expectedError: ErrNoCurrentSuggestion,
	},
	"When there is a current suggestion, then the suggestion value is updated": {
		setup: func(manager *SuggestionManager) {
			manager.CurrentSuggestionState =
				&manager.FileSuggestionData[0].Suggestions[0][0]
		},
		newValue:                  "updated value",
		expectedCurrentSuggestion: "updated value",
		expectedOtherSuggestion:   "other-current",
	},
	"When updating one suggestion, then other suggestions are unchanged": {
		setup: func(manager *SuggestionManager) {
			manager.CurrentSuggestionState =
				&manager.FileSuggestionData[0].Suggestions[0][0]
		},
		newValue:                  "new first suggestion",
		expectedCurrentSuggestion: "new first suggestion",
		expectedOtherSuggestion:   "other-current",
	},
	"When the current suggestion value is updated multiple times, then the latest value is used": {
		setup: func(manager *SuggestionManager) {
			manager.CurrentSuggestionState =
				&manager.FileSuggestionData[0].Suggestions[0][0]
		},
		newValue:                  "final value",
		expectedCurrentSuggestion: "final value",
		expectedOtherSuggestion:   "other-current",
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
				assertManagerIndices(t, manager, tc.expectedCurrentFileIndex, tc.expectedCurrentIssueIndex, tc.expectedCurrentSuggestionIndex)

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
				assertManagerIndices(t, manager, tc.expectedFileIndex, tc.expectedCurrentIssueIndex, 0)
				assert.Equal(t, tc.expectedCurrentFileName, manager.CurrentFileName)
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
				assertManagerIndices(t, manager, tc.expectedCurrentFileIndex, tc.expectedCurrentIssueIndex, tc.expectedSuggestionIndex)
				assertCurrentSuggestion(t, manager, tc.expectedFileName, tc.expectedSuggestionName, tc.expectedSuggestionState)
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
				assertManagerIndices(t, manager, tc.expectedCurrentFileIndex, tc.expectedCurrentIssueIndex, tc.expectedSuggestionIndex)
				assertCurrentSuggestion(t, manager, tc.expectedFileName, tc.expectedSuggestionName, tc.expectedSuggestionState)
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
				assertManagerIndices(t, manager, tc.expectedCurrentFileIndex, tc.expectedCurrentIssueIndex, tc.expectedSuggestionIndex)
				assertCurrentSuggestion(t, manager, tc.expectedFileName, tc.expectedSuggestionName, tc.expectedSuggestionState)
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
				assertManagerIndices(t, manager, tc.expectedCurrentFileIndex, tc.expectedCurrentIssueIndex, tc.expectedSuggestionIndex)

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
				assertManagerIndices(t, manager, tc.expectedCurrentFileIndex, tc.expectedCurrentIssueIndex, tc.expectedSuggestionIndex)
				assertCurrentSuggestion(t, manager, tc.expectedFileName, tc.expectedSuggestionName, tc.expectedSuggestionState)
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
				assertManagerIndices(t, manager, tc.expectedCurrentFileIndex, tc.expectedCurrentIssueIndex, tc.expectedSuggestionIndex)
				assertCurrentSuggestion(t, manager, tc.expectedFileName, tc.expectedSuggestionName, tc.expectedSuggestionState)
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

	t.Run("UpdateCurrentSuggestionValue", func(t *testing.T) {
		t.Parallel()

		for name, tc := range updateCurrentSuggestionValueTestCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				manager := newManagerForUpdateCurrentSuggestionValueTests()

				tc.setup(manager)

				err := manager.UpdateCurrentSuggestionValue(tc.newValue)

				if tc.expectedError != nil {
					require.ErrorIs(t, err, tc.expectedError)
					return
				}

				require.NoError(t, err)
				assert.Equal(t, tc.expectedCurrentSuggestion, manager.CurrentSuggestionState.CurrentSuggestion)
				assert.Equal(t, tc.expectedOtherSuggestion, manager.FileSuggestionData[0].Suggestions[0][1].CurrentSuggestion)
			})
		}
	})
}

// Helpers

// Assertion helpers

// assertManagerIndices checks the current file, issue, and suggestion indices
func assertManagerIndices(t *testing.T, manager *SuggestionManager, expectedFileIndex, expectedIssueIndex, expectedSuggestionIndex int) {
	t.Helper()

	assert.Equal(t, expectedFileIndex, manager.CurrentFileIndex, "current file index mismatch")
	assert.Equal(t, expectedIssueIndex, manager.CurrentIssueIndex, "current issue index mismatch")
	assert.Equal(t, expectedSuggestionIndex, manager.CurrentSuggestionIndex, "current suggestion index mismatch")
}

// assertCurrentSuggestion checks the current file name, suggestion name, and suggestion state
func assertCurrentSuggestion(t *testing.T, manager *SuggestionManager, expectedFileName, expectedSuggestionName string, expectedState *SuggestionState) {
	t.Helper()

	assert.Equal(t, expectedFileName, manager.CurrentFileName, "current file name mismatch")
	assert.Equal(t, expectedSuggestionName, manager.CurrentSuggestionName, "current suggestion name mismatch")

	if expectedState == nil {
		assert.Nil(t, manager.CurrentSuggestionState, "expected nil suggestion state")
	} else {
		require.NotNil(t, manager.CurrentSuggestionState, "expected non-nil suggestion state")
		assert.Equal(t, *expectedState, *manager.CurrentSuggestionState, "suggestion state mismatch")
	}
}

// Common test manager factory functions
func newManagerForAcceptSuggestionTests() *SuggestionManager {
	return &SuggestionManager{
		Suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "Issue 1"},
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

func newManagerForPreviousSuggestionTests() *SuggestionManager {
	return createStandardTestManager(2, 3)
}

func newManagerForPreviousIssueTests() *SuggestionManager {
	return createStandardTestManager(3, 3)
}

func newManagerForNextSuggestionTests() *SuggestionManager {
	return createStandardTestManager(1, 3)
}

func newManagerForNextIssueTests() *SuggestionManager {
	manager := createStandardTestManager(3, 2)
	manager.runAll = true
	return manager
}

func newManagerForUpdateCurrentSuggestionValueTests() *SuggestionManager {
	return &SuggestionManager{
		FileSuggestionData: []FileSuggestionInfo{
			{
				Name: "file1.html",
				Suggestions: [][]SuggestionState{
					{
						{
							Original:          "original",
							CurrentSuggestion: "current",
						},
						{
							Original:          "other-original",
							CurrentSuggestion: "other-current",
						},
					},
				},
			},
		},
		CurrentFileIndex:       0,
		CurrentIssueIndex:      0,
		CurrentSuggestionIndex: 0,
	}
}

func newManagerForNextFileTests() *SuggestionManager {
	manager := createStandardTestManager(3, 3)
	manager.runAll = true
	return manager
}

func newManagerForPreviousFileTests() *SuggestionManager {
	return createStandardTestManager(2, 3)
}

// createStandardTestManager creates a manager with the given number of issues and given number of files
// and 3 files (file1.html, file2.html, file3.html)
func createStandardTestManager(numIssues, numFiles int) *SuggestionManager {
	suggestions := make([]potentiallyfixableissue.PotentiallyFixableIssue, numIssues)
	for i := range numIssues {
		suggestions[i] = potentiallyfixableissue.PotentiallyFixableIssue{Name: "Issue " + fmt.Sprint(i+1), GetSuggestions: func(s string) (map[string]string, error) { return map[string]string{}, nil }}
	}

	files := make([]FileSuggestionInfo, numFiles)
	for i := range numFiles {
		files[i] = FileSuggestionInfo{
			Name:        fmt.Sprintf("file%d.html", i+1),
			Suggestions: make([][]SuggestionState, numIssues),
		}
	}

	return &SuggestionManager{
		Suggestions:        suggestions,
		FileSuggestionData: files,
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

func pointerToBool(value bool) *bool {
	return &value
}
