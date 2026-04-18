package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (m FixableIssuesModel) sectionBreakView() string {
	var sectionBreakDisplay = m.sectionBreakInfo.input.View()
	if !m.sectionBreakInfo.pasteFailed {
		return sectionBreakDisplay
	}

	return warningStyle.Render(warningIcon+" The previous paste attempt failed. Please try again.") + "\n" + sectionBreakDisplay
}

func (m *FixableIssuesModel) handleSectionBreakMsgs(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	oldValue := m.sectionBreakInfo.input.Value()
	m.sectionBreakInfo.input, cmd = m.sectionBreakInfo.input.Update(msg)
	newValue := m.sectionBreakInfo.input.Value()
	cmds = append(cmds, cmd)

	if oldValue != newValue { // only chide the warning once an actual change has been made to the text or another paste happens
		m.sectionBreakInfo.pasteFailed = false
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.sectionBreakInfo.pasteFailed && m.logFile != nil {
			fmt.Fprintln(m.logFile, "Next change after failed paste")
		}

		m.sectionBreakInfo.isPasting = false

		switch msg.String() {
		case "enter":
			*m.sectionBreakInfo.contextBreak = strings.TrimSpace(m.sectionBreakInfo.input.Value())
			if *m.sectionBreakInfo.contextBreak != "" {
				m.currentStage = suggestionsProcessing

				cmd, err := m.setupForNextSuggestions()
				if err != nil {
					m.Err = err

					return tea.Quit
				}

				cmds = append(cmds, cmd)
			}
		case "ctrl+v":
			if m.logFile != nil {
				fmt.Fprintln(m.logFile, "Starting Paste")
			}

			m.sectionBreakInfo.pasteFailed = false
			m.sectionBreakInfo.isPasting = true
		}
	}

	return tea.Batch(cmds...)
}
