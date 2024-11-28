package tui

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	stringdiff "github.com/pjkaufman/go-go-gadgets/pkg/string-diff"
)

type SuggestionsModel struct {
	file, group, OriginalText, FinalText string
	sectionSuggestionStates              []suggestionState
	editSuggestion                       textinput.Model
	currentSuggestIndex                  int
	Done, editMode                       bool
	Err                                  error
}

type suggestionState struct {
	isAccepted                                               bool
	original, originalSuggestion, currentSuggestion, display string
}

func NewSuggestionsModel(title, subtitle string, suggestions map[string]string) SuggestionsModel {
	ti := textinput.New()
	ti.Placeholder = "Enter an edited version of the original string"
	ti.CharLimit = 2000
	ti.Width = 20
	var (
		sectionSuggestionStates = make([]suggestionState, len(suggestions))
		i                       = 0
	)

	for original, suggestion := range suggestions {
		sectionSuggestionStates[i] = suggestionState{
			original:           original,
			originalSuggestion: suggestion,
			currentSuggestion:  suggestion,
			display:            getStringDiff(original, suggestion),
		}
	}

	return SuggestionsModel{
		editSuggestion:          ti,
		file:                    title,
		group:                   subtitle,
		sectionSuggestionStates: sectionSuggestionStates,
	}
}

func (f SuggestionsModel) Init() tea.Cmd {
	return nil
}

func (f SuggestionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd               tea.Cmd
		currentSuggestion = f.sectionSuggestionStates[f.currentSuggestIndex]
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return f, tea.Quit
		case "e":
			if !f.editMode && !currentSuggestion.isAccepted {
				f.editMode = true
				f.editSuggestion.SetValue(currentSuggestion.currentSuggestion)

				return f, nil
			}

		case "ctrl+s":
			if f.editMode {
				currentSuggestion.currentSuggestion = f.editSuggestion.Value()
				f.editMode = false

				return f, nil
			}

		case "enter":
			if !f.editMode {
				if !currentSuggestion.isAccepted {
					// TODO: the replace count should be -1 in some instances
					f.FinalText = strings.Replace(f.FinalText, currentSuggestion.original, currentSuggestion.currentSuggestion, 1)
				}

				f.currentSuggestIndex++
			}

			return f, nil
		case "c":
			if !f.editMode {
				// Copy original value to the clipboard
				clipboard.WriteAll(currentSuggestion.original)
				return f, nil
			}

		case "right", "l":
			if !f.editMode && f.currentSuggestIndex+1 < len(f.sectionSuggestionStates) {
				f.currentSuggestIndex++

				return f, nil
			}

		case "left", "h":
			if !f.editMode && f.currentSuggestIndex > 0 {
				f.currentSuggestIndex--

				return f, nil
			}
		}

		// Handle edit mode text input
		if f.editMode {
			f.editSuggestion, cmd = f.editSuggestion.Update(msg)

			return f, cmd
		}
	}

	return f, nil
}

func (f SuggestionsModel) View() string {
	return f.getSectionBreakView()
}

func (f SuggestionsModel) getSectionBreakView() string {
	var (
		s                 strings.Builder
		currentSuggestion = f.sectionSuggestionStates[f.currentSuggestIndex]
	)

	s.WriteString(titleStyle.Render(fmt.Sprintf("Current File: %s", f.file)) + "\n")
	s.WriteString(subtitleStyle.Render(fmt.Sprintf("Issue Group: %s", f.group)) + "\n\n")

	if f.editMode {
		s.WriteString(f.editSuggestion.View())
	} else {
		s.WriteString(getStringDiff(currentSuggestion.original, currentSuggestion.currentSuggestion))
	}

	s.WriteString("\n\n")

	f.displaySectionBreakControls(&s)

	return s.String()
}

func (f SuggestionsModel) displaySectionBreakControls(s *strings.Builder) {
	s.WriteString(subtitleStyle.Render("Controls:") + "\n")
	s.WriteString("← / → : Previous/Next Suggestion   ")
	s.WriteString("Enter: Accept   ")
	s.WriteString("E: Edit   ")
	s.WriteString("C: Copy   ")
	s.WriteString("Q: Quit\n")
}

func getStringDiff(original, new string) string {
	diffString, err := stringdiff.GetPrettyDiffString(strings.TrimLeft(original, "\n"), strings.TrimLeft(new, "\n"))
	if err != nil {
		logger.WriteError(err.Error())
	}

	return diffString
}
