package ui

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pjkaufman/go-go-gadgets/epub-lint/cmd/tui"
	potentiallyfixableissue "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/potentially-fixable-issue"
	stringdiff "github.com/pjkaufman/go-go-gadgets/pkg/string-diff"
)

var (
	maxDisplayHeight = 20
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
	runAll, ready                bool
	height, width                int
	logFile                      io.Writer
	Err                          error
}

type sectionBreakStageInfo struct {
	input        textinput.Model
	contextBreak *string
}

type PotentiallyFixableStageInfo struct {
	FilePaths                                                                           []string
	FileTexts                                                                           []string
	currentFile, currentSuggestionName                                                  string
	currentSuggestionIndex, potentialFixableIssueIndex, currentFileIndex                int
	suggestions                                                                         []potentiallyfixableissue.PotentiallyFixableIssue
	currentSuggestion                                                                   *potentiallyfixableissue.PotentiallyFixableIssue
	CssUpdateRequired, AddCssSectionBreakIfMissing, AddCssPageBreakIfMissing, isEditing bool
	sectionSuggestionStates                                                             []suggestionState
	currentSuggestionState                                                              *suggestionState
	suggestionEdit                                                                      textarea.Model
	suggestionDisplay                                                                   viewport.Model
	scrollbar                                                                           tea.Model

	// leftStatusWidth int
}

type CssSelectionStageInfo struct {
	cssFiles        []string
	SelectedCssFile string
	currentCssIndex int
}

type suggestionState struct {
	isAccepted                                               bool
	original, originalSuggestion, currentSuggestion, display string
}

