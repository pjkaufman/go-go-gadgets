//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . PotentiallyFixableIssue

package suggestionmanager

import (
	"fmt"
	"io"
	"sort"
	"strings"

	potentiallyfixableissue "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/potentially-fixable-issue"
)

// SuggestionManager manages the state and navigation of suggestions across files and issue types.
type SuggestionManager struct {
	suggestions            []potentiallyfixableissue.PotentiallyFixableIssue
	fileSuggestionData     []FileSuggestionInfo
	currentFileIndex       int
	currentIssueIndex      int
	currentSuggestionIndex int
	currentSuggestionName  string
	runAll                 bool
	skipCss                bool
	logFile                io.Writer
}

// FileSuggestionInfo contains the suggestion data for a single file.
type FileSuggestionInfo struct {
	Name        string
	Text        string
	Suggestions [][]SuggestionState
}

// SuggestionState represents the state of a single suggestion.
type SuggestionState struct {
	IsAccepted                           bool
	OriginallyHadHalfwidthCircleKatakana bool
	Original                             string
	OriginalSuggestion                   string
	CurrentSuggestion                    string
	Display                              string
}

// NewSuggestionManager creates a new SuggestionManager instance.
func NewSuggestionManager(
	suggestions []potentiallyfixableissue.PotentiallyFixableIssue,
	files []FileSuggestionInfo,
	runAll bool,
	skipCss bool,
	logFile io.Writer,
) *SuggestionManager {
	return &SuggestionManager{
		suggestions:            suggestions,
		fileSuggestionData:     files,
		currentFileIndex:       0,
		currentIssueIndex:      0,
		currentSuggestionIndex: 0,
		runAll:                 runAll,
		skipCss:                skipCss,
		logFile:                logFile,
	}
}

// GetCurrentFile returns the name of the current file.
func (sm *SuggestionManager) GetCurrentFile() string {
	if sm.currentFileIndex < len(sm.fileSuggestionData) {
		return sm.fileSuggestionData[sm.currentFileIndex].Name
	}
	return ""
}

// GetCurrentFileIndex returns the current file index.
func (sm *SuggestionManager) GetCurrentFileIndex() int {
	return sm.currentFileIndex
}

// GetFileCount returns the total number of files.
func (sm *SuggestionManager) GetFileCount() int {
	return len(sm.fileSuggestionData)
}

// GetCurrentIssue returns the current issue type.
func (sm *SuggestionManager) GetCurrentIssue() *potentiallyfixableissue.PotentiallyFixableIssue {
	if sm.currentIssueIndex < len(sm.suggestions) {
		return &sm.suggestions[sm.currentIssueIndex]
	}
	return nil
}

// GetCurrentIssueName returns the name of the current issue type.
func (sm *SuggestionManager) GetCurrentIssueName() string {
	return sm.currentSuggestionName
}

// GetCurrentSuggestion returns the current suggestion state.
func (sm *SuggestionManager) GetCurrentSuggestion() *SuggestionState {
	if sm.currentFileIndex >= len(sm.fileSuggestionData) ||
		sm.currentIssueIndex >= len(sm.suggestions) ||
		sm.currentSuggestionIndex >= len(sm.fileSuggestionData[sm.currentFileIndex].Suggestions[sm.currentIssueIndex]) {
		return nil
	}
	return &sm.fileSuggestionData[sm.currentFileIndex].Suggestions[sm.currentIssueIndex][sm.currentSuggestionIndex]
}

// GetCurrentSuggestionIndex returns the current suggestion index.
func (sm *SuggestionManager) GetCurrentSuggestionIndex() int {
	return sm.currentSuggestionIndex
}

// GetCurrentSuggestionCount returns the count of suggestions for the current issue type in the current file.
func (sm *SuggestionManager) GetCurrentSuggestionCount() int {
	if sm.currentFileIndex >= len(sm.fileSuggestionData) ||
		sm.currentIssueIndex >= len(sm.suggestions) {
		return 0
	}
	return len(sm.fileSuggestionData[sm.currentFileIndex].Suggestions[sm.currentIssueIndex])
}

// UpdateFileContent updates the content of a file at the specified index.
func (sm *SuggestionManager) UpdateFileContent(fileIndex int, newContent string) error {
	if fileIndex < 0 || fileIndex >= len(sm.fileSuggestionData) {
		return fmt.Errorf("invalid file index: %d", fileIndex)
	}
	sm.fileSuggestionData[fileIndex].Text = newContent
	return nil
}

