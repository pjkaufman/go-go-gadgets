package epub

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pjkaufman/go-go-gadgets/ebook-lint/cmd/epub/tui"
)

type FixableTuiModel struct {
	currentStage                                        stage
	currentFile                                         string
	sectionBreakInput                                   tui.SectionBreakModel
	potentiallyFixableIssues                            []potentiallyFixableIssue
	sectionSuggestionStates                             []suggestionState
	fileTexts                                           map[string]string
	cssFiles, handledFiles                              []string
	runAll, addCssSectionIfMissing, addCssPageIfMissing bool
	err                                                 error
}

type suggestionState struct {
	isAccepted, isEditing                                    bool
	original, originalSuggestion, currentSuggestion, display string
}

var (
	titleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	subtitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("246"))
	// activeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("190"))
	// inactiveStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	// diffAddStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	// diffRemoveStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)

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
		cssFiles:                 cssFiles,
		runAll:                   runAll,
	}
}

func (f FixableTuiModel) Init() tea.Cmd {
	return nil
}

func (f FixableTuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch f.currentStage {
	case sectionContextBreak:
		f.sectionBreakInput, cmd = f.sectionBreakInput.Update(msg)

		return f, cmd
	case suggestionsProcessing:
		return f.handlePotentialSuggestionsKeys(msg)
	}

	// TODO: handle this differently
	f.sectionBreakInput, cmd = f.sectionBreakInput.Update(msg)

	return f, cmd
}

func (f FixableTuiModel) handlePotentialSuggestionsKeys(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		// case tea.KeyEnter:
		// 	contextBreak = f.sectionBreakInput.Value()

		// 	if strings.TrimSpace(contextBreak) != "" {
		// 		f.currentStage = suggestionsProcessing
		// 	}
		case tea.KeyCtrlC, tea.KeyEsc:
			return f, tea.Quit
		}

	case error:
		f.err = msg
		return f, nil
	}

	f.sectionBreakInput, cmd = f.sectionBreakInput.Update(msg)
	return f, cmd
}

func (f FixableTuiModel) View() string {
	switch f.currentStage {
	case sectionContextBreak:
		return f.sectionBreakInput.View()
	case suggestionsProcessing:
		return f.getPotentialSuggestionsView()
	}

	return f.sectionBreakInput.View()
}

func (f FixableTuiModel) getPotentialSuggestionsView() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render(fmt.Sprintf("Current File: %s", f.currentFile)) + "\n")

	s.WriteString("Suggested Change:\n\n")
	s.WriteString("This is a place holder")
	s.WriteString("\n\n")

	f.displayPotentialSuggestionsControls(&s)

	return s.String()
}

func (f FixableTuiModel) displayPotentialSuggestionsControls(s *strings.Builder) {
	s.WriteString(subtitleStyle.Render("Controls:") + "\n")
	s.WriteString("Enter: Continue   Ctrl+C/Esc: Quit\n")
}
