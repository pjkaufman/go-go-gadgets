package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/acarl005/stripansi"
	"github.com/atotto/clipboard"
	"github.com/muesli/reflow/wordwrap"
)

func (m *FixableIssuesModel) suggestionsView() string {
	return m.suggestionView()
}

func (m FixableIssuesModel) getSuggestionWidth() int {
	return m.body.Width() - (suggestionBorderStyle.GetBorderLeftSize() + suggestionBorderStyle.GetBorderRightSize())
}

func (m *FixableIssuesModel) suggestionHeaderView() string {
	var (
		maxTextWidth           = m.getSuggestionWidth()
		fileName               = m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentFileName
		fileStatus             string
		fileNumberInfo         = fmt.Sprintf("(%d/%d)", m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentFileIndex+1, len(m.PotentiallyFixableIssuesInfo.SuggestionManager.FileSuggestionData))
		remainingFileNameWidth = maxTextWidth - (lipgloss.Width(documentIcon) + 2 + lipgloss.Width(fileNumberInfo) + lipgloss.Width(m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentFileName))
	)

	// truncate file name
	if remainingFileNameWidth < 0 {
		fileName = "..." + fileName[remainingFileNameWidth*-1+3:]
		fileStatus = documentIcon + " " + fileNameStyle.Render(fileName+" "+fileNumberInfo)
	} else {
		fileStatus = fillLine(documentIcon+" "+fileNameStyle.Render(fileName+" "+fileNumberInfo), maxTextWidth)
	}

	var suggestionStatus = wordwrap.String(sectionIcon+" "+m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestionName+fmt.Sprintf(" (%d/%d)", m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentIssueIndex+1, len(m.PotentiallyFixableIssuesInfo.SuggestionManager.Suggestions)), maxTextWidth)
	var lines = strings.Split(suggestionStatus, "\n")

	var afterIcon = strings.Index(lines[0], " ") + 1
	lines[0] = lines[0][:afterIcon] + suggestionNameStyle.Render(lines[0][afterIcon:])
	for i := 1; i < len(lines); i++ {
		lines[i] = suggestionNameStyle.Render(lines[i])
	}

	suggestionStatus = strings.Join(lines, "\n")

	var warningStatus = ""
	if m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestionState != nil && m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestionState.OriginallyHadHalfwidthCircleKatakana {
		warningStatus = wordwrap.String(warningIcon+" "+warningStyle.Render(` At least one instance of "°" in the displayed text may be the Japanese handakuten.`), maxTextWidth)
		lines = strings.Split(warningStatus, "\n")

		var afterIcon = strings.Index(lines[0], " ") + 1
		lines[0] = lines[0][:afterIcon] + warningStyle.Render(lines[0][afterIcon:])
		for i := 1; i < len(lines); i++ {
			lines[i] = warningStyle.Render(lines[i])
		}

		warningStatus = "\n" + strings.Join(lines, "\n")
	}

	return fmt.Sprintf("%s\n%s%s\n%s\n", fileStatus, suggestionStatus, warningStatus, hrStyle.Render(strings.Repeat("─", maxTextWidth)))
}

