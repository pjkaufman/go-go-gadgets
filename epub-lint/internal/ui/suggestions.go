package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

const maxLeftStatusWidth = 40

func (m *FixableIssuesModel) suggestionsView() string {
	return lipgloss.JoinHorizontal(lipgloss.Left, m.leftStatusView(), m.suggestionView())
}

func (m FixableIssuesModel) getSuggestionWidth(statusWidth int) int {
	return m.body.Width - (statusWidth + leftStatusBorderStyle.GetBorderLeftSize() + suggestionBorderStyle.GetBorderRightSize())
}

func (m FixableIssuesModel) getLeftStatusWidth() int {
	return min(lipgloss.Width(m.leftStatusView()), maxLeftStatusWidth)
}

func (m *FixableIssuesModel) leftStatusView() string {
	var (
		maxTextWidth           = maxLeftStatusWidth - leftStatusBorderStyle.GetHorizontalBorderSize()
		fileName               = m.PotentiallyFixableIssuesInfo.currentFile
		fileStatus             string
		fileNumberInfo         = fmt.Sprintf("(%d/%d)", m.PotentiallyFixableIssuesInfo.currentFileIndex+1, len(m.PotentiallyFixableIssuesInfo.FileSuggestionData))
		remainingFileNameWidth = maxTextWidth - (lipgloss.Width(documentIcon) + 2 + lipgloss.Width(fileNumberInfo) + lipgloss.Width(m.PotentiallyFixableIssuesInfo.currentFile))
	)

	// truncate file name
	if remainingFileNameWidth < 0 {
		fileName = "..." + fileName[remainingFileNameWidth*-1+3:]
		fileStatus = documentIcon + " " + fileNameStyle.Render(fileName+" "+fileNumberInfo)
	} else {
		fileStatus = fillLine(documentIcon+" "+fileNameStyle.Render(fileName+" "+fileNumberInfo), maxTextWidth)
	}

	var suggestionStatus = wordwrap.String(sectionIcon+" "+m.PotentiallyFixableIssuesInfo.currentSuggestionName+fmt.Sprintf(" (%d/%d)", m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex+1, len(m.PotentiallyFixableIssuesInfo.suggestions)), maxTextWidth)
	var lines = strings.Split(suggestionStatus, "\n")

	var afterIcon = strings.Index(lines[0], " ") + 1
	lines[0] = lines[0][:afterIcon] + suggestionNameStyle.Render(lines[0][afterIcon:])
	for i := 1; i < len(lines); i++ {
		lines[i] = suggestionNameStyle.Render(lines[i])
	}

	suggestionStatus = strings.Join(lines, "\n")

	var (
		statusView      = fileStatus + "\n" + suggestionStatus
		remainingHeight int
		statusPadding   string
	)

	remainingHeight = m.body.Height - (lipgloss.Height(statusView) + leftStatusBorderStyle.GetVerticalBorderSize())

	if remainingHeight > 0 {
		statusPadding = strings.Repeat("\n", remainingHeight)
	}

	return leftStatusBorderStyle.Render(statusView + statusPadding)
}

func (m *FixableIssuesModel) suggestionView() string {
	var s strings.Builder

	if m.PotentiallyFixableIssuesInfo.currentSuggestionState == nil {
		s.WriteString("\nNo current file is selected. Something may have gone wrong...\n\n")
	} else {
		modeIcon := viewIcon
		modeName := "View"
		if m.PotentiallyFixableIssuesInfo.isEditing {
			modeIcon = editIcon
			modeName = "Edit"
		}

		if m.PotentiallyFixableIssuesInfo.isEditing {
			var customBorderStyle = m.createCustomBorder(modeIcon, modeName, 0)

			s.WriteString(customBorderStyle.Render(m.PotentiallyFixableIssuesInfo.suggestionEdit.View()) + "\n\n")
		} else if m.PotentiallyFixableIssuesInfo.suggestionDisplay.TotalLineCount() > m.PotentiallyFixableIssuesInfo.suggestionDisplay.VisibleLineCount() {
			var customBorderStyle = m.createCustomBorder(modeIcon, modeName, scrollbarPadding)

			s.WriteString(customBorderStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top,
				m.PotentiallyFixableIssuesInfo.suggestionDisplay.View(),
				m.PotentiallyFixableIssuesInfo.scrollbar.View(),
			)))
		} else {
			var customBorderStyle = m.createCustomBorder(modeIcon, modeName, 0)

			s.WriteString(customBorderStyle.Render(m.PotentiallyFixableIssuesInfo.suggestionDisplay.View()))
		}
	}

	return s.String()
}

