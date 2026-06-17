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
	potentiallyfixableissue "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/potentially-fixable-issue"
)

var (
	ErrUserKilledProgram    = errors.New("user killed program")
	potentialFailedPasteMsg = "exit status 1"
)

const (
	borderWidth             = 2
	scrollbarPadding        = 3
	minHeight               = 21
	minWidth                = 41
	minLargeLayoutThreshold = 73
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
	pasteFailed  bool
	isPasting    bool // keep track of this so we know when to ignore an error that happens
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
	isAccepted, originallyHadHalfwidthCircleKatakana         bool
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

		m.recalculateElementSizes(false)

		cmds = append(cmds, tea.ClearScreen)
	case error:
		if m.logFile != nil {
			fmt.Fprintf(m.logFile, "Unexpected error encountered %q. Stage is %d. Is pasting: %t. Is potential failed paste error: %t.\n", msg, m.currentStage, m.sectionBreakInfo.isPasting, potentialFailedPasteMsg == msg.Error())
		}
		if m.sectionBreakInfo.isPasting && potentialFailedPasteMsg == msg.Error() {
			m.sectionBreakInfo.pasteFailed = true
		} else {
			m.Err = msg

			return m, tea.Quit
		}
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
		if m.height < minHeight || m.width < minWidth {
			view.SetContent(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, fmt.Sprintf("Terminal size is too small:\nWidth = %d Height = %d\n\nMinimum Needed:\nWidth = %d Height = %d", m.width, m.height, minWidth, minHeight)))

			return view
		}

		var (
			header = m.headerView()
			footer = m.footerView()
		)
		m.body.SetContent(m.bodyView())

		view.SetContent(lipgloss.JoinVertical(lipgloss.Left, header, m.body.View(), footer))
	}

	return view
}

func (m *FixableIssuesModel) recalculateElementSizes(resetSuggestionYOffset bool) {
	m.body.SetWidth(m.width)
	m.body.SetHeight(max(0, m.height-(m.headerHeight()+m.footerHeight())))

	m.sectionBreakInfo.input.SetWidth(m.width)
	m.setSuggestionDisplay(resetSuggestionYOffset)
}

func (m FixableIssuesModel) headerView() string {
	return headerBorderStyle.Render(fillLine(m.headerText(), m.width-headerBorderStyle.GetHorizontalBorderSize()))
}

func (m FixableIssuesModel) headerText() string {
	return lipgloss.JoinHorizontal(lipgloss.Center, titleStyle.Render(m.title), " | ", m.getStageHeaders())
}

func (m FixableIssuesModel) getStageHeaders() string {
	if m.width < minLargeLayoutThreshold {
		return activeStyle.Render(m.stages[min(int(m.currentStage), len(m.stages)-1)])
	}

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
