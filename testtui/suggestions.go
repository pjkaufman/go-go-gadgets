package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type suggestions struct {
	height, width                      *int
	currentFile, currentSuggestionName string
	isEditing                          bool
	// suggestionData                     []fileSuggestionInfo

	// currentSuggestionIndex, potentialFixableIssueIndex, currentFileIndex int
	// potentialIssues []potentiallyFixableIssue
	// currentIssue                                                                  *potentiallyFixableIssue
	// cssUpdateRequired, addCssSectionBreakIfMissing, addCssPageBreakIfMissing, isEditing bool
	// currentSuggestionState *suggestionState
	// suggestionEdit         textarea.Model
	suggestionDisplay viewport.Model
	// scrollbar              tea.Model
}

// type fileSuggestionInfo struct {
// 	fileName    string
// 	fileText    string
// 	suggestions []suggestionState
// }

// type suggestionState struct {
// 	isAccepted                                               bool
// 	original, originalSuggestion, currentSuggestion, display string
// }

func newSuggestions(height, width *int) suggestions {
	v := viewport.New(0, 0)

	return suggestions{
		height:                height,
		width:                 width,
		currentFile:           "OEBS/Text/file.html",
		currentSuggestionName: "Suggestion Name",
		suggestionDisplay:     v,
	}
}

func (m suggestions) Init() tea.Cmd {
	return nil
}

func (m suggestions) View() string {
	var status = m.leftStatusView()

	if m.height != nil {
		m.suggestionDisplay.Height = *m.height - leftStatusBorderStyle.GetVerticalBorderSize()
	}

	if m.width != nil {
		m.suggestionDisplay.Width = *m.width - (lipgloss.Width(status) + leftStatusBorderStyle.GetBorderLeftSize() + suggestionBorderStyle.GetBorderRightSize())
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, status, m.suggestionView())
}

func (m suggestions) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m suggestions) leftStatusView() string {
	var (
		statusView      = fmt.Sprintf("%s %s\n%s %s\n", documentIcon, fileNameStyle.Render(m.currentFile), suggestionIcon, suggestionNameStyle.Render(m.currentSuggestionName))
		remainingHeight int
		statusPadding   string
	)

	if m.height != nil {
		remainingHeight = *m.height - (lipgloss.Height(statusView) + leftStatusBorderStyle.GetVerticalBorderSize())
	}

	if remainingHeight > 0 {
		statusPadding = strings.Repeat("\n", remainingHeight)
	}

	return leftStatusBorderStyle.Render(statusView + statusPadding)
}

func (m suggestions) suggestionView() string {
	if !m.isEditing {
		m.suggestionDisplay.SetContent("Hello I am the suggestion....")
	}

	return suggestionBorderStyle.Render(m.suggestionDisplay.View())
}
