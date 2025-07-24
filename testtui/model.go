package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	title                       string
	currentStage, width, height int
	bodyHeight, bodyWidth       *int
	stages                      []string
	bodyContent                 []tea.Model
	ready                       bool
	body                        viewport.Model
	help                        help.Model
}

func newModel() model {
	var height, width int

	return model{
		title: "Epub Linter Manually Fixable Issues",
		bodyContent: []tea.Model{
			newSectionBreak(true),
			newSuggestions(&height, &width),
		},
		help:       help.New(),
		body:       viewport.New(0, 0),
		bodyHeight: &height,
		bodyWidth:  &width,
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
			header       = m.headerView()
			footer       = m.footerView()
			headerHeight = lipgloss.Height(header) + headerBorderStyle.GetBorderBottomSize()
			footerHeight = lipgloss.Height(footer) + footerBorderStyle.GetBorderTopSize()
		)

		m.body.Width = m.width
		m.body.Height = max(0, m.height-(headerHeight+footerHeight)+2)
		// preserve pointer value to allow for proper seinding of the value to any views that need it
		// using a new pointer would break the references in the body views
		(*m.bodyHeight) = m.body.Height
		(*m.bodyWidth) = m.body.Width
		m.body.SetContent(m.bodyContent[m.currentStage].View())

		return lipgloss.JoinVertical(lipgloss.Center, header, m.body.View(), footer)
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

	var cmds []tea.Cmd
	currentStage, cmd := m.bodyContent[m.currentStage].Update(msg)
	m.bodyContent[m.currentStage] = currentStage
	cmds = append(cmds, cmd)

	m.help, cmd = m.help.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func (m model) headerView() string {
	return headerBorderStyle.Render(fillLine(lipgloss.JoinHorizontal(lipgloss.Center, titleStyle.Render(m.title), " | ", m.getStageHeaders()), m.width-headerBorderStyle.GetHorizontalBorderSize()))
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
	// return footerBorderStyle.Render(fillLine(controlsStyle.Render("Controls:"), m.width-footerBorderStyle.GetHorizontalBorderSize()) + "\n" + m.help.View(m.bodyContent[m.currentStage].HelpKeys()))
	var s strings.Builder
	s.WriteString(fillLine(controlsStyle.Render("Controls:"), m.width-footerBorderStyle.GetHorizontalBorderSize()) + "\n")

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
		// "Ctrl+R: Reset",
		// 	"Ctrl+E: Cancel edit",
		// 	"Ctrl+S: Accept",
		// 	"Esc: Quit",
		// 	"Ctrl+C: Exit without saving",
		// 	}
		// } else if m.potentiallyFixableIssuesInfo.currentSuggestionState != nil && m.potentiallyFixableIssuesInfo.currentSuggestionState.isAccepted {
		// controls = []string{
		// 	"← / → : Previous/Next Suggestion",
		// 	"C: Copy",
		// 	"Esc: Quit",
		// 	"Ctrl+C: Exit without saving",
		// }
		// } else {
		controls = []string{
			"← / → : Previous/Next Suggestion",
			"E: Edit",
			"C: Copy",
			"Enter: Accept",
			"Esc: Quit",
			"Ctrl+C: Exit without saving",
		}
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

	return footerBorderStyle.Render(s.String())
}

// type modelKeyMap struct {
// 	// base keys
// 	Accept     key.Binding
// 	Quit       key.Binding
// 	ExitNoSave key.Binding
// 	// edit suggestion keys
// 	ResetSuggestion key.Binding
// 	CancelEdit      key.Binding
// 	AcceptEdit      key.Binding
// 	// "Ctrl+R: Reset",
// 	// "Ctrl+E: Cancel edit",
// 	// "Ctrl+S: Accept",
// }

// // ShortHelp returns keybindings to be shown in the mini help view. It's part
// // of the key.Map interface.
// func (k modelKeyMap) ShortHelp() []key.Binding {
// 	return []key.Binding{k.Accept, k.Quit, k.ExitNoSave}
// }

// // FullHelp returns keybindings for the expanded help view. It's part of the
// // key.Map interface.
// func (k modelKeyMap) FullHelp() [][]key.Binding {
// 	return [][]key.Binding{
// 		{k.Accept, k.Quit, k.ExitNoSave}, // second column
// 	}
// }

// var keys = modelKeyMap{
// 	Quit: key.NewBinding(
// 		key.WithKeys("esc"),
// 		key.WithHelp("esc", "quit"),
// 	),
// 	Accept: key.NewBinding(
// 		key.WithKeys("enter"),
// 		key.WithHelp("enter", "accept"),
// 	),
// 	ExitNoSave: key.NewBinding(
// 		key.WithKeys("ctrl+c"),
// 		key.WithHelp("ctrl+c", "exit without saving"),
// 	),
// 	ResetSuggestion: key.NewBinding(
// 		key.WithKeys("ctrl+r"),
// 		key.WithHelp("ctrl+r", "reset"),
// 	),
// 	CancelEdit: key.NewBinding(
// 		key.WithKeys("ctrl+e"),
// 		key.WithHelp("ctrl+e", "cancel edit"),
// 	),
// 	AcceptEdit: key.NewBinding(
// 		key.WithKeys("ctrl+s"),
// 		key.WithHelp("ctrl+s", "accpet"),
// 	),
// }
