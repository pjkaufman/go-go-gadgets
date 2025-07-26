package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// type suggestions struct {
// 	// state                                                                *State
// 	currentFile, currentSuggestionName                                   string
// 	isEditing                                                            bool
// 	suggestionData                                                       [][]suggestionState
// 	currentSuggestion                                                    *suggestionState
// 	currentSuggestionIndex, currentFileIndex, potentialFixableIssueIndex int

// 	// currentSuggestionIndex, potentialFixableIssueIndex, currentFileIndex int
// 	potentialIssues []potentiallyfixableissue.PotentiallyFixableIssue
// 	currentIssue    *potentiallyfixableissue.PotentiallyFixableIssue
// 	// cssUpdateRequired, addCssSectionBreakIfMissing, addCssPageBreakIfMissing, isEditing bool
// 	// currentSuggestionState *suggestionState
// 	// suggestionEdit         textarea.Model
// 	suggestionDisplay viewport.Model
// 	// scrollbar              tea.Model
// }

// type suggestionState struct {
// 	isAccepted                                               bool
// 	original, originalSuggestion, currentSuggestion, display string
// }

// func newSuggestions(state *State) suggestions {
// 	v := viewport.New(0, 0)

// 	return suggestions{
// 		state:                 state,
// 		currentFile:           "OEBS/Text/file.html",
// 		currentSuggestionName: "Suggestion Name",
// 		suggestionDisplay:     v,
// 		// suggestionData: []fileSuggestionInfo{
// 		// 	{
// 		// 		fileName: "OEBS/Text/file.html",
// 		// 		suggestions: [][]suggestionState{
// 		// 			{
// 		// 				{
// 		// 					original:           "This is the original",
// 		// 					originalSuggestion: "This is the new display value. How do you like them apples?",
// 		// 					currentSuggestion:  "This is the new display value. How do you like them apples?",
// 		// 					display:            "This is the new display value. How do you like them apples?",
// 		// 				},
// 		// 				{
// 		// 					original:           "Suggestion 2 is even longer than original. How does this play?",
// 		// 					originalSuggestion: "New suggestion is here to stay and play. How are things going to look?",
// 		// 					currentSuggestion:  "New suggestion is here to stay and play. How are things going to look?",
// 		// 					display:            "New suggestion is here to stay and play. How are things going to look?",
// 		// 				},
// 		// 			},
// 		// 		},
// 		// 	},
// 		// },
// 	}
// }

func (m FixableIssuesModel) suggestionsView() string {
	var (
		status         = m.leftStatusView()
		height         = m.body.Height - leftStatusBorderStyle.GetVerticalBorderSize()
		remainingWidth = m.body.Width - (lipgloss.Width(status) + leftStatusBorderStyle.GetBorderLeftSize() + suggestionBorderStyle.GetBorderRightSize())
	)

	m.PotentiallyFixableIssuesInfo.suggestionDisplay.Height = height
	m.PotentiallyFixableIssuesInfo.suggestionDisplay.Width = remainingWidth
	m.PotentiallyFixableIssuesInfo.suggestionEdit.SetWidth(remainingWidth)
	m.PotentiallyFixableIssuesInfo.suggestionEdit.SetHeight(height)

	return lipgloss.JoinHorizontal(lipgloss.Left, status, m.suggestionView())
}

func (m FixableIssuesModel) leftStatusView() string {
	var (
		statusView      = fmt.Sprintf("%s %s\n%s %s\n", documentIcon, fileNameStyle.Render(m.PotentiallyFixableIssuesInfo.currentFile), suggestionIcon, suggestionNameStyle.Render(m.PotentiallyFixableIssuesInfo.currentSuggestionName))
		remainingHeight int
		statusPadding   string
	)

	remainingHeight = m.body.Height - (lipgloss.Height(statusView) + leftStatusBorderStyle.GetVerticalBorderSize())

	if remainingHeight > 0 {
		statusPadding = strings.Repeat("\n", remainingHeight)
	}

	return leftStatusBorderStyle.Render(statusView + statusPadding)
}

func (m FixableIssuesModel) suggestionView() string {
	var s strings.Builder
	// s.WriteString(titleStyle.Render(fmt.Sprintf("Current File (%d/%d): %s ", m.PotentiallyFixableIssuesInfo.currentFileIndex+1, len(m.PotentiallyFixableIssuesInfo.FilePaths), m.PotentiallyFixableIssuesInfo.currentFile)) + "\n")
	// s.WriteString(groupStyle.Render(fmt.Sprintf("Issue Group: %s", m.PotentiallyFixableIssuesInfo.currentSuggestionName) + "\n"))

	if m.PotentiallyFixableIssuesInfo.currentSuggestionState == nil {
		s.WriteString("\nNo current file is selected. Something may have gone wrong...\n\n")
	} else {
		// var (
		// 	suggestion = displayStyle.Width(m.width - columnPadding).Render(fmt.Sprintf(`"%s"`, m.PotentiallyFixableIssuesInfo.currentSuggestionState.display))
		// 	// expectedHeight = strings.Count(suggestion, "\n") + 1
		// )
		// if m.PotentiallyFixableIssuesInfo.isEditing {
		// 	s.WriteString("\nEditing suggestion:\n\n")
		// } else if m.PotentiallyFixableIssuesInfo.currentSuggestionState.isAccepted {
		// 	s.WriteString("\n" + acceptedChangeTitleStyle.Render("Accepted change:") + "\n")
		// } else {
		// 	s.WriteString(fmt.Sprintf("\nSuggested change (%d/%d):\n", expectedHeight, maxDisplayHeight))
		// }

		if m.PotentiallyFixableIssuesInfo.isEditing {
			s.WriteString(suggestionBorderStyle.Render(m.PotentiallyFixableIssuesInfo.suggestionEdit.View()) + "\n\n")
		} else {
			if m.PotentiallyFixableIssuesInfo.suggestionDisplay.TotalLineCount() > m.PotentiallyFixableIssuesInfo.suggestionDisplay.VisibleLineCount() {
				s.WriteString(suggestionBorderStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top,
					m.PotentiallyFixableIssuesInfo.suggestionDisplay.View(),
					m.PotentiallyFixableIssuesInfo.scrollbar.View(),
				)))
			} else {
				s.WriteString(suggestionBorderStyle.Render(m.PotentiallyFixableIssuesInfo.suggestionDisplay.View()))
			}

			s.WriteString(fmt.Sprintf("\033[0m\n\nSuggestion %d of %d.\n\n", m.PotentiallyFixableIssuesInfo.currentSuggestionIndex+1, len(m.PotentiallyFixableIssuesInfo.sectionSuggestionStates)))
		}
	}

	return s.String()
}

