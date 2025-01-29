package linter

import (
	"fmt"
	"strings"
)

// RemoveEmptyOpfElement processes the specified element in the OPF contents
func RemoveEmptyOpfElements(elementName string, lineNum int, opfContents string) (string, bool, error) {
	lines := strings.Split(opfContents, "\n")
	if lineNum < 0 || lineNum >= len(lines) {
		return opfContents, false, fmt.Errorf("line number out of range")
	}

	// Locate the specified line
	line := lines[lineNum]
	elementStart := "<" + elementName
	selfClosingIndicator := "/>"
	elementEnd := "</"

	// Find the start of the element
	startIndex := strings.Index(line, elementStart)
	if startIndex == -1 {
		return opfContents, false, fmt.Errorf("element not found on the given line")
	}

	endIndex := strings.Index(line[startIndex:], selfClosingIndicator)
	if endIndex == -1 {
		initialEndIndex := strings.Index(line[startIndex:], elementEnd)
		if initialEndIndex == -1 {
			return opfContents, false, fmt.Errorf("end of element not found on the given line")
		}

		endIndex = initialEndIndex + strings.Index(line[startIndex+initialEndIndex:], ">") + 1
	} else {
		endIndex += len(selfClosingIndicator)
	}

	endIndex += startIndex

	// Strip out the content of the element
	line = line[:startIndex] + line[endIndex:]

	// Check if the remaining line content is whitespace
	if strings.TrimSpace(line) == "" {
		// Remove the line
		lines = append(lines[:lineNum], lines[lineNum+1:]...)
		return strings.Join(lines, "\n"), true, nil
	}

	// Update the line in the original content
	lines[lineNum] = line
	return strings.Join(lines, "\n"), false, nil
}