func NewFixableIssuesModel(runAll, runSectionBreak bool, potentiallyFixableIssues []potentiallyfixableissue.PotentiallyFixableIssue, cssFiles []string, logFile io.Writer, contextBreak *string) FixableIssuesModel {
	ti := textinput.New()
	ti.Width = 20
	ti.CharLimit = 200
	ti.Placeholder = "Section break"

	ta := textarea.New()
	ta.Placeholder = "Enter an edited version of the original string"
	ta.CharLimit = 10000
	ta.ShowLineNumbers = false

	var currentStage = sectionBreak
	if runAll || runSectionBreak {
		ti.Focus()
	} else {
		currentStage = suggestionsProcessing
	}

	v := viewport.New(80, 20)
	v.MouseWheelEnabled = true

	sb := tui.NewVertical()
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
		currentStage: currentStage,
		logFile:      logFile,
		stages: []string{
			"Section Break",
			"Suggestions",
			"CSS File Selection",
		},
		title: "Epub Linter Manually Fixable Issues",
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

	m.sectionBreakInfo.input, cmd = m.sectionBreakInfo.input.Update(msg)
	cmds = append(cmds, cmd)

	m.PotentiallyFixableIssuesInfo.suggestionEdit, cmd = m.PotentiallyFixableIssuesInfo.suggestionEdit.Update(msg)
	cmds = append(cmds, cmd)

	m.PotentiallyFixableIssuesInfo.suggestionDisplay, cmd = m.PotentiallyFixableIssuesInfo.suggestionDisplay.Update(msg)
	cmds = append(cmds, cmd)

	m.PotentiallyFixableIssuesInfo.scrollbar, cmd = m.PotentiallyFixableIssuesInfo.scrollbar.Update(m.PotentiallyFixableIssuesInfo.suggestionDisplay)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch m.currentStage {
		case sectionBreak:
			switch key {
			case "enter":
				*m.sectionBreakInfo.contextBreak = strings.TrimSpace(m.sectionBreakInfo.input.Value())
				if *m.sectionBreakInfo.contextBreak != "" {
					m.currentStage = suggestionsProcessing

					cmd, err := m.setupForNextSuggestions()
					if err != nil {
						m.Err = err

						return m, tea.Quit
					}

					cmds = append(cmds, cmd)
				}
			case "ctrl+c":
				m.Err = ErrUserKilledProgram

				return m, tea.Quit
			case "esc":
				return m, tea.Quit
			}
		case suggestionsProcessing:
			if m.PotentiallyFixableIssuesInfo.isEditing {
				switch key {
				case "ctrl+c":
					m.Err = ErrUserKilledProgram

					return m, tea.Quit
				case "esc":
					return m, m.exitOrMoveToCssSelection()
				case "ctrl+s":
					m.PotentiallyFixableIssuesInfo.currentSuggestionState.currentSuggestion = alignWhitespace(m.PotentiallyFixableIssuesInfo.currentSuggestionState.original, m.PotentiallyFixableIssuesInfo.suggestionEdit.Value())
					m.PotentiallyFixableIssuesInfo.isEditing = false

					var err error
					m.PotentiallyFixableIssuesInfo.currentSuggestionState.display, err = getStringDiff(m.PotentiallyFixableIssuesInfo.currentSuggestionState.original, m.PotentiallyFixableIssuesInfo.currentSuggestionState.currentSuggestion)
					if err != nil {
						m.Err = err

						return m, tea.Quit
					}

					return m, m.setSuggestionDisplay()
				case "ctrl+e":
					m.PotentiallyFixableIssuesInfo.isEditing = false

					var err error
					m.PotentiallyFixableIssuesInfo.currentSuggestionState.display, err = getStringDiff(m.PotentiallyFixableIssuesInfo.currentSuggestionState.original, m.PotentiallyFixableIssuesInfo.currentSuggestionState.currentSuggestion)
					if err != nil {
						m.Err = err

						return m, tea.Quit
					}
				case "ctrl+r":
					m.PotentiallyFixableIssuesInfo.suggestionEdit.SetValue(m.PotentiallyFixableIssuesInfo.currentSuggestionState.originalSuggestion)
				}
			} else {
				switch key {
				case "ctrl+c":
					m.Err = ErrUserKilledProgram

					return m, tea.Quit
				case "esc":
					return m, m.exitOrMoveToCssSelection()
				case "right":
					cmd, err := m.moveToNextSuggestion()
					if err != nil {
						m.Err = err

						return m, tea.Quit
					}

					cmds = append(cmds, cmd)
				case "left":
					cmd = m.moveToPreviousSuggestion()
					cmds = append(cmds, cmd)
				case "c":
					// Copy original value to the clipboard
					// original, err := repairUnicode(m.currentFileState.original)
					// if err != nil {
					// 	m.Err = err

					// 	return m, tea.Quit
					// }

					// err = clipboard.WriteAll(original)
					// TODO: make sure values are utf-8 compliant
					err := clipboard.WriteAll(m.PotentiallyFixableIssuesInfo.currentSuggestionState.original)
					if err != nil {
						m.Err = err

						return m, tea.Quit
					}
				case "enter":
					if !m.PotentiallyFixableIssuesInfo.currentSuggestionState.isAccepted && m.PotentiallyFixableIssuesInfo.currentSuggestion != nil {
						var replaceCount = 1
						if m.PotentiallyFixableIssuesInfo.currentSuggestion.UpdateAllInstances {
							replaceCount = -1
						}

						m.PotentiallyFixableIssuesInfo.FileTexts[m.PotentiallyFixableIssuesInfo.currentFileIndex] = strings.Replace(m.PotentiallyFixableIssuesInfo.FileTexts[m.PotentiallyFixableIssuesInfo.currentFileIndex], m.PotentiallyFixableIssuesInfo.currentSuggestionState.original, m.PotentiallyFixableIssuesInfo.currentSuggestionState.currentSuggestion, replaceCount)

						m.PotentiallyFixableIssuesInfo.currentSuggestionState.isAccepted = true

						if m.PotentiallyFixableIssuesInfo.currentSuggestion.AddCssSectionBreakIfMissing {
							m.PotentiallyFixableIssuesInfo.AddCssSectionBreakIfMissing = true
							m.PotentiallyFixableIssuesInfo.CssUpdateRequired = true
						} else if m.PotentiallyFixableIssuesInfo.currentSuggestion.AddCssPageBreakIfMissing {
							m.PotentiallyFixableIssuesInfo.AddCssPageBreakIfMissing = true
							m.PotentiallyFixableIssuesInfo.CssUpdateRequired = true
						}

						cmd, err := m.moveToNextSuggestion()
						if err != nil {
							m.Err = err

							return m, tea.Quit
						}

						cmds = append(cmds, cmd)
					}
				case "e":
					if m.PotentiallyFixableIssuesInfo.currentSuggestionState != nil && !m.PotentiallyFixableIssuesInfo.currentSuggestionState.isAccepted {
						m.PotentiallyFixableIssuesInfo.isEditing = true
						m.PotentiallyFixableIssuesInfo.suggestionEdit.SetValue(m.PotentiallyFixableIssuesInfo.currentSuggestionState.currentSuggestion)

						cmd = m.PotentiallyFixableIssuesInfo.suggestionEdit.Focus()
						cmds = append(cmds, cmd)
					}
				}
			}
		case stageCssSelection:
			switch key {
			case "ctrl+c":
				m.Err = ErrUserKilledProgram

				return m, tea.Quit
			case "up":
				if m.CssSelectionInfo.currentCssIndex > 0 {
					m.CssSelectionInfo.currentCssIndex--
					m.CssSelectionInfo.SelectedCssFile = m.CssSelectionInfo.cssFiles[m.CssSelectionInfo.currentCssIndex]
				}
			case "down":
				if m.CssSelectionInfo.currentCssIndex+1 < len(m.CssSelectionInfo.cssFiles) {
					m.CssSelectionInfo.currentCssIndex++
					m.CssSelectionInfo.SelectedCssFile = m.CssSelectionInfo.cssFiles[m.CssSelectionInfo.currentCssIndex]
				}
			case "enter":
				m.currentStage = finalStage
			}
		default:
			switch key {
			case "ctrl+c":
				m.Err = ErrUserKilledProgram

				return m, tea.Quit
			case "esc":
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.ready = true
		m.height = msg.Height
		// TODO: see about removing this value...
		maxDisplayHeight = msg.Height / 3
		m.width = msg.Width

		// TODO: replace with the logic for finding the body width and height
		var maxWidth = m.width - columnPadding

		m.PotentiallyFixableIssuesInfo.suggestionEdit.SetWidth(maxWidth)
		m.PotentiallyFixableIssuesInfo.suggestionDisplay.Width = maxWidth - scrollbarPadding
		if m.PotentiallyFixableIssuesInfo.suggestionDisplay.Height > maxDisplayHeight {
			m.PotentiallyFixableIssuesInfo.suggestionDisplay.Height = maxDisplayHeight
		}

		m.PotentiallyFixableIssuesInfo.scrollbar, cmd = m.PotentiallyFixableIssuesInfo.scrollbar.Update(tui.HeightMsg(m.PotentiallyFixableIssuesInfo.suggestionDisplay.Height))
		cmds = append(cmds, cmd)

		cmds = append(cmds, tea.ClearScreen)
	case error:
		m.Err = msg
		return m, tea.Quit
	}

	if m.currentStage == finalStage {
		return m, tea.Quit
	} else if m.currentStage != initialStage {
		cmds = append(cmds, tea.ClearScreen)
	}

	return m, tea.Batch(cmds...)
}

// func (m FixableIssuesModel) View() string {
// 	var (
// 		header = m.headerView()
// 		s      strings.Builder
// 	)

// 	switch m.currentStage {
// 	case sectionBreak:
// 		s.WriteString("\n" + m.sectionBreakInfo.input.View() + "\n\n")
// 	case suggestionsProcessing:
// 		// s.WriteString(titleStyle.Render(fmt.Sprintf("Current File (%d/%d): %s ", m.PotentiallyFixableIssuesInfo.currentFileIndex+1, len(m.PotentiallyFixableIssuesInfo.FilePaths), m.PotentiallyFixableIssuesInfo.currentFile)) + "\n")
// 		// s.WriteString(groupStyle.Render(fmt.Sprintf("Issue Group: %s", m.PotentiallyFixableIssuesInfo.currentSuggestionName) + "\n"))

// 		if m.PotentiallyFixableIssuesInfo.currentSuggestionState == nil {
// 			s.WriteString("\nNo current file is selected. Something may have gone wrong...\n\n")
// 		} else {
// 			var (
// 				suggestion     = displayStyle.Width(m.width - columnPadding).Render(fmt.Sprintf(`"%s"`, m.PotentiallyFixableIssuesInfo.currentSuggestionState.display))
// 				expectedHeight = strings.Count(suggestion, "\n") + 1
// 			)
// 			if m.PotentiallyFixableIssuesInfo.isEditing {
// 				s.WriteString("\nEditing suggestion:\n\n")
// 			} else if m.PotentiallyFixableIssuesInfo.currentSuggestionState.isAccepted {
// 				s.WriteString("\n" + acceptedChangeTitleStyle.Render("Accepted change:") + "\n")
// 			} else {
// 				s.WriteString(fmt.Sprintf("\nSuggested change (%d/%d):\n", expectedHeight, maxDisplayHeight))
// 			}

// 			if m.PotentiallyFixableIssuesInfo.isEditing {
// 				s.WriteString(suggestionBorderStyle.Render(m.PotentiallyFixableIssuesInfo.suggestionEdit.View()) + "\n\n")
// 			} else {
// 				if m.PotentiallyFixableIssuesInfo.suggestionDisplay.TotalLineCount() > m.PotentiallyFixableIssuesInfo.suggestionDisplay.VisibleLineCount() {
// 					s.WriteString(suggestionBorderStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top,
// 						m.PotentiallyFixableIssuesInfo.suggestionDisplay.View(),
// 						m.PotentiallyFixableIssuesInfo.scrollbar.View(),
// 					)))
// 				} else {
// 					s.WriteString(suggestionBorderStyle.Render(m.PotentiallyFixableIssuesInfo.suggestionDisplay.View()))
// 				}

// 				s.WriteString(fmt.Sprintf("\033[0m\n\nSuggestion %d of %d.\n\n", m.PotentiallyFixableIssuesInfo.currentSuggestionIndex+1, len(m.PotentiallyFixableIssuesInfo.sectionSuggestionStates)))
// 			}
// 		}
// 	case stageCssSelection:
// 		s.WriteString("\nSelect the CSS file to modify:\n\n")
// 		for i, cssFile := range m.cssSelectionInfo.cssFiles {
// 			cursor := " "
// 			if m.cssSelectionInfo.currentCssIndex == i {
// 				cursor = ">"
// 			}

// 			s.WriteString(fmt.Sprintf("%s %d. %s\n", cursor, i+1, cssFile))
// 		}

// 		s.WriteString("\n")
// 	}

// 	m.displayControls(&s)

// 	return header + s.String()
// }

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
	var s strings.Builder
	s.WriteString(fillLine(controlsStyle.Render("Controls:"), m.width-footerBorderStyle.GetHorizontalBorderSize()) + "\n")

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
		line     strings.Builder
		maxWidth = m.width - columnPadding
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
	var s strings.Builder

	switch m.currentStage {
	case sectionBreak:
		return m.sectionBreakView()
	case suggestionsProcessing:
		return m.suggestionsView()
		// // s.WriteString(titleStyle.Render(fmt.Sprintf("Current File (%d/%d): %s ", m.PotentiallyFixableIssuesInfo.currentFileIndex+1, len(m.PotentiallyFixableIssuesInfo.FilePaths), m.PotentiallyFixableIssuesInfo.currentFile)) + "\n")
		// // s.WriteString(groupStyle.Render(fmt.Sprintf("Issue Group: %s", m.PotentiallyFixableIssuesInfo.currentSuggestionName) + "\n"))

		// if m.PotentiallyFixableIssuesInfo.currentSuggestionState == nil {
		// 	s.WriteString("\nNo current file is selected. Something may have gone wrong...\n\n")
		// } else {
		// 	var (
		// 		suggestion     = displayStyle.Width(m.width - columnPadding).Render(fmt.Sprintf(`"%s"`, m.PotentiallyFixableIssuesInfo.currentSuggestionState.display))
		// 		expectedHeight = strings.Count(suggestion, "\n") + 1
		// 	)
		// 	if m.PotentiallyFixableIssuesInfo.isEditing {
		// 		s.WriteString("\nEditing suggestion:\n\n")
		// 	} else if m.PotentiallyFixableIssuesInfo.currentSuggestionState.isAccepted {
		// 		s.WriteString("\n" + acceptedChangeTitleStyle.Render("Accepted change:") + "\n")
		// 	} else {
		// 		s.WriteString(fmt.Sprintf("\nSuggested change (%d/%d):\n", expectedHeight, maxDisplayHeight))
		// 	}

		// 	if m.PotentiallyFixableIssuesInfo.isEditing {
		// 		s.WriteString(suggestionBorderStyle.Render(m.PotentiallyFixableIssuesInfo.suggestionEdit.View()) + "\n\n")
		// 	} else {
		// 		if m.PotentiallyFixableIssuesInfo.suggestionDisplay.TotalLineCount() > m.PotentiallyFixableIssuesInfo.suggestionDisplay.VisibleLineCount() {
		// 			s.WriteString(suggestionBorderStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top,
		// 				m.PotentiallyFixableIssuesInfo.suggestionDisplay.View(),
		// 				m.PotentiallyFixableIssuesInfo.scrollbar.View(),
		// 			)))
		// 		} else {
		// 			s.WriteString(suggestionBorderStyle.Render(m.PotentiallyFixableIssuesInfo.suggestionDisplay.View()))
		// 		}

		// 		s.WriteString(fmt.Sprintf("\033[0m\n\nSuggestion %d of %d.\n\n", m.PotentiallyFixableIssuesInfo.currentSuggestionIndex+1, len(m.PotentiallyFixableIssuesInfo.sectionSuggestionStates)))
		// 	}
		// }
	case stageCssSelection:
		s.WriteString("\nSelect the CSS file to modify:\n\n")
		for i, cssFile := range m.CssSelectionInfo.cssFiles {
			cursor := " "
			if m.CssSelectionInfo.currentCssIndex == i {
				cursor = ">"
			}

			s.WriteString(fmt.Sprintf("%s %d. %s\n", cursor, i+1, cssFile))
		}

		s.WriteString("\n")
	}

	return s.String()
}

