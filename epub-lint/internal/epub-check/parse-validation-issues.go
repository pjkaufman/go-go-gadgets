package epubcheck

import (
	"strconv"
	"strings"
)

// ParseEPUBCheckOutput parses the contents of an EPUBCheck output from a string.
func ParseEPUBCheckOutput(logContents string) (ValidationErrors, error) {
	var validationErrors ValidationErrors
	lines := strings.Split(logContents, "\n")

	for _, line := range lines {
		// Find the code (between first '(' and ')')
		start := strings.Index(line, "(")
		end := strings.Index(line, ")")
		if start == -1 || end == -1 || end < start {
			continue
		}
		code := line[start+1 : end]

		// Find the .epub/ marker for file path
		epubIdx := strings.Index(line, ".epub/")
		if epubIdx == -1 {
			continue
		}
		// Find the next '(' after .epub/ which marks the start of (line,column)
		pathStart := epubIdx + len(".epub/")
		pathEnd := strings.Index(line[pathStart:], "(")
		if pathEnd == -1 {
			continue
		}
		filePath := line[pathStart : pathStart+pathEnd]

		// Get line and column (between '(' and ')' after file path)
		locStart := pathStart + pathEnd + 1
		locEnd := strings.Index(line[locStart:], ")")
		if locEnd == -1 {
			continue
		}
		locStr := line[locStart : locStart+locEnd]
		locParts := strings.SplitN(locStr, ",", 2)
		lineNum, colNum := -1, -1
		if len(locParts) == 2 {
			lineNum, _ = strconv.Atoi(strings.TrimSpace(locParts[0]))
			colNum, _ = strconv.Atoi(strings.TrimSpace(locParts[1]))
		}

		// Message: after the final "):"
		afterLoc := line[locStart+locEnd:]
		colonIdx := strings.Index(afterLoc, ":")
		if colonIdx == -1 {
			continue
		}
		message := strings.TrimSpace(afterLoc[colonIdx+1:])

		var pos *Position
		if lineNum != -1 && colNum != -1 {
			pos = &Position{Line: lineNum, Column: colNum}
		}

		validationErrors.ValidationIssues = append(validationErrors.ValidationIssues, ValidationError{
			Code:     code,
			FilePath: filePath,
			Location: pos,
			Message:  message,
		})
	}

	return validationErrors, nil
}
