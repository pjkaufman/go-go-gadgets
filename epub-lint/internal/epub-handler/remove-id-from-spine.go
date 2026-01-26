package epubhandler

import (
	"fmt"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
)

const (
	spineStartTag = "<spine"
	spineEndTag   = "</spine>"
)

var ErrNoSpine = fmt.Errorf("spine tag not found in OPF contents")

func RemoveIdFromSpine(opfContents, fileId string) (positions.TextEdit, error) {
	var edit positions.TextEdit

	startIndex, spineContent, err := getSpineContents(opfContents)
	if err != nil {
		return edit, err
	}

	idRef := fmt.Sprintf(`idref="%s"`, fileId)
	idRefIndex := strings.Index(spineContent, idRef)
	if idRefIndex == -1 {
		return edit, nil
	}

	startItemRefIndex := strings.LastIndex(spineContent[:idRefIndex], itemRefStartTag)
	if startItemRefIndex == -1 {
		return edit, fmt.Errorf("failed to parse itemref out of spine content for id %q due to missing %q", fileId, itemRefStartTag)
	}

	endItemRefIndex := strings.Index(spineContent[idRefIndex:], "/>")
	if endItemRefIndex == -1 {
		return edit, fmt.Errorf("failed to parse itemref out of spine content for id %q due to missing %q", fileId, "/>")
	}

	endItemRefIndex += 2

	startOfLineIndex := strings.LastIndex(spineContent[:idRefIndex], "\n")
	if startOfLineIndex == -1 {
		startOfLineIndex = 0
	}

	if startItemRefIndex < startOfLineIndex {
		return edit, fmt.Errorf("failed to parse itemref line out of spine content for id %q due to the start of itemref being on a different line from the itemref's href", fileId)
	}

	endOfLineIndex := strings.LastIndex(spineContent[idRefIndex:], "\n")
	if endOfLineIndex == -1 {
		endOfLineIndex = len(spineContent) - 1
	}

	if endItemRefIndex > endOfLineIndex {
		return edit, fmt.Errorf("failed to parse itemref line out of spine content for id %q due to the end of itemref being on a different line from the itemref's href", fileId)
	}

	var startPos, endPos positions.Position
	if strings.TrimSpace(spineContent[startOfLineIndex:startItemRefIndex]+spineContent[idRefIndex+endItemRefIndex:idRefIndex+endOfLineIndex]) == "" {
		// no other line content other than whitespace is on the line...
		startPos = positions.IndexToPosition(opfContents, startIndex+startOfLineIndex)
		endPos = positions.IndexToPosition(opfContents, startIndex+idRefIndex+endOfLineIndex)
	} else {
		startPos = positions.IndexToPosition(opfContents, startIndex+startItemRefIndex)
		endPos = positions.IndexToPosition(opfContents, startIndex+idRefIndex+endItemRefIndex)
	}

	edit.Range.Start = startPos
	edit.Range.End = endPos

	return edit, nil
}

func getSpineContents(opfContents string) (int, string, error) {
	startIndex := strings.Index(opfContents, spineStartTag)
	endIndex := strings.Index(opfContents, spineEndTag)

	if startIndex == -1 || endIndex == -1 {
		return 0, "", ErrNoSpine
	}

	return startIndex + len(spineStartTag) + 1, opfContents[startIndex+len(spineStartTag)+1 : endIndex], nil
}
