package ui

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/davecgh/go-spew/spew"
	potentiallyfixableissue "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/potentially-fixable-issue"
)

var ErrUserKilledProgram = errors.New("user killed program")

const (
	borderWidth      = 2
	scrollbarPadding = 3
)

type stage int

const (
	sectionBreak stage = iota
	suggestionsProcessing
	stageCssSelection
	finalStage
)

type FixableIssuesModel struct {
	sectionBreakInfo             sectionBreakStageInfo
	PotentiallyFixableIssuesInfo PotentiallyFixableStageInfo
	CssSelectionInfo             CssSelectionStageInfo
	currentStage                 stage
	body                         viewport.Model
	title                        string
	stages                       []string
	runAll, skipCss, ready       bool
	height, width                int
	logFile                      io.Writer
	Err                          error
}

type sectionBreakStageInfo struct {
	input        textinput.Model
	contextBreak *string
}

type PotentiallyFixableStageInfo struct {
	FileSuggestionData                                                                  []FileSuggestionInfo
	currentFile, currentSuggestionName                                                  string
	currentSuggestionIndex, potentialFixableIssueIndex, currentFileIndex                int
	suggestions                                                                         []potentiallyfixableissue.PotentiallyFixableIssue
	currentSuggestion                                                                   *potentiallyfixableissue.PotentiallyFixableIssue
	CssUpdateRequired, AddCssSectionBreakIfMissing, AddCssPageBreakIfMissing, isEditing bool
	currentSuggestionState                                                              *SuggestionState
	suggestionEdit                                                                      textarea.Model
	suggestionDisplay                                                                   viewport.Model
	scrollbar                                                                           tea.Model
}

type FileSuggestionInfo struct {
	Name        string
	Text        string
	Suggestions [][]SuggestionState
}

type CssSelectionStageInfo struct {
	cssFiles        []string
	SelectedCssFile string
	currentCssIndex int
}

type SuggestionState struct {
	isAccepted                                               bool
	original, originalSuggestion, currentSuggestion, display string
}

func NewFixableIssuesModel(runAll, skipCss, runSectionBreak bool, potentiallyFixableIssues []potentiallyfixableissue.PotentiallyFixableIssue, cssFiles []string, logFile io.Writer, contextBreak *string) FixableIssuesModel {
	ti := textinput.New()
	ti.SetWidth(20)
	ti.CharLimit = 200
	ti.Placeholder = "Section break"

	ta := textarea.New()
	ta.Prompt = ""
	ta.Placeholder = "Enter an edited version of the original string"
	ta.CharLimit = 10000
	ta.ShowLineNumbers = false

	var currentStage = sectionBreak
	if !skipCss && (runAll || runSectionBreak) {
		ti.Focus()
	} else {
		currentStage = suggestionsProcessing
	}

	v := viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
	v.MouseWheelEnabled = true

	sb := NewVertical()
	sb.Style = sb.Style.Border(lipgloss.RoundedBorder(), true)

	return FixableIssuesModel{
		sectionBreakInfo: sectionBreakStageInfo{
			input:        ti,
			contextBreak: contextBreak,
		},
		PotentiallyFixableIssuesInfo: PotentiallyFixableStageInfo{
			suggestions:       potentiallyFixableIssues,
			suggestionEdit:    ta,
			suggestionDisplay: v,
			scrollbar:         sb,
		},
		CssSelectionInfo: CssSelectionStageInfo{
			cssFiles: cssFiles,
		},
		runAll:       runAll,
		skipCss:      skipCss,
		currentStage: currentStage,
		logFile:      logFile,
		stages: []string{
			"Section Break",
			"Suggestions",
			"Select CSS File",
		},
		title: "Manually Fixable Issues",
	}
}

func (m FixableIssuesModel) Init() tea.Cmd {
	return nil
}

func (m FixableIssuesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.logFile != nil {
		spew.Fdump(m.logFile, msg)
	}
	if m.Err != nil {
		return m, tea.Quit
	}

	var (
		cmds         []tea.Cmd
		cmd          tea.Cmd
		initialStage = m.currentStage
	)

	if !m.ready && m.currentStage == suggestionsProcessing {
		cmd, m.Err = m.setupForNextSuggestions()

		cmds = append(cmds, cmd)
	}

	// general logic for handling keys here
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.Err = ErrUserKilledProgram

			return m, tea.Quit
		case "esc":
			cmd = m.exitOrMoveToCssSelection()

			return m, cmd
		}
	case tea.WindowSizeMsg:
		m.ready = true
		m.height = msg.Height
		m.width = msg.Width

		m.body.SetWidth(m.width)
		m.body.SetHeight(max(0, m.height-(m.headerHeight()+m.footerHeight())))

		m.setSuggestionDisplay(false)

		cmds = append(cmds, tea.ClearScreen)
	case error:
		m.Err = msg
		return m, tea.Quit
	}

	switch m.currentStage {
	case sectionBreak:
		cmd = m.handleSectionBreakMsgs(msg)
		cmds = append(cmds, cmd)
	case suggestionsProcessing:
		cmd = m.handleSuggestionMsgs(msg)
		cmds = append(cmds, cmd)
	case stageCssSelection:
		m.handleCssSelectionMsgs(msg)
	}

	if m.currentStage == finalStage {
		return m, tea.Quit
	} else if m.currentStage != initialStage {
		cmds = append(cmds, tea.ClearScreen)
	}

	return m, tea.Batch(cmds...)
}

