package epubhandler

import (
	"fmt"
	"strings"
)

const (
	ManifestStartTag = "<manifest>"
	ManifestEndTag   = "</manifest>"
)

func GetManifestContents(opfContents string) (int, int, string, error) {
	startIndex := strings.Index(opfContents, ManifestStartTag)
	endIndex := strings.Index(opfContents, ManifestEndTag)

	if startIndex == -1 || endIndex == -1 {
		return 0, 0, "", ErrNoManifest
	}

	return startIndex + len(ManifestStartTag), endIndex, opfContents[startIndex+len(ManifestStartTag) : endIndex], nil
}

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

// GetLineBoundsIfEmpty returns the start and end el indexes passed in
// if the line is not otherwise empty. If it is empty, it will return
// the starting and ending indexes of the current line
func GetLineBoundsIfEmpty(contents string, startOfEl, endOfEl int) (int, int) {
	var startOfLine = strings.LastIndex(contents[:startOfEl], "\n")
	if startOfLine == -1 {
		startOfLine = 0
	}

	var (
		removeLine = strings.TrimSpace(contents[startOfLine:startOfEl]) == ""
		endOfLine  = strings.Index(contents[endOfEl:], "\n")
	)
	if endOfLine == -1 {
		endOfLine = len(contents)
	} else {
		endOfLine += endOfEl
	}

	removeLine = removeLine && strings.TrimSpace(contents[endOfEl:endOfLine]) == ""
	if removeLine {
		return startOfLine, endOfLine
	}

	return startOfEl, endOfEl
}

// ExtractAttribute returns the attribute value in the specified text.
// Note: text is meant to be a single line string and attribute is just the attribute name.
func ExtractAttribute(text, attribute string) (string, int, int, error) {
	var (
		attributeIndicator = attribute + "="
		startOfAttribute   = strings.Index(text, attributeIndicator)
	)
	if startOfAttribute == -1 {
		return "", -1, -1, fmt.Errorf("%s attribute not found", attribute)
	}

	startOfAttribute += len(attributeIndicator)
	var quote = text[startOfAttribute : startOfAttribute+1]
	startOfAttribute++

	var endOfAttribute = strings.Index(text[startOfAttribute:], quote)
	if endOfAttribute == -1 {
		return "", -1, -1, fmt.Errorf("%s attribute value not found", attribute)
	}

	endOfAttribute += startOfAttribute

	return text[startOfAttribute:endOfAttribute], startOfAttribute, endOfAttribute, nil
}
