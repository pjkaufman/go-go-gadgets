package main

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// custom messages/commands
type advanceStage struct{}

func nextStage() tea.Msg {
	return advanceStage{}
}

// model
type model struct {
	title                       string
	currentStage, width, height int
	stages                      []string
	bodyContent                 []tea.Model
	ready                       bool
	body                        viewport.Model
	help                        help.Model
}

func newModel() model {
	return model{
		title: "Epub Linter Manually Fixable Issues",
		// title: "EL MFI",
		bodyContent: []tea.Model{
			newSectionBreak(true),
			newSuggestions(),
		},
		help: help.New(),
		body: viewport.New(0, 0),
		stages: []string{
			"Section Break",
			"Suggestions",
			"CSS File Selection",
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	if m.ready {
		var (
			header = m.headerView() + "\n"
			footer = m.footerView()
		)

		m.body.SetYOffset(lipgloss.Height(header))
		m.body.Width = m.width
		m.body.Height = max(0, m.height-(lipgloss.Height(header)+lipgloss.Height(footer))+1)
		m.body.SetContent(m.bodyContent[m.currentStage].View())

		return header + m.body.View() + footer
	}

	return ""
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tea.Quit
		case "ctrl+c":
			// TODO: make sure this is an error once ready
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
	case advanceStage:
		// TODO: add more in depth logic
		m.currentStage++
	}

	var cmd tea.Cmd
	m.bodyContent[m.currentStage], cmd = m.bodyContent[m.currentStage].Update(msg)

	return m, cmd
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// styles
var (
	titleStyle         = lipgloss.NewStyle().Bold(true)
	inactiveStyle      = lipgloss.NewStyle().Faint(true)
	activeStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	titleBorder        = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderBottom(true)
	controlsStyle      = lipgloss.NewStyle().Faint(true).Bold(true)
	controlBorderStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderTop(true)
)

func (m model) headerView() string {
	return titleBorder.Render(fillLine(lipgloss.JoinHorizontal(lipgloss.Center, titleStyle.Render(m.title), " | ", m.getStageHeaders()), m.width)) + "\n"
}

func (m model) getStageHeaders() string {
	var stageHeaders = make([]string, len(m.stages))

	var style lipgloss.Style
	for i, header := range m.stages {
		if i == m.currentStage {
			style = activeStyle
		} else {
			style = inactiveStyle
		}

		stageHeaders[i] = style.Render(header)
	}

	return strings.Join(stageHeaders, inactiveStyle.Render(" > "))
}

func (m model) footerView() string {
	var s strings.Builder
	s.WriteString(fillLine(controlsStyle.Render("Controls:"), m.width) + "\n")

	var controls []string
	switch m.currentStage {
	case 0:
		controls = []string{
			"Enter: Accept",
			"Esc: Quit",
			"Ctrl+C: Exit without saving",
		}
	case 1:
		// if m.potentiallyFixableIssuesInfo.isEditing {
		// 	controls = []string{
		// 		"Ctrl+R: Reset",
		// 		"Ctrl+E: Cancel edit",
		// 		"Ctrl+S: Accept",
		// 		"Esc: Quit",
		// 		"Ctrl+C: Exit without saving",
		// 	}
		// } else if m.potentiallyFixableIssuesInfo.currentSuggestionState != nil && m.potentiallyFixableIssuesInfo.currentSuggestionState.isAccepted {
		controls = []string{
			"← / → : Previous/Next Suggestion",
			"C: Copy",
			"Esc: Quit",
			"Ctrl+C: Exit without saving",
		}
		// } else {
		// 	controls = []string{
		// 		"← / → : Previous/Next Suggestion",
		// 		"E: Edit",
		// 		"C: Copy",
		// 		"Enter: Accept",
		// 		"Esc: Quit",
		// 		"Ctrl+C: Exit without saving",
		// 	}
		// }
	case 2:
		controls = []string{
			"↑ / ↓ : Previous/Next Suggestion",
			"Enter: Accept",
			"Ctrl+C: Exit without saving",
		}
	}

	var (
		line     strings.Builder
		maxWidth = m.width
	)
	for _, help := range controls {
		if line.Len() == 0 {
			line.WriteString(help)
			s.WriteString(help)
		} else if line.Len()+len(help)+3 <= maxWidth {
			s.WriteString(" • " + help)
			line.WriteString(" • " + help)
		} else {
			s.WriteString("\n")
			line.Reset()

			line.WriteString(help)
			s.WriteString(help)
		}
	}

	if line.Len() != 0 {
		s.WriteString("\n")
	}

	return controlBorderStyle.Render(s.String())
}

func fillLine(currentValue string, width int) string {
	var amountToFill = width - lipgloss.Width(currentValue)
	if amountToFill < 1 {
		return currentValue
	}

	return currentValue + strings.Repeat(" ", amountToFill)
}

// Section Break
type sectionBreak struct {
	input        textinput.Model
	contextBreak string
}

func newSectionBreak(focus bool) sectionBreak {
	ti := textinput.New()
	ti.Width = 20
	ti.CharLimit = 200
	ti.Placeholder = "Section break"

	if focus {
		ti.Focus()
	}

	return sectionBreak{
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
			m.contextBreak = strings.TrimSpace(m.input.Value())
			if m.contextBreak != "" {
				cmds = append(cmds, nextStage)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

// Suggestions
type suggestions struct {
}

func newSuggestions() suggestions {
	return suggestions{}
}

func (m suggestions) Init() tea.Cmd {
	return nil
}

func (m suggestions) View() string {
	return "Suggestions"
}

func (m suggestions) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
