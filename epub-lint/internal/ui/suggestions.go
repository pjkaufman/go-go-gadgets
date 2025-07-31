package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pjkaufman/go-go-gadgets/epub-lint/cmd/tui"
)

func (m FixableIssuesModel) suggestionsView() string {
	var (
		status         = m.leftStatusView()
		height         = m.body.Height - leftStatusBorderStyle.GetVerticalBorderSize()
		remainingWidth = m.getSuggestionWidth(status)
	)

	m.PotentiallyFixableIssuesInfo.suggestionDisplay.Height = height
	m.PotentiallyFixableIssuesInfo.suggestionDisplay.Width = remainingWidth
	m.PotentiallyFixableIssuesInfo.suggestionEdit.SetWidth(remainingWidth)
	m.PotentiallyFixableIssuesInfo.suggestionEdit.SetHeight(height)

	m.setSuggestionDisplay()

	return lipgloss.JoinHorizontal(lipgloss.Left, status, m.suggestionView())
}

func (m FixableIssuesModel) getSuggestionWidth(status string) int {
	return m.body.Width - (lipgloss.Width(status) + leftStatusBorderStyle.GetBorderLeftSize() + suggestionBorderStyle.GetBorderRightSize())
}

func (m FixableIssuesModel) leftStatusView() string {
	var (
		statusView      = fmt.Sprintf("%s %s\n%s %s\n", documentIcon, fileNameStyle.Render(m.PotentiallyFixableIssuesInfo.currentFile), suggestionIcon, suggestionNameStyle.Render(m.PotentiallyFixableIssuesInfo.currentSuggestionName))
		remainingHeight int
		statusPadding   string
	)

	remainingHeight = m.body.Height - (lipgloss.Height(statusView) + leftStatusBorderStyle.GetVerticalBorderSize())

	if remainingHeight > 0 {
		statusPadding = strings.Repeat("\n", remainingHeight)
	}

	return leftStatusBorderStyle.Render(statusView + statusPadding)
}

func (m FixableIssuesModel) suggestionView() string {
	var s strings.Builder
	// s.WriteString(titleStyle.Render(fmt.Sprintf("Current File (%d/%d): %s ", m.PotentiallyFixableIssuesInfo.currentFileIndex+1, len(m.PotentiallyFixableIssuesInfo.FilePaths), m.PotentiallyFixableIssuesInfo.currentFile)) + "\n")

	if m.PotentiallyFixableIssuesInfo.currentSuggestionState == nil {
		s.WriteString("\nNo current file is selected. Something may have gone wrong...\n\n")
	} else {
		modeIcon := viewIcon
		modeName := "View"
		if m.PotentiallyFixableIssuesInfo.isEditing {
			modeIcon = editIcon
			modeName = "Edit"
		}

		borderConfig := NewBorderConfig(m.PotentiallyFixableIssuesInfo.suggestionDisplay.Height+2, m.PotentiallyFixableIssuesInfo.suggestionDisplay.Width+2) // +2 for border width/height
		modeInfo := fmt.Sprintf("%s %s", modeIcon, modeName)
		statusInfo := fmt.Sprintf("%d/%d", m.PotentiallyFixableIssuesInfo.currentSuggestionIndex+1, len(m.PotentiallyFixableIssuesInfo.sectionSuggestionStates))
		borderConfig.SetInfoItems(modeInfo, statusInfo)

		baseBorder := lipgloss.RoundedBorder()
		customBorder := borderConfig.GetBorder(baseBorder)
		customBorderStyle := lipgloss.NewStyle().Border(customBorder)
		customBorderStyle = customBorderStyle.BorderForeground(suggestionBorderStyle.GetBorderTopForeground())

		if m.PotentiallyFixableIssuesInfo.isEditing {
			s.WriteString(customBorderStyle.Render(m.PotentiallyFixableIssuesInfo.suggestionEdit.View()) + "\n\n")
		} else if m.PotentiallyFixableIssuesInfo.suggestionDisplay.TotalLineCount() > m.PotentiallyFixableIssuesInfo.suggestionDisplay.VisibleLineCount() {
			s.WriteString(customBorderStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top,
				m.PotentiallyFixableIssuesInfo.suggestionDisplay.View(),
				m.PotentiallyFixableIssuesInfo.scrollbar.View(),
			)))
		} else {
			s.WriteString(customBorderStyle.Render(m.PotentiallyFixableIssuesInfo.suggestionDisplay.View()))
		}
	}

	return s.String()
}

func (m *FixableIssuesModel) handleSuggestionMsgs(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.CssSelectionInfo.currentCssIndex > 0 {
				m.CssSelectionInfo.currentCssIndex--
				m.CssSelectionInfo.SelectedCssFile = m.CssSelectionInfo.cssFiles[m.CssSelectionInfo.currentCssIndex]
			}
		case "down":
			if m.CssSelectionInfo.currentCssIndex+1 < len(m.CssSelectionInfo.cssFiles) {
				m.CssSelectionInfo.currentCssIndex++
				m.CssSelectionInfo.SelectedCssFile = m.CssSelectionInfo.cssFiles[m.CssSelectionInfo.currentCssIndex]
			}
		case "enter":
			m.currentStage = finalStage
		}
	}

	return nil
}

