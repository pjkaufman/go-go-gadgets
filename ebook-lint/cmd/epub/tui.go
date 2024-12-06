package epub

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	stringdiff "github.com/pjkaufman/go-go-gadgets/pkg/string-diff"
)

// var (
// 	appNameStyle = lipgloss.NewStyle().Background(lipgloss.Color("99")).Padding(0, 1)
// )

// type fixableTuiModel struct {
// 	currentStage stage
// 	currentFile  string
// 	// sectionBreakInput                                   tui.SectionBreakModel
// 	// suggestionHandler                                   tui.SuggestionsModel
// 	sectionBreakInput                                   textinput.Model
// 	editableInput                                       textarea.Model
// 	potentiallyFixableIssues                            []potentiallyFixableIssue
// filePaths                                           []string
// fileTexts                                           map[string]string
// 	currentFilePathIndex, currentSuggestIndex           int
// 	width, height                                       int
// 	cssFiles, handledFiles                              []string
// 	runAll, addCssSectionIfMissing, addCssPageIfMissing bool
// 	Err                                                 error
// }

// type stage int

// const (
// 	sectionContextBreak stage = iota
// 	suggestionsProcessing
// 	stageCssSelection
// 	finalStage
// )

// func newFixableTuiModel(runAll, runSectionBreak bool, potentiallyFixableIssues []potentiallyFixableIssue, cssFiles []string) fixableTuiModel {
// 	var startingStage = sectionContextBreak
// 	if !runAll && !runSectionBreak {
// 		startingStage = suggestionsProcessing
// 	}

// 	ti := textinput.New()
// 	ti.Placeholder = "Section Break"
// 	ti.CharLimit = 100
// 	ti.Width = 20

// 	ta := textarea.New()
// 	ta.Placeholder = "Enter an edited version of the original string"
// 	ta.CharLimit = math.MaxInt
// 	ta.ShowLineNumbers = false
// 	// var input = tui.NewSectionBreakModel()

// 	return fixableTuiModel{
// 		// sectionBreakInput:        input,
// 		editableInput:            ta,
// 		sectionBreakInput:        ti,
// 		currentStage:             startingStage,
// 		potentiallyFixableIssues: potentiallyFixableIssues,
// 		fileTexts:                make(map[string]string),
// 		cssFiles:                 cssFiles,
// 		runAll:                   runAll,
// 	}
// }

// func (f fixableTuiModel) Init() tea.Cmd {
// 	return nil
// }

// // func (f fixableTuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// // 	var cmd tea.Cmd

// // 	if f.Err != nil {
// // 		return f, tea.Quit
// // 	}

// // 	switch msg := msg.(type) {
// // 	case tea.WindowSizeMsg:
// // 		f.width = msg.Width
// // 		f.height = msg.Height
// // 	}

// // 	switch f.currentStage {
// // 	case sectionContextBreak:
// // 		if f.sectionBreakInput.Err != nil {
// // 			f.Err = f.sectionBreakInput.Err

// // 			return f, tea.Quit
// // 		}

// // 		f.sectionBreakInput, cmd = f.sectionBreakInput.Update(msg)
// // 		f = f.advanceStageIfNeeded()

// // 		return f, cmd
// // 	case suggestionsProcessing:
// // 		if f.suggestionHandler.Err != nil {
// // 			f.Err = f.suggestionHandler.Err

// // 			return f, tea.Quit
// // 		}

// // 		f.suggestionHandler, cmd = f.suggestionHandler.Update(msg)
// // 		f = f.advanceStageIfNeeded()

// // 		return f, cmd
// // 	}

// // 	return f, cmd
// // }

// func (f fixableTuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	if f.Err != nil {
// 		return f, tea.Quit
// 	}

// 	var (
// 		cmd  tea.Cmd
// 		cmds []tea.Cmd
// 	)

// 	f.sectionBreakInput, cmd = f.sectionBreakInput.Update(cmd)
// 	cmds = append(cmds, cmd)

// 	f.editableInput, cmd = f.editableInput.Update(cmd)
// 	cmds = append(cmds, cmd)

// 	switch msg := msg.(type) {
// 	case tea.WindowSizeMsg:
// 		f.width = msg.Width
// 		f.height = msg.Height

