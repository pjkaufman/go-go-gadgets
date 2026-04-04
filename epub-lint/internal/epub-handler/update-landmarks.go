package epubhandler

import (
	"fmt"
	"strings"
)

func UpdateLandmarks(contents, relativeFilePath, relativeCoverPath, relativeTocPath string) string {
	landmarksIndicator := strings.Index(contents, `epub:type="landmarks"`)
	if landmarksIndicator == -1 {
		return contents
	}

	startOfEl := strings.LastIndex(contents[:landmarksIndicator], "<")
	if startOfEl == -1 {
		return contents
	}

	var (
		startingElContent = contents[startOfEl+1 : landmarksIndicator]
		endOfElName       = strings.Index(startingElContent, " ")
	)
	if endOfElName == -1 { // this should not happen, but just in case, we will not make any changes
		return contents
	}

	endOfLandMarkEl := strings.Index(contents[landmarksIndicator:], fmt.Sprintf("</%s>", startingElContent[:endOfElName]))
	if endOfLandMarkEl == -1 { // another scenario that should not happen
		return contents
	}

	var (
		remainingLandmarkContents = contents[landmarksIndicator : landmarksIndicator+endOfLandMarkEl]
		nextRelativePathIndex     int
		currentActualIndex        = landmarksIndicator
		pathToFindHref            = fmt.Sprintf(`href="%s"`, relativeFilePath)
	)
	for nextRelativePathIndex != -1 {
		nextRelativePathIndex = strings.Index(remainingLandmarkContents, pathToFindHref)
		if nextRelativePathIndex == -1 {
			break
		}

		var (
			absoluteIndex = currentActualIndex + nextRelativePathIndex
			tagStart      = strings.LastIndex(contents[:absoluteIndex], "<")
			tagEnd        = strings.Index(contents[absoluteIndex+len(pathToFindHref):], ">")
		)
		if tagStart == -1 || tagEnd == -1 {
			currentActualIndex += nextRelativePathIndex + len(pathToFindHref)
			remainingLandmarkContents = remainingLandmarkContents[nextRelativePathIndex+len(pathToFindHref):]
			continue
		}

		tagEnd += absoluteIndex + len(pathToFindHref)

		var (
			tag                 = contents[tagStart : tagEnd+1]
			epubType, _, _, err = ExtractAttribute(tag, "epub:type")
		)
		if err != nil {
			currentActualIndex += nextRelativePathIndex + len(pathToFindHref)
			remainingLandmarkContents = remainingLandmarkContents[nextRelativePathIndex+len(pathToFindHref):]
			continue
		}

		var replacement string
		switch epubType {
		case "cover":
			if relativeCoverPath == "" {
				currentActualIndex += nextRelativePathIndex + len(pathToFindHref)
				remainingLandmarkContents = remainingLandmarkContents[nextRelativePathIndex+len(pathToFindHref):]
				continue
			}

			replacement = fmt.Sprintf(`href="%s"`, relativeCoverPath)
		case "toc":
			if relativeTocPath == "" {
				currentActualIndex += nextRelativePathIndex + len(pathToFindHref)
				remainingLandmarkContents = remainingLandmarkContents[nextRelativePathIndex+len(pathToFindHref):]
				continue
			}

			replacement = fmt.Sprintf(`href="%s"`, relativeTocPath)
		default:
			currentActualIndex += nextRelativePathIndex + len(pathToFindHref)
			remainingLandmarkContents = remainingLandmarkContents[nextRelativePathIndex+len(pathToFindHref):]
			continue
		}

		contents = contents[:absoluteIndex] + replacement + contents[absoluteIndex+len(pathToFindHref):]

		delta := len(replacement) - len(pathToFindHref)
		currentActualIndex += nextRelativePathIndex + len(replacement)
		remainingLandmarkContents = contents[currentActualIndex : landmarksIndicator+endOfLandMarkEl+delta]
	}

	return contents
}
