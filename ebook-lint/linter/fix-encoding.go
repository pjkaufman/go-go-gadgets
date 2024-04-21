package linter

import (
	"regexp"
	"strings"
)

var xmlElRegex = regexp.MustCompile(`<\?xml([^\n>?])*\?>`)
var hasUpdatedFirstXmlEl = false

const defaultXmlEl = `<?xml version="1.0" encoding="utf-8"?>` + "\n"

func EnsureEncodingIsPresent(text string) string {
	hasUpdatedFirstXmlEl = false

	if xmlElRegex.MatchString(text) {
		return xmlElRegex.ReplaceAllStringFunc(text, addEncodingIfMissing)
	}

	return defaultXmlEl + text
}

func addEncodingIfMissing(part string) string {
	if hasUpdatedFirstXmlEl || strings.Contains(part, "encoding=") {
		hasUpdatedFirstXmlEl = true
		return part
	}

	hasUpdatedFirstXmlEl = true
	return strings.Replace(part, "?>", ` encoding="utf-8"?>`, 1)
}
