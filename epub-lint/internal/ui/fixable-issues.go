package ui

import (
	"errors"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	potentiallyfixableissue "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/potentially-fixable-issue"
)

var ErrUserKilledProgram = errors.New("user killed program")

const (
	columnPadding    = 3
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
	ti.Width = 20
	ti.CharLimit = 200
	ti.Placeholder = "Section break"

	ta := textarea.New()
	ta.Placeholder = "Enter an edited version of the original string"
	ta.CharLimit = 10000
	ta.ShowLineNumbers = false

	var currentStage = sectionBreak
	if !skipCss && (runAll || runSectionBreak) {
		ti.Focus()
	} else {
		currentStage = suggestionsProcessing
	}

	v := viewport.New(80, 20)
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
	if m.Err != nil {
		return m, tea.Quit
	}

	var (
		cmds         []tea.Cmd
		cmd          tea.Cmd
		initialStage = m.currentStage
	)

	// general logic for handling keys here
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.Err = ErrUserKilledProgram

			return m, tea.Quit
		case "esc":
			return m, m.exitOrMoveToCssSelection()
		}
	case tea.WindowSizeMsg:
		m.ready = true
		m.height = msg.Height
		m.width = msg.Width

		m.body.Width = m.width
		m.body.Height = max(0, m.height-(m.headerHeight()+m.footerHeight()))

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
		cmd = m.handleCssSelectionMsgs(msg)
		cmds = append(cmds, cmd)
	}

	if m.currentStage == finalStage {
		return m, tea.Quit
	} else if m.currentStage != initialStage {
		cmds = append(cmds, tea.ClearScreen)
	}

	return m, tea.Batch(cmds...)
}

func (m FixableIssuesModel) View() string {
	if m.ready {
		var (
			header = m.headerView()
			footer = m.footerView()
		)
		m.body.SetContent(m.bodyView())

		return lipgloss.JoinVertical(lipgloss.Center, header, m.body.View(), footer)
	}

	return ""
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
				"Ctrl+E: Cancel edit",
				"Ctrl+S: Accept",
				"Esc: Quit",
				"Ctrl+C: Exit without saving",
			}
		} else if m.PotentiallyFixableIssuesInfo.currentSuggestionState != nil && m.PotentiallyFixableIssuesInfo.currentSuggestionState.isAccepted {
			controls = []string{
				"← / → : Previous/Next Suggestion",
				"C: Copy",
				"Esc: Quit",
				"Ctrl+C: Exit without saving",
			}
		} else {
			controls = []string{
				"← / → : Previous/Next Suggestion",
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