func (m *FixableIssuesModel) handleSuggestionMsgs(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.PotentiallyFixableIssuesInfo.suggestionEdit, cmd = m.PotentiallyFixableIssuesInfo.suggestionEdit.Update(msg)
	cmds = append(cmds, cmd)

	m.PotentiallyFixableIssuesInfo.suggestionDisplay, cmd = m.PotentiallyFixableIssuesInfo.suggestionDisplay.Update(msg)
	cmds = append(cmds, cmd)

	m.PotentiallyFixableIssuesInfo.scrollbar, cmd = m.PotentiallyFixableIssuesInfo.scrollbar.Update(msg)
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

				m.setSuggestionDisplay(true)

				return tea.Batch(cmds...)
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
			case "up":
				if m.PotentiallyFixableIssuesInfo.suggestionDisplay.YOffset <= 0 {
					m.PotentiallyFixableIssuesInfo.suggestionDisplay.YOffset = 0
				} else {
					m.PotentiallyFixableIssuesInfo.suggestionDisplay.YOffset--
				}
			case "down":
				if m.logFile != nil {
					fmt.Fprintf(m.logFile, "YOffset before update %d vs total lines %d for content %q\n", m.PotentiallyFixableIssuesInfo.suggestionDisplay.YOffset, m.PotentiallyFixableIssuesInfo.suggestionDisplay.TotalLineCount(), m.PotentiallyFixableIssuesInfo.suggestionDisplay.View())
				}
				if m.PotentiallyFixableIssuesInfo.suggestionDisplay.YOffset >= m.PotentiallyFixableIssuesInfo.suggestionDisplay.TotalLineCount() {
					m.PotentiallyFixableIssuesInfo.suggestionDisplay.YOffset = m.PotentiallyFixableIssuesInfo.suggestionDisplay.TotalLineCount()
				} else {
					m.PotentiallyFixableIssuesInfo.suggestionDisplay.YOffset++
				}
			case "right":
				cmd, err := m.moveToNextSuggestion()
				if err != nil {
					m.Err = err

					return tea.Quit
				}

				cmds = append(cmds, cmd)
			case "left":
				m.moveToPreviousSuggestion()
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

					m.PotentiallyFixableIssuesInfo.FileSuggestionData[m.PotentiallyFixableIssuesInfo.currentFileIndex].Text = strings.Replace(m.PotentiallyFixableIssuesInfo.FileSuggestionData[m.PotentiallyFixableIssuesInfo.currentFileIndex].Text, m.PotentiallyFixableIssuesInfo.currentSuggestionState.original, m.PotentiallyFixableIssuesInfo.currentSuggestionState.currentSuggestion, replaceCount)

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

func (m *FixableIssuesModel) setupForNextSuggestions() (tea.Cmd, error) {
	if m.logFile != nil {
		fmt.Fprintln(m.logFile, "Getting next suggestions")
	}

	m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex++

	for m.PotentiallyFixableIssuesInfo.currentFileIndex < len(m.PotentiallyFixableIssuesInfo.FileSuggestionData) {
		m.PotentiallyFixableIssuesInfo.currentFile = m.PotentiallyFixableIssuesInfo.FileSuggestionData[m.PotentiallyFixableIssuesInfo.currentFileIndex].Name
		if m.logFile != nil {
			fmt.Fprintf(m.logFile, "Current file is %q is %d of %d\n", m.PotentiallyFixableIssuesInfo.currentFile, m.PotentiallyFixableIssuesInfo.currentFileIndex+1, len(m.PotentiallyFixableIssuesInfo.FileSuggestionData))
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
			} else if m.skipCss && (potentialFixableIssue.AddCssPageBreakIfMissing || potentialFixableIssue.AddCssSectionBreakIfMissing) {
				if m.logFile != nil {
					fmt.Fprintf(m.logFile, "Skipping possible fixable issue %q because css related rules are to be skipped\n", potentialFixableIssue.Name)
				}

				m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex++
				continue
			}

			if len(m.PotentiallyFixableIssuesInfo.FileSuggestionData[m.PotentiallyFixableIssuesInfo.currentFileIndex].Suggestions[m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex]) != 0 {
				if m.logFile != nil {
					fmt.Fprintf(m.logFile, "Possible fixable issue %q has %d suggestion(s) already\n", potentialFixableIssue.Name, len(m.PotentiallyFixableIssuesInfo.FileSuggestionData[m.PotentiallyFixableIssuesInfo.currentFileIndex].Suggestions[m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex]))
				}

				m.PotentiallyFixableIssuesInfo.currentSuggestionIndex = 0
				m.PotentiallyFixableIssuesInfo.currentSuggestionState = &m.PotentiallyFixableIssuesInfo.FileSuggestionData[m.PotentiallyFixableIssuesInfo.currentFileIndex].Suggestions[m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex][0]
				m.PotentiallyFixableIssuesInfo.currentSuggestionName = potentialFixableIssue.Name
				m.setSuggestionDisplay(true)

				return nil, nil
			}

			var (
				suggestions = potentialFixableIssue.GetSuggestions(m.PotentiallyFixableIssuesInfo.FileSuggestionData[m.PotentiallyFixableIssuesInfo.currentFileIndex].Text)
			)

			if m.logFile != nil {
				fmt.Fprintf(m.logFile, "Possible fixable issue %q has %d suggestion(s) found\n", potentialFixableIssue.Name, len(suggestions))
			}

			if len(suggestions) != 0 {
				m.PotentiallyFixableIssuesInfo.currentSuggestion = &potentialFixableIssue
				m.PotentiallyFixableIssuesInfo.FileSuggestionData[m.PotentiallyFixableIssuesInfo.currentFileIndex].Suggestions[m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex] = make([]SuggestionState, len(suggestions))

				var i = 0
				for original, suggestion := range suggestions {
					var display, err = getStringDiff(original, suggestion)
					if err != nil {
						return nil, err
					}

					m.PotentiallyFixableIssuesInfo.FileSuggestionData[m.PotentiallyFixableIssuesInfo.currentFileIndex].Suggestions[m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex][i] = SuggestionState{
						original:           original,
						originalSuggestion: suggestion,
						currentSuggestion:  suggestion,
						display:            display,
					}

					i++
				}

				sort.Slice(m.PotentiallyFixableIssuesInfo.FileSuggestionData[m.PotentiallyFixableIssuesInfo.currentFileIndex].Suggestions[m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex], func(i, j int) bool {
					return m.PotentiallyFixableIssuesInfo.FileSuggestionData[m.PotentiallyFixableIssuesInfo.currentFileIndex].Suggestions[m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex][i].original < m.PotentiallyFixableIssuesInfo.FileSuggestionData[m.PotentiallyFixableIssuesInfo.currentFileIndex].Suggestions[m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex][j].original
				})

				m.PotentiallyFixableIssuesInfo.currentSuggestionIndex = 0
				m.PotentiallyFixableIssuesInfo.currentSuggestionState = &m.PotentiallyFixableIssuesInfo.FileSuggestionData[m.PotentiallyFixableIssuesInfo.currentFileIndex].Suggestions[m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex][0]
				m.PotentiallyFixableIssuesInfo.currentSuggestionName = potentialFixableIssue.Name
				m.setSuggestionDisplay(true)

				return nil, nil
			}

			m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex++
		}

		m.PotentiallyFixableIssuesInfo.currentFileIndex++
		m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex = 0
	}

	return m.exitOrMoveToCssSelection(), nil
}