func (m *FixableIssuesModel) setupForNextSuggestions() (tea.Cmd, error) {
	if m.logFile != nil {
		fmt.Fprintln(m.logFile, "Getting next suggestions")
	}

	for m.PotentiallyFixableIssuesInfo.currentFileIndex < len(m.PotentiallyFixableIssuesInfo.FilePaths) {
		m.PotentiallyFixableIssuesInfo.currentFile = m.PotentiallyFixableIssuesInfo.FilePaths[m.PotentiallyFixableIssuesInfo.currentFileIndex]
		if m.logFile != nil {
			fmt.Fprintf(m.logFile, "Current file is %q is %d of %d\n", m.PotentiallyFixableIssuesInfo.currentFile, m.PotentiallyFixableIssuesInfo.currentFileIndex+1, len(m.PotentiallyFixableIssuesInfo.FilePaths))
		}

		for m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex < len(m.PotentiallyFixableIssuesInfo.suggestions) {
			var potentialFixableIssue = m.PotentiallyFixableIssuesInfo.suggestions[m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex]
			if m.logFile != nil {
				fmt.Fprintf(m.logFile, "Possible fixable issue %q is %d of %d issues.\n", potentialFixableIssue.Name, m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex+1, len(m.PotentiallyFixableIssuesInfo.suggestions))
			}

			if !m.runAll && (potentialFixableIssue.IsEnabled == nil || *potentialFixableIssue.IsEnabled) {
				if m.logFile != nil {
					fmt.Fprintf(m.logFile, "Skipping possible fixable issue %q with isEnabled set to %v\n", potentialFixableIssue.Name, potentialFixableIssue.IsEnabled)
				}

				m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex++
				continue
			}

			var (
				suggestions = potentialFixableIssue.GetSuggestions(m.PotentiallyFixableIssuesInfo.FileTexts[m.PotentiallyFixableIssuesInfo.currentFileIndex])
			)

			if m.logFile != nil {
				fmt.Fprintf(m.logFile, "Possible fixable issue %q has %d suggestion(s) found\n", potentialFixableIssue.Name, len(suggestions))
			}

			if len(suggestions) != 0 {
				m.PotentiallyFixableIssuesInfo.currentSuggestion = &potentialFixableIssue
				m.PotentiallyFixableIssuesInfo.sectionSuggestionStates = make([]suggestionState, len(suggestions))

				var i = 0
				for original, suggestion := range suggestions {
					var display, err = getStringDiff(original, suggestion)
					if err != nil {
						return nil, err
					}

					m.PotentiallyFixableIssuesInfo.sectionSuggestionStates[i] = suggestionState{
						original:           original,
						originalSuggestion: suggestion,
						currentSuggestion:  suggestion,
						display:            display,
					}

					i++
				}

				m.PotentiallyFixableIssuesInfo.currentSuggestionIndex = 0
				m.PotentiallyFixableIssuesInfo.currentSuggestionState = &m.PotentiallyFixableIssuesInfo.sectionSuggestionStates[0]
				cmd := m.setSuggestionDisplay()
				m.PotentiallyFixableIssuesInfo.currentSuggestionName = potentialFixableIssue.Name

				m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex++

				return cmd, nil
			}

			m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex++
		}

		m.PotentiallyFixableIssuesInfo.currentFileIndex++
		m.PotentiallyFixableIssuesInfo.potentialFixableIssueIndex = 0
	}

	if m.PotentiallyFixableIssuesInfo.CssUpdateRequired {
		m.currentStage = stageCssSelection
	} else {
		m.currentStage = finalStage
	}

	return nil, nil
}

