package linter

import (
	"strings"
)

const (
	defaultXmlEl = `<?xml version="1.0" encoding="utf-8"?>` + "\n"
	openingXmlEl = `<?xml`
	closingXmlEl = `?>`
	encodingAttr = `encoding=`
)

func EnsureEncodingIsPresent(text string) string {
	var xmlOpeningElTag = strings.Index(text, openingXmlEl)
	if xmlOpeningElTag == -1 {
		return defaultXmlEl + text
	}

	var xmlEndOfTagIndex = strings.Index(text[xmlOpeningElTag:], closingXmlEl)
	if xmlEndOfTagIndex == -1 {
		return text
	}

	var attributeEnd = xmlOpeningElTag + xmlEndOfTagIndex
	var xmlEl = text[xmlOpeningElTag:attributeEnd]
	var encodingAttrIndex = strings.Index(xmlEl, encodingAttr)
	if encodingAttrIndex == -1 {
		return text[0:attributeEnd] + " " + encodingAttr + "\"utf-8\"" + text[attributeEnd:]
	}

	startOfAttr := encodingAttrIndex + len(encodingAttr)
	endOfAttr := startOfAttr + 1
	endingIndicator := xmlEl[startOfAttr:endOfAttr]
	var encodingVal, char string
	for endOfAttr < len(xmlEl) {
		char = xmlEl[endOfAttr : endOfAttr+1]
		if char == endingIndicator {
			break
		}

		encodingVal += char
		endOfAttr++
	}

	if strings.TrimSpace(encodingVal) == "" {
		return text[0:startOfAttr+1] + "utf-8" + text[endOfAttr:]
	} else {
		return text
	}
}