func (m *FixableIssuesModel) suggestionView() string {
	if m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestionState == nil {
		return "\nNo current file is selected. Something may have gone wrong...\n\n"
	}

	var (
		modeIcon = viewIcon
		modeName = "View"
		s        strings.Builder
	)
	if m.PotentiallyFixableIssuesInfo.isEditing {
		modeIcon = editIcon
		modeName = "Edit"
	}

	var (
		customBorderPadding = 0
		customBorderStyle   lipgloss.Style
		isUsingScrollbar    bool
	)
	if !m.PotentiallyFixableIssuesInfo.isEditing && m.PotentiallyFixableIssuesInfo.suggestionDisplay.TotalLineCount() > m.PotentiallyFixableIssuesInfo.suggestionDisplay.VisibleLineCount() {
		isUsingScrollbar = true
		customBorderPadding = scrollbarPadding
	}

	s.WriteString(m.suggestionHeaderView())

	customBorderStyle = m.createCustomBorder(modeIcon, modeName, customBorderPadding)

	if m.PotentiallyFixableIssuesInfo.isEditing {
		s.WriteString(m.PotentiallyFixableIssuesInfo.suggestionEdit.View())
	} else if isUsingScrollbar {
		s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
			m.PotentiallyFixableIssuesInfo.suggestionDisplay.View(),
			m.PotentiallyFixableIssuesInfo.scrollbar.View().Content,
		))
	} else {
		s.WriteString(m.PotentiallyFixableIssuesInfo.suggestionDisplay.View())
	}

	return customBorderStyle.Render(s.String())
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
				m.Err = m.PotentiallyFixableIssuesInfo.SuggestionManager.UpdateCurrentSuggestionValue(alignWhitespace(m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestionState.Original, m.PotentiallyFixableIssuesInfo.suggestionEdit.Value()))
				if m.Err != nil {
					return tea.Quit
				}

				m.PotentiallyFixableIssuesInfo.isEditing = false
				m.PotentiallyFixableIssuesInfo.suggestionEdit.Blur()
				m.Err = m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestionState.GetStringDiffAsDisplay()
				if m.Err != nil {
					return tea.Quit
				}

				m.recalculateElementSizes(true)

				return tea.Batch(cmds...)
			case "ctrl+e":
				m.PotentiallyFixableIssuesInfo.isEditing = false
				m.PotentiallyFixableIssuesInfo.suggestionEdit.Blur()

				m.Err = m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestionState.GetStringDiffAsDisplay()
				if m.Err != nil {
					return tea.Quit
				}

				m.recalculateElementSizes(false)
			case "ctrl+r":
				var originalSuggestion = m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestionState.ReplaceBrokenDisplayCharacters(m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestionState.OriginalSuggestion)
				m.PotentiallyFixableIssuesInfo.suggestionEdit.SetValue(originalSuggestion)
			case "ctrl+o":
				var original = m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestionState.ReplaceBrokenDisplayCharacters(m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestionState.Original)
				m.PotentiallyFixableIssuesInfo.suggestionEdit.SetValue(original)
			}
		} else {
			switch msg.String() {
			case "ctrl+c":
				m.Err = ErrUserKilledProgram

				return tea.Quit
			case "esc":
				return m.exitOrMoveToCssSelection()
			case "up":
				var yOffset = m.PotentiallyFixableIssuesInfo.suggestionDisplay.YOffset()
				if yOffset <= 0 {
					m.PotentiallyFixableIssuesInfo.suggestionDisplay.SetYOffset(0)
				} else {
					m.PotentiallyFixableIssuesInfo.suggestionDisplay.SetYOffset(yOffset - 1)
				}
			case "down":
				var yOffset = m.PotentiallyFixableIssuesInfo.suggestionDisplay.YOffset()
				if m.logFile != nil {
					fmt.Fprintf(m.logFile, "YOffset before update %d vs total lines %d for content %q\n", yOffset, m.PotentiallyFixableIssuesInfo.suggestionDisplay.TotalLineCount(), m.PotentiallyFixableIssuesInfo.suggestionDisplay.View())
				}
				if yOffset >= m.PotentiallyFixableIssuesInfo.suggestionDisplay.TotalLineCount() {
					m.PotentiallyFixableIssuesInfo.suggestionDisplay.SetYOffset(m.PotentiallyFixableIssuesInfo.suggestionDisplay.TotalLineCount())
				} else {
					m.PotentiallyFixableIssuesInfo.suggestionDisplay.SetYOffset(yOffset + 1)
				}
			case "right":
				cmd, err := m.moveToNextSuggestion()
				if err != nil {
					m.Err = err

					return tea.Quit
				}

				m.recalculateElementSizes(false)

				cmds = append(cmds, cmd)
			case "left":
				m.moveToPreviousSuggestion()
				m.recalculateElementSizes(false)
			case "c":
				// TODO: make sure values are utf-8 compliant
				err := clipboard.WriteAll(m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestionState.Original)
				if err != nil {
					m.Err = err

					return tea.Quit
				}
			case "enter":
				if !m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestionState.IsAccepted && m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestion != nil {
					m.Err = m.PotentiallyFixableIssuesInfo.SuggestionManager.AcceptSuggestion()
					if m.Err != nil {
						return tea.Quit
					}

					if m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestion.AddCssSectionBreakIfMissing {
						m.PotentiallyFixableIssuesInfo.AddCssSectionBreakIfMissing = true
						m.PotentiallyFixableIssuesInfo.CssUpdateRequired = true
					} else if m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestion.AddCssPageBreakIfMissing {
						m.PotentiallyFixableIssuesInfo.AddCssPageBreakIfMissing = true
						m.PotentiallyFixableIssuesInfo.CssUpdateRequired = true
					}

					cmd, err := m.moveToNextSuggestion()
					if err != nil {
						m.Err = err

						return tea.Quit
					}

					m.recalculateElementSizes(false)

					cmds = append(cmds, cmd)
				}
			case "e":
				if m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestionState != nil && !m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestionState.IsAccepted {
					m.PotentiallyFixableIssuesInfo.isEditing = true
					var currentSuggestion = m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestionState.ReplaceBrokenDisplayCharacters(m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestionState.CurrentSuggestion)

					m.PotentiallyFixableIssuesInfo.suggestionEdit.SetValue(currentSuggestion)

					cmd = m.PotentiallyFixableIssuesInfo.suggestionEdit.Focus()
					cmds = append(cmds, cmd)

					m.recalculateElementSizes(false)
				}
			case "ctrl+d":
				cmd, err := m.handleForwardSuggestionUpdate(m.PotentiallyFixableIssuesInfo.SuggestionManager.MoveToNextIssue)
				if err != nil {
					m.Err = err

					return tea.Quit
				}

				cmds = append(cmds, cmd)
			case "ctrl+u":
				if m.PotentiallyFixableIssuesInfo.SuggestionManager.MoveToPreviousIssue() {
					m.recalculateElementSizes(true)
				}
			case "pgdown":
				cmd, err := m.handleForwardSuggestionUpdate(m.PotentiallyFixableIssuesInfo.SuggestionManager.MoveToNextFile)
				if err != nil {
					m.Err = err

					return tea.Quit
				}

				cmds = append(cmds, cmd)
			case "pgup":
				if m.PotentiallyFixableIssuesInfo.SuggestionManager.MoveToPreviousFile() {
					m.recalculateElementSizes(true)
				}
			}
		}
	}

	return tea.Batch(cmds...)
}

