package epubhandler

import (
	"fmt"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
)

const (
	hrefAttribute   = "href="
	itemStartTag    = "<item"
	itemRefStartTag = "<itemref"
)

func RemoveFileFromOpf(opfContents, fileName string) ([]positions.TextEdit, error) {
	var edits []positions.TextEdit

	startIndex, _, manifestContent, err := GetManifestContents(opfContents)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(manifestContent, "\n")
	var fileID string

	manifestGlobalOffset := startIndex + len(ManifestStartTag)

	for i, line := range lines {
		lineSubset := line
		localOffset := 0

		for {
			startOfItem := strings.Index(lineSubset, itemStartTag)
			if startOfItem == -1 {
				break
			}

			endOfItem := strings.Index(lineSubset, "/>")
			if endOfItem == -1 {
				return nil, fmt.Errorf("failed to parse item out of line contents %q due to missing %q", lineSubset, "/>")
			}

			itemStart := localOffset + startOfItem
			itemEnd := localOffset + endOfItem + 2
			itemEl := line[itemStart:itemEnd]

			// find href
			hrefIdx := strings.Index(itemEl, hrefAttribute)
			if hrefIdx == -1 {
				localOffset += endOfItem
				lineSubset = lineSubset[endOfItem:]
				continue
			}

			hrefIdx += len(hrefAttribute)
			quote := itemEl[hrefIdx : hrefIdx+1]
			hrefIdx++

			hrefEnd := strings.Index(itemEl[hrefIdx:], quote)
			if hrefEnd == -1 {
				localOffset += endOfItem
				lineSubset = lineSubset[endOfItem:]
				continue
			}

			hrefContent := itemEl[hrefIdx : hrefIdx+hrefEnd]
			if !strings.HasSuffix(hrefContent, fileName) {
				localOffset += endOfItem
				lineSubset = lineSubset[endOfItem:]
				continue
			}

			// false positive check
			trimmed := strings.TrimSuffix(hrefContent, fileName)
			var prev rune
			if len(trimmed) == 0 {
				prev = rune(quote[0])
			} else {
				prev = rune(trimmed[len(trimmed)-1])
			}
			if prev != '\'' && prev != '"' && prev != '\\' && prev != '/' {
				localOffset += endOfItem
				lineSubset = lineSubset[endOfItem:]
				continue
			}

			fileID = ExtractID(itemEl)

			// compute global delete range
			globalStart := manifestGlobalOffset + computeLineOffset(lines, i) + itemStart
			globalEnd := manifestGlobalOffset + computeLineOffset(lines, i) + itemEnd

			edits = append(edits, positions.TextEdit{
				Range: positions.Range{
					Start: positions.IndexToPosition(opfContents, globalStart),
					End:   positions.IndexToPosition(opfContents, globalEnd),
				},
			})

			break
		}
	}

	if fileID == "" {
		return edits, nil
	}

	// append spine edits
	spineEdit, err := RemoveIdFromSpine(opfContents, fileID)
	if err != nil {
		return nil, err
	}

	if !spineEdit.IsEmpty() {

		edits = append(edits, spineEdit)
	}

	return edits, nil
}

// computeLineOffset returns the byte offset of the start of line i within manifestContent
func computeLineOffset(lines []string, lineIndex int) int {
	offset := 0
	for j := 0; j < lineIndex; j++ {
		offset += len(lines[j]) + 1 // +1 for newline
	}
	return offset
}

func ExtractID(line string) string {
	const idAttr = `id="`
	startIndex := strings.Index(line, idAttr)
	if startIndex == -1 {
		return ""
	}

	startIndex += len(idAttr)
	endIndex := strings.Index(line[startIndex:], `"`)
	if endIndex == -1 {
		return ""
	}

	return line[startIndex : startIndex+endIndex]
}
