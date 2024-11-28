package tui

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	stringdiff "github.com/pjkaufman/go-go-gadgets/pkg/string-diff"
)

type SuggestionsModel struct {
	file, group, fileStatus, OriginalText, FinalText string
	sectionSuggestionStates                          []suggestionState
	suggestionInput                                  textarea.Model
	currentSuggestIndex                              int
	Done, ChangeMade, replaceAll, editMode           bool
	Err                                              error
}

type suggestionState struct {
	isAccepted                                               bool
	original, originalSuggestion, currentSuggestion, display string
}

func NewSuggestionsModel(title, subtitle, fileStatus string, suggestions map[string]string) (SuggestionsModel, error) {
	ti := textarea.New()
	ti.Placeholder = "Enter an edited version of the original string"
	ti.CharLimit = 2000
	ti.ShowLineNumbers = false
	var (
		sectionSuggestionStates = make([]suggestionState, len(suggestions))
		i                       = 0
	)

	for original, suggestion := range suggestions {
		var display, err = getStringDiff(original, suggestion)
		if err != nil {
			return SuggestionsModel{}, err
		}

		sectionSuggestionStates[i] = suggestionState{
			original:           original,
			originalSuggestion: suggestion,
			currentSuggestion:  suggestion,
			display:            display,
		}

		i++
	}

	if len(sectionSuggestionStates) > 0 {
		ti.SetValue(sectionSuggestionStates[0].display)
	}

	return SuggestionsModel{
		suggestionInput:         ti,
		file:                    title,
		group:                   subtitle,
		fileStatus:              fileStatus,
		sectionSuggestionStates: sectionSuggestionStates,
	}, nil
}

func (f SuggestionsModel) Init() tea.Cmd {
	return nil
}

func (f SuggestionsModel) Update(msg tea.Msg) (SuggestionsModel, tea.Cmd) {
	if f.editMode {
		return f.handleEditKeys(msg)
	}

	return f.handleNonEditKeys(msg)
}

func (f SuggestionsModel) handleEditKeys(msg tea.Msg) (SuggestionsModel, tea.Cmd) {
	var (
		cmd               tea.Cmd
		err               error
		currentSuggestion suggestionState
	)

	if f.currentSuggestIndex+1 < len(f.sectionSuggestionStates) {
		currentSuggestion = f.sectionSuggestionStates[f.currentSuggestIndex]
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			// TODO: populate this up without breaking the CLI
			f.Err = fmt.Errorf("User killed program")
			return f, nil
		case "esc":
			return f, tea.Quit
		case "ctrl+s":
			currentSuggestion.currentSuggestion = f.suggestionInput.Value()
			f.editMode = false

			f.sectionSuggestionStates[f.currentSuggestIndex].display, err = getStringDiff(strings.TrimSpace(currentSuggestion.original), strings.TrimSpace(currentSuggestion.currentSuggestion))
			if err != nil {
				f.Err = err
				return f, nil
			}

			return f, nil
		case "ctrl+e":
			f.editMode = false

			f.sectionSuggestionStates[f.currentSuggestIndex].display, err = getStringDiff(strings.TrimSpace(currentSuggestion.original), strings.TrimSpace(currentSuggestion.currentSuggestion))
			if err != nil {
				f.Err = err
				return f, nil
			}

			return f, nil
		case "ctrl+r":
			f.suggestionInput.SetValue(currentSuggestion.originalSuggestion)

			return f, nil
		}
	case tea.WindowSizeMsg:
		return f, tea.ClearScreen
	case error:
		f.Err = msg

		return f, nil
	}

	f.suggestionInput, cmd = f.suggestionInput.Update(msg)
	return f, cmd
}

