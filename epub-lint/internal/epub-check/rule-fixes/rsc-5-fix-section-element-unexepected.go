package rulefixes

import (
	"strings"
)

func FixSectionElementUnexpected(line, column int, contents string) string {
	offset := GetPositionOffset(contents, line, column) // gets the index that actually represents the line and column in the current file
	if offset == -1 {
		return contents
	}

	openSection := "<section"
	openIdx := strings.LastIndex(contents[:offset], openSection)
	if openIdx == -1 {
		return contents
	}

	endSection := "</section>"
	endIdx := strings.Index(contents[offset:], endSection)
	if endIdx == -1 {
		return contents
	}

	openingSectionEl := contents[openIdx:offset]

	// for performance and simplicity we will assume this is on a single line and not multiple...
	var currentLine string
	lineStart := strings.LastIndex(contents[:offset], "\n")
	lineStart++ // we will not include the \n and if it is -1, we will be at 0

	lineEnd := strings.Index(contents[offset:], "\n")
	if lineEnd == -1 {
		if lineStart == 0 {
			currentLine = contents
		} else {
			currentLine = contents[lineStart:]
		}
	} else {
		currentLine = contents[lineStart : offset+lineEnd]
	}

	var (
		lineSubset                     = currentLine[:openIdx-lineStart]
		elementEndStart, indexToMoveTo int
		movedBeforeElements            []string
	)

	for elementEndStart != -1 {
		elementEndStart = strings.LastIndex(lineSubset, "<")
		if elementEndStart == -1 {
			break
		}

		// if it is an end tag, we skip it
		if lineSubset[elementEndStart+1] == '/' {
			lineSubset = lineSubset[:elementEndStart]
			continue
		}

		elementEndEnd := strings.Index(lineSubset[elementEndStart:], " ")
		if elementEndEnd == -1 {
			elementEndEnd = strings.Index(lineSubset[elementEndStart:], ">")
			if elementEndEnd == -1 {
				lineSubset = lineSubset[:elementEndStart]
				continue
			}
		}

		var tagName = lineSubset[elementEndStart+1 : elementEndStart+elementEndEnd]
		if tagName != "p" && tagName != "span" {
			break
		}

		indexToMoveTo = elementEndStart
		movedBeforeElements = append(movedBeforeElements, tagName)

		lineSubset = lineSubset[:elementEndStart]
	}

	if len(movedBeforeElements) == 0 {
		return contents
	}

	var (
		endLineContent = contents[offset+endIdx+len(endSection) : offset+lineEnd]
		updatedLine    = contents[lineStart:openIdx] + contents[offset:offset+endIdx] + endLineContent
	)
	if indexToMoveTo == 0 {
		updatedLine = openingSectionEl + updatedLine + endSection
	} else if strings.TrimSpace(currentLine[:indexToMoveTo]) == "" {
		updatedLine = updatedLine[:indexToMoveTo] + openingSectionEl + updatedLine[indexToMoveTo:] + endSection
	} else {
		var endIndexToMoveTo = strings.Index(updatedLine, endLineContent)
		for _, tagName := range movedBeforeElements {
			var (
				endTag      = "</" + tagName + ">"
				endTagIndex = strings.Index(endLineContent, endTag)
			)
			if endTagIndex == -1 { // something is wrong, so we are skipping this one...
				return contents
			}

			endIndexToMoveTo += endTagIndex + len(endTag)
			endLineContent = endLineContent[endTagIndex+len(endTag):]
		}

		// add end of section first to prevent accounting for that shift as well
		updatedLine = updatedLine[:endIndexToMoveTo] + endSection + updatedLine[endIndexToMoveTo:]
		updatedLine = updatedLine[:indexToMoveTo] + openingSectionEl + updatedLine[indexToMoveTo:]
	}

	return contents[:lineStart] + updatedLine + contents[offset+lineEnd:]
}
