package ui

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
	rw "github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
)

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

	if m.PotentiallyFixableIssuesInfo.currentSuggestionState == nil {
		s.WriteString("\nNo current file is selected. Something may have gone wrong...\n\n")
	} else {
		modeIcon := viewIcon
		modeName := "View"
		if m.PotentiallyFixableIssuesInfo.isEditing {
			modeIcon = editIcon
			modeName = "Edit"
		}

		borderConfig := NewBorderConfig(m.PotentiallyFixableIssuesInfo.suggestionDisplay.Height+2, m.PotentiallyFixableIssuesInfo.suggestionDisplay.Width+2) // +2 for border width/height
		modeInfo := fmt.Sprintf("%s %s", modeIcon, modeName)
		statusInfo := fmt.Sprintf("%d/%d", m.PotentiallyFixableIssuesInfo.currentSuggestionIndex+1, len(m.PotentiallyFixableIssuesInfo.sectionSuggestionStates))
		borderConfig.SetInfoItems(modeInfo, statusInfo)

		baseBorder := lipgloss.RoundedBorder()
		customBorder := borderConfig.GetBorder(baseBorder)
		customBorderStyle := lipgloss.NewStyle().Border(customBorder)
		customBorderStyle = customBorderStyle.BorderForeground(suggestionBorderStyle.GetBorderTopForeground())

		if m.PotentiallyFixableIssuesInfo.isEditing {
			s.WriteString(customBorderStyle.Render(m.PotentiallyFixableIssuesInfo.suggestionEdit.View()) + "\n\n")
		} else if m.PotentiallyFixableIssuesInfo.suggestionDisplay.TotalLineCount() > m.PotentiallyFixableIssuesInfo.suggestionDisplay.VisibleLineCount() {
			s.WriteString(customBorderStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top,
				m.PotentiallyFixableIssuesInfo.suggestionDisplay.View(),
				m.PotentiallyFixableIssuesInfo.scrollbar.View(),
			)))
		} else {
			s.WriteString(customBorderStyle.Render(m.PotentiallyFixableIssuesInfo.suggestionDisplay.View()))
		}
	}

	return s.String()
}

func wrapLines(text string, width int) string {
	var (
		s             strings.Builder
		originalLines = strings.Split(text, "\n")
	)

	for i, line := range originalLines {
		wrappedLines := wrap([]rune(line), width)

		for j, wrappedLine := range wrappedLines {
			s.WriteString(string(wrappedLine))
			if i+1 != len(originalLines) || j+i != len(wrappedLine) {
				s.WriteString("\n")
			}
		}
	}

	return s.String()
}

func wrap(runes []rune, width int) [][]rune {
	var (
		lines  = [][]rune{{}}
		word   = []rune{}
		row    int
		spaces int
	)

	// Word wrap the runes
	for _, r := range runes {
		if unicode.IsSpace(r) {
			spaces++
		} else {
			word = append(word, r)
		}

		if spaces > 0 { //nolint:nestif
			if uniseg.StringWidth(string(lines[row]))+uniseg.StringWidth(string(word))+spaces > width {
				row++
				lines = append(lines, []rune{})
				lines[row] = append(lines[row], word...)
				lines[row] = append(lines[row], repeatSpaces(spaces)...)
				spaces = 0
				word = nil
			} else {
				lines[row] = append(lines[row], word...)
				lines[row] = append(lines[row], repeatSpaces(spaces)...)
				spaces = 0
				word = nil
			}
		} else {
			// If the last character is a double-width rune, then we may not be able to add it to this line
			// as it might cause us to go past the width.
			lastCharLen := rw.RuneWidth(word[len(word)-1])
			if uniseg.StringWidth(string(word))+lastCharLen > width {
				// If the current line has any content, let's move to the next
				// line because the current word fills up the entire line.
				if len(lines[row]) > 0 {
					row++
					lines = append(lines, []rune{})
				}
				lines[row] = append(lines[row], word...)
				word = nil
			}
		}
	}

	if uniseg.StringWidth(string(lines[row]))+uniseg.StringWidth(string(word))+spaces >= width {
		lines = append(lines, []rune{})
		lines[row+1] = append(lines[row+1], word...)
		// We add an extra space at the end of the line to account for the
		// trailing space at the end of the previous soft-wrapped lines so that
		// behaviour when navigating is consistent and so that we don't need to
		// continually add edges to handle the last line of the wrapped input.
		spaces++
		lines[row+1] = append(lines[row+1], repeatSpaces(spaces)...)
	} else {
		lines[row] = append(lines[row], word...)
		spaces++
		lines[row] = append(lines[row], repeatSpaces(spaces)...)
	}

	return lines
}

func repeatSpaces(n int) []rune {
	return []rune(strings.Repeat(string(' '), n))
}
