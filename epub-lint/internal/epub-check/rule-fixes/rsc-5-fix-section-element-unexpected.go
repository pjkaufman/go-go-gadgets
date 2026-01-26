package rulefixes

import (
	"strings"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
)

func FixSectionElementUnexpected(line, column int, contents string) (edits []positions.TextEdit) {
	offset := positions.GetPositionOffset(contents, line, column) // gets the index that actually represents the line and column in the current file
	if offset == -1 {
		return
	}

	openSection := "<section"
	openIdx := strings.LastIndex(contents[:offset], openSection)
	if openIdx == -1 {
		return
	}

	const endSection = "</section>"
	endIdx := strings.Index(contents[offset:], endSection)
	if endIdx == -1 {
		return
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
		return
	}

	// remove opening and closing els
	var (
		endStart      = offset + endIdx
		endEnd        = offset + endIdx + len(endSection)
		lineEndOffset = offset + lineEnd
	)
	edits = append(edits, positions.TextEdit{
		Range: positions.Range{
			Start: positions.IndexToPosition(contents, openIdx),
			End:   positions.IndexToPosition(contents, offset),
		},
	},
		positions.TextEdit{
			Range: positions.Range{
				Start: positions.IndexToPosition(contents, endStart),
				End:   positions.IndexToPosition(contents, endEnd),
			},
		})

	var (
		endLineContent               = contents[endEnd:lineEndOffset]
		lineContents                 = contents[lineStart:lineEndOffset]
		insertStartPos, insertEndPos positions.Position
	)
	if indexToMoveTo == 0 {
		insertStartPos = positions.IndexToPosition(contents, lineStart)
		insertEndPos = positions.IndexToPosition(contents, lineEndOffset)
	} else if strings.TrimSpace(currentLine[:indexToMoveTo]) == "" {
		insertStartPos = positions.IndexToPosition(contents, lineStart+indexToMoveTo)
		insertEndPos = positions.IndexToPosition(contents, lineEndOffset)
	} else {
		var endIndexToMoveTo = strings.Index(lineContents, endLineContent)
		for _, tagName := range movedBeforeElements {
			var (
				endTag      = "</" + tagName + ">"
				endTagIndex = strings.Index(endLineContent, endTag)
			)
			if endTagIndex == -1 { // something is wrong, so we are skipping this one...
				return
			}

			endIndexToMoveTo += endTagIndex + len(endTag)
			endLineContent = endLineContent[endTagIndex+len(endTag):]
		}

		insertStartPos = positions.IndexToPosition(contents, lineStart+indexToMoveTo)
		insertEndPos = positions.IndexToPosition(contents, lineStart+endIndexToMoveTo)
	}

	edits = append(edits, positions.TextEdit{
		Range: positions.Range{
			Start: insertStartPos,
			End:   insertStartPos,
		},
		NewText: openingSectionEl,
	},
		positions.TextEdit{
			Range: positions.Range{
				Start: insertEndPos,
				End:   insertEndPos,
			},
			NewText: endSection,
		})

	return
}
