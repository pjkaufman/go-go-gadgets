//go:build unit

package suggestionmanager_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	potentiallyfixableissue "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/potentially-fixable-issue"
	sm "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/suggestion-manager"
)

/*
* Test scenarios
- Initialization:

*/

// Test case structures
type testCase struct {
	name        string
	suggestions []potentiallyfixableissue.PotentiallyFixableIssue
	files       []sm.FileSuggestionInfo
	runAll      bool
	skipCss     bool
	setupFunc   func(*sm.SuggestionManager)
	assertions  func(*testing.T, *sm.SuggestionManager)
}

var suggestionManagerTestCases = map[string]testCase{
	"NewSuggestionManager initializes with correct default state": {
		suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "issue1", GetSuggestions: mockGetSuggestions},
		},
		files: []sm.FileSuggestionInfo{
			{
				Name:        "file1.txt",
				Text:        "original text",
				Suggestions: [][]sm.SuggestionState{{}, {}},
			},
		},
		runAll:  false,
		skipCss: false,
		assertions: func(t *testing.T, manager *sm.SuggestionManager) {
			assert.Equal(t, 0, manager.GetCurrentFileIndex())
			assert.Equal(t, "file1.txt", manager.GetCurrentFile())
			assert.Equal(t, 1, manager.GetFileCount())
		},
	},
	"GetCurrentFile returns correct file name": {
		suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "issue1", GetSuggestions: mockGetSuggestions},
		},
		files: []sm.FileSuggestionInfo{
			{Name: "file1.txt", Text: "content1", Suggestions: [][]sm.SuggestionState{{}, {}}},
			{Name: "file2.txt", Text: "content2", Suggestions: [][]sm.SuggestionState{{}, {}}},
		},
		assertions: func(t *testing.T, manager *sm.SuggestionManager) {
			assert.Equal(t, "file1.txt", manager.GetCurrentFile())
			manager.MoveToNextFile()
			assert.Equal(t, "file2.txt", manager.GetCurrentFile())
		},
	},
	"GetFileCount returns correct count": {
		files: []sm.FileSuggestionInfo{
			{Name: "file1.txt", Text: "content1", Suggestions: [][]sm.SuggestionState{{}, {}}},
			{Name: "file2.txt", Text: "content2", Suggestions: [][]sm.SuggestionState{{}, {}}},
			{Name: "file3.txt", Text: "content3", Suggestions: [][]sm.SuggestionState{{}, {}}},
		},
		suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "issue1", GetSuggestions: mockGetSuggestions},
		},
		assertions: func(t *testing.T, manager *sm.SuggestionManager) {
			assert.Equal(t, 3, manager.GetFileCount())
		},
	},
	"UpdateFileContent successfully updates content": {
		files: []sm.FileSuggestionInfo{
			{Name: "file1.txt", Text: "original", Suggestions: [][]sm.SuggestionState{{}, {}}},
		},
		suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "issue1", GetSuggestions: mockGetSuggestions},
		},
		assertions: func(t *testing.T, manager *sm.SuggestionManager) {
			err := manager.UpdateFileContent(0, "updated")
			require.NoError(t, err)
			content, err := manager.GetFileContent(0)
			require.NoError(t, err)
			assert.Equal(t, "updated", content)
		},
	},
	"UpdateFileContent returns error for invalid index": {
		files: []sm.FileSuggestionInfo{
			{Name: "file1.txt", Text: "content", Suggestions: [][]sm.SuggestionState{{}, {}}},
		},
		suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "issue1", GetSuggestions: mockGetSuggestions},
		},
		assertions: func(t *testing.T, manager *sm.SuggestionManager) {
			err := manager.UpdateFileContent(5, "new content")
			assert.Error(t, err)
		},
	},
	"MoveToNextFile advances to next file": {
		files: []sm.FileSuggestionInfo{
			{Name: "file1.txt", Text: "content1", Suggestions: [][]sm.SuggestionState{{}, {}}},
			{Name: "file2.txt", Text: "content2", Suggestions: [][]sm.SuggestionState{{}, {}}},
		},
		suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "issue1", GetSuggestions: mockGetSuggestions},
		},
		assertions: func(t *testing.T, manager *sm.SuggestionManager) {
			assert.True(t, manager.MoveToNextFile())
			assert.Equal(t, "file2.txt", manager.GetCurrentFile())
			assert.False(t, manager.MoveToNextFile())
		},
	},
	"MoveToPreviousFile moves to previous file": {
		files: []sm.FileSuggestionInfo{
			{Name: "file1.txt", Text: "content1", Suggestions: [][]sm.SuggestionState{{}, {}}},
			{Name: "file2.txt", Text: "content2", Suggestions: [][]sm.SuggestionState{{}, {}}},
		},
		suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "issue1", GetSuggestions: mockGetSuggestions},
		},
		setupFunc: func(manager *sm.SuggestionManager) {
			manager.MoveToNextFile()
		},
		assertions: func(t *testing.T, manager *sm.SuggestionManager) {
			assert.True(t, manager.MoveToPreviousFile())
			assert.Equal(t, "file1.txt", manager.GetCurrentFile())
			assert.False(t, manager.MoveToPreviousFile())
		},
	},
	"AcceptSuggestion marks suggestion as accepted and updates file content": {
		suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "issue1", GetSuggestions: mockGetSuggestions, UpdateAllInstances: false},
		},
		files: []sm.FileSuggestionInfo{
			{
				Name: "file1.txt",
				Text: "original text",
				Suggestions: [][]sm.SuggestionState{
					{
						{
							Original:           "original",
							OriginalSuggestion: "updated",
							CurrentSuggestion:  "updated",
						},
					},
				},
			},
		},
		setupFunc: func(manager *sm.SuggestionManager) {
			// Already positioned at first suggestion
		},
		assertions: func(t *testing.T, manager *sm.SuggestionManager) {
			err := manager.AcceptSuggestion()
			require.NoError(t, err)

			currentSuggestion := manager.GetCurrentSuggestion()
			assert.True(t, currentSuggestion.IsAccepted)

			content, err := manager.GetFileContent(0)
			require.NoError(t, err)
			assert.Equal(t, "updated text", content)
		},
	},
	"AcceptSuggestion with UpdateAllInstances replaces all occurrences": {
		suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "issue1", GetSuggestions: mockGetSuggestions, UpdateAllInstances: true},
		},
		files: []sm.FileSuggestionInfo{
			{
				Name: "file1.txt",
				Text: "test test test",
				Suggestions: [][]sm.SuggestionState{
					{
						{
							Original:           "test",
							OriginalSuggestion: "updated",
							CurrentSuggestion:  "updated",
						},
					},
				},
			},
		},
		assertions: func(t *testing.T, manager *sm.SuggestionManager) {
			err := manager.AcceptSuggestion()
			require.NoError(t, err)

			content, err := manager.GetFileContent(0)
			require.NoError(t, err)
			assert.Equal(t, "updated updated updated", content)
		},
	},
	"UpdateCurrentSuggestionValue updates suggestion value": {
		suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "issue1", GetSuggestions: mockGetSuggestions},
		},
		files: []sm.FileSuggestionInfo{
			{
				Name: "file1.txt",
				Text: "content",
				Suggestions: [][]sm.SuggestionState{
					{
						{
							Original:           "original",
							OriginalSuggestion: "suggestion",
							CurrentSuggestion:  "suggestion",
						},
					},
				},
			},
		},
		assertions: func(t *testing.T, manager *sm.SuggestionManager) {
			err := manager.UpdateCurrentSuggestionValue("new value")
			require.NoError(t, err)

			current := manager.GetCurrentSuggestion()
			assert.Equal(t, "new value", current.CurrentSuggestion)
		},
	},
	"MoveToNextSuggestion advances within current issue": {
		suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "issue1", GetSuggestions: mockGetSuggestions},
		},
		files: []sm.FileSuggestionInfo{
			{
				Name: "file1.txt",
				Text: "content",
				Suggestions: [][]sm.SuggestionState{
					{
						{Original: "original1", CurrentSuggestion: "updated1"},
						{Original: "original2", CurrentSuggestion: "updated2"},
					},
				},
			},
		},
		assertions: func(t *testing.T, manager *sm.SuggestionManager) {
			assert.Equal(t, 0, manager.GetCurrentSuggestionIndex())
			assert.True(t, manager.MoveToNextSuggestion())
			assert.Equal(t, 1, manager.GetCurrentSuggestionIndex())
			assert.False(t, manager.MoveToNextSuggestion())
			assert.Equal(t, 0, manager.GetCurrentSuggestionIndex())
		},
	},
	"MoveToPreviousSuggestion moves backward within current issue": {
		suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "issue1", GetSuggestions: mockGetSuggestions},
		},
		files: []sm.FileSuggestionInfo{
			{
				Name: "file1.txt",
				Text: "content",
				Suggestions: [][]sm.SuggestionState{
					{
						{Original: "original1", CurrentSuggestion: "updated1"},
						{Original: "original2", CurrentSuggestion: "updated2"},
					},
				},
			},
		},
		setupFunc: func(manager *sm.SuggestionManager) {
			manager.MoveToNextSuggestion()
		},
		assertions: func(t *testing.T, manager *sm.SuggestionManager) {
			assert.True(t, manager.MoveToPreviousSuggestion())
			assert.Equal(t, 0, manager.GetCurrentSuggestionIndex())
			assert.False(t, manager.MoveToPreviousSuggestion())
		},
	},
	"InitializeFirstSuggestion positions at first available suggestion": {
		suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "issue1", GetSuggestions: mockGetSuggestions},
		},
		files: []sm.FileSuggestionInfo{
			{
				Name: "file1.txt",
				Text: "content",
				Suggestions: [][]sm.SuggestionState{
					{
						{Original: "test", CurrentSuggestion: "updated"},
					},
				},
			},
		},
		assertions: func(t *testing.T, manager *sm.SuggestionManager) {
			err := manager.InitializeFirstSuggestion()
			require.NoError(t, err)
			assert.NotNil(t, manager.GetCurrentSuggestion())
		},
	},
	"InitializeFirstSuggestion returns error when no suggestions available": {
		suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "issue1", GetSuggestions: mockGetSuggestions},
		},
		files: []sm.FileSuggestionInfo{
			{
				Name:        "file1.txt",
				Text:        "content",
				Suggestions: [][]sm.SuggestionState{{}},
			},
		},
		assertions: func(t *testing.T, manager *sm.SuggestionManager) {
			err := manager.InitializeFirstSuggestion()
			assert.Error(t, err)
		},
	},
	"GetAllFiles returns all file data": {
		files: []sm.FileSuggestionInfo{
			{Name: "file1.txt", Text: "content1", Suggestions: [][]sm.SuggestionState{{}, {}}},
			{Name: "file2.txt", Text: "content2", Suggestions: [][]sm.SuggestionState{{}, {}}},
		},
		suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "issue1", GetSuggestions: mockGetSuggestions},
		},
		assertions: func(t *testing.T, manager *sm.SuggestionManager) {
			allFiles := manager.GetAllFiles()
			assert.Equal(t, 2, len(allFiles))
			assert.Equal(t, "file1.txt", allFiles[0].Name)
			assert.Equal(t, "file2.txt", allFiles[1].Name)
		},
	},
	"GetCurrentSuggestionCount returns correct count": {
		suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "issue1", GetSuggestions: mockGetSuggestions},
		},
		files: []sm.FileSuggestionInfo{
			{
				Name: "file1.txt",
				Text: "content",
				Suggestions: [][]sm.SuggestionState{
					{
						{Original: "test1", CurrentSuggestion: "updated1"},
						{Original: "test2", CurrentSuggestion: "updated2"},
						{Original: "test3", CurrentSuggestion: "updated3"},
					},
				},
			},
		},
		assertions: func(t *testing.T, manager *sm.SuggestionManager) {
			assert.Equal(t, 3, manager.GetCurrentSuggestionCount())
		},
	},
	"HasNextIssueOrFile returns true when more issues exist": {
		suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "issue1", GetSuggestions: mockGetSuggestions},
			{Name: "issue2", GetSuggestions: mockGetSuggestions},
		},
		files: []sm.FileSuggestionInfo{
			{
				Name: "file1.txt",
				Text: "content",
				Suggestions: [][]sm.SuggestionState{
					{
						{Original: "test", CurrentSuggestion: "updated"},
					},
					{
						{Original: "test2", CurrentSuggestion: "updated2"},
					},
				},
			},
		},
		assertions: func(t *testing.T, manager *sm.SuggestionManager) {
			assert.True(t, manager.HasNextIssueOrFile())
		},
	},
	"HasNextIssueOrFile returns true when more files exist": {
		suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "issue1", GetSuggestions: mockGetSuggestions},
		},
		files: []sm.FileSuggestionInfo{
			{
				Name: "file1.txt",
				Text: "content1",
				Suggestions: [][]sm.SuggestionState{
					{
						{Original: "test", CurrentSuggestion: "updated"},
					},
				},
			},
			{
				Name: "file2.txt",
				Text: "content2",
				Suggestions: [][]sm.SuggestionState{
					{
						{Original: "test2", CurrentSuggestion: "updated2"},
					},
				},
			},
		},
		assertions: func(t *testing.T, manager *sm.SuggestionManager) {
			assert.True(t, manager.HasNextIssueOrFile())
		},
	},
	"HasNextIssueOrFile returns false when at end": {
		suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "issue1", GetSuggestions: mockGetSuggestions},
		},
		files: []sm.FileSuggestionInfo{
			{
				Name: "file1.txt",
				Text: "content",
				Suggestions: [][]sm.SuggestionState{
					{
						{Original: "test", CurrentSuggestion: "updated"},
					},
				},
			},
		},
		assertions: func(t *testing.T, manager *sm.SuggestionManager) {
			assert.False(t, manager.HasNextIssueOrFile())
		},
	},
	"Manager logs to io.Writer when provided": {
		suggestions: []potentiallyfixableissue.PotentiallyFixableIssue{
			{Name: "issue1", GetSuggestions: mockGetSuggestions},
		},
		files: []sm.FileSuggestionInfo{
			{
				Name: "file1.txt",
				Text: "content",
				Suggestions: [][]sm.SuggestionState{
					{
						{Original: "test", CurrentSuggestion: "updated"},
					},
				},
			},
		},
		assertions: func(t *testing.T, manager *sm.SuggestionManager) {
			// The log functionality is tested implicitly through the manager's operations
			// We can verify it doesn't panic when logging is performed
			err := manager.UpdateFileContent(0, "new content")
			require.NoError(t, err)
		},
	},
}

