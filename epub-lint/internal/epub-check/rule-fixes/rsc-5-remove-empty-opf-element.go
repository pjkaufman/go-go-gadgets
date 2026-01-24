package rulefixes

import (
	"fmt"
	"strings"
)

// RemoveEmptyOpfElement processes the specified element in the OPF contents
func RemoveEmptyOpfElements(elementName string, lineNum int, opfContents string) (TextEdit, bool, error) {
	var edit TextEdit
	lineNum--
	lines := strings.Split(opfContents, "\n")
	if lineNum < 0 || lineNum >= len(lines) {
		return edit, false, fmt.Errorf("line number out of range")
	}

	// Locate the specified line
	line := lines[lineNum]
	elementStart := "<" + elementName
	selfClosingIndicator := "/>"
	elementEnd := "</"

	// Find the start of the element
	startIndex := strings.Index(line, elementStart)
	if startIndex == -1 {
		return edit, false, fmt.Errorf("element not found on the given line")
	}

	endIndex := strings.Index(line[startIndex:], selfClosingIndicator)
	if endIndex == -1 {
		initialEndIndex := strings.Index(line[startIndex:], elementEnd)
		if initialEndIndex == -1 {
			return edit, false, fmt.Errorf("end of element not found on the given line")
		}

		endIndex = initialEndIndex + strings.Index(line[startIndex+initialEndIndex:], ">") + 1
	} else {
		endIndex += len(selfClosingIndicator)
	}

	endIndex += startIndex

	// Strip out the content of the element
	updatedLine := line[:startIndex] + line[endIndex:]

	edit.Range.Start.Line = lineNum + 1
	edit.Range.End.Line = lineNum + 1

	// Check if the remaining line content is whitespace
	if strings.TrimSpace(updatedLine) == "" {
		// Remove the line
		edit.Range.End.Line += 1

		return edit, true, nil
	}

	edit.Range.Start.Column = getColumnForLine(line, startIndex)
	edit.Range.End.Column = getColumnForLine(line, endIndex)

	return edit, false, nil
}