// AcceptSuggestion marks the current suggestion as accepted and applies it to the file content.
// If the suggestion has UpdateAllInstances set to true, all instances will be replaced.
func (sm *SuggestionManager) AcceptSuggestion() error {
	currentSuggestion := sm.GetCurrentSuggestion()
	if currentSuggestion == nil {
		return fmt.Errorf("no current suggestion available")
	}

	if currentSuggestion.IsAccepted {
		return fmt.Errorf("suggestion already accepted")
	}

	currentIssue := sm.GetCurrentIssue()
	if currentIssue == nil {
		return fmt.Errorf("no current issue available")
	}

	replaceCount := 1
	if currentIssue.UpdateAllInstances {
		replaceCount = -1
	}

	sm.fileSuggestionData[sm.currentFileIndex].Text = strings.Replace(
		sm.fileSuggestionData[sm.currentFileIndex].Text,
		currentSuggestion.Original,
		currentSuggestion.CurrentSuggestion,
		replaceCount,
	)

	currentSuggestion.IsAccepted = true
	sm.fileSuggestionData[sm.currentFileIndex].Suggestions[sm.currentIssueIndex][sm.currentSuggestionIndex] = *currentSuggestion

	return nil
}

// UpdateCurrentSuggestionValue updates the current suggestion's value.
func (sm *SuggestionManager) UpdateCurrentSuggestionValue(newValue string) error {
	currentSuggestion := sm.GetCurrentSuggestion()
	if currentSuggestion == nil {
		return fmt.Errorf("no current suggestion available")
	}

	currentSuggestion.CurrentSuggestion = newValue
	sm.fileSuggestionData[sm.currentFileIndex].Suggestions[sm.currentIssueIndex][sm.currentSuggestionIndex] = *currentSuggestion

	return nil
}

// MoveToNextSuggestion advances to the next suggestion within the current issue type.
// Returns true if a next suggestion exists, false if we need to move to the next issue.
func (sm *SuggestionManager) MoveToNextSuggestion() bool {
	if sm.currentSuggestionIndex+1 < len(sm.fileSuggestionData[sm.currentFileIndex].Suggestions[sm.currentIssueIndex]) {
		sm.currentSuggestionIndex++
		return true
	}
	sm.currentSuggestionIndex = 0
	return false
}

// MoveToPreviousSuggestion moves to the previous suggestion within the current issue type.
// Returns true if a previous suggestion exists, false if we need to move to the previous issue.
func (sm *SuggestionManager) MoveToPreviousSuggestion() bool {
	if sm.currentSuggestionIndex > 0 {
		sm.currentSuggestionIndex--
		return true
	}
	return false
}

// MoveToNextIssue advances to the next issue type that has suggestions.
// Returns true if a next issue with suggestions exists, false otherwise.
func (sm *SuggestionManager) MoveToNextIssue() bool {
	sm.currentSuggestionIndex = 0

	for sm.currentIssueIndex < len(sm.suggestions) {
		issue := &sm.suggestions[sm.currentIssueIndex]
		sm.log("Checking issue %q (index %d)", issue.Name, sm.currentIssueIndex)

		// Skip disabled issues if not running all
		if !sm.runAll && (issue.IsEnabled == nil || *issue.IsEnabled) {
			sm.log("Skipping issue %q because it is disabled", issue.Name)
			sm.currentIssueIndex++
			continue
		}

		// Skip CSS-related issues if skipCss is set
		if sm.skipCss && (issue.AddCssPageBreakIfMissing || issue.AddCssSectionBreakIfMissing) {
			sm.log("Skipping issue %q because CSS-related rules are to be skipped", issue.Name)
			sm.currentIssueIndex++
			continue
		}

		// Check if we already have suggestions for this issue
		if len(sm.fileSuggestionData[sm.currentFileIndex].Suggestions[sm.currentIssueIndex]) > 0 {
			sm.log("Found existing suggestions for issue %q", issue.Name)
			sm.currentSuggestionName = issue.Name
			return true
		}

		// Generate new suggestions
		suggestions, err := issue.GetSuggestions(sm.fileSuggestionData[sm.currentFileIndex].Text)
		if err != nil {
			sm.log("Error generating suggestions for issue %q: %v", issue.Name, err)
			sm.currentIssueIndex++
			continue
		}

		if len(suggestions) > 0 {
			sm.log("Generated %d suggestions for issue %q", len(suggestions), issue.Name)
			sm.fileSuggestionData[sm.currentFileIndex].Suggestions[sm.currentIssueIndex] = sm.createSuggestionStates(suggestions)
			sm.currentSuggestionName = issue.Name
			sort.Slice(sm.fileSuggestionData[sm.currentFileIndex].Suggestions[sm.currentIssueIndex], func(i, j int) bool {
				return sm.fileSuggestionData[sm.currentFileIndex].Suggestions[sm.currentIssueIndex][i].Original <
					sm.fileSuggestionData[sm.currentFileIndex].Suggestions[sm.currentIssueIndex][j].Original
			})
			return true
		}

		sm.currentIssueIndex++
	}

	return false
}