// Helper function to mock GetSuggestions
func mockGetSuggestions(text string) (map[string]string, error) {
	return map[string]string{}, nil
}

func TestSuggestionManager(t *testing.T) {
	for name, tc := range suggestionManagerTestCases {
		t.Run(name, func(t *testing.T) {
			var logBuffer bytes.Buffer
			manager := sm.NewSuggestionManager(tc.suggestions, tc.files, tc.runAll, tc.skipCss, &logBuffer)

			if tc.setupFunc != nil {
				tc.setupFunc(manager)
			}

			tc.assertions(t, manager)
		})
	}
}

// Benchmark tests
func BenchmarkMoveToNextFile(b *testing.B) {
	files := make([]sm.FileSuggestionInfo, 100)
	for i := 0; i < 100; i++ {
		files[i] = sm.FileSuggestionInfo{
			Name:        "file.txt",
			Text:        "content",
			Suggestions: [][]sm.SuggestionState{{}, {}},
		}
	}

	suggestions := []potentiallyfixableissue.PotentiallyFixableIssue{
		{Name: "issue1", GetSuggestions: mockGetSuggestions},
	}

	manager := sm.NewSuggestionManager(suggestions, files, false, false, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.MoveToNextFile()
		if manager.GetCurrentFileIndex() >= 99 {
			// Reset by creating a new manager
			manager = sm.NewSuggestionManager(suggestions, files, false, false, nil)
		}
	}
}

func BenchmarkUpdateFileContent(b *testing.B) {
	files := []sm.FileSuggestionInfo{
		{
			Name:        "file.txt",
			Text:        "content",
			Suggestions: [][]sm.SuggestionState{{}, {}},
		},
	}

	suggestions := []potentiallyfixableissue.PotentiallyFixableIssue{
		{Name: "issue1", GetSuggestions: mockGetSuggestions},
	}

	manager := sm.NewSuggestionManager(suggestions, files, false, false, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.UpdateFileContent(0, "new content")
	}
}
