package epubhandler

import (
	"fmt"
	"strings"
)

// Gets the start of the toc element contents up to and including the closing el
func GetNavTOCContentPositionInfo(navFileContents string) (int, int) {
	const tocEpubType = `epub:type="toc"`
	epubTypeIndex := strings.Index(navFileContents, tocEpubType)
	if epubTypeIndex == -1 {
		return -1, -1
	}

	var (
		startOfEl = strings.LastIndex(navFileContents[:epubTypeIndex], "<")
		endOfEl   = strings.Index(navFileContents[epubTypeIndex:], ">")
	)
	if startOfEl == -1 || endOfEl == -1 {
		return -1, -1
	}

	endOfEl += 1 + epubTypeIndex

	elNameEndIndex := strings.Index(navFileContents[startOfEl:], " ")
	if elNameEndIndex == -1 {
		return -1, -1
	}

	elNameEndIndex += startOfEl
	elName := navFileContents[startOfEl+1 : elNameEndIndex]

	endOfElIndex := strings.Index(navFileContents[endOfEl:], fmt.Sprintf("</%s>", elName))
	if endOfElIndex == -1 {
		return -1, -1
	}

	endOfElIndex += endOfEl

	return endOfEl, endOfElIndex
}