func (m FixableIssuesModel) View() tea.View {
	view := tea.NewView("")
	view.AltScreen = true
	if m.ready {
		var (
			header = m.headerView()
			footer = m.footerView()
		)
		m.body.SetContent(m.bodyView())

		var content = lipgloss.JoinVertical(lipgloss.Center, header, m.body.View(), footer)

		if m.logFile != nil {
			fmt.Fprintf(m.logFile, "Suggestion content (%d, %d): %q\n", m.body.Width(), m.body.Height(), content)
		}

		view.SetContent(content)
	}

	return view
}

func (m FixableIssuesModel) headerView() string {
	return headerBorderStyle.Render(fillLine(lipgloss.JoinHorizontal(lipgloss.Center, titleStyle.Render(m.title), " | ", m.getStageHeaders()), m.width-headerBorderStyle.GetHorizontalBorderSize()))
}

func (m FixableIssuesModel) getStageHeaders() string {
	var stageHeaders = make([]string, len(m.stages))

	var style lipgloss.Style
	for i, header := range m.stages {
		if i == int(m.currentStage) {
			style = activeStyle
		} else {
			style = inactiveStyle
		}

		stageHeaders[i] = style.Render(header)
	}

	return strings.Join(stageHeaders, inactiveStyle.Render(" > "))
}

func (m FixableIssuesModel) footerView() string {
	var (
		s        strings.Builder
		maxWidth = m.width - footerBorderStyle.GetHorizontalBorderSize()
	)
	s.WriteString(fillLine(controlsStyle.Render("Controls:"), maxWidth) + "\n")

	var controls []string
	switch m.currentStage {
	case sectionBreak:
		controls = []string{
			"Enter: Accept",
			"Esc: Quit",
			"Ctrl+C: Exit without saving",
		}
	case suggestionsProcessing:
		if m.PotentiallyFixableIssuesInfo.isEditing {
			controls = []string{
				"Ctrl+R: Reset",
				"Ctrl+O: Original content",
				"Ctrl+E: Cancel edit",
				"Ctrl+S: Accept",
				"Esc: Quit",
				"Ctrl+C: Exit without saving",
			}
		} else if m.PotentiallyFixableIssuesInfo.currentSuggestionState != nil && m.PotentiallyFixableIssuesInfo.currentSuggestionState.isAccepted {
			controls = []string{
				"← / → : Previous/Next Suggestion",
				"Ctrl+U / Ctrl+D: Previous/Next Potential Issue Type",
				"Ctrl+PgUp / Ctrl+PgDn: Previous/Next File",
				"C: Copy",
				"Esc: Quit",
				"Ctrl+C: Exit without saving",
			}
		} else {
			controls = []string{
				"← / → : Previous/Next Suggestion",
				"Ctrl+U / Ctrl+D: Previous/Next Potential Issue Type",
				"Ctrl+PgUp / Ctrl+PgDn: Previous/Next File",
				"E: Edit",
				"C: Copy",
				"Enter: Accept",
				"Esc: Quit",
				"Ctrl+C: Exit without saving",
			}
		}
	case stageCssSelection:
		controls = []string{
			"↑ / ↓ : Previous/Next Suggestion",
			"Enter: Accept",
			"Ctrl+C: Exit without saving",
		}
	}

	var (
		line strings.Builder
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

	return footerBorderStyle.Render(s.String())
}

func (m FixableIssuesModel) bodyView() string {
	switch m.currentStage {
	case sectionBreak:
		return m.sectionBreakView()
	case suggestionsProcessing:
		return m.suggestionsView()
	case stageCssSelection:
		return m.cssSelectionView()
	}

	return ""
}

func (m *FixableIssuesModel) exitOrMoveToCssSelection() tea.Cmd {
	if m.PotentiallyFixableIssuesInfo.CssUpdateRequired {
		m.currentStage = stageCssSelection

		if len(m.CssSelectionInfo.cssFiles) != 0 {
			m.CssSelectionInfo.SelectedCssFile = m.CssSelectionInfo.cssFiles[m.CssSelectionInfo.currentCssIndex]
		}
	} else {
		m.currentStage = finalStage
		return tea.Quit
	}

	return nil
}

func (m FixableIssuesModel) headerHeight() int {
	return lipgloss.Height(m.headerView()) + headerBorderStyle.GetBorderBottomSize()
}

func (m FixableIssuesModel) footerHeight() int {
	return lipgloss.Height(m.footerView()) + footerBorderStyle.GetBorderTopSize()
}
