package main

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/help"
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

type model struct {
	title                       string
	currentStage, width, height int
	stages                      []string
	ready                       bool
	help                        help.Model
}

func newModel() model {
	return model{
		title: "Epub Linter Manually Fixable Issues",
		// title: "EL MFI",
		help: help.New(),
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
			body   = strings.Repeat("\n", max(0, m.height-(lipgloss.Height(header)+lipgloss.Height(footer))+1))
		)

		return header + body + footer
	}

	return ""
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
	}

	return m, nil
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
