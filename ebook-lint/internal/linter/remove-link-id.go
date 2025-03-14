package linter

import (
	"strings"
)

func RemoveLinkId(fileContents string, lineToUpdate, startOfFragment int) string {
	var lines = strings.Split(fileContents, "\n")
	if len(lines) < lineToUpdate {
		return fileContents
	}

	var indicatedLine = lines[lineToUpdate]
	if len(indicatedLine) < startOfFragment {
		return fileContents
	}

	// Work backwards from startOfFragment to find href or src attribute
	var (
		hrefAttributeIndicator = "href="
		srcAttributeIndicator  = "src="
		fragmentStart          = strings.LastIndex(indicatedLine[:startOfFragment], hrefAttributeIndicator)
		srcStart               = strings.LastIndex(indicatedLine[:startOfFragment], srcAttributeIndicator)
		startAttr              int
	)
	if fragmentStart == -1 && srcStart == -1 {
		return fileContents
	} else if fragmentStart > srcStart {
		startAttr = fragmentStart + len(hrefAttributeIndicator)
	} else {
		startAttr = srcStart + len(srcAttributeIndicator)
	}

	var endingQuote = string(indicatedLine[startAttr])
	endOfFragment := startAttr + 1 + strings.Index(indicatedLine[startAttr+1:], endingQuote)
	if endOfFragment == -1 {
		return fileContents
	}

	fragment := indicatedLine[startAttr+1 : endOfFragment]
	idIndicatorStart := strings.Index(fragment, "#")
	if idIndicatorStart == -1 {
		return fileContents
	}

	lines[lineToUpdate] = indicatedLine[:startAttr+1] + fragment[:idIndicatorStart] + indicatedLine[endOfFragment:]

	return strings.Join(lines, "\n")
}
