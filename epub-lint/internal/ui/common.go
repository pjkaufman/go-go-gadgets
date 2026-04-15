package ui

import (
	"fmt"
	"strings"
	"unicode"

	"charm.land/lipgloss/v2"
	stringdiff "github.com/pjkaufman/go-go-gadgets/pkg/string-diff"
)

// icons
var (
	documentIcon = string([]byte{0xF0, 0x9F, 0x97, 0x8E}) // UTF-8 encoding for "🗎"
	sectionIcon  = string([]byte{0xC2, 0xA7})             // UTF-8 encoding for "§"
	viewIcon     = string([]byte{0xF0, 0x9F, 0x91, 0x81}) // UTF-8 encoding for "👁"
	editIcon     = string([]byte{0xE2, 0x9C, 0x8E})       // UTF-8 encoding for "✎"
	warningIcon  = string([]byte{0xE2, 0x9A, 0xA0})       // UTF-8 encoding for "⚠"
)

func unreachableCode() {
	fmt.Println("Unreachable")
}

func fillLine(currentValue string, width int) string {
	var amountToFill = width - lipgloss.Width(currentValue)
	if amountToFill < 1 {
		return currentValue
	}

	return currentValue + strings.Repeat(" ", amountToFill)
}

func getStringDiff(original, updated string) (string, error) {
	return stringdiff.GetPrettyDiffString(strings.TrimLeft(original, "\n"), strings.TrimLeft(updated, "\n"))
}

// textarea gets rid of tabs when creating changes, so in order to preserve tabs in the starting whitespace of a line
// we will use the value of original as the template for what whitespace is needed for each line present
func alignWhitespace(original, updated string) string {
	var (
		origLines = strings.Split(original, "\n")
		newLines  = strings.Split(updated, "\n")
		minLines  = min(len(newLines), len(origLines))
	)

	for i := range minLines {
		origPrefix := strings.Builder{}
		for j := range len(origLines[i]) {
			if !unicode.IsSpace(rune(origLines[i][j])) {
				break
			}

			origPrefix.WriteByte(origLines[i][j])
		}

		newLines[i] = origPrefix.String() + strings.TrimLeft(newLines[i], " \t")
	}

	return strings.Join(newLines, "\n")
}