// MoveToPreviousIssue moves to the previous issue type that has suggestions.
// Returns true if a previous issue with suggestions exists, false otherwise.
func (sm *SuggestionManager) MoveToPreviousIssue() bool {
	for sm.currentIssueIndex > 0 {
		sm.currentIssueIndex--

		if len(sm.fileSuggestionData[sm.currentFileIndex].Suggestions[sm.currentIssueIndex]) > 0 {
			sm.currentSuggestionIndex = 0
			sm.currentSuggestionName = sm.suggestions[sm.currentIssueIndex].Name
			return true
		}
	}

	return false
}

// MoveToNextFile advances to the next file and resets issue/suggestion indices.
// Returns true if a next file exists, false otherwise.
func (sm *SuggestionManager) MoveToNextFile() bool {
	if sm.currentFileIndex+1 < len(sm.fileSuggestionData) {
		sm.currentFileIndex++
		sm.currentIssueIndex = 0
		sm.currentSuggestionIndex = 0
		return true
	}
	return false
}

// MoveToPreviousFile moves to the previous file.
// Returns true if a previous file exists, false otherwise.
func (sm *SuggestionManager) MoveToPreviousFile() bool {
	if sm.currentFileIndex == 0 {
		return false
	}

	sm.currentFileIndex--
	sm.currentIssueIndex = 0
	sm.currentSuggestionIndex = 0

	// Find the first issue in the previous file that has suggestions
	for i := 0; i < len(sm.suggestions); i++ {
		if len(sm.fileSuggestionData[sm.currentFileIndex].Suggestions[i]) > 0 {
			sm.currentIssueIndex = i
			sm.currentSuggestionName = sm.suggestions[i].Name
			return true
		}
	}

	return false
}

// InitializeFirstSuggestion initializes the manager to show the first available suggestion.
// Returns an error if no suggestions are available.
func (sm *SuggestionManager) InitializeFirstSuggestion() error {
	sm.currentFileIndex = 0
	sm.currentIssueIndex = 0
	sm.currentSuggestionIndex = 0

	if !sm.MoveToNextIssue() {
		return fmt.Errorf("no suggestions available")
	}

	return nil
}

// HasNextIssueOrFile checks if there are more issues to process or files to move to.
// This is useful for determining if we've reached the end of all suggestions.
func (sm *SuggestionManager) HasNextIssueOrFile() bool {
	// Check if there are more issues in the current file
	for i := sm.currentIssueIndex + 1; i < len(sm.suggestions); i++ {
		if len(sm.fileSuggestionData[sm.currentFileIndex].Suggestions[i]) > 0 {
			return true
		}
	}

	// Check if there are more files
	if sm.currentFileIndex+1 < len(sm.fileSuggestionData) {
		return true
	}

	return false
}

// GetFileContent returns the current content of the specified file.
func (sm *SuggestionManager) GetFileContent(fileIndex int) (string, error) {
	if fileIndex < 0 || fileIndex >= len(sm.fileSuggestionData) {
		return "", fmt.Errorf("invalid file index: %d", fileIndex)
	}
	return sm.fileSuggestionData[fileIndex].Text, nil
}

// GetAllFiles returns all file suggestion data.
func (sm *SuggestionManager) GetAllFiles() []FileSuggestionInfo {
	return sm.fileSuggestionData
}

// createSuggestionStates converts a map of suggestions into SuggestionState objects.
func (sm *SuggestionManager) createSuggestionStates(suggestions map[string]string) []SuggestionState {
	states := make([]SuggestionState, 0, len(suggestions))
	for original, suggestion := range suggestions {
		states = append(states, SuggestionState{
			Original:           original,
			OriginalSuggestion: suggestion,
			CurrentSuggestion:  suggestion,
			Display:            "", // Display will be set by the caller
		})
	}
	return states
}

// log logs a message to the configured log file if one exists.
func (sm *SuggestionManager) log(format string, args ...interface{}) {
	if sm.logFile != nil {
		fmt.Fprintf(sm.logFile, format+"\n", args...)
	}
}
