package rulefixes

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// FixFailedBlockquoteParsing takes a line and column (both 1-based), and
// file contents, and returns updated contents
// It inserts </p> before the ending blockquote element at the given position,
// and <p> after the opening blockquote element corresponding to that blockquote.
func FixFailedBlockquoteParsing(line, column int, contents string) string {
	if line < 1 {
		return contents
	}

	offset := getColumnOffset(contents, line, column)
	if offset == -1 {
		return contents
	}

	closeTag := "</blockquote>"
	openTag := "<blockquote"
	pOpen := "<p>"
	pClose := "</p>"

	if offset < len(closeTag) {
		return contents
	}

	closeIdx := offset - len(closeTag)

	// Scan backwards for previous non-whitespace chars before </blockquote>
	beforeClose := contents[:closeIdx]
	trimmed := trimRightSpace(beforeClose)
	if !strings.HasSuffix(trimmed, "</span>") &&
		!strings.HasSuffix(trimmed, "/>") &&
		trimmed[len(trimmed)-1] == '>' {
		return contents
	}

	// Find the opening <blockquote before closeIdx
	openIdx := strings.LastIndex(contents[:closeIdx], openTag)
	if openIdx == -1 {
		return contents
	}
	openEnd := strings.Index(contents[openIdx:], ">")
	if openEnd == -1 || openIdx+openEnd > closeIdx {
		return contents
	}
	openEndIdx := openIdx + openEnd

	// Insert <p> after opening tag and </p> before closing tag
	var sb strings.Builder
	sb.WriteString(contents[:openEndIdx+1]) // up to and including opening tag
	sb.WriteString(pOpen)
	sb.WriteString(contents[openEndIdx+1 : closeIdx]) // content between
	sb.WriteString(pClose)
	sb.WriteString(contents[closeIdx:]) // from closing tag onward

	return sb.String()
}

func trimRightSpace(s string) string {
	i := len(s)
	for i > 0 && unicode.IsSpace(rune(s[i-1])) {
		i--
	}

	return s[:i]
}

func getColumnOffset(contents string, line, column int) int {
	lines := strings.Split(contents, "\n")
	if line > len(lines) {
		return -1
	}

	byteOffset := 0
	for i := 0; i < line-1; i++ {
		byteOffset += len(lines[i]) + 1
	}

	curLine := lines[line-1]
	colByte := 0
	remainingRunes := column - 1
	for remainingRunes > 0 && colByte < len(curLine) {
		_, size := utf8.DecodeRuneInString(curLine[colByte:])
		if size == 0 {
			break
		}
		colByte += size
		remainingRunes--
	}
	if remainingRunes > 0 {
		colByte = len(curLine)
	}

	byteOffset += colByte
	if byteOffset > len(contents) {
		byteOffset = len(contents)
	}

	return byteOffset
}
