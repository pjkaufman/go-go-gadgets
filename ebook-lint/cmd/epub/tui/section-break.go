package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

type SectionBreakModel struct {
	sectionBreakInput textinput.Model
	sectionBreak      string
	Done              bool
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
	if f.Done {
		return f, nil
	}

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

	case error:
		logger.WriteError(msg.Error())
	}

	f.sectionBreakInput, cmd = f.sectionBreakInput.Update(msg)
	return f, cmd
}

func (f SectionBreakModel) View() string {
	if f.Done {
		return ""
	}

	var s strings.Builder

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