func (m *FixableIssuesModel) moveToNextSuggestion() (tea.Cmd, error) {
	m.PotentiallyFixableIssuesInfo.suggestionDisplay.YOffset = 0
	if m.PotentiallyFixableIssuesInfo.currentSuggestionIndex+1 < len(m.PotentiallyFixableIssuesInfo.FileSuggestionData[m.PotentiallyFixableIssuesInfo.currentFileIndex].Suggestions[m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex]) {
		m.PotentiallyFixableIssuesInfo.currentSuggestionIndex++
		m.PotentiallyFixableIssuesInfo.currentSuggestionState = &m.PotentiallyFixableIssuesInfo.FileSuggestionData[m.PotentiallyFixableIssuesInfo.currentFileIndex].Suggestions[m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex][m.PotentiallyFixableIssuesInfo.currentSuggestionIndex]
		m.setSuggestionDisplay(true)

		return nil, nil
	}

	return m.setupForNextSuggestions()
}

func (m *FixableIssuesModel) moveToPreviousSuggestion() {
	if m.PotentiallyFixableIssuesInfo.currentSuggestionIndex > 0 {
		m.PotentiallyFixableIssuesInfo.currentSuggestionIndex--
		m.PotentiallyFixableIssuesInfo.currentSuggestionState = &m.PotentiallyFixableIssuesInfo.FileSuggestionData[m.PotentiallyFixableIssuesInfo.currentFileIndex].Suggestions[m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex][m.PotentiallyFixableIssuesInfo.currentSuggestionIndex]
		m.setSuggestionDisplay(true)

		return
	}

	var (
		originalCurrentFileIndex           = m.PotentiallyFixableIssuesInfo.currentFileIndex
		originalPotentialFixableIssueIndex = m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex
	)
	for m.PotentiallyFixableIssuesInfo.currentFileIndex != 0 || m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex != 0 {
		if m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex == 0 {
			m.PotentiallyFixableIssuesInfo.currentFileIndex--

			m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex = len(m.PotentiallyFixableIssuesInfo.suggestions) - 1
		} else {
			m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex--
		}

		for m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex > 0 && len(m.PotentiallyFixableIssuesInfo.FileSuggestionData[m.PotentiallyFixableIssuesInfo.currentFileIndex].Suggestions[m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex]) == 0 {
			m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex--
		}

		var numSuggestions = len(m.PotentiallyFixableIssuesInfo.FileSuggestionData[m.PotentiallyFixableIssuesInfo.currentFileIndex].Suggestions[m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex])
		if numSuggestions != 0 {
			m.PotentiallyFixableIssuesInfo.currentSuggestionIndex = numSuggestions - 1
			m.PotentiallyFixableIssuesInfo.currentSuggestionState = &m.PotentiallyFixableIssuesInfo.FileSuggestionData[m.PotentiallyFixableIssuesInfo.currentFileIndex].Suggestions[m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex][m.PotentiallyFixableIssuesInfo.currentSuggestionIndex]
			m.setSuggestionDisplay(true)

			return
		}
	}

	m.PotentiallyFixableIssuesInfo.currentFileIndex = originalCurrentFileIndex
	m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex = originalPotentialFixableIssueIndex
}

