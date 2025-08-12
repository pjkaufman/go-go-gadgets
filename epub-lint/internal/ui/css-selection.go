package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m FixableIssuesModel) cssSelectionView() string {
	var s strings.Builder
	s.WriteString("\nSelect the CSS file to modify:\n\n")
	for i, cssFile := range m.CssSelectionInfo.cssFiles {
		cursor := " "
		if m.CssSelectionInfo.currentCssIndex == i {
			cursor = ">"
		}

		s.WriteString(fmt.Sprintf("%s %d. %s\n", cursor, i+1, cssFile))
	}

	s.WriteString("\n")

	return s.String()
}

func (m *FixableIssuesModel) handleCssSelectionMsgs(msg tea.Msg) tea.Cmd {
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
