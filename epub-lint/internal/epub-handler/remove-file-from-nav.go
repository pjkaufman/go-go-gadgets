package epubhandler

import "strings"

func RemoveFileFromNav(text, file string) string {
	hrefAttributeEnd := strings.Index(text, file+`"`)
	if hrefAttributeEnd == -1 {
		return text
	}

	// check if it truly is in an href...
	startOfAttribute := strings.LastIndex(text[:hrefAttributeEnd], `"`)
	if startOfAttribute == -1 { // something is wrong, so we are going to skip this...
		return text
	}

	if !strings.HasSuffix(text[:startOfAttribute], "href=") {
		return text
	}

	tagName, tagStart := getPreviousHtmlTagIfNotClosingTag(text, startOfAttribute)
	if tagStart == -1 {
		return text
	}

	if tagName != "a" { // this is not a scenario we are handling at this time...
		return text
	}

	tagName, tagStart = getPreviousHtmlTagIfNotClosingTag(text, tagStart)
	if tagStart == -1 {
		return text
	}

	if tagName != "li" { // this is not a scenario we are handling at this time...
		return text
	}

	// it should be safe to just excise the content until the next `</li>` at this point
	const endingListItem = "</li>"
	endingTagIndex := strings.Index(text[startOfAttribute:], endingListItem)
	if endingTagIndex == -1 {
		return text
	}

	endingTagIndex += len(endingListItem) + startOfAttribute

	// check if line just has starting whitespace
	startOfLine := strings.LastIndex(text[:tagStart], "\n")
	if startOfLine == -1 {
		startOfLine = 0
	}

	var (
		startOfRemoval            = tagStart
		startingWhitespaceRemoved bool
	)
	if strings.TrimSpace(text[startOfLine:tagStart]) == "" {
		startOfRemoval = startOfLine
		startingWhitespaceRemoved = true
	}

	endOfLine := strings.Index(text[endingTagIndex:], "\n")
	if endOfLine == -1 {
		endOfLine = len(text)
	} else {
		endOfLine += endingTagIndex
	}

	var endOfRemoval = endingTagIndex
	if startingWhitespaceRemoved && strings.TrimSpace(text[endingTagIndex:endOfLine]) == "" {
		endOfRemoval = endOfLine
	}

	return text[:startOfRemoval] + text[endOfRemoval:]
}

// getPreviousHtmlTagIfNotClosingTag gets the nearest prior opening html tag and the start of its position
// if there is no closing tag between it and the start index
func getPreviousHtmlTagIfNotClosingTag(text string, start int) (string, int) {
	elStart := strings.LastIndex(text[:start], "<")
	if elStart == -1 { // if this happens there is no prior HTML tag
		return "", -1
	}

	if text[elStart+1] == '/' { // if the element is a closing tag/is self closing return an empty string
		return "", -1
	}

	elementEndEnd := strings.Index(text[elStart:], ">")
	if elementEndEnd == -1 {
		return "", -1
	}

	var elementEndIndex = elStart + elementEndEnd
	elementEndEnd = strings.Index(text[elStart:elStart+elementEndEnd], " ")
	if elementEndEnd != -1 {
		elementEndIndex = elStart + elementEndEnd
	}

	return text[elStart+1 : elementEndIndex], elStart
}
