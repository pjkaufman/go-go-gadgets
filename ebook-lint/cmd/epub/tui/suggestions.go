package tui

import (
	"fmt"
	"regexp"
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
	width, height                                    int
	Err                                              error
}

type suggestionState struct {
	isAccepted                                               bool
	original, originalSuggestion, currentSuggestion, display string
}

func NewSuggestionsModel(title, subtitle, fileStatus string, suggestions map[string]string, width, height int) (SuggestionsModel, error) {
	ti := textarea.New()
	ti.Placeholder = "Enter an edited version of the original string"
	ti.CharLimit = 2000
	ti.ShowLineNumbers = false
	ti.SetWidth(width - 4)

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
		width:                   width,
		height:                  height,
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

			f.sectionSuggestionStates[f.currentSuggestIndex].display, err = getStringDiff(currentSuggestion.original, currentSuggestion.currentSuggestion)
			if err != nil {
				f.Err = err
				return f, nil
			}

			return f, nil
		case "ctrl+e":
			f.editMode = false

			f.sectionSuggestionStates[f.currentSuggestIndex].display, err = getStringDiff(currentSuggestion.original, currentSuggestion.currentSuggestion)
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
		f.width = msg.Width
		f.height = msg.Height

		f.suggestionInput.SetWidth(msg.Width - 4)

		return f, tea.Batch(tea.ClearScreen, f.suggestionInput.Focus())
	case error:
		f.Err = msg

		return f, nil
	}

	f.suggestionInput, cmd = f.suggestionInput.Update(msg)
	return f, cmd
}

func (f SuggestionsModel) handleNonEditKeys(msg tea.Msg) (SuggestionsModel, tea.Cmd) {
	var (
		cmds              []tea.Cmd
		cmd               tea.Cmd
		currentSuggestion suggestionState
	)

	f.suggestionInput, cmd = f.suggestionInput.Update(msg)
	cmds = append(cmds, cmd)

	if f.currentSuggestIndex+1 < len(f.sectionSuggestionStates) {
		currentSuggestion = f.sectionSuggestionStates[f.currentSuggestIndex]
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			f.Err = fmt.Errorf("user killed program")
			return f, nil
		case "q", "esc":
			return f, tea.Quit
		case "e":
			if !currentSuggestion.isAccepted {
				f.editMode = true
				f.suggestionInput.SetValue(currentSuggestion.currentSuggestion)

				cmd = f.suggestionInput.Focus()
				cmds = append(cmds, cmd)
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
		case "c":
			// Copy original value to the clipboard
			err := clipboard.WriteAll(currentSuggestion.original)
			if err != nil {
				f.Err = err
				return f, tea.Quit
			}
		case "right", "l":
			f = f.moveToNextSuggestion()
		case "left", "h":
			f = f.moveToPreviousSuggestion()
		}
	case tea.WindowSizeMsg:
		f.width = msg.Width
		f.height = msg.Height
		f.suggestionInput.SetWidth(msg.Width - 4)

		cmds = append(cmds, tea.ClearScreen)
	case error:
		f.Err = msg

		return f, tea.Quit
	}

	return f, tea.Batch(cmds...)
}

func (f SuggestionsModel) View() string {
	var s strings.Builder
	clearScreen(&s)
	// s.WriteString("\n")
	if len(f.sectionSuggestionStates) == 0 {
		s.WriteString(warningStyle.Width(f.width).Render("No suggestions found"))
		return s.String()
	}

	s.WriteString(titleStyle.Width(f.width).Render(fmt.Sprintf("Current File: %s", f.file)) + "\n")
	s.WriteString(fileStatusStyle.Width(f.width).Render(f.fileStatus) + "\n")
	s.WriteString(groupStyle.Width(f.width).Render(fmt.Sprintf("Issue Group: %s", f.group) + "\n\n"))

	// if f.editMode {
	// 	s.WriteString(f.suggestionInput.View() + "\n\n")
	// } else {
	// 	// TODO: if need be write \r here
	// 	// s.WriteString("\r")
	// 	// s.WriteString(wordwrap.String(f.sectionSuggestionStates[f.currentSuggestIndex].display, f.width-10))
	// 	s.WriteString(suggestionStyle.Width(f.width).Render(f.sectionSuggestionStates[f.currentSuggestIndex].display + "\n\n"))
	// }

	s.WriteString(generalStyle.Width(f.width).Render(fmt.Sprintf("Suggestion %d of %d.\n\n", f.currentSuggestIndex+1, len(f.sectionSuggestionStates))))

	f.displaySuggestionControls(&s)
	s.WriteString(generalStyle.Width(f.width).Render(fmt.Sprintf("height: %d | width: %d\n", f.height, f.width)))

	return s.String()
}

func (f SuggestionsModel) displaySuggestionControls(s *strings.Builder) {
	s.WriteString(groupStyle.Width(f.width).Render("Controls:") + "\n")

	var controls []string
	if !f.editMode {
		controls = []string{
			"← / → : Previous/Next Suggestion",
			"E: Edit",
			"C: Copy",
			"Enter: Accept",
			"Q/Esc: Quit",
			"Ctrl+C: Exit without saving",
		}
	} else {
		controls = []string{
			"Ctrl+R: Reset",
			"Ctrl+E: Cancel edit",
			"Ctrl+S: Accept",
			"Esc: Quit",
			"Ctrl+C: Exit without saving",
		}
	}

	s.WriteString(generalStyle.Width(f.width).Render(strings.Join(controls, "   ") + "\n"))
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

var removeStartingLineWhitespace = regexp.MustCompile(`(^|\n)[ \t]+`)

func getStringDiff(original, new string) (string, error) {
	// fmt.Printf("Original (before trim): '%s'\n", original)
	// fmt.Printf("New (before trim): '%s'\n", new)

	original = strings.TrimSpace(removeStartingLineWhitespace.ReplaceAllString(original, "\n"))
	new = strings.TrimSpace(removeStartingLineWhitespace.ReplaceAllString(new, "\n"))
	// fmt.Printf("Original (after trim): '%s'\n", original)
	// fmt.Printf("New (after trim): '%s'\n", new)

	// fmt.Printf("DIFF DEBUG: Original after regex: '%s'\n", original)
	// fmt.Printf("DIFF DEBUG: New after regex: '%s'\n", new)

	var diffString, err = stringdiff.GetPrettyDiffString(original, new)
	diffString = strings.TrimSpace(diffString)

	// fmt.Printf("DIFF DEBUG: Result: '%s'\n", diffString)
	// fmt.Printf("DIFF DEBUG: Result Representation: %q\n", diffString)

	// fmt.Printf("Diff String (after trim): '%s'\n", strings.TrimSpace(diffString))

	return diffString, err
}
