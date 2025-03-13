package linter

import (
	"fmt"
	"strings"
)

func RemoveLinkId(fileContents string, lineToUpdate, startOfFragment int) string {
	var lines = strings.Split(fileContents, "\n")
	if len(lines) <= lineToUpdate {
		fmt.Println("Bail on lines")
		return fileContents
	}

	var indicatedLine = lines[lineToUpdate]
	if len(indicatedLine) <= startOfFragment {
		fmt.Println("Bail on chars")
		return fileContents
	}

	// get the fragment
	// TODO: make it tolerate single quotes as well
	var (
		endOfFragment    = startOfFragment + 1 + strings.Index(indicatedLine[startOfFragment+1:], "\"")
		fragment         = indicatedLine[startOfFragment:endOfFragment]
		idIndicatorStart = strings.Index(fragment, "#")
	)
	fmt.Println(endOfFragment, fragment, idIndicatorStart)
	if idIndicatorStart == -1 {
		return fileContents
	}

	lines[lineToUpdate] = indicatedLine[:startOfFragment] + fragment[:idIndicatorStart] + indicatedLine[endOfFragment:]

	return strings.Join(lines, "\n")
}
