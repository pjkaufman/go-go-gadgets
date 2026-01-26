package epubhandler

import (
	"fmt"
	"slices"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
)

const (
	hrefAttribute   = "href="
	itemStartTag    = "<item"
	itemRefStartTag = "<itemref"
)

func RemoveFileFromOpf(opfContents, fileName string) (string, error) {
	startIndex, endIndex, manifestContent, err := GetManifestContents(opfContents)
	if err != nil {
		return "", err
	}

	lines := strings.Split(manifestContent, "\n")
	var fileID string

	for i, line := range lines {
		// find each item in the line and check if the href is the one you are looking for
		// if it is, remove it from that line. If the line is now blank, remove it.
		var (
			lineSubset  = line
			startOfItem int
		)
		for startOfItem != -1 {
			startOfItem = strings.Index(lineSubset, itemStartTag)
			if startOfItem == -1 {
				break
			}

			// for now we will assume the items are self-closing
			var endOfItem = strings.Index(lineSubset, "/>")
			if endOfItem == -1 {
				return "", fmt.Errorf("failed to parse item out of line contents %q due to missing %q", lineSubset, "/>")
			}

			var (
				itemEndIndex = endOfItem + 2
				itemEl       = lineSubset[startOfItem:itemEndIndex]
			)
			// to make sure we are only operating on the href
			var startOfHref = strings.Index(itemEl, hrefAttribute)
			if startOfHref == -1 {
				lineSubset = lineSubset[itemEndIndex:]
				continue
			}

			startOfHref += len(hrefAttribute)
			var hrefQuote = itemEl[startOfHref : startOfHref+1]
			startOfHref++

			var endOfHrefIndex = strings.Index(itemEl[startOfHref:], hrefQuote)
			if endOfHrefIndex == -1 {
				lineSubset = lineSubset[itemEndIndex:]
				continue // something went wrong, so we have to ignore the el...
			}

			var hrefContent = itemEl[startOfHref : startOfHref+endOfHrefIndex]
			if !strings.HasSuffix(hrefContent, fileName) {
				lineSubset = lineSubset[endOfItem:] // start over with any other items on this line
				continue
			}

			hrefContent = strings.TrimSuffix(hrefContent, fileName)

			// check for a false positive by checking that previous char is not a slash or a quote
			var previousChar rune
			if len(hrefContent) == 0 {
				previousChar = rune(hrefQuote[0])
			} else {
				previousChar = rune(hrefContent[len(hrefContent)-1])
			}
			if previousChar != '\'' && previousChar != '"' && previousChar != '\\' && previousChar != '/' {
				lineSubset = lineSubset[endOfItem:]
				continue
			}

			fileID = ExtractID(itemEl)
			line = strings.Replace(line, itemEl, "", 1)

			if strings.TrimSpace(line) == "" {
				lines = slices.Delete(lines, i, i+1)
			} else {
				lines[i] = line
			}

			break
		}
	}

	updatedManifestContent := strings.Join(lines, "\n")
	updatedOpfContents := opfContents[:startIndex] + updatedManifestContent + opfContents[endIndex:]

	if fileID == "" {
		return updatedOpfContents, nil
	}

	edit, err := RemoveIdFromSpine(updatedOpfContents, fileID)
	if err != nil {
		return "", err
	}

	if edit.IsEmpty() {
		return updatedOpfContents, nil
	}

	return positions.ApplyEdits("", updatedOpfContents, []positions.TextEdit{edit})
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