func (m *FixableIssuesModel) handleForwardSuggestionUpdate(moveFunc func() (foundSuggestion bool, err error)) (tea.Cmd, error) {
	foundSuggestion, err := moveFunc()
	if err != nil {
		return nil, err
	}

	var cmd tea.Cmd
	if !foundSuggestion {
		cmd = m.exitOrMoveToCssSelection()
	} else {
		m.recalculateElementSizes(true)
	}

	return cmd, nil
}

func (m *FixableIssuesModel) moveToNextSuggestion() (tea.Cmd, error) {
	m.PotentiallyFixableIssuesInfo.suggestionDisplay.SetYOffset(0)
	var hasNextSuggestion = m.PotentiallyFixableIssuesInfo.SuggestionManager.MoveToNextSuggestion()
	if hasNextSuggestion {
		m.setSuggestionDisplay(true)

		return nil, nil
	}

	return m.handleForwardSuggestionUpdate(m.PotentiallyFixableIssuesInfo.SuggestionManager.SetupForNextSuggestions)
}

func (m *FixableIssuesModel) moveToPreviousSuggestion() {
	if m.PotentiallyFixableIssuesInfo.SuggestionManager.MoveToPreviousSuggestion() {
		m.setSuggestionDisplay(true)
	}
}

func (m FixableIssuesModel) createCustomBorder(modeIcon, modeName string, extraWidthPadding int) lipgloss.Style {
	borderConfig := NewBorderConfig(m.PotentiallyFixableIssuesInfo.suggestionDisplay.Height()+borderWidth, m.PotentiallyFixableIssuesInfo.suggestionDisplay.Width()+borderWidth+extraWidthPadding) // +2 for border width/height
	modeInfo := fmt.Sprintf("%s %s", modeIcon, modeName)
	statusInfo := fmt.Sprintf("%d/%d", m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestionIndex+1, len(m.PotentiallyFixableIssuesInfo.SuggestionManager.FileSuggestionData[m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentFileIndex].Suggestions[m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentIssueIndex]))
	borderConfig.SetInfoItems(modeInfo, statusInfo)

	baseBorder := lipgloss.RoundedBorder()
	customBorder := borderConfig.GetBorder(baseBorder)
	customBorderStyle := lipgloss.NewStyle().Border(customBorder)

	return customBorderStyle.BorderForeground(suggestionBorderStyle.GetBorderTopForeground())
}

