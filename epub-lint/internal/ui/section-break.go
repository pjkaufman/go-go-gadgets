package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m FixableIssuesModel) sectionBreakView() string {
	return m.sectionBreakInfo.input.View()
}

func (m *FixableIssuesModel) handleSectionBreakMsgs(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	m.sectionBreakInfo.input, cmd = m.sectionBreakInfo.input.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
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
		}
	}

	return tea.Batch(cmds...)
}
