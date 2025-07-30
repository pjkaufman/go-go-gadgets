package ui

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	rw "github.com/mattn/go-runewidth"
	stringdiff "github.com/pjkaufman/go-go-gadgets/pkg/string-diff"
)

// icons
var (
	documentIcon   = string([]byte{0xF0, 0x9F, 0x97, 0x8E}) // UTF-8 encoding for "🗎"
	suggestionIcon = string([]byte{0xE2, 0x91, 0x82})       // UTF-8 encoding for "⑂"
	viewIcon       = string([]byte{0xF0, 0x9F, 0x91, 0x81}) // UTF-8 encoding for "👁"
	editIcon       = string([]byte{0xE2, 0x9C, 0x8E})       // UTF-8 encoding for "✎"
)

func fillLine(currentValue string, width int) string {
	var amountToFill = width - lipgloss.Width(currentValue)
	if amountToFill < 1 {
		return currentValue
	}

	return currentValue + strings.Repeat(" ", amountToFill)
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

// from https://github.com/charmbracelet/bubbles/blob/d66fddf5e780b2bf30e386dbf4e65b55b258197f/textarea/textarea.go#L1398-L1461
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
			if lipgloss.Width(string(lines[row]))+lipgloss.Width(string(word))+spaces > width {
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
			if ansi.StringWidth(string(word))+lastCharLen > width {
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

	if lipgloss.Width(string(lines[row]))+lipgloss.Width(string(word))+spaces >= width {
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

// from https://github.com/charmbracelet/bubbles/blob/d66fddf5e780b2bf30e386dbf4e65b55b258197f/textarea/textarea.go#L1463-L1465
func repeatSpaces(n int) []rune {
	return []rune(strings.Repeat(string(' '), n))
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
