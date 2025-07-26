package ui

import (
	"errors"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	potentiallyfixableissue "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/potentially-fixable-issue"
)

// type stage int

// const (
// 	sectionBreak stage = iota
// 	suggestionsProcessing
// 	stageCssSelection
// 	finalStage
// )

type FixableIssuesModel struct {
	// general
	CurrentStage          int
	BodyHeight, BodyWidth int
	Ready                 bool
	RunAll                bool
	// file data
	FilePaths []string
	FileTexts []string
	// body data
	ContextBreak string

	title         string
	width, height int
	stages        []string
	ready         bool
	// components
	body              viewport.Model
	help              help.Model
	sectionBreakInput textinput.Model
	suggestionDisplay viewport.Model

	// TODO: sort out where this data goes
	currentFile, currentIssueName                                        string
	isEditing                                                            bool
	suggestionData                                                       [][]suggestionState
	currentSuggestion                                                    *suggestionState
	currentSuggestionIndex, currentFileIndex, potentialFixableIssueIndex int

	// currentSuggestionIndex, potentialFixableIssueIndex, currentFileIndex int
	potentialIssues []potentiallyfixableissue.PotentiallyFixableIssue
	currentIssue    *potentiallyfixableissue.PotentiallyFixableIssue
	// cssUpdateRequired, addCssSectionBreakIfMissing, addCssPageBreakIfMissing, isEditing bool
	// currentSuggestionState *suggestionState
	// suggestionEdit         textarea.Model
	// scrollbar              tea.Model

	logFile io.Writer
	Err     error
}

func NewFixableIssuesModel(runAll, runSectionBreak bool, potentiallyFixableIssues []potentiallyfixableissue.PotentiallyFixableIssue, cssFiles []string, logFile io.Writer) FixableIssuesModel {
	ti := textinput.New()
	ti.Width = 20
	ti.CharLimit = 200
	ti.Placeholder = "Section break"

	var currentStage = 0
	if runAll || runSectionBreak {
		ti.Focus()
	} else {
		currentStage = 1
	}

	suggestionDisplay := viewport.New(0, 0)
	suggestionDisplay.MouseWheelEnabled = true

	return FixableIssuesModel{
		title: "Epub Linter Manually Fixable Issues",
		stages: []string{
			"Section Break",
			"Suggestions",
			"CSS File Selection",
		},
		CurrentStage: currentStage,

		help:              help.New(),
		body:              viewport.New(0, 0),
		suggestionDisplay: suggestionDisplay,
		sectionBreakInput: ti,

		logFile: logFile,
	}
}

func (m FixableIssuesModel) Init() tea.Cmd {
	return nil
}

func (m FixableIssuesModel) View() string {
	if m.ready {
		var (
			header       = m.headerView()
			footer       = m.footerView()
			headerHeight = lipgloss.Height(header) + headerBorderStyle.GetBorderBottomSize()
			footerHeight = lipgloss.Height(footer) + footerBorderStyle.GetBorderTopSize()
		)

		m.body.Width = m.width
		m.body.Height = max(0, m.height-(headerHeight+footerHeight)+2)
		m.BodyHeight = m.body.Height
		m.BodyWidth = m.body.Width
		m.body.SetContent(m.bodyView())

		return lipgloss.JoinVertical(lipgloss.Center, header, m.body.View(), footer)
	}

	return ""
}

func (m FixableIssuesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.Err = errors.New("must die....")
			return m, tea.Quit
		case "ctrl+c":
			m.Err = errors.New("must die....")
			// TODO: make sure this is an error once ready
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch m.CurrentStage {
	case 0:
		cmd = m.HandleSectionBreakKeys(msg)
	}

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

func (m FixableIssuesModel) headerView() string {
	return headerBorderStyle.Render(fillLine(lipgloss.JoinHorizontal(lipgloss.Center, titleStyle.Render(m.title), " | ", m.getStageHeaders()), m.width-headerBorderStyle.GetHorizontalBorderSize()))
}

func (m FixableIssuesModel) getStageHeaders() string {
	var stageHeaders = make([]string, len(m.stages))

	var style lipgloss.Style
	for i, header := range m.stages {
		if i == m.CurrentStage {
			style = activeStyle
		} else {
			style = inactiveStyle
		}

		stageHeaders[i] = style.Render(header)
	}

	return strings.Join(stageHeaders, inactiveStyle.Render(" > "))
}

func (m FixableIssuesModel) footerView() string {
	// return footerBorderStyle.Render(fillLine(controlsStyle.Render("Controls:"), m.width-footerBorderStyle.GetHorizontalBorderSize()) + "\n" + m.help.View(m.bodyContent[m.State.CurrentStage].HelpKeys()))
	var s strings.Builder
	s.WriteString(fillLine(controlsStyle.Render("Controls:"), m.width-footerBorderStyle.GetHorizontalBorderSize()) + "\n")

	var controls []string
	switch m.CurrentStage {
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

	// if line.Len() != 0 {
	// 	s.WriteString("\n")
	// }

	return footerBorderStyle.Render(s.String())
}

func (m FixableIssuesModel) bodyView() string {
	switch m.CurrentStage {
	case 0:
		return m.SectionBreakView()
	case 1:
		return m.SuggestionsView()
	}

	return ""
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