func (m *FixableIssuesModel) moveToNextSuggestion() (tea.Cmd, error) {
	if m.PotentiallyFixableIssuesInfo.currentSuggestionIndex+1 < len(m.PotentiallyFixableIssuesInfo.sectionSuggestionStates) {
		m.PotentiallyFixableIssuesInfo.currentSuggestionIndex++
		m.PotentiallyFixableIssuesInfo.currentSuggestionState = &m.PotentiallyFixableIssuesInfo.sectionSuggestionStates[m.PotentiallyFixableIssuesInfo.currentSuggestionIndex]

		return m.setSuggestionDisplay(), nil
	}

	return m.setupForNextSuggestions()
}

func (m *FixableIssuesModel) moveToPreviousSuggestion() tea.Cmd {
	if m.PotentiallyFixableIssuesInfo.currentSuggestionIndex > 0 {
		m.PotentiallyFixableIssuesInfo.currentSuggestionIndex--
		m.PotentiallyFixableIssuesInfo.currentSuggestionState = &m.PotentiallyFixableIssuesInfo.sectionSuggestionStates[m.PotentiallyFixableIssuesInfo.currentSuggestionIndex]
		return m.setSuggestionDisplay()
	}

	return nil
}

func (m *FixableIssuesModel) setSuggestionDisplay() tea.Cmd {
	if m.PotentiallyFixableIssuesInfo.currentSuggestionState == nil {
		return nil
	}

	if m.logFile != nil {
		fmt.Fprintf(m.logFile, "current width %d; border width: %d; scrollbar padding: %d\n", m.PotentiallyFixableIssuesInfo.suggestionDisplay.Width, suggestionBorderStyle.GetHorizontalBorderSize(), scrollbarPadding)
	}

	var (
		expectedSuggestionWidth = m.PotentiallyFixableIssuesInfo.suggestionDisplay.Width
		suggestion              = displayStyle.Render(wrapLines(fmt.Sprintf(`"%s"`, m.PotentiallyFixableIssuesInfo.currentSuggestionState.display), expectedSuggestionWidth))
	)

	if m.logFile != nil {
		fmt.Fprintf(m.logFile, "New suggestion getting set with width %d and a value of %q\n", expectedSuggestionWidth, suggestion)
	}

	m.PotentiallyFixableIssuesInfo.suggestionDisplay.SetContent(suggestion)

	// TODO: how do I handle the resizing of the suggestion display
	// one option would be to
	var cmd tea.Cmd
	m.PotentiallyFixableIssuesInfo.scrollbar, cmd = m.PotentiallyFixableIssuesInfo.scrollbar.Update(tui.HeightMsg(m.PotentiallyFixableIssuesInfo.suggestionDisplay.Height))

	return cmd
}

func (m *FixableIssuesModel) exitOrMoveToCssSelection() tea.Cmd {
	if m.PotentiallyFixableIssuesInfo.CssUpdateRequired {
		m.currentStage = stageCssSelection
	} else {
		return tea.Quit
	}

	return nil
}

func getStringDiff(original, new string) (string, error) {
	return stringdiff.GetPrettyDiffString(strings.TrimLeft(original, "\n"), strings.TrimLeft(new, "\n"))
}

// textarea gets rid of tabs when creating changes, so in order to preserve tabs in the starting whitespace of a line
// we will use the value of original as the template for what whitespace is needed for each line present
func alignWhitespace(original, new string) string {
	origLines := strings.Split(original, "\n")
	newLines := strings.Split(new, "\n")

	var min = len(newLines)
	if len(origLines) < min {
		min = len(origLines)
	}

	for i := 0; i < min; i++ {
		origPrefix := ""
		for j := 0; j < len(origLines[i]); j++ {
			if !unicode.IsSpace(rune(origLines[i][j])) {
				break
			}
			origPrefix += string(origLines[i][j])
		}
		newLines[i] = origPrefix + strings.TrimLeft(newLines[i], " \t")
	}

	return strings.Join(newLines, "\n")
}
