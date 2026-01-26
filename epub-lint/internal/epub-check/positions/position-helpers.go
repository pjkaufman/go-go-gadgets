package positions

import (
	"strings"
	"unicode/utf8"
)

func GetColumnFromIndex(contents string, line, index int) int {
	lines := strings.Split(contents, "\n")
	if line < 1 || line > len(lines) {
		return -1
	}

	// Calculate start-of-line byte offset
	byteOffset := 0
	for i := 0; i < line-1; i++ {
		byteOffset += len(lines[i]) + 1 // +1 for '\n'
	}

	// Clamp index to end of contents
	if index > len(contents) {
		index = len(contents)
	}

	// Offset in the line
	lineStart := byteOffset
	lineEnd := lineStart + len(lines[line-1])
	if index < lineStart {
		return 0
	}
	if index > lineEnd {
		index = lineEnd
	}

	lineBytes := contents[lineStart:index]
	col := 0
	for len(lineBytes) > 0 {
		_, size := utf8.DecodeRuneInString(lineBytes)
		lineBytes = lineBytes[size:]
		col++
	}
	return col + 1 // Columns are 1-based
}

func IndexToPosition(contents string, index int) Position {
	if index < 0 {
		return Position{Line: 1, Column: 1}
	}

	// Clamp to end of contents
	if index > len(contents) {
		index = len(contents)
	}

	// Count lines up to (but NOT including) index
	line := strings.Count(contents[:index], "\n") + 1

	// Find start of this line
	lineStart := strings.LastIndex(contents[:index], "\n")
	if lineStart == -1 {
		lineStart = 0
	} else {
		lineStart++ // move past '\n'
	}

	// Column = number of runes BEFORE the cursor
	col := utf8.RuneCountInString(contents[lineStart:index]) + 1

	return Position{Line: line, Column: col}
}

func GetColumnForLine(line string, index int) int {
	if index > len(line) {
		index = len(line)
	}

	var col, byteOffset int
	for byteOffset < index {
		_, size := utf8.DecodeRuneInString(line[byteOffset:])
		if size == 0 {
			break
		}

		byteOffset += size
		col++
	}

	return col + 1 // Columns are 1-based, like your other function
}

func GetPositionOffset(contents string, line, column int) int {
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
