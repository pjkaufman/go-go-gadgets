package epubhandler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
)

const (
	spineStartTag = "<spine"
	spineEndTag   = "</spine>"
)

var ErrNoSpine = errors.New("spine tag not found in OPF contents")

func RemoveIdFromSpine(opfContents, fileId string) (positions.TextEdit, error) {
	var edit positions.TextEdit

	startIndex, spineContent, err := getSpineContents(opfContents)
	if err != nil {
		return edit, err
	}

	idRef := fmt.Sprintf(`idref=%q`, fileId)
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

	var (
		startOfContentToRemove, endOfContentToRemove int
		startOfEl                                    = startItemRefIndex
		endOfEl                                      = idRefIndex + endItemRefIndex
	)
	startOfContentToRemove, endOfContentToRemove = GetLineBoundsIfEmpty(spineContent, startOfEl, endOfEl)
	if strings.Contains(spineContent[startOfContentToRemove+1:endOfContentToRemove], "\n") { // adding one gets past the initial newline character
		return edit, fmt.Errorf("failed to parse itemref line out of spine content for id %q due to the content having a newline in it somewhere", fileId)
	}

	var startPos, endPos positions.Position
	if startOfContentToRemove != startOfEl || endOfContentToRemove != endOfEl {
		// no other line content other than whitespace is on the line...
		startPos = positions.IndexToPosition(opfContents, startIndex+startOfContentToRemove)
		endPos = positions.IndexToPosition(opfContents, startIndex+endOfContentToRemove)
	} else {
		startPos = positions.IndexToPosition(opfContents, startIndex+startOfEl)
		endPos = positions.IndexToPosition(opfContents, startIndex+endOfEl)
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