// 		f.editableInput.SetWidth(f.width)

// 		cmds = append(cmds, tea.ClearScreen)
// 	case tea.Key:
// 		switch msg.Type {
// 		case tea.KeyCtrlC, tea.KeyEsc:
// 			return f, tea.Quit
// 		}

// 	case error:
// 		f.Err = msg

// 		return f, tea.Quit
// 	}

// 	// switch f.currentStage {
// 	// case sectionContextBreak:
// 	// 	if f.sectionBreakInput.Err != nil {
// 	// 		f.Err = f.sectionBreakInput.Err

// 	// 		return f, tea.Quit
// 	// 	}

// 	// 	f.sectionBreakInput, cmd = f.sectionBreakInput.Update(msg)
// 	// 	f = f.advanceStageIfNeeded()

// 	// 	return f, cmd
// 	// case suggestionsProcessing:
// 	// 	if f.suggestionHandler.Err != nil {
// 	// 		f.Err = f.suggestionHandler.Err

// 	// 		return f, tea.Quit
// 	// 	}

// 	// 	f.suggestionHandler, cmd = f.suggestionHandler.Update(msg)
// 	// 	f = f.advanceStageIfNeeded()

// 	// 	return f, cmd
// 	// }

// 	return f, tea.Batch(cmds...)
// }

// func (f fixableTuiModel) View() string {
// 	var s strings.Builder
// 	s.WriteString(appNameStyle.Render("Ebook Linter Manually Fixable Issues") + "\n\n")

// 	switch f.currentStage {
// 	// case sectionContextBreak:
// 	// 	return f.sectionBreakInput.View()
// 	// case suggestionsProcessing:
// 	// 	return f.suggestionHandler.View()
// 	}

// 	return s.String()
// }

// // func (f fixableTuiModel) advanceStageIfNeeded() fixableTuiModel {
// // 	switch f.currentStage {
// // 	case sectionContextBreak:
// // 		if f.sectionBreakInput.Done {
// // 			contextBreak = f.sectionBreakInput.Value()
// // 			f.currentStage = suggestionsProcessing
// // 			f = f.getNextSuggestion()
// // 		}
// // 	case suggestionsProcessing:
// // 		if f.suggestionHandler.Done {
// // 			if f.currentFilePathIndex+1 < len(f.filePaths) {
// // 				f = f.getNextSuggestion()

// // 				return f
// // 			}

// // 			f.currentStage = stageCssSelection
// // 			f.Err = fmt.Errorf("Not implemented stage yet")
// // 			// TODO: figure out how to close the program
// // 		}
// // 	}

// // 	return f
// // }

// func (f fixableTuiModel) getNextSuggestion() fixableTuiModel {
// 	for f.currentFilePathIndex+1 < len(f.filePaths) {
// 		for f.currentSuggestIndex+1 < len(f.potentiallyFixableIssues) {
// 			var (
// 				currentFilePath = f.filePaths[f.currentFilePathIndex]
// 				suggestions     = f.potentiallyFixableIssues[f.currentSuggestIndex].getSuggestions(f.fileTexts[currentFilePath])
// 			)

// 			if len(suggestions) != 0 {
// 				var err error
// 				f.suggestionHandler, err = tui.NewSuggestionsModel(currentFilePath, f.potentiallyFixableIssues[f.currentSuggestIndex].name, fmt.Sprintf("File %d of %d", f.currentFilePathIndex+1, len(f.filePaths)), suggestions, f.width, f.height)
// 				if err != nil {
// 					f.Err = err

// 					return f
// 				}

// 				f.currentSuggestIndex++

// 				return f
// 			}

// 			f.currentSuggestIndex++
// 		}

// 		f.currentFilePathIndex++
// 		f.currentSuggestIndex = 0
// 	}

// 	return f
// }

