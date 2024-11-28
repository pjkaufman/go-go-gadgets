package epub

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pjkaufman/go-go-gadgets/ebook-lint/cmd/epub/tui"
)

type FixableTuiModel struct {
	currentStage                                        stage
	currentFile                                         string
	sectionBreakInput                                   tui.SectionBreakModel
	suggestionHandler                                   tui.SuggestionsModel
	potentiallyFixableIssues                            []potentiallyFixableIssue
	filePaths                                           []string
	fileTexts                                           map[string]string
	currentFilePathIndex, currentSuggestIndex           int
	cssFiles, handledFiles                              []string
	runAll, addCssSectionIfMissing, addCssPageIfMissing bool
	Err                                                 error
}

type stage int

const (
	sectionContextBreak stage = iota
	suggestionsProcessing
	stageCssSelection
	finalStage
)

func NewFixableTuiModel(runAll, runSectionBreak bool, potentiallyFixableIssues []potentiallyFixableIssue, cssFiles []string) FixableTuiModel {
	var startingStage = sectionContextBreak
	if !runAll && !runSectionBreak {
		startingStage = suggestionsProcessing
	}

	return FixableTuiModel{
		sectionBreakInput:        tui.NewSectionBreakModel(),
		currentStage:             startingStage,
		potentiallyFixableIssues: potentiallyFixableIssues,
		fileTexts:                make(map[string]string),
		cssFiles:                 cssFiles,
		runAll:                   runAll,
	}
}

func (f FixableTuiModel) Init() tea.Cmd {
	return nil
}

func (f FixableTuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if f.Err != nil {
		return f, tea.Quit
	}

	switch f.currentStage {
	case sectionContextBreak:
		if f.sectionBreakInput.Err != nil {
			f.Err = f.sectionBreakInput.Err

			return f, tea.Quit
		}

		f.sectionBreakInput, cmd = f.sectionBreakInput.Update(msg)
		f = f.advanceStageIfNeeded()

		return f, cmd
	case suggestionsProcessing:
		if f.suggestionHandler.Err != nil {
			f.Err = f.suggestionHandler.Err

			return f, tea.Quit
		}

		f.suggestionHandler, cmd = f.suggestionHandler.Update(msg)
		f = f.advanceStageIfNeeded()

		return f, cmd
	}

	return f, cmd
}

func (f FixableTuiModel) View() string {
	switch f.currentStage {
	case sectionContextBreak:
		return f.sectionBreakInput.View()
	case suggestionsProcessing:
		return f.suggestionHandler.View()
	}

	return ""
}

func (f FixableTuiModel) advanceStageIfNeeded() FixableTuiModel {
	switch f.currentStage {
	case sectionContextBreak:
		if f.sectionBreakInput.Done {
			contextBreak = f.sectionBreakInput.Value()
			f.currentStage = suggestionsProcessing
			f = f.getNextSuggestion()
		}
	case suggestionsProcessing:
		if f.suggestionHandler.Done {
			if f.currentFilePathIndex+1 < len(f.filePaths) {
				f = f.getNextSuggestion()

				return f
			}

			f.currentStage = stageCssSelection
			f.Err = fmt.Errorf("Not implemented stage yet")
			// TODO: figure out how to close the program
		}
	}

	return f
}

func (f FixableTuiModel) getNextSuggestion() FixableTuiModel {
	for f.currentFilePathIndex+1 < len(f.filePaths) {
		for f.currentSuggestIndex+1 < len(f.potentiallyFixableIssues) {
			var (
				currentFilePath = f.filePaths[f.currentFilePathIndex]
				suggestions     = f.potentiallyFixableIssues[f.currentSuggestIndex].getSuggestions(f.fileTexts[currentFilePath])
			)

			if len(suggestions) != 0 {
				var err error
				f.suggestionHandler, err = tui.NewSuggestionsModel(currentFilePath, f.potentiallyFixableIssues[f.currentSuggestIndex].name, fmt.Sprintf("File %d of %d", f.currentFilePathIndex+1, len(f.filePaths)), suggestions)
				if err != nil {
					f.Err = err

					return f
				}

				f.currentSuggestIndex++

				return f
			}

			f.currentSuggestIndex++
		}

		f.currentFilePathIndex++
		f.currentSuggestIndex = 0
	}

	return f
}
