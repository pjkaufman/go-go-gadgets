package ui

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
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
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.PotentiallyFixableIssuesInfo.suggestionEdit, cmd = m.PotentiallyFixableIssuesInfo.suggestionEdit.Update(msg)
	cmds = append(cmds, cmd)

	m.PotentiallyFixableIssuesInfo.suggestionDisplay, cmd = m.PotentiallyFixableIssuesInfo.suggestionDisplay.Update(msg)
	cmds = append(cmds, cmd)

	m.PotentiallyFixableIssuesInfo.scrollbar, cmd = m.PotentiallyFixableIssuesInfo.scrollbar.Update(m.PotentiallyFixableIssuesInfo.suggestionDisplay)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.PotentiallyFixableIssuesInfo.isEditing {
			switch msg.String() {
			case "ctrl+s":
				m.PotentiallyFixableIssuesInfo.currentSuggestionState.currentSuggestion = alignWhitespace(m.PotentiallyFixableIssuesInfo.currentSuggestionState.original, m.PotentiallyFixableIssuesInfo.suggestionEdit.Value())
				m.PotentiallyFixableIssuesInfo.isEditing = false

				var err error
				m.PotentiallyFixableIssuesInfo.currentSuggestionState.display, err = getStringDiff(m.PotentiallyFixableIssuesInfo.currentSuggestionState.original, m.PotentiallyFixableIssuesInfo.currentSuggestionState.currentSuggestion)
				if err != nil {
					m.Err = err

					return tea.Quit
				}

				return m.setSuggestionDisplay()
			case "ctrl+e":
				m.PotentiallyFixableIssuesInfo.isEditing = false

				var err error
				m.PotentiallyFixableIssuesInfo.currentSuggestionState.display, err = getStringDiff(m.PotentiallyFixableIssuesInfo.currentSuggestionState.original, m.PotentiallyFixableIssuesInfo.currentSuggestionState.currentSuggestion)
				if err != nil {
					m.Err = err

					return tea.Quit
				}
			case "ctrl+r":
				m.PotentiallyFixableIssuesInfo.suggestionEdit.SetValue(m.PotentiallyFixableIssuesInfo.currentSuggestionState.originalSuggestion)
			}
		} else {
			switch msg.String() {
			case "ctrl+c":
				m.Err = ErrUserKilledProgram

				return tea.Quit
			case "esc":
				return m.exitOrMoveToCssSelection()
			case "right":
				cmd, err := m.moveToNextSuggestion()
				if err != nil {
					m.Err = err

					return tea.Quit
				}

				cmds = append(cmds, cmd)
			case "left":
				cmd = m.moveToPreviousSuggestion()
				cmds = append(cmds, cmd)
			case "c":
				// Copy original value to the clipboard
				// original, err := repairUnicode(m.currentFileState.original)
				// if err != nil {
				// 	m.Err = err

				// 	return m, tea.Quit
				// }

				// err = clipboard.WriteAll(original)
				// TODO: make sure values are utf-8 compliant
				err := clipboard.WriteAll(m.PotentiallyFixableIssuesInfo.currentSuggestionState.original)
				if err != nil {
					m.Err = err

					return tea.Quit
				}
			case "enter":
				if !m.PotentiallyFixableIssuesInfo.currentSuggestionState.isAccepted && m.PotentiallyFixableIssuesInfo.currentSuggestion != nil {
					var replaceCount = 1
					if m.PotentiallyFixableIssuesInfo.currentSuggestion.UpdateAllInstances {
						replaceCount = -1
					}

					m.PotentiallyFixableIssuesInfo.FileTexts[m.PotentiallyFixableIssuesInfo.currentFileIndex] = strings.Replace(m.PotentiallyFixableIssuesInfo.FileTexts[m.PotentiallyFixableIssuesInfo.currentFileIndex], m.PotentiallyFixableIssuesInfo.currentSuggestionState.original, m.PotentiallyFixableIssuesInfo.currentSuggestionState.currentSuggestion, replaceCount)

					m.PotentiallyFixableIssuesInfo.currentSuggestionState.isAccepted = true

					if m.PotentiallyFixableIssuesInfo.currentSuggestion.AddCssSectionBreakIfMissing {
						m.PotentiallyFixableIssuesInfo.AddCssSectionBreakIfMissing = true
						m.PotentiallyFixableIssuesInfo.CssUpdateRequired = true
					} else if m.PotentiallyFixableIssuesInfo.currentSuggestion.AddCssPageBreakIfMissing {
						m.PotentiallyFixableIssuesInfo.AddCssPageBreakIfMissing = true
						m.PotentiallyFixableIssuesInfo.CssUpdateRequired = true
					}

					cmd, err := m.moveToNextSuggestion()
					if err != nil {
						m.Err = err

						return tea.Quit
					}

					cmds = append(cmds, cmd)
				}
			case "e":
				if m.PotentiallyFixableIssuesInfo.currentSuggestionState != nil && !m.PotentiallyFixableIssuesInfo.currentSuggestionState.isAccepted {
					m.PotentiallyFixableIssuesInfo.isEditing = true
					m.PotentiallyFixableIssuesInfo.suggestionEdit.SetValue(m.PotentiallyFixableIssuesInfo.currentSuggestionState.currentSuggestion)

					cmd = m.PotentiallyFixableIssuesInfo.suggestionEdit.Focus()
					cmds = append(cmds, cmd)
				}
			}
		}
	}

	return tea.Batch(cmds...)
}
