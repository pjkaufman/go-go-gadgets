package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type sectionBreak struct {
	state *State
	input textinput.Model
}

func newSectionBreak(focus bool, state *State) sectionBreak {
	ti := textinput.New()
	ti.Width = 20
	ti.CharLimit = 200
	ti.Placeholder = "Section break"

	if focus {
		ti.Focus()
	}

	return sectionBreak{
		state: state,
		input: ti,
	}
}

func (m sectionBreak) Init() tea.Cmd {
	return nil
}

func (m sectionBreak) View() string {
	return m.input.View()
}

func (m sectionBreak) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.state.ContextBreak = strings.TrimSpace(m.input.Value())
			if m.state.ContextBreak != "" {
				cmds = append(cmds, nextStage)
			}
		}
	}

	return m, tea.Batch(cmds...)
}
