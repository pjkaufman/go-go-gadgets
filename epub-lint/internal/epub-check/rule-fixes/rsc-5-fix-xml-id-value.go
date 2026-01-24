package rulefixes

import (
	"strings"
	"unicode"
)

func isValidStart(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

// isValidChar checks if the given character is valid in the ID.
func isValidChar(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' || ch == '-'
}

// convertToValidID converts an invalid OPF/NCX ID to a valid ID.
func convertToValidID(id string) string {
	id = strings.TrimSpace(id)
	id = strings.ReplaceAll(id, " ", "_")

	if len(id) == 0 {
		return id
	}

	var firstChar = rune(id[0])
	// Ensure the ID starts with a valid character
	if !isValidChar(firstChar) {
		id = "_" + id[1:]
	} else if !isValidStart(firstChar) {
		id = "_" + id
	}

	// Replace invalid characters with underscores
	var builder strings.Builder
	for _, ch := range id {
		if isValidChar(ch) {
			builder.WriteRune(ch)
		} else {
			builder.WriteRune('_')
		}
	}

	return builder.String()
}

func FixXmlIdValue(original string, lineNumber int, attribute string) (edit TextEdit) {
	lines := strings.Split(original, "\n")
	if lineNumber <= 0 || lineNumber > len(lines) {
		return
	}

	line := lines[lineNumber-1] // lineNumber is 1-based
	attributePrefix := attribute + `="`

	startIndex := strings.Index(line, attributePrefix)
	if startIndex != -1 {
		startIndex += len(attributePrefix)
		endIndex := startIndex
		for endIndex < len(line) && line[endIndex] != '"' {
			endIndex++
		}

		if endIndex < len(line) {
			invalidID := line[startIndex:endIndex]
			validID := convertToValidID(invalidID)
			edit.Range.Start.Line = lineNumber
			edit.Range.End.Line = lineNumber
			edit.Range.Start.Column = getColumnForLine(line, startIndex)
			edit.Range.End.Column = getColumnForLine(line, endIndex)

			edit.NewText = validID
			// line = line[:startIndex] + validID + line[endIndex:]

			// lines[lineNumber-1] = line
		}
	}

	return
}