// func (m FixableIssuesModel) suggestionView() string {
// 	if !m.isEditing {
// 		m.suggestionDisplay.SetContent(m.currentSuggestion.currentSuggestion)
// 	}

// 	return suggestionBorderStyle.Render(m.suggestionDisplay.View())
// }

// func (m *FixableIssuesModel) setupForNextSuggestions() error {
// 	if m.logFile != nil {
// 		fmt.Fprintln(m.logFile, "Getting next suggestions")
// 	}

// 	for m.currentFileIndex < len(m.FilePaths) {
// 		m.currentFile = m.FilePaths[m.currentFileIndex]
// 		if m.logFile != nil {
// 			fmt.Fprintf(m.logFile, "Current file is %q is %d of %d\n", m.currentFile, m.currentFileIndex+1, len(m.FilePaths))
// 		}

// 		for m.potentialFixableIssueIndex < len(m.suggestionData) {
// 			var potentialFixableIssue = m.potentialIssues[m.potentialFixableIssueIndex]
// 			if m.logFile != nil {
// 				fmt.Fprintf(m.logFile, "Possible fixable issue %q is %d of %d issues.\n", potentialFixableIssue.Name, m.potentialFixableIssueIndex+1, len(m.potentialIssues))
// 			}

// 			if !m.RunAll && (potentialFixableIssue.IsEnabled == nil || *potentialFixableIssue.IsEnabled) {
// 				if m.logFile != nil {
// 					fmt.Fprintf(m.logFile, "Skipping possible fixable issue %q with isEnabled set to %v\n", potentialFixableIssue.Name, potentialFixableIssue.IsEnabled)
// 				}

// 				m.potentialFixableIssueIndex++
// 				continue
// 			}

// 			var (
// 				suggestions = potentialFixableIssue.GetSuggestions(m.FileTexts[m.currentFileIndex])
// 			)

// 			if m.logFile != nil {
// 				fmt.Fprintf(m.logFile, "Possible fixable issue %q has %d suggestion(s) found\n", potentialFixableIssue.Name, len(suggestions))
// 			}

// 			if len(suggestions) != 0 {
// 				m.currentIssue = &potentialFixableIssue
// 				m.currentIssueName = potentialFixableIssue.Name
// 				m.suggestionData[m.currentFileIndex] = make([]suggestionState, len(suggestions))

// 				var i = 0
// 				for original, suggestion := range suggestions {
// 					var display, err = getStringDiff(original, suggestion)
// 					if err != nil {
// 						return err
// 					}

// 					m.suggestionData[m.currentFileIndex][i] = suggestionState{
// 						original:           original,
// 						originalSuggestion: suggestion,
// 						currentSuggestion:  suggestion,
// 						display:            display,
// 					}

// 					i++
// 				}

// 				m.currentSuggestionIndex = 0
// 				m.currentSuggestion = &m.suggestionData[m.currentFileIndex][0]
// 				// cmd := m.setSuggestionDisplay()
// 				m.currentIssueName = potentialFixableIssue.Name

// 				m.potentialFixableIssueIndex++

// 				return nil
// 			}

// 			m.potentialFixableIssueIndex++
// 		}

// 		m.currentFileIndex++
// 		m.potentialFixableIssueIndex = 0
// 	}

// 	// TODO: add logic for advancing to the next stage here...
// 	// if m.state.CssUpdateRequired {
// 	// 	m.currentStage = stageCssSelection
// 	// } else {
// 	// 	m.currentStage = finalStage
// 	// }

// 	return nil
// }

// func getStringDiff(original, new string) (string, error) {
// 	return stringdiff.GetPrettyDiffString(strings.TrimLeft(original, "\n"), strings.TrimLeft(new, "\n"))
// }

// // textarea gets rid of tabs when creating changes, so in order to preserve tabs in the starting whitespace of a line
// // we will use the value of original as the template for what whitespace is needed for each line present
// func alignWhitespace(original, new string) string {
// 	origLines := strings.Split(original, "\n")
// 	newLines := strings.Split(new, "\n")

// 	var min = len(newLines)
// 	if len(origLines) < min {
// 		min = len(origLines)
// 	}

// 	for i := 0; i < min; i++ {
// 		origPrefix := ""
// 		for j := 0; j < len(origLines[i]); j++ {
// 			if !unicode.IsSpace(rune(origLines[i][j])) {
// 				break
// 			}
// 			origPrefix += string(origLines[i][j])
// 		}
// 		newLines[i] = origPrefix + strings.TrimLeft(newLines[i], " \t")
// 	}

// 	return strings.Join(newLines, "\n")
// }
