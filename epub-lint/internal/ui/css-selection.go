package ui

import (
	"fmt"
	"strings"
)

func (m FixableIssuesModel) cssSelectionView() string {
	var s strings.Builder
	s.WriteString("\nSelect the CSS file to modify:\n\n")
	for i, cssFile := range m.CssSelectionInfo.cssFiles {
		cursor := " "
		if m.CssSelectionInfo.currentCssIndex == i {
			cursor = ">"
		}

		s.WriteString(fmt.Sprintf("%s %d. %s\n", cursor, i+1, cssFile))
	}

	s.WriteString("\n")

	return s.String()
}