func (m *FixableIssuesModel) setupForNextSuggestions() (tea.Cmd, error) {
	if m.logFile != nil {
		fmt.Fprintln(m.logFile, "Getting next suggestions")
	}

	for m.PotentiallyFixableIssuesInfo.currentFileIndex < len(m.PotentiallyFixableIssuesInfo.FilePaths) {
		m.PotentiallyFixableIssuesInfo.currentFile = m.PotentiallyFixableIssuesInfo.FilePaths[m.PotentiallyFixableIssuesInfo.currentFileIndex]
		if m.logFile != nil {
			fmt.Fprintf(m.logFile, "Current file is %q is %d of %d\n", m.PotentiallyFixableIssuesInfo.currentFile, m.PotentiallyFixableIssuesInfo.currentFileIndex+1, len(m.PotentiallyFixableIssuesInfo.FilePaths))
		}

		for m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex < len(m.PotentiallyFixableIssuesInfo.suggestions) {
			var potentialFixableIssue = m.PotentiallyFixableIssuesInfo.suggestions[m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex]
			if m.logFile != nil {
				fmt.Fprintf(m.logFile, "Possible fixable issue %q is %d of %d issues.\n", potentialFixableIssue.Name, m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex+1, len(m.PotentiallyFixableIssuesInfo.suggestions))
			}

			if !m.runAll && (potentialFixableIssue.IsEnabled == nil || *potentialFixableIssue.IsEnabled) {
				if m.logFile != nil {
					fmt.Fprintf(m.logFile, "Skipping possible fixable issue %q with isEnabled set to %v\n", potentialFixableIssue.Name, potentialFixableIssue.IsEnabled)
				}

				m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex++
				continue
			}

			var (
				suggestions = potentialFixableIssue.GetSuggestions(m.PotentiallyFixableIssuesInfo.FileTexts[m.PotentiallyFixableIssuesInfo.currentFileIndex])
			)

			if m.logFile != nil {
				fmt.Fprintf(m.logFile, "Possible fixable issue %q has %d suggestion(s) found\n", potentialFixableIssue.Name, len(suggestions))
			}

			if len(suggestions) != 0 {
				m.PotentiallyFixableIssuesInfo.currentSuggestion = &potentialFixableIssue
				m.PotentiallyFixableIssuesInfo.sectionSuggestionStates = make([]suggestionState, len(suggestions))

				var i = 0
				for original, suggestion := range suggestions {
					var display, err = getStringDiff(original, suggestion)
					if err != nil {
						return nil, err
					}

					m.PotentiallyFixableIssuesInfo.sectionSuggestionStates[i] = suggestionState{
						original:           original,
						originalSuggestion: suggestion,
						currentSuggestion:  suggestion,
						display:            display,
					}

					i++
				}

				m.PotentiallyFixableIssuesInfo.currentSuggestionIndex = 0
				m.PotentiallyFixableIssuesInfo.currentSuggestionState = &m.PotentiallyFixableIssuesInfo.sectionSuggestionStates[0]
				cmd := m.setSuggestionDisplay()
				m.PotentiallyFixableIssuesInfo.currentSuggestionName = potentialFixableIssue.Name

				m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex++

				return cmd, nil
			}

			m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex++
		}

		m.PotentiallyFixableIssuesInfo.currentFileIndex++
		m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex = 0
	}

	return m.exitOrMoveToCssSelection(), nil
}

func (m *FixableIssuesModel) moveToNextSuggestion() (tea.Cmd, error) {
	if m.PotentiallyFixableIssuesInfo.currentSuggestionIndex+1 < len(m.PotentiallyFixableIssuesInfo.sectionSuggestionStates) {
		m.PotentiallyFixableIssuesInfo.currentSuggestionIndex++
		m.PotentiallyFixableIssuesInfo.currentSuggestionState = &m.PotentiallyFixableIssuesInfo.sectionSuggestionStates[m.PotentiallyFixableIssuesInfo.currentSuggestionIndex]

		return m.setSuggestionDisplay(), nil
	}

	return m.setupForNextSuggestions()
}

func (m *FixableIssuesModel) moveToPreviousSuggestion() tea.Cmd {
	if m.PotentiallyFixableIssuesInfo.currentSuggestionIndex > 0 {
		m.PotentiallyFixableIssuesInfo.currentSuggestionIndex--
		m.PotentiallyFixableIssuesInfo.currentSuggestionState = &m.PotentiallyFixableIssuesInfo.sectionSuggestionStates[m.PotentiallyFixableIssuesInfo.currentSuggestionIndex]
		return m.setSuggestionDisplay()
	}

	return nil
}

func (m *FixableIssuesModel) setSuggestionDisplay() tea.Cmd {
	if m.PotentiallyFixableIssuesInfo.currentSuggestionState == nil {
		return nil
	}

	if m.logFile != nil {
		fmt.Fprintf(m.logFile, "current width %d; border width: %d; scrollbar padding: %d\n", m.PotentiallyFixableIssuesInfo.suggestionDisplay.Width, suggestionBorderStyle.GetHorizontalBorderSize(), scrollbarPadding)
	}

	var (
		expectedSuggestionWidth = m.getSuggestionWidth(m.leftStatusView())
		suggestion              = displayStyle.Render(wrapLines(fmt.Sprintf(`"%s"`, m.PotentiallyFixableIssuesInfo.currentSuggestionState.display), expectedSuggestionWidth))
	)

	if m.logFile != nil {
		fmt.Fprintf(m.logFile, "New suggestion getting set with width %d and a value of %q\n", expectedSuggestionWidth, suggestion)
	}

	m.PotentiallyFixableIssuesInfo.suggestionDisplay.SetContent(suggestion)

	// TODO: how do I handle the resizing of the suggestion display
	// one option would be to
	var cmd tea.Cmd
	m.PotentiallyFixableIssuesInfo.scrollbar, cmd = m.PotentiallyFixableIssuesInfo.scrollbar.Update(tui.HeightMsg(m.PotentiallyFixableIssuesInfo.suggestionDisplay.Height))

	return cmd
}