func (m FixableIssuesModel) createCustomBorder(modeIcon, modeName string, extraWidthPadding int) lipgloss.Style {
	borderConfig := NewBorderConfig(m.PotentiallyFixableIssuesInfo.suggestionDisplay.Height+2, m.PotentiallyFixableIssuesInfo.suggestionDisplay.Width+2+extraWidthPadding) // +2 for border width/height
	modeInfo := fmt.Sprintf("%s %s", modeIcon, modeName)
	statusInfo := fmt.Sprintf("%d/%d", m.PotentiallyFixableIssuesInfo.currentSuggestionIndex+1, len(m.PotentiallyFixableIssuesInfo.FileSuggestionData[m.PotentiallyFixableIssuesInfo.currentFileIndex].Suggestions[m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex]))
	borderConfig.SetInfoItems(modeInfo, statusInfo)

	baseBorder := lipgloss.RoundedBorder()
	customBorder := borderConfig.GetBorder(baseBorder)
	customBorderStyle := lipgloss.NewStyle().Border(customBorder)

	return customBorderStyle.BorderForeground(suggestionBorderStyle.GetBorderTopForeground())
}

func (m *FixableIssuesModel) setSuggestionDisplay(resetYOffset bool) {
	if m.PotentiallyFixableIssuesInfo.currentSuggestionState == nil {
		return
	}

	if resetYOffset {
		m.PotentiallyFixableIssuesInfo.suggestionDisplay.YOffset = 0
	}

	var (
		height         = m.body.Height - leftStatusBorderStyle.GetVerticalBorderSize()
		remainingWidth = m.getSuggestionWidth(m.getLeftStatusWidth())
	)

	if m.logFile != nil {
		fmt.Fprintf(m.logFile, "Status width %d, body width %d, and a height of %d\n", m.getLeftStatusWidth(), remainingWidth, height)
	}

	m.PotentiallyFixableIssuesInfo.suggestionDisplay.Height = height
	m.PotentiallyFixableIssuesInfo.suggestionDisplay.Width = remainingWidth
	m.PotentiallyFixableIssuesInfo.suggestionEdit.SetWidth(remainingWidth)
	m.PotentiallyFixableIssuesInfo.suggestionEdit.SetHeight(height)

	var (
		expectedSuggestionWidth = m.PotentiallyFixableIssuesInfo.suggestionDisplay.Width
		// this does not currently handle ansi code wrapping
		suggestion = displayStyle.Render(wordwrap.String(fmt.Sprintf(`"%s"`, m.PotentiallyFixableIssuesInfo.currentSuggestionState.display), expectedSuggestionWidth))
	)

	m.PotentiallyFixableIssuesInfo.suggestionDisplay.SetContent(suggestion)

	// recalculate pertinent values if the scrollbar does indeed get displayed
	if m.PotentiallyFixableIssuesInfo.suggestionDisplay.TotalLineCount() > m.PotentiallyFixableIssuesInfo.suggestionDisplay.VisibleLineCount() {
		m.PotentiallyFixableIssuesInfo.suggestionDisplay.Width -= scrollbarPadding

		expectedSuggestionWidth = m.PotentiallyFixableIssuesInfo.suggestionDisplay.Width
		// this does not currently handle ansi code wrapping
		suggestion = displayStyle.Render(wordwrap.String(fmt.Sprintf(`"%s"`, m.PotentiallyFixableIssuesInfo.currentSuggestionState.display), expectedSuggestionWidth))

		m.PotentiallyFixableIssuesInfo.suggestionDisplay.SetContent(suggestion)
	}

	if m.logFile != nil {
		fmt.Fprintf(m.logFile, "New suggestion getting set with width %d and a value of %q\n", expectedSuggestionWidth, suggestion)
	}

	// the cmd here is always nil, so we never need to actually handle it
	m.PotentiallyFixableIssuesInfo.scrollbar, _ = m.PotentiallyFixableIssuesInfo.scrollbar.Update(m.PotentiallyFixableIssuesInfo.suggestionDisplay)
	m.PotentiallyFixableIssuesInfo.scrollbar, _ = m.PotentiallyFixableIssuesInfo.scrollbar.Update(HeightMsg(m.PotentiallyFixableIssuesInfo.suggestionDisplay.Height))
}
