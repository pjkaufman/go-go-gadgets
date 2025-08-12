package ui

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
	stringdiff "github.com/pjkaufman/go-go-gadgets/pkg/string-diff"
)

// icons
var (
	documentIcon = string([]byte{0xF0, 0x9F, 0x97, 0x8E}) // UTF-8 encoding for "🗎"
	sectionIcon  = string([]byte{0xC2, 0xA7})             // UTF-8 encoding for "§"
	viewIcon     = string([]byte{0xF0, 0x9F, 0x91, 0x81}) // UTF-8 encoding for "👁"
	editIcon     = string([]byte{0xE2, 0x9C, 0x8E})       // UTF-8 encoding for "✎"
)

func fillLine(currentValue string, width int) string {
	var amountToFill = width - lipgloss.Width(currentValue)
	if amountToFill < 1 {
		return currentValue
	}

	return currentValue + strings.Repeat(" ", amountToFill)
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
