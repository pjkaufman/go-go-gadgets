package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
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
