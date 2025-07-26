package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m FixableIssuesModel) suggestionsView() string {
	var (
		status         = m.leftStatusView()
		height         = m.body.Height - leftStatusBorderStyle.GetVerticalBorderSize()
		remainingWidth = m.body.Width - (lipgloss.Width(status) + leftStatusBorderStyle.GetBorderLeftSize() + suggestionBorderStyle.GetBorderRightSize())
	)

	m.PotentiallyFixableIssuesInfo.suggestionDisplay.Height = height
	m.PotentiallyFixableIssuesInfo.suggestionDisplay.Width = remainingWidth
	m.PotentiallyFixableIssuesInfo.suggestionEdit.SetWidth(remainingWidth)
	m.PotentiallyFixableIssuesInfo.suggestionEdit.SetHeight(height)

	return lipgloss.JoinHorizontal(lipgloss.Left, status, m.suggestionView())
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
		// var (
		// 	suggestion = displayStyle.Width(m.width - columnPadding).Render(fmt.Sprintf(`"%s"`, m.PotentiallyFixableIssuesInfo.currentSuggestionState.display))
		// 	// expectedHeight = strings.Count(suggestion, "\n") + 1
		// )
		// if m.PotentiallyFixableIssuesInfo.isEditing {
		// 	s.WriteString("\nEditing suggestion:\n\n")
		// } else if m.PotentiallyFixableIssuesInfo.currentSuggestionState.isAccepted {
		// 	s.WriteString("\n" + acceptedChangeTitleStyle.Render("Accepted change:") + "\n")
		// } else {
		// 	s.WriteString(fmt.Sprintf("\nSuggested change (%d/%d):\n", expectedHeight, maxDisplayHeight))
		// }

		if m.PotentiallyFixableIssuesInfo.isEditing {
			s.WriteString(suggestionBorderStyle.Render(m.PotentiallyFixableIssuesInfo.suggestionEdit.View()) + "\n\n")
		} else if m.PotentiallyFixableIssuesInfo.suggestionDisplay.TotalLineCount() > m.PotentiallyFixableIssuesInfo.suggestionDisplay.VisibleLineCount() {
			s.WriteString(suggestionBorderStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top,
				m.PotentiallyFixableIssuesInfo.suggestionDisplay.View(),
				m.PotentiallyFixableIssuesInfo.scrollbar.View(),
			)))
		} else {
			s.WriteString(suggestionBorderStyle.Render(m.PotentiallyFixableIssuesInfo.suggestionDisplay.View()))
		}

		// s.WriteString(fmt.Sprintf("\033[0m\n\nSuggestion %d of %d.\n\n", m.PotentiallyFixableIssuesInfo.currentSuggestionIndex+1, len(m.PotentiallyFixableIssuesInfo.sectionSuggestionStates)))
		// }
	}

	return s.String()
}
