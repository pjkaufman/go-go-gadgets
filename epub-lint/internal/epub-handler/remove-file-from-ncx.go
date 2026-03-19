package epubhandler

import (
	"fmt"
	"strings"
)

func RemoveFileFromNcx(contents, relativeFilePath string) string {
	var (
		remainingNcx           = contents
		nextRelativePathIndex  int
		currentActualIndex     int
		pathToFindSrc          = fmt.Sprintf(`src="%s"`, relativeFilePath) // we may want to account for references to ids in the file as well, but for now this should work
		navPointStartIndicator = "<navPoint"
		navPointEndIndicator   = "</navPoint>"
	)
	for nextRelativePathIndex != -1 {
		nextRelativePathIndex = strings.Index(remainingNcx, pathToFindSrc)
		if nextRelativePathIndex == -1 {
			break
		}

		var (
			absoluteIndex = currentActualIndex + nextRelativePathIndex
			navPointStart = strings.LastIndex(contents[:absoluteIndex], navPointStartIndicator)
			navPointEnd   = strings.Index(contents[absoluteIndex+len(pathToFindSrc):], navPointEndIndicator)
		)
		if navPointStart == -1 || navPointEnd == -1 {
			currentActualIndex += nextRelativePathIndex + len(pathToFindSrc)
			remainingNcx = remainingNcx[nextRelativePathIndex+len(pathToFindSrc):]
			continue
		}

		navPointEnd += absoluteIndex + len(pathToFindSrc) + len(navPointEndIndicator)

		var startOfLine = strings.LastIndex(contents[:navPointStart], "\n")
		if startOfLine == -1 {
			startOfLine = 0
		}

		var removeLine = strings.TrimSpace(contents[startOfLine:navPointStart]) == ""

		var endOfLine = strings.Index(contents[navPointEnd:], "\n")
		if endOfLine == -1 {
			endOfLine = len(contents)
		} else {
			endOfLine += navPointEnd
		}

		removeLine = removeLine && strings.TrimSpace(contents[navPointEnd:endOfLine]) == ""
		if removeLine {
			navPointStart = startOfLine
			navPointEnd = endOfLine
		}

		contents = contents[:navPointStart] + contents[navPointEnd:]
		remainingNcx = contents[navPointStart:]
	}

	return contents
}
