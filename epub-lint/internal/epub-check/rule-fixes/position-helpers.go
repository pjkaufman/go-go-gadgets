package rulefixes

import (
	"strings"
	"unicode/utf8"
)

func getColumnFromIndex(contents string, line, index int) int {
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

func indexToPosition(contents string, index int) Position {
	if index < 0 {
		return Position{Line: 1, Column: 1}
	}

	// Early clamp to avoid unnecessary work
	if index >= len(contents) {
		// Place at the very end (after last character of last line)
		lines := strings.Split(contents, "\n")
		if len(lines) == 0 {
			return Position{Line: 1, Column: 1}
		}

		lastLine := lines[len(lines)-1]

		return Position{
			Line:   len(lines),
			Column: utf8.RuneCountInString(lastLine) + 1,
		}
	}

	var (
		upToIndex = contents[:index+1]
		line      = strings.Count(upToIndex, "\n") + 1
		col       int
	)
	if line == 1 {
		col = utf8.RuneCountInString(upToIndex) + 1
	} else {
		col = utf8.RuneCountInString(upToIndex[strings.LastIndex(upToIndex, "\n")+1:])
	}

	return Position{Line: line, Column: col}
}

func getColumnForLine(line string, index int) int {
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
