package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type SectionBreakModel struct {
	sectionBreakInput textinput.Model
	sectionBreak      string
	Done              bool
	Err               error
}

func NewSectionBreakModel() SectionBreakModel {
	ti := textinput.New()
	ti.Placeholder = "Section Break"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 20

	return SectionBreakModel{
		sectionBreakInput: ti,
	}
}

func (f SectionBreakModel) Init() tea.Cmd {
	return nil
}

func (f SectionBreakModel) Update(msg tea.Msg) (SectionBreakModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			f.sectionBreak = f.sectionBreakInput.Value()

			if strings.TrimSpace(f.sectionBreak) != "" {
				f.Done = true
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			return f, tea.Quit
		}
	case tea.WindowSizeMsg:
		return f, tea.ClearScreen
	case error:
		f.Err = msg

		return f, nil
	}

	f.sectionBreakInput, cmd = f.sectionBreakInput.Update(msg)
	return f, cmd
}

func (f SectionBreakModel) View() string {
	var s strings.Builder
	clearScreen(&s)

	s.WriteString("What is the section break for the epub?\n\n")
	s.WriteString(f.sectionBreakInput.View())
	s.WriteString("\n\n")

	f.displaySectionBreakControls(&s)

	return s.String()
}

func (f SectionBreakModel) displaySectionBreakControls(s *strings.Builder) {
	s.WriteString(groupStyle.Render("Controls:") + "\n")
	s.WriteString("Enter: Continue   Ctrl+C/Esc: Quit\n")
}

func (f SectionBreakModel) Value() string {
	return f.sectionBreak
}