func (m *FixableIssuesModel) setSuggestionDisplay(resetYOffset bool) {
	if m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestionState == nil {
		return
	}

	if resetYOffset {
		m.PotentiallyFixableIssuesInfo.suggestionDisplay.SetYOffset(0)
	}

	var (
		height         = max(m.body.Height()-leftStatusBorderStyle.GetVerticalBorderSize()-strings.Count(m.suggestionHeaderView(), "\n"), 0)
		remainingWidth = m.getSuggestionWidth()
	)

	if m.logFile != nil {
		fmt.Fprintf(m.logFile, "Body width %d, and a height of %d\n", remainingWidth, height)
	}

	m.PotentiallyFixableIssuesInfo.suggestionDisplay.SetHeight(height)
	m.PotentiallyFixableIssuesInfo.suggestionDisplay.SetWidth(remainingWidth)
	m.PotentiallyFixableIssuesInfo.suggestionEdit.SetWidth(remainingWidth)
	m.PotentiallyFixableIssuesInfo.suggestionEdit.SetHeight(height)

	var (
		expectedSuggestionWidth = m.PotentiallyFixableIssuesInfo.suggestionDisplay.Width()
		suggestion              = m.buildSuggestion(m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestionState.Display, expectedSuggestionWidth)
	)

	m.PotentiallyFixableIssuesInfo.suggestionDisplay.SetContent(suggestion)

	// recalculate pertinent values if the scrollbar does indeed get displayed
	if m.PotentiallyFixableIssuesInfo.suggestionDisplay.TotalLineCount() > m.PotentiallyFixableIssuesInfo.suggestionDisplay.VisibleLineCount() {
		m.PotentiallyFixableIssuesInfo.suggestionDisplay.SetWidth(remainingWidth - scrollbarPadding)

		expectedSuggestionWidth = m.PotentiallyFixableIssuesInfo.suggestionDisplay.Width()
		suggestion = m.buildSuggestion(m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestionState.Display, expectedSuggestionWidth)

		m.PotentiallyFixableIssuesInfo.suggestionDisplay.SetContent(suggestion)
	}

	if m.logFile != nil {
		var (
			original = stripansi.Strip(suggestion)
		)
		fmt.Fprintf(m.logFile, "New suggestion getting set with width %d and a value of %q vs original %q\n", expectedSuggestionWidth, suggestion, original)
		var lines = strings.Split(original, "\n")
		for i, line := range lines {
			fmt.Fprintf(m.logFile, "  Line %d: %q Size: %d\n", i+1, line, lipgloss.Width(line))
		}
	}

	// the cmd here is always nil, so we never need to actually handle it
	m.PotentiallyFixableIssuesInfo.scrollbar, _ = m.PotentiallyFixableIssuesInfo.scrollbar.Update(m.PotentiallyFixableIssuesInfo.suggestionDisplay)
	m.PotentiallyFixableIssuesInfo.scrollbar, _ = m.PotentiallyFixableIssuesInfo.scrollbar.Update(HeightMsg(m.PotentiallyFixableIssuesInfo.suggestionDisplay.Height()))
}

func (m *FixableIssuesModel) buildSuggestion(displayText string, expectedSuggestionWidth int) string {
	//nolint:gocritic // can include ansi escape codes so we should ignore this issue here
	text := fmt.Sprintf(`"%s"`, displayText) // includes ANSI
	text = m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestionState.ReplaceBrokenDisplayCharacters(text)
	return displayStyle.Width(expectedSuggestionWidth).Render(text)
}

func replaceBrokenDisplayCharacters(text string) (string, bool) {
	// text with handakuten in them are not having their width calculated correctly, so I will just remove them
	// and we can display a warning if need bee
	if strings.Contains(text, "ﾟ") {
		return strings.ReplaceAll(text, "ﾟ", "°"), true
	}

	return text, false
}