var (
	appNameStyle    = lipgloss.NewStyle().Background(lipgloss.Color("99")).Padding(0, 1)
	titleStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	groupStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("246"))
	fileStatusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("190"))
	displayStyle    = lipgloss.NewStyle()
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
	currentStage                 stage
	sectionBreakInfo             sectionBreakStageInfo
	potentiallyFixableIssuesInfo potentiallyFixableStageInfo
	cssSelectionInfo             cssSelectionStageInfo
	// filePaths, cssFiles                                  []string
	// fileTexts                                            map[string]string
	runAll bool
	// potentiallyFixableIssues                             []potentiallyFixableIssue
	// potentialFixableIssueIndex, currentFileIndex         int
	// sectionSuggestionStates                              []suggestionState
	// currentFile, currentSuggestionTitle, SelectedCssFile string
	// currentSuggestionIndex, currentCssIndex              int
	// currentFileState                                     *suggestionState
	height, width int
	Err           error
}

type sectionBreakStageInfo struct {
	sectionBreakInput textinput.Model
}

type potentiallyFixableStageInfo struct {
	currentFileState                             *suggestionState
	currentSuggestionIndex                       int
	currentFile, currentSuggestionTitle          string
	potentiallyFixableIssues                     []potentiallyFixableIssue
	currentSuggestion                            *potentiallyFixableIssue
	potentialFixableIssueIndex, currentFileIndex int
	cssUpdateRequired                            bool
	sectionSuggestionStates                      []suggestionState
	filePaths                                    []string
	fileTexts                                    map[string]string
}

