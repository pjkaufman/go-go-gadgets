package epub

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	stringdiff "github.com/pjkaufman/go-go-gadgets/pkg/string-diff"
)

var (
	appNameStyle             = lipgloss.NewStyle().Background(lipgloss.Color("99")).Padding(0, 1)
	titleStyle               = lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	groupStyle               = lipgloss.NewStyle().Foreground(lipgloss.Color("246"))
	fileStatusStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("190"))
	acceptedChangeTitleStyle = lipgloss.NewStyle().Bold(true)
	displayStyle             = lipgloss.NewStyle()
)

const columnPadding = 10

type stage int

const (
	sectionBreak stage = iota
	suggestionsProcessing
	stageCssSelection
	finalStage
)

type fixableIssuesModel struct {
	sectionBreakInfo             sectionBreakStageInfo
	potentiallyFixableIssuesInfo potentiallyFixableStageInfo
	cssSelectionInfo             cssSelectionStageInfo
	currentStage                 stage
	runAll                       bool
	height, width                int
	Err                          error
}

type sectionBreakStageInfo struct {
	input textinput.Model
}

type potentiallyFixableStageInfo struct {
	filePaths                                                                           []string
	fileTexts                                                                           map[string]string
	currentFile, currentSuggestionName                                                  string
	currentSuggestionIndex, potentialFixableIssueIndex, currentFileIndex                int
	suggestions                                                                         []potentiallyFixableIssue
	currentSuggestion                                                                   *potentiallyFixableIssue
	cssUpdateRequired, addCssSectionBreakIfMissing, addCssPageBreakIfMissing, isEditing bool
	sectionSuggestionStates                                                             []suggestionState
	currentSuggestionState                                                              *suggestionState
	suggestionEdit                                                                      textarea.Model
}

type cssSelectionStageInfo struct {
	cssFiles        []string
	selectedCssFile string
	currentCssIndex int
}

type suggestionState struct {
	isAccepted                                               bool
	original, originalSuggestion, currentSuggestion, display string
}

func newModel(runAll, runSectionBreak bool, potentiallyFixableIssues []potentiallyFixableIssue, cssFiles []string) fixableIssuesModel {
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

	return fixableIssuesModel{
		sectionBreakInfo: sectionBreakStageInfo{
			input: ti,
		},
		potentiallyFixableIssuesInfo: potentiallyFixableStageInfo{
			fileTexts:      map[string]string{},
			suggestions:    potentiallyFixableIssues,
			suggestionEdit: ta,
		},
		cssSelectionInfo: cssSelectionStageInfo{
			cssFiles: cssFiles,
		},
		runAll:       runAll,
		currentStage: currentStage,
	}
}

func (m fixableIssuesModel) Init() tea.Cmd {
	return nil
}

