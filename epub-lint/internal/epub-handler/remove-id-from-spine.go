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

	startIndex, _, spineContent, err := getSpineContents(opfContents)
	if err != nil {
		return edit, err
	}

	lines := strings.Split(spineContent, "\n")
	idRef := fmt.Sprintf(`idref="%s"`, fileId)

	// Global offset where <spine ...> content begins
	spineGlobalOffset := startIndex + len(spineStartTag)

	for i, line := range lines {
		if !strings.Contains(line, idRef) {
			continue
		}

		lineSubset := line
		localOffset := 0

		for {
			startOfItemref := strings.Index(lineSubset, itemRefStartTag)
			if startOfItemref == -1 {
				break
			}

			endOfItemref := strings.Index(lineSubset, "/>")
			if endOfItemref == -1 {
				return edit, fmt.Errorf("failed to parse itemref out of line contents %q due to missing %q", lineSubset, "/>")
			}

			itemStart := localOffset + startOfItemref
			itemEnd := localOffset + endOfItemref + 2
			itemrefEl := line[itemStart:itemEnd]

			if strings.Contains(itemrefEl, idRef) {
				// Compute global delete range
				globalStart := spineGlobalOffset + computeLineOffset(lines, i) + itemStart
				globalEnd := spineGlobalOffset + computeLineOffset(lines, i) + itemEnd

				edit = positions.TextEdit{
					Range: positions.Range{
						Start: positions.IndexToPosition(opfContents, globalStart),
						End:   positions.IndexToPosition(opfContents, globalEnd),
					},
				}

				return edit, nil
			}

			// Continue scanning the rest of the line
			localOffset += endOfItemref + 2
			lineSubset = lineSubset[endOfItemref+2:]
		}
	}

	return edit, nil
}

func getSpineContents(opfContents string) (int, int, string, error) {
	startIndex := strings.Index(opfContents, spineStartTag)
	endIndex := strings.Index(opfContents, spineEndTag)

	if startIndex == -1 || endIndex == -1 {
		return 0, 0, "", ErrNoSpine
	}

	return startIndex, endIndex, opfContents[startIndex+len(spineStartTag) : endIndex], nil
}
