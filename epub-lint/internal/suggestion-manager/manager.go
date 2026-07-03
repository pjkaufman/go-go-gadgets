package suggestionmanager

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	potentiallyfixableissue "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/potentially-fixable-issue"
)

var (
	ErrNoCurrentSuggestion       = errors.New("no current suggestion available")
	ErrSuggestionAlreadyAccepted = errors.New("suggestion already accepted")
	ErrNoCurrentIssueAvailable   = errors.New("no current issue available")
)

// SuggestionManager manages the state and navigation of suggestions across files and issue types.
type SuggestionManager struct {
	Suggestions            []potentiallyfixableissue.PotentiallyFixableIssue
	FileSuggestionData     []FileSuggestionInfo
	CurrentFileIndex       int
	CurrentIssueIndex      int
	CurrentSuggestionIndex int
	CurrentSuggestionName  string
	CurrentFileName        string
	CurrentSuggestionState *SuggestionState
	CurrentSuggestion      *potentiallyfixableissue.PotentiallyFixableIssue
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

// NewSuggestionManager creates a new SuggestionManager instance.
// filePathToText maps file paths to their cleaned/processed text content.
// This allows the caller to handle file I/O, cleanup, and sorting before initialization.
func NewSuggestionManager(
	suggestions []potentiallyfixableissue.PotentiallyFixableIssue,
	filePathToText map[string]string,
	runAll bool,
	skipCss bool,
	logFile io.Writer,
) *SuggestionManager {
	// Convert map to sorted FileSuggestionInfo slice
	fileSuggestionData := make([]FileSuggestionInfo, 0, len(filePathToText))
	numFixableIssues := len(suggestions)

	for filePath, text := range filePathToText {
		fileSuggestionData = append(fileSuggestionData, FileSuggestionInfo{
			Name:        filePath,
			Text:        text,
			Suggestions: make([][]SuggestionState, numFixableIssues),
		})
	}

	// Sort by file path to ensure consistent ordering
	sort.Slice(fileSuggestionData, func(i, j int) bool {
		return fileSuggestionData[i].Name < fileSuggestionData[j].Name
	})

	return &SuggestionManager{
		Suggestions:            suggestions,
		FileSuggestionData:     fileSuggestionData,
		CurrentFileIndex:       0,
		CurrentIssueIndex:      -1, // when setup is called the first time this will get set to 0
		CurrentSuggestionIndex: 0,
		runAll:                 runAll,
		skipCss:                skipCss,
		logFile:                logFile,
	}
}

// GetCurrentSuggestion returns the current suggestion state.
func (sm *SuggestionManager) GetCurrentSuggestion() *SuggestionState {
	if sm.CurrentFileIndex >= len(sm.FileSuggestionData) ||
		sm.CurrentIssueIndex >= len(sm.Suggestions) ||
		sm.CurrentSuggestionIndex >= len(sm.FileSuggestionData[sm.CurrentFileIndex].Suggestions[sm.CurrentIssueIndex]) {
		return nil
	}

	return &sm.FileSuggestionData[sm.CurrentFileIndex].Suggestions[sm.CurrentIssueIndex][sm.CurrentSuggestionIndex]
}

// AcceptSuggestion marks the current suggestion as accepted and applies it to the file content.
// If the suggestion has UpdateAllInstances set to true, all instances will be replaced.
func (sm *SuggestionManager) AcceptSuggestion() error {
	if sm.CurrentSuggestionState == nil {
		return ErrNoCurrentSuggestion
	}

	if sm.CurrentSuggestionState.IsAccepted {
		return ErrSuggestionAlreadyAccepted
	}

	if sm.CurrentSuggestion == nil {
		return ErrNoCurrentIssueAvailable
	}

	replaceCount := 1
	if sm.CurrentSuggestion.UpdateAllInstances {
		replaceCount = -1
	}

	sm.FileSuggestionData[sm.CurrentFileIndex].Text = strings.Replace(
		sm.FileSuggestionData[sm.CurrentFileIndex].Text,
		sm.CurrentSuggestionState.Original,
		sm.CurrentSuggestionState.CurrentSuggestion,
		replaceCount,
	)

	sm.CurrentSuggestionState.IsAccepted = true

	return nil
}

// UpdateCurrentSuggestionValue updates the current suggestion's value.
func (sm *SuggestionManager) UpdateCurrentSuggestionValue(newValue string) error {
	if sm.CurrentSuggestionState == nil {
		return ErrNoCurrentSuggestion
	}

	sm.CurrentSuggestionState.CurrentSuggestion = newValue
	sm.CurrentSuggestionState.undoReplaceBrokenDisplayCharacters()

	return nil
}

// MoveToNextSuggestion advances to the next suggestion within the current issue type.
// Returns true if a next suggestion exists, false if we need to search for the next set of issues.
func (sm *SuggestionManager) MoveToNextSuggestion() bool {
	if sm.CurrentSuggestionIndex+1 < len(sm.FileSuggestionData[sm.CurrentFileIndex].Suggestions[sm.CurrentIssueIndex]) {
		sm.CurrentSuggestionIndex++
		sm.CurrentSuggestionState = &sm.FileSuggestionData[sm.CurrentFileIndex].Suggestions[sm.CurrentIssueIndex][sm.CurrentSuggestionIndex]

		return true
	}

	return false
}

// MoveToPreviousSuggestion moves to the previous suggestion.
// Returns true if a previous suggestion exists, false if no prior suggestion exists.
func (sm *SuggestionManager) MoveToPreviousSuggestion() bool {
	if sm.CurrentSuggestionIndex > 0 {
		sm.CurrentSuggestionIndex--
		sm.CurrentSuggestionState = &sm.FileSuggestionData[sm.CurrentFileIndex].Suggestions[sm.CurrentIssueIndex][sm.CurrentSuggestionIndex]

		return true
	}

	var (
		originalCurrentFileIndex           = sm.CurrentFileIndex
		originalPotentialFixableIssueIndex = sm.CurrentIssueIndex
	)
	for sm.CurrentFileIndex != 0 || sm.CurrentIssueIndex != 0 {
		if sm.CurrentIssueIndex == 0 {
			sm.CurrentFileIndex--

			sm.CurrentSuggestionIndex = len(sm.Suggestions) - 1
		} else {
			sm.CurrentIssueIndex--
		}

		for sm.CurrentIssueIndex > 0 && len(sm.FileSuggestionData[sm.CurrentFileIndex].Suggestions[sm.CurrentIssueIndex]) == 0 {
			sm.CurrentIssueIndex--
		}

		var numSuggestions = len(sm.FileSuggestionData[sm.CurrentFileIndex].Suggestions[sm.CurrentIssueIndex])
		if numSuggestions != 0 {
			sm.CurrentSuggestionIndex = numSuggestions - 1
			sm.CurrentSuggestionState = &sm.FileSuggestionData[sm.CurrentFileIndex].Suggestions[sm.CurrentIssueIndex][sm.CurrentSuggestionIndex]
			sm.CurrentFileName = sm.FileSuggestionData[sm.CurrentFileIndex].Name
			sm.CurrentSuggestionName = sm.Suggestions[sm.CurrentIssueIndex].Name

			return true
		}
	}

	sm.CurrentFileIndex = originalCurrentFileIndex
	sm.CurrentIssueIndex = originalPotentialFixableIssueIndex

	return false
}

// MoveToNextIssue advances to the next issue type that has suggestions.
// Returns true if there is another suggestion left otherwise it returns false.
func (sm *SuggestionManager) MoveToNextIssue() (bool, error) {
	sm.CurrentSuggestionIndex = len(sm.FileSuggestionData[sm.CurrentFileIndex].Suggestions[sm.CurrentIssueIndex])

	return sm.SetupForNextSuggestions()
}

// MoveToPreviousIssue moves to the previous issue type that has suggestions.
// Returns true if a previous issue with suggestions exists, false otherwise.
func (sm *SuggestionManager) MoveToPreviousIssue() bool {
	var (
		currentFileIndex           = sm.CurrentFileIndex
		potentialFixableIssueIndex = sm.CurrentIssueIndex
	)
	for currentFileIndex != 0 || potentialFixableIssueIndex != 0 {
		if potentialFixableIssueIndex == 0 {
			currentFileIndex--

			potentialFixableIssueIndex = len(sm.Suggestions) - 1
		} else {
			potentialFixableIssueIndex--
		}

		for potentialFixableIssueIndex > 0 && len(sm.FileSuggestionData[currentFileIndex].Suggestions[potentialFixableIssueIndex]) == 0 {
			potentialFixableIssueIndex--
		}

		var numSuggestions = len(sm.FileSuggestionData[currentFileIndex].Suggestions[potentialFixableIssueIndex])
		if numSuggestions != 0 {
			sm.CurrentFileIndex = currentFileIndex
			sm.CurrentIssueIndex = potentialFixableIssueIndex
			sm.CurrentSuggestionIndex = 0
			sm.CurrentSuggestionName = sm.Suggestions[potentialFixableIssueIndex].Name
			sm.CurrentFileName = sm.FileSuggestionData[currentFileIndex].Name
			sm.CurrentSuggestionState = &sm.FileSuggestionData[sm.CurrentFileIndex].Suggestions[sm.CurrentIssueIndex][sm.CurrentSuggestionIndex]

			return true
		}
	}

	return false
}

// MoveToNextFile advances to the next file that has suggestions if possible.
func (sm *SuggestionManager) MoveToNextFile() (bool, error) {
	if sm.CurrentFileIndex+1 < len(sm.FileSuggestionData) {
		sm.CurrentFileIndex++
		sm.CurrentIssueIndex = 0
		sm.CurrentSuggestionIndex = 0
		sm.CurrentFileName = sm.FileSuggestionData[sm.CurrentFileIndex].Name
		sm.CurrentSuggestionName = sm.Suggestions[sm.CurrentIssueIndex].Name
	} else {
		sm.CurrentIssueIndex = len(sm.Suggestions) - 1
		sm.CurrentSuggestionIndex = len(sm.FileSuggestionData[sm.CurrentFileIndex].Suggestions[sm.CurrentIssueIndex])
	}

	return sm.SetupForNextSuggestions()
}

// MoveToPreviousFile moves to the previous file.
// Returns true if a previous file with a suggestion exists, false otherwise.
func (sm *SuggestionManager) MoveToPreviousFile() bool {
	if sm.CurrentFileIndex == 0 {
		return false
	}

	var (
		currentFileIndex           = sm.CurrentFileIndex
		potentialFixableIssueIndex int
	)
	for currentFileIndex != 0 {
		currentFileIndex--
		potentialFixableIssueIndex = 0
		// skipping files forwards can cause a gap to happen in the potentially fixable issue data, but for now, we will ignore it
		// since if the previous file had any suggestions the first one that has data should be present
		// thus this should work fine, but if it does not we can tweak this
		for potentialFixableIssueIndex < len(sm.Suggestions) {
			if len(sm.FileSuggestionData[currentFileIndex].Suggestions[potentialFixableIssueIndex]) != 0 {
				sm.CurrentFileIndex = currentFileIndex
				sm.CurrentIssueIndex = potentialFixableIssueIndex
				sm.CurrentSuggestionIndex = 0
				sm.CurrentSuggestionName = sm.Suggestions[potentialFixableIssueIndex].Name
				sm.CurrentFileName = sm.FileSuggestionData[currentFileIndex].Name
				sm.CurrentSuggestionState = &sm.FileSuggestionData[sm.CurrentFileIndex].Suggestions[sm.CurrentIssueIndex][sm.CurrentSuggestionIndex]

				return true
			}

			potentialFixableIssueIndex++
		}
	}

	return false
}

func (sm *SuggestionManager) SetupForNextSuggestions() (bool, error) {
	sm.logf("Getting next suggestions")
	sm.CurrentIssueIndex++

	for sm.CurrentFileIndex < len(sm.FileSuggestionData) {
		sm.CurrentFileName = sm.FileSuggestionData[sm.CurrentFileIndex].Name
		sm.logf("Current file is %q is %d of %d\n", sm.CurrentFileName, sm.CurrentFileIndex+1, len(sm.FileSuggestionData))

		for sm.CurrentIssueIndex < len(sm.Suggestions) {
			var potentialFixableIssue = sm.Suggestions[sm.CurrentIssueIndex]
			sm.logf("Possible fixable issue %q is %d of %d issues.", potentialFixableIssue.Name, sm.CurrentIssueIndex+1, len(sm.Suggestions))

			if !sm.runAll && (potentialFixableIssue.IsEnabled == nil || !*potentialFixableIssue.IsEnabled) {
				sm.logf("Skipping possible fixable issue %q with isEnabled set to %v", potentialFixableIssue.Name, potentialFixableIssue.IsEnabled)

				sm.CurrentIssueIndex++
				continue
			} else if sm.skipCss && (potentialFixableIssue.AddCssPageBreakIfMissing || potentialFixableIssue.AddCssSectionBreakIfMissing) {
				sm.logf("Skipping possible fixable issue %q because css related rules are to be skipped", potentialFixableIssue.Name)
				sm.CurrentIssueIndex++

				continue
			}

			if len(sm.FileSuggestionData[sm.CurrentFileIndex].Suggestions[sm.CurrentIssueIndex]) != 0 {
				sm.logf("Possible fixable issue %q has %d suggestion(s) already\n", potentialFixableIssue.Name, len(sm.FileSuggestionData[sm.CurrentFileIndex].Suggestions[sm.CurrentIssueIndex]))

				sm.CurrentSuggestionIndex = 0
				sm.CurrentSuggestionState = &sm.FileSuggestionData[sm.CurrentFileIndex].Suggestions[sm.CurrentIssueIndex][0]
				sm.CurrentSuggestionName = potentialFixableIssue.Name
				sm.CurrentFileName = sm.FileSuggestionData[sm.CurrentFileIndex].Name

				return true, nil
			}

			suggestions, err := potentialFixableIssue.GetSuggestions(sm.FileSuggestionData[sm.CurrentFileIndex].Text)
			if err != nil {
				return false, err
			}

			sm.logf("Possible fixable issue %q has %d suggestion(s) found\n", potentialFixableIssue.Name, len(suggestions))

			if len(suggestions) != 0 {
				sm.CurrentSuggestion = &potentialFixableIssue
				sm.FileSuggestionData[sm.CurrentFileIndex].Suggestions[sm.CurrentIssueIndex] = make([]SuggestionState, len(suggestions))

				var (
					i             = 0
					newSuggestion SuggestionState
				)
				for original, suggestion := range suggestions {
					newSuggestion = SuggestionState{
						Original:           original,
						OriginalSuggestion: suggestion,
						CurrentSuggestion:  suggestion,
					}

					err = newSuggestion.GetStringDiffAsDisplay()
					if err != nil {
						return false, err
					}

					sm.FileSuggestionData[sm.CurrentFileIndex].Suggestions[sm.CurrentIssueIndex][i] = newSuggestion
					i++
				}

				sort.Slice(sm.FileSuggestionData[sm.CurrentFileIndex].Suggestions[sm.CurrentIssueIndex], func(i, j int) bool {
					return sm.FileSuggestionData[sm.CurrentFileIndex].Suggestions[sm.CurrentIssueIndex][i].Original < sm.FileSuggestionData[sm.CurrentFileIndex].Suggestions[sm.CurrentIssueIndex][j].Original
				})

				sm.CurrentSuggestionIndex = 0
				sm.CurrentSuggestionState = &sm.FileSuggestionData[sm.CurrentFileIndex].Suggestions[sm.CurrentIssueIndex][0]
				sm.CurrentSuggestionName = potentialFixableIssue.Name

				return true, nil
			}

			sm.CurrentIssueIndex++
		}

		sm.CurrentFileIndex++
		sm.CurrentIssueIndex = 0
	}

	return false, nil
}

// logf logs a message to the configured logf file if one exists.
func (sm *SuggestionManager) logf(format string, args ...any) {
	if sm.logFile != nil {
		fmt.Fprintf(sm.logFile, format+"\n", args...)
	}
}