type cssSelectionStageInfo struct {
	cssFiles        []string
	SelectedCssFile string
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

	var currentStage = sectionBreak
	if runAll || runSectionBreak {
		ti.Focus()
	} else {
		currentStage = suggestionsProcessing
	}

	return fixableIssuesModel{
		sectionBreakInfo: sectionBreakStageInfo{
			sectionBreakInput: ti,
		},
		potentiallyFixableIssuesInfo: potentiallyFixableStageInfo{
			fileTexts:                map[string]string{},
			potentiallyFixableIssues: potentiallyFixableIssues,
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

	m.sectionBreakInfo.sectionBreakInput, cmd = m.sectionBreakInfo.sectionBreakInput.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch m.currentStage {
		case sectionBreak:
			switch key {
			case "enter":
				contextBreak = strings.TrimSpace(m.sectionBreakInfo.sectionBreakInput.Value())
				if contextBreak != "" {
					m.currentStage = suggestionsProcessing

					err := m.setupForNextSuggestions()
					if err != nil {
						m.Err = err

						return m, tea.Quit
					}

					// cmds = append(cmds, tea.ClearScreen)
				}
			case "ctrl+c", "esc":
				return m, tea.Quit
			}
		case suggestionsProcessing:
			switch key {
			case "ctrl+c", "esc":
				return m, tea.Quit
			case "right":
				err := m.moveToNextSuggestion()
				if err != nil {
					m.Err = err

					return m, tea.Quit
				}

				// cmds = append(cmds, tea.ClearScreen)
			case "left":
				m.moveToPreviousSuggestion()

				// cmds = append(cmds, tea.ClearScreen)
			case "c":
				// Copy original value to the clipboard
				// original, err := repairUnicode(m.currentFileState.original)
				// if err != nil {
				// 	m.Err = err

				// 	return m, tea.Quit
				// }

				// err = clipboard.WriteAll(original)
				// TODO: make sure values are utf-8 compliant
				err := clipboard.WriteAll(m.potentiallyFixableIssuesInfo.currentFileState.original)
				if err != nil {
					m.Err = err

					return m, tea.Quit
				}
			case "enter":
				if !m.potentiallyFixableIssuesInfo.currentFileState.isAccepted && m.potentiallyFixableIssuesInfo.currentSuggestion != nil {
					var replaceCount = 1
					if m.potentiallyFixableIssuesInfo.currentSuggestion.updateAllInstances {
						replaceCount = -1
					}

					m.potentiallyFixableIssuesInfo.fileTexts[m.potentiallyFixableIssuesInfo.currentFile] = strings.Replace(m.potentiallyFixableIssuesInfo.fileTexts[m.potentiallyFixableIssuesInfo.currentFile], m.potentiallyFixableIssuesInfo.currentFileState.original, m.potentiallyFixableIssuesInfo.currentFileState.currentSuggestion, replaceCount)

					if m.potentiallyFixableIssuesInfo.currentSuggestion.addCssIfMissing {
						m.potentiallyFixableIssuesInfo.cssUpdateRequired = true
					}

					err := m.moveToNextSuggestion()
					if err != nil {
						m.Err = err

						return m, tea.Quit
					}

					// cmds = append(cmds, tea.ClearScreen)
				}
			}
		case stageCssSelection:
			switch key {
			case "ctrl+c", "esc":
				return m, tea.Quit
			case "up":
				if m.cssSelectionInfo.currentCssIndex > 0 {
					m.cssSelectionInfo.currentCssIndex--
					m.cssSelectionInfo.SelectedCssFile = m.cssSelectionInfo.cssFiles[m.cssSelectionInfo.currentCssIndex]
				}
			case "down":
				if m.cssSelectionInfo.currentCssIndex+1 < len(m.cssSelectionInfo.cssFiles) {
					m.cssSelectionInfo.currentCssIndex++
					m.cssSelectionInfo.SelectedCssFile = m.cssSelectionInfo.cssFiles[m.cssSelectionInfo.currentCssIndex]
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
		s.WriteString("\n" + m.sectionBreakInfo.sectionBreakInput.View() + "\n\n")
	case suggestionsProcessing:
		s.WriteString(titleStyle.Render(fmt.Sprintf("Current File: %s", m.potentiallyFixableIssuesInfo.currentFile)) + "\n")
		s.WriteString(fileStatusStyle.Render(fmt.Sprintf("File %d of %d", m.potentiallyFixableIssuesInfo.currentFileIndex+1, len(m.potentiallyFixableIssuesInfo.filePaths))) + "\n")
		s.WriteString(groupStyle.Render(fmt.Sprintf("Issue Group: %s", m.potentiallyFixableIssuesInfo.currentSuggestionTitle) + "\n"))

		if m.potentiallyFixableIssuesInfo.currentFileState == nil {
			s.WriteString("\nNo current file is selected. Something may have gone wrong...\n\n")
		} else {
			s.WriteString("\nSuggested change:\n")
			s.WriteString(displayStyle.Width(m.width-columnPadding).Render(fmt.Sprintf(`"%s"`, m.potentiallyFixableIssuesInfo.currentFileState.display)) + "\n\n")
			s.WriteString(fmt.Sprintf("Suggestion %d of %d.\n\n", m.potentiallyFixableIssuesInfo.currentSuggestionIndex+1, len(m.potentiallyFixableIssuesInfo.sectionSuggestionStates)))
		}
	case stageCssSelection:
		var startingChar = "\n"
		for i, cssFile := range m.cssSelectionInfo.cssFiles {

			cursor := " "
			if m.cssSelectionInfo.currentCssIndex == i {
				cursor = ">"
			}

			s.WriteString(fmt.Sprintf("%s%s %d. %s\n", startingChar, cursor, i+1, cssFile))

			startingChar = ""
		}

		s.WriteString("\n")
	}

	m.displayControls(&s)

	s.WriteString("\n")

	switch m.currentStage {
	case sectionBreak:
		s.WriteString("Stage: section break")
	case suggestionsProcessing:
		s.WriteString("Stage: suggestion processing")
	default:
		s.WriteString("Stage: other")
	}

	return s.String()
}

func (m fixableIssuesModel) displayControls(s *strings.Builder) {
	s.WriteString(groupStyle.Render("Controls:") + "\n")

	var controls []string
	switch m.currentStage {
	case sectionBreak:
		controls = []string{
			"Enter: Accept",
			"Ctrl+C/Esc: Quit",
		}
	case suggestionsProcessing:
		// TODO: handle edit mode
		controls = []string{
			"← / → : Previous/Next Suggestion",
			// "E: Edit", // TODO: decide how to add this
			"C: Copy",
			"Enter: Accept",
			"Ctrl+C/Esc: Quit",
			// "Ctrl+C: Exit without saving",
		}
	case stageCssSelection:
		controls = []string{
			"↑ / ↓ : Previous/Next Suggestion",
			"Enter: Accept",
			"Ctrl+C/Esc: Quit",
			// "Ctrl+C: Exit without saving",
		}
	}

	// if !f.editMode {
	// controls = []string{
	// 	"← / → : Previous/Next Suggestion",
	// 	"E: Edit",
	// 	"C: Copy",
	// 	"Enter: Accept",
	// 	"Q/Esc: Quit",
	// 	"Ctrl+C: Exit without saving",
	// 	}
	// } else {
	// 	controls = []string{
	// 		"Ctrl+R: Reset",
	// 		"Ctrl+E: Cancel edit",
	// 		"Ctrl+S: Accept",
	// 		"Esc: Quit",
	// 		"Ctrl+C: Exit without saving",
	// 	}
	// }

	s.WriteString(strings.Join(controls, " • ") + "\n")
}

func (m *fixableIssuesModel) setupForNextSuggestions() error {
	for m.potentiallyFixableIssuesInfo.currentFileIndex+1 < len(m.potentiallyFixableIssuesInfo.filePaths) {
		for m.potentiallyFixableIssuesInfo.potentialFixableIssueIndex+1 < len(m.potentiallyFixableIssuesInfo.potentiallyFixableIssues) {
			var potentialFixableIssue = m.potentiallyFixableIssuesInfo.potentiallyFixableIssues[m.potentiallyFixableIssuesInfo.potentialFixableIssueIndex]
			if !m.runAll && (potentialFixableIssue.isEnabled == nil || *potentialFixableIssue.isEnabled) {
				m.potentiallyFixableIssuesInfo.potentialFixableIssueIndex++
				continue
			}

			var (
				currentFilePath = m.potentiallyFixableIssuesInfo.filePaths[m.potentiallyFixableIssuesInfo.currentFileIndex]
				suggestions     = potentialFixableIssue.getSuggestions(m.potentiallyFixableIssuesInfo.fileTexts[currentFilePath])
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
				m.potentiallyFixableIssuesInfo.currentFileState = &m.potentiallyFixableIssuesInfo.sectionSuggestionStates[0]
				m.potentiallyFixableIssuesInfo.currentSuggestionTitle = potentialFixableIssue.name

				m.potentiallyFixableIssuesInfo.potentialFixableIssueIndex++

				return nil
			}

			m.potentiallyFixableIssuesInfo.potentialFixableIssueIndex++
		}

		m.potentiallyFixableIssuesInfo.currentFileIndex++
		m.potentiallyFixableIssuesInfo.potentialFixableIssueIndex = 0
	}

	// if m.potentiallyFixableIssuesInfo.cssUpdateRequired {
	m.currentStage = stageCssSelection
	// } else {
	// 	m.currentStage = finalStage
	// }

	return nil
}

func (m *fixableIssuesModel) moveToNextSuggestion() error {
	if m.potentiallyFixableIssuesInfo.currentSuggestionIndex+1 < len(m.potentiallyFixableIssuesInfo.sectionSuggestionStates) {
		m.potentiallyFixableIssuesInfo.currentSuggestionIndex++
		m.potentiallyFixableIssuesInfo.currentFileState = &m.potentiallyFixableIssuesInfo.sectionSuggestionStates[m.potentiallyFixableIssuesInfo.currentSuggestionIndex]

		return nil
	}

	return m.setupForNextSuggestions()
}

func (m *fixableIssuesModel) moveToPreviousSuggestion() {
	if m.potentiallyFixableIssuesInfo.currentSuggestionIndex > 0 {
		m.potentiallyFixableIssuesInfo.currentSuggestionIndex--
		m.potentiallyFixableIssuesInfo.currentFileState = &m.potentiallyFixableIssuesInfo.sectionSuggestionStates[m.potentiallyFixableIssuesInfo.currentSuggestionIndex]
	}
}

// var removeStartingLineWhitespace = regexp.MustCompile(`(^|\n)[ \t]+`)

func getStringDiff(original, new string) (string, error) {
	// original = strings.TrimSpace(removeStartingLineWhitespace.ReplaceAllString(original, "\n"))
	// new = strings.TrimSpace(removeStartingLineWhitespace.ReplaceAllString(new, "\n"))

	// var diffString, err = stringdiff.GetPrettyDiffString(original, new)
	// diffString = strings.TrimSpace(diffString)

	return stringdiff.GetPrettyDiffString(strings.TrimLeft(original, "\n"), strings.TrimLeft(new, "\n"))
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