func (f SuggestionsModel) handleNonEditKeys(msg tea.Msg) (SuggestionsModel, tea.Cmd) {
	var (
		cmd               tea.Cmd
		currentSuggestion suggestionState
	)

	if f.currentSuggestIndex+1 < len(f.sectionSuggestionStates) {
		currentSuggestion = f.sectionSuggestionStates[f.currentSuggestIndex]
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			f.Err = fmt.Errorf("User killed program")
			return f, nil
		case "q", "esc":
			return f, tea.Quit
		case "e":
			if !currentSuggestion.isAccepted {
				f.editMode = true
				f.suggestionInput.SetValue(currentSuggestion.currentSuggestion)

				return f, nil
			}
		case "enter":
			if !currentSuggestion.isAccepted {
				var replaceCount = 1
				if f.replaceAll {
					replaceCount = -1
				}

				f.FinalText = strings.Replace(f.FinalText, currentSuggestion.original, currentSuggestion.currentSuggestion, replaceCount)

				f.ChangeMade = true
				currentSuggestion.isAccepted = true
			}

			f = f.moveToNextSuggestion()

			return f, nil
		case "c":
			// Copy original value to the clipboard
			err := clipboard.WriteAll(currentSuggestion.original)
			if err != nil {
				f.Err = err
				return f, nil
			}

			return f, nil
		case "right", "l":
			f = f.moveToNextSuggestion()

			return f, nil

		case "left", "h":
			f = f.moveToPreviousSuggestion()

			return f, nil
		}
	case tea.WindowSizeMsg:
		return f, tea.ClearScreen
	case error:
		f.Err = msg

		return f, nil
	default:
		f.suggestionInput, cmd = f.suggestionInput.Update(msg)
	}

	return f, cmd
}

func (f SuggestionsModel) View() string {
	var s strings.Builder
	clearScreen(&s)

	if len(f.sectionSuggestionStates) == 0 {
		s.WriteString("No suggestions found")

		return s.String()
	}

	s.WriteString(titleStyle.Render(fmt.Sprintf("Current File: %s", f.file)) + "\n")
	s.WriteString(fileStatusStyle.Render(f.fileStatus))
	s.WriteString("\n")
	s.WriteString(groupStyle.Render(fmt.Sprintf("Issue Group: %s", f.group)) + "\n\n")

	if f.editMode {
		s.WriteString(f.suggestionInput.View())
	} else {
		s.WriteString(f.sectionSuggestionStates[f.currentSuggestIndex].display)
	}

	s.WriteString("\n\n")

	s.WriteString("Suggestion ")
	s.WriteString(fmt.Sprint(f.currentSuggestIndex + 1))
	s.WriteString(" of ")
	s.WriteString(fmt.Sprint(len(f.sectionSuggestionStates)))
	s.WriteString(".\n\n")

	f.displaySuggestionControls(&s)

	// s.WriteString("\nOriginal:")
	// s.WriteString(f.sectionSuggestionStates[f.currentSuggestIndex].original)
	// s.WriteString("\nOriginal Suggestion:")
	// s.WriteString(f.sectionSuggestionStates[f.currentSuggestIndex].originalSuggestion)
	// s.WriteString("\nCurrent Suggestion:")
	// s.WriteString(f.sectionSuggestionStates[f.currentSuggestIndex].currentSuggestion)
	// s.WriteString("\nDisplay Value:")
	// s.WriteString(f.suggestionInput.Value())

	return s.String()
}

func (f SuggestionsModel) displaySuggestionControls(s *strings.Builder) {
	s.WriteString(groupStyle.Render("Controls:") + "\n")

	if !f.editMode {
		s.WriteString("← / → : Previous/Next Suggestion   ")
		s.WriteString("E: Edit   ")
		s.WriteString("C: Copy   ")
		s.WriteString("Enter: Accept   ")
		s.WriteString("Q/Esc: Quit")
	} else {
		s.WriteString("Ctrl+R: Reset   ")
		s.WriteString("Ctrl+E: Cancel edit   ")
		s.WriteString("Ctrl+S: Accept   ")
		s.WriteString("Esc: Quit")
		s.WriteString("Ctrl+C: Exit without saving\n")
	}
}

func (f SuggestionsModel) moveToNextSuggestion() SuggestionsModel {
	// TODO: determine how to tell if th user accidentally tried to move to the next suggestion on the last one or if they actually wanted them reset.
	if f.currentSuggestIndex+1 < len(f.sectionSuggestionStates) {
		f.currentSuggestIndex++
		f.suggestionInput.SetValue(f.sectionSuggestionStates[f.currentSuggestIndex].display)
	} else {
		f.Done = true
	}

	return f
}

func (f SuggestionsModel) moveToPreviousSuggestion() SuggestionsModel {
	if f.currentSuggestIndex > 0 {
		f.currentSuggestIndex--
		f.suggestionInput.SetValue(f.sectionSuggestionStates[f.currentSuggestIndex].display)
	}

	return f
}

func getStringDiff(original, new string) (string, error) {
	return stringdiff.GetPrettyDiffString(strings.TrimLeft(original, "\n"), strings.TrimLeft(new, "\n"))
}
