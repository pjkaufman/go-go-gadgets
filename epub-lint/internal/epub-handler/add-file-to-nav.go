package epubhandler

import (
	"fmt"
	"strings"
)

// AddFileToNav adds the specified file at the end of the nav contents
func AddFileToNav(navContents, filePath, title string) string {
	endOfEl, endOfElIndex := GetNavTOCContentPositionInfo(navContents)
	if endOfEl == -1 || endOfElIndex == -1 {
		return navContents
	}

	// for now we shall assume that all nav TOCs are made up of an ol element.
	const endingOlEl = "</ol>"
	insertIndex := strings.LastIndex(navContents[:endOfElIndex], endingOlEl)
	if insertIndex == -1 {
		return navContents
	}

	return navContents[:insertIndex] + fmt.Sprintf(`<li><a href=%q>%s</a></li>`, filePath, title) + "\n" + navContents[insertIndex:]
}
