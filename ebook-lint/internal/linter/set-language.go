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

// const (
// 	langAttribute    = "lang="
// 	xmlLangAttribute = "xml:lang="
// 	openingHtmlTag   = "<html"
// )

// TODO: see about fixing the issues with this
// func EnsureLanguageIsSet(text, lang string) string {
// 	var htmlOpenStart = strings.Index(text, openingHtmlTag)
// 	if htmlOpenStart == -1 {
// 		return text
// 	}

// 	var htmlOpenEnd = strings.Index(text[htmlOpenStart:], `>`)
// 	if htmlOpenEnd == -1 {
// 		return text
// 	}

// 	var attributeEnd = htmlOpenStart + htmlOpenEnd
// 	var htmlEl = text[htmlOpenStart:attributeEnd]
// 	var newHtmlEl strings.Builder
// 	var langAttrIndex = strings.Index(htmlEl, langAttribute)
// 	if langAttrIndex == -1 {
// 		return text[:htmlOpenEnd] + " " + langAttribute + "\"" + lang + "\" " + xmlLangAttribute + "\"" + lang + "\"" + text[htmlOpenEnd:]
// 	}

// 	var (
// 		regularLangAttrIsHandled, xmlLangAttributeIsHandled bool
// 		startOfAttr, endOfAttr                              int
// 		langVal, endingIndicator, char                      string
// 	)

// 	// fmt.Println("1", text[htmlOpenStart+len(openingHtmlTag):htmlOpenStart+langAttrIndex])
// 	// fmt.Println("2", htmlEl[len(openingHtmlTag):langAttrIndex])
// 	newHtmlEl.WriteString(strings.TrimSpace(htmlEl[len(openingHtmlTag):langAttrIndex]))
// 	for langAttrIndex != -1 {
// 		startOfAttr = langAttrIndex + len(langAttribute)
// 		endOfAttr = startOfAttr + 1
// 		endingIndicator = htmlEl[startOfAttr:endOfAttr]
// 		langVal = ""
// 		for endOfAttr < len(htmlEl) {
// 			char = htmlEl[endOfAttr : endOfAttr+1]
// 			if char == endingIndicator {
// 				break
// 			}

// 			langVal += char
// 			endOfAttr++
// 		}

// 		if langAttrIndex != 0 {
// 			newHtmlEl.WriteString(htmlEl[0:langAttrIndex])
// 		}

// 		if strings.TrimSpace(langVal) == "" {
// 			newHtmlEl.WriteString(htmlEl[langAttrIndex : startOfAttr+1])
// 			newHtmlEl.WriteString(lang)
// 			newHtmlEl.WriteString(endingIndicator)
// 		} else {
// 			newHtmlEl.WriteString(htmlEl[langAttrIndex : endOfAttr+1])
// 		}

// 		if langAttrIndex == 0 || htmlEl[langAttrIndex-1:langAttrIndex] != ":" {
// 			regularLangAttrIsHandled = true
// 		} else {
// 			xmlLangAttributeIsHandled = true
// 		}

// 		if endOfAttr == len(htmlEl) {
// 			break
// 		}

// 		htmlEl = htmlEl[endOfAttr+1:]
// 		langAttrIndex = strings.Index(htmlEl, langAttribute)
// 		if langAttrIndex != -1 {
// 			newHtmlEl.WriteString(htmlEl[0:langAttrIndex])
// 		}
// 	}

// 	if !regularLangAttrIsHandled {
// 		newHtmlEl.WriteString(" ")
// 		newHtmlEl.WriteString(langAttribute)
// 		newHtmlEl.WriteString("\"")
// 		newHtmlEl.WriteString(lang)
// 		newHtmlEl.WriteString("\"")
// 	}

// 	if !xmlLangAttributeIsHandled {
// 		newHtmlEl.WriteString(" ")
// 		newHtmlEl.WriteString(xmlLangAttribute)
// 		newHtmlEl.WriteString("\"")
// 		newHtmlEl.WriteString(lang)
// 		newHtmlEl.WriteString("\"")
// 	}

// 	// if previousIndex < htmlOpenEnd {
// 	newHtmlEl.WriteString(htmlEl)
// 	// }

// 	return text[0:htmlOpenStart] + newHtmlEl.String() + text[attributeEnd:]
// }