func (m fixableIssuesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	m.potentiallyFixableIssuesInfo.suggestionEdit, cmd = m.potentiallyFixableIssuesInfo.suggestionEdit.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch m.currentStage {
		case sectionBreak:
			switch key {
			case "enter":
				contextBreak = strings.TrimSpace(m.sectionBreakInfo.input.Value())
				if contextBreak != "" {
					m.currentStage = suggestionsProcessing

					err := m.setupForNextSuggestions()
					if err != nil {
						m.Err = err

						return m, tea.Quit
					}
				}
			case "ctrl+c", "esc":
				return m, tea.Quit
			}
		case suggestionsProcessing:
			if m.potentiallyFixableIssuesInfo.isEditing {
				switch key {
				case "ctrl+c", "esc":
					return m, tea.Quit
				case "ctrl+s":
					m.potentiallyFixableIssuesInfo.currentSuggestionState.currentSuggestion = alignWhitespace(m.potentiallyFixableIssuesInfo.currentSuggestionState.original, m.potentiallyFixableIssuesInfo.suggestionEdit.Value())
					m.potentiallyFixableIssuesInfo.isEditing = false

					var err error
					m.potentiallyFixableIssuesInfo.currentSuggestionState.display, err = getStringDiff(m.potentiallyFixableIssuesInfo.currentSuggestionState.original, m.potentiallyFixableIssuesInfo.currentSuggestionState.currentSuggestion)
					if err != nil {
						m.Err = err

						return m, tea.Quit
					}
				case "ctrl+e":
					m.potentiallyFixableIssuesInfo.isEditing = false

					var err error
					m.potentiallyFixableIssuesInfo.currentSuggestionState.display, err = getStringDiff(m.potentiallyFixableIssuesInfo.currentSuggestionState.original, m.potentiallyFixableIssuesInfo.currentSuggestionState.currentSuggestion)
					if err != nil {
						m.Err = err

						return m, tea.Quit
					}
				case "ctrl+r":
					m.potentiallyFixableIssuesInfo.suggestionEdit.SetValue(m.potentiallyFixableIssuesInfo.currentSuggestionState.originalSuggestion)
				}
			} else {
				switch key {
				case "ctrl+c", "esc":
					return m, tea.Quit
				case "right":
					err := m.moveToNextSuggestion()
					if err != nil {
						m.Err = err

						return m, tea.Quit
					}
				case "left":
					m.moveToPreviousSuggestion()
				case "c":
					// Copy original value to the clipboard
					// original, err := repairUnicode(m.currentFileState.original)
					// if err != nil {
					// 	m.Err = err

					// 	return m, tea.Quit
					// }

					// err = clipboard.WriteAll(original)
					// TODO: make sure values are utf-8 compliant
					err := clipboard.WriteAll(m.potentiallyFixableIssuesInfo.currentSuggestionState.original)
					if err != nil {
						m.Err = err

						return m, tea.Quit
					}
				case "enter":
					if !m.potentiallyFixableIssuesInfo.currentSuggestionState.isAccepted && m.potentiallyFixableIssuesInfo.currentSuggestion != nil {
						var replaceCount = 1
						if m.potentiallyFixableIssuesInfo.currentSuggestion.updateAllInstances {
							replaceCount = -1
						}

						m.potentiallyFixableIssuesInfo.fileTexts[m.potentiallyFixableIssuesInfo.currentFile] = strings.Replace(m.potentiallyFixableIssuesInfo.fileTexts[m.potentiallyFixableIssuesInfo.currentFile], m.potentiallyFixableIssuesInfo.currentSuggestionState.original, m.potentiallyFixableIssuesInfo.currentSuggestionState.currentSuggestion, replaceCount)

						m.potentiallyFixableIssuesInfo.currentSuggestionState.isAccepted = true

						if m.potentiallyFixableIssuesInfo.currentSuggestion.addCssSectionBreakIfMissing {
							m.potentiallyFixableIssuesInfo.addCssSectionBreakIfMissing = true
							m.potentiallyFixableIssuesInfo.cssUpdateRequired = true
						} else if m.potentiallyFixableIssuesInfo.currentSuggestion.addCssPageBreakIfMissing {
							m.potentiallyFixableIssuesInfo.addCssPageBreakIfMissing = true
							m.potentiallyFixableIssuesInfo.cssUpdateRequired = true
						}

						err := m.moveToNextSuggestion()
						if err != nil {
							m.Err = err

							return m, tea.Quit
						}
					}
				case "e":
					if m.potentiallyFixableIssuesInfo.currentSuggestionState != nil && !m.potentiallyFixableIssuesInfo.currentSuggestionState.isAccepted {
						m.potentiallyFixableIssuesInfo.isEditing = true
						m.potentiallyFixableIssuesInfo.suggestionEdit.SetValue(m.potentiallyFixableIssuesInfo.currentSuggestionState.currentSuggestion)

						cmd = m.potentiallyFixableIssuesInfo.suggestionEdit.Focus()
						cmds = append(cmds, cmd)
					}
				}
			}
		case stageCssSelection:
			switch key {
			case "ctrl+c", "esc":
				return m, tea.Quit
			case "up":
				if m.cssSelectionInfo.currentCssIndex > 0 {
					m.cssSelectionInfo.currentCssIndex--
					m.cssSelectionInfo.selectedCssFile = m.cssSelectionInfo.cssFiles[m.cssSelectionInfo.currentCssIndex]
				}
			case "down":
				if m.cssSelectionInfo.currentCssIndex+1 < len(m.cssSelectionInfo.cssFiles) {
					m.cssSelectionInfo.currentCssIndex++
					m.cssSelectionInfo.selectedCssFile = m.cssSelectionInfo.cssFiles[m.cssSelectionInfo.currentCssIndex]
				}
			case "enter":
				m.currentStage = finalStage
			}
		default:
			switch key {
			case "ctrl+c", "esc":
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

		m.potentiallyFixableIssuesInfo.suggestionEdit.SetWidth(m.width - columnPadding)

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

func (m fixableIssuesModel) View() string {
	var s strings.Builder
	s.WriteString(appNameStyle.Render("Ebook Linter Manually Fixable Issues") + "\n")

	switch m.currentStage {
	case sectionBreak:
		s.WriteString("\n" + m.sectionBreakInfo.input.View() + "\n\n")
	case suggestionsProcessing:
		s.WriteString(titleStyle.Render(fmt.Sprintf("Current File: %s", m.potentiallyFixableIssuesInfo.currentFile)) + "\n")
		s.WriteString(fileStatusStyle.Render(fmt.Sprintf("File %d of %d", m.potentiallyFixableIssuesInfo.currentFileIndex+1, len(m.potentiallyFixableIssuesInfo.filePaths))) + "\n")
		s.WriteString(groupStyle.Render(fmt.Sprintf("Issue Group: %s", m.potentiallyFixableIssuesInfo.currentSuggestionName) + "\n"))

		if m.potentiallyFixableIssuesInfo.currentSuggestionState == nil {
			s.WriteString("\nNo current file is selected. Something may have gone wrong...\n\n")
		} else {
			if m.potentiallyFixableIssuesInfo.isEditing {
				s.WriteString("\nEditing suggestion:\n\n")
			} else if m.potentiallyFixableIssuesInfo.currentSuggestionState.isAccepted {
				s.WriteString("\n" + acceptedChangeTitleStyle.Render("Accepted change:") + "\n")
			} else {
				s.WriteString("\nSuggested change:\n")
			}

			if m.potentiallyFixableIssuesInfo.isEditing {
				s.WriteString(m.potentiallyFixableIssuesInfo.suggestionEdit.View() + "\n\n")
			} else {
				s.WriteString(displayStyle.Width(m.width-columnPadding).Render(fmt.Sprintf(`"%s"`, m.potentiallyFixableIssuesInfo.currentSuggestionState.display)) + "\n\n\033[0m")
				s.WriteString(fmt.Sprintf("Suggestion %d of %d.\n\n", m.potentiallyFixableIssuesInfo.currentSuggestionIndex+1, len(m.potentiallyFixableIssuesInfo.sectionSuggestionStates)))
			}
		}
	case stageCssSelection:
		s.WriteString("\nSelect the CSS file to modify:\n\n")
		for i, cssFile := range m.cssSelectionInfo.cssFiles {
			cursor := " "
			if m.cssSelectionInfo.currentCssIndex == i {
				cursor = ">"
			}

			s.WriteString(fmt.Sprintf("%s %d. %s\n", cursor, i+1, cssFile))
		}

		s.WriteString("\n")
	}

	m.displayControls(&s)

	return s.String()
}

func (m fixableIssuesModel) displayControls(s *strings.Builder) {
	s.WriteString(groupStyle.Render("Controls:") + "\n")

	// TODO: swap to esc as exit and ctrl+c as kill
	var controls []string
	switch m.currentStage {
	case sectionBreak:
		controls = []string{
			"Enter: Accept",
			"Ctrl+C/Esc: Quit",
		}
	case suggestionsProcessing:
		// TODO: handle edit mode
		if m.potentiallyFixableIssuesInfo.isEditing {
			controls = []string{
				"Ctrl+R: Reset",
				"Ctrl+E: Cancel edit",
				"Ctrl+S: Accept",
				"Ctrl+C/Esc: Quit",
				// "Esc: Quit",
				// "Ctrl+C: Exit without saving",
			}
		} else if m.potentiallyFixableIssuesInfo.currentSuggestionState != nil && m.potentiallyFixableIssuesInfo.currentSuggestionState.isAccepted {
			controls = []string{
				"← / → : Previous/Next Suggestion",
				"C: Copy",
				"Ctrl+C/Esc: Quit",
				// "Ctrl+C: Exit without saving",
			}
		} else {
			controls = []string{
				"← / → : Previous/Next Suggestion",
				"E: Edit",
				"C: Copy",
				"Enter: Accept",
				"Ctrl+C/Esc: Quit",
				// "Ctrl+C: Exit without saving",
			}
		}
	case stageCssSelection:
		controls = []string{
			"↑ / ↓ : Previous/Next Suggestion",
			"Enter: Accept",
			"Ctrl+C/Esc: Quit",
			// "Ctrl+C: Exit without saving",
		}
	}

	s.WriteString(strings.Join(controls, " • ") + "\n")
}

func (m *fixableIssuesModel) setupForNextSuggestions() error {
	for m.potentiallyFixableIssuesInfo.currentFileIndex+1 < len(m.potentiallyFixableIssuesInfo.filePaths) {
		m.potentiallyFixableIssuesInfo.currentFile = m.potentiallyFixableIssuesInfo.filePaths[m.potentiallyFixableIssuesInfo.currentFileIndex]

		for m.potentiallyFixableIssuesInfo.potentialFixableIssueIndex+1 < len(m.potentiallyFixableIssuesInfo.suggestions) {
			var potentialFixableIssue = m.potentiallyFixableIssuesInfo.suggestions[m.potentiallyFixableIssuesInfo.potentialFixableIssueIndex]
			if !m.runAll && (potentialFixableIssue.isEnabled == nil || *potentialFixableIssue.isEnabled) {
				m.potentiallyFixableIssuesInfo.potentialFixableIssueIndex++
				continue
			}

			var (
				suggestions = potentialFixableIssue.getSuggestions(m.potentiallyFixableIssuesInfo.fileTexts[m.potentiallyFixableIssuesInfo.currentFile])
			)

			if len(suggestions) != 0 {
				m.potentiallyFixableIssuesInfo.currentSuggestion = &potentialFixableIssue
				m.potentiallyFixableIssuesInfo.sectionSuggestionStates = make([]suggestionState, len(suggestions))

				var i = 0
				for original, suggestion := range suggestions {
					var display, err = getStringDiff(original, suggestion)
					if err != nil {
						return err
					}

					m.potentiallyFixableIssuesInfo.sectionSuggestionStates[i] = suggestionState{
						original:           original,
						originalSuggestion: suggestion,
						currentSuggestion:  suggestion,
						display:            display,
					}

					i++
				}

				m.potentiallyFixableIssuesInfo.currentSuggestionIndex = 0
				m.potentiallyFixableIssuesInfo.currentSuggestionState = &m.potentiallyFixableIssuesInfo.sectionSuggestionStates[0]
				m.potentiallyFixableIssuesInfo.currentSuggestionName = potentialFixableIssue.name

				m.potentiallyFixableIssuesInfo.potentialFixableIssueIndex++

				return nil
			}

			m.potentiallyFixableIssuesInfo.potentialFixableIssueIndex++
		}

		m.potentiallyFixableIssuesInfo.currentFileIndex++
		m.potentiallyFixableIssuesInfo.potentialFixableIssueIndex = 0
	}

	if m.potentiallyFixableIssuesInfo.cssUpdateRequired {
		m.currentStage = stageCssSelection
	} else {
		m.currentStage = finalStage
	}

	return nil
}

func (m *fixableIssuesModel) moveToNextSuggestion() error {
	if m.potentiallyFixableIssuesInfo.currentSuggestionIndex+1 < len(m.potentiallyFixableIssuesInfo.sectionSuggestionStates) {
		m.potentiallyFixableIssuesInfo.currentSuggestionIndex++
		m.potentiallyFixableIssuesInfo.currentSuggestionState = &m.potentiallyFixableIssuesInfo.sectionSuggestionStates[m.potentiallyFixableIssuesInfo.currentSuggestionIndex]

		return nil
	}

	return m.setupForNextSuggestions()
}

func (m *fixableIssuesModel) moveToPreviousSuggestion() {
	if m.potentiallyFixableIssuesInfo.currentSuggestionIndex > 0 {
		m.potentiallyFixableIssuesInfo.currentSuggestionIndex--
		m.potentiallyFixableIssuesInfo.currentSuggestionState = &m.potentiallyFixableIssuesInfo.sectionSuggestionStates[m.potentiallyFixableIssuesInfo.currentSuggestionIndex]
	}
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

// func repairUTF8(s string) (string, error) {
// 	buf := make([]byte, 0, len(s))
// 	for i, r := range s {
// 		b, size := utf8.DecodeRune([]byte(r)
// 		// if err != nil {
// 		// 	return "", fmt.Errorf("character at index %d is not valid UTF-8", i)
// 		// }
// 		buf = append(buf, b)
// 		i += size - 1
// 	}
// 	return string(buf), nil
// }

// func repairUnicode(s string) (string, error) {
// 	buf := make([]byte, 0, len(s))
// 	for i, r := range s {
// 		if r > 0x10FFFF {
// 			return "", fmt.Errorf("character %q at index %d is not part of Unicode", string(r), i)
// 		}
// 		buf = append(buf, byte(r))
// 	}
// 	return string(buf), nil
// }
