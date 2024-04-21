package linter

import (
	"fmt"
	"regexp"
	"strings"
)

var htmlElRegex = regexp.MustCompile(`<html([^\n>])*>`)
var langAttributeRegex = regexp.MustCompile(`([^:]lang=["'])([^"'\n>]*)(["'])`)
var xmlLangAttributeRegex = regexp.MustCompile(`(xml:lang=["'])([^"'\n>]*)(["'])`)

const (
	langAttribute    = "lang="
	xmlLangAttribute = "xml:lang="
)

func EnsureLanguageIsSet(text, lang string) string {
	var htmlElInfo = htmlElRegex.FindAllString(text, 1)
	if len(htmlElInfo) == 0 {
		return text
	}

	var htmlEl = htmlElInfo[0]
	var newHtmlEl = addAttributeIfMissingOrUpdateIfEmpty(htmlEl, langAttribute, lang, langAttributeRegex)
	newHtmlEl = addAttributeIfMissingOrUpdateIfEmpty(newHtmlEl, xmlLangAttribute, lang, xmlLangAttributeRegex)

	return strings.Replace(text, htmlEl, newHtmlEl, 1)
}

func addAttributeIfMissingOrUpdateIfEmpty(htmlEl, attribute, lang string, attributeRegex *regexp.Regexp) string {
	var attributeInfo = attributeRegex.FindAllStringSubmatch(htmlEl, 1)
	if len(attributeInfo) != 0 {
		var attributeValue = attributeInfo[0][2]
		if strings.Trim(attributeValue, " ") == "" {
			return strings.Replace(htmlEl, attributeInfo[0][0], fmt.Sprintf("%s%s%s", attributeInfo[0][1], lang, attributeInfo[0][3]), 1)
		}

		return htmlEl
	}

	return strings.Replace(htmlEl, ">", fmt.Sprintf(" %s\"%s\">", attribute, lang), 1)
}
