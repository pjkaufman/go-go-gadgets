package linter

import (
	"strings"
)

const (
	langAttribute    = "lang="
	xmlLangAttribute = "xml:lang="
	openingHtmlTag   = "<html"
)

func EnsureLanguageIsSet(text, lang string) string {
	var htmlOpenStart = strings.Index(text, openingHtmlTag)
	if htmlOpenStart == -1 {
		return text
	}

	var htmlOpenEnd = strings.Index(text[htmlOpenStart:], `>`)
	if htmlOpenEnd == -1 {
		return text
	}

	var htmlEl = text[htmlOpenStart:htmlOpenEnd]
	var newHtmlEl strings.Builder
	var langAttrIndex = strings.Index(htmlEl, langAttribute)
	if langAttrIndex == -1 {
		return text[htmlOpenEnd-1:] + " " + langAttribute + "=\"" + lang + "\" " + xmlLangAttribute + "=\"" + lang + "\"" + text[htmlOpenEnd:]
	}

	var (
		regularLangAttrIsHandled, xmlLangAttributeIsHandled bool
		startOfAttr, endOfAttr                              int
		langVal, endingIndicator, char                      string
	)
	for langAttrIndex != -1 {
		startOfAttr = langAttrIndex + len(langAttribute)
		endOfAttr = startOfAttr + 1
		endingIndicator = htmlEl[startOfAttr:endOfAttr]
		langVal = ""
		for endOfAttr < len(htmlEl) {
			char = htmlEl[endOfAttr : endOfAttr+1]
			if char == endingIndicator {
				break
			}

			langVal += char
			endOfAttr++
		}

		if langAttrIndex != 0 {
			newHtmlEl.WriteString(htmlEl[0:langAttrIndex])
		}

		if strings.TrimSpace(langVal) == "" {
			newHtmlEl.WriteString(htmlEl[langAttrIndex : startOfAttr+1])
			newHtmlEl.WriteString(lang)
			newHtmlEl.WriteString(endingIndicator)
		} else {
			newHtmlEl.WriteString(htmlEl[langAttrIndex : endOfAttr+1])
		}

		if langAttrIndex == 0 || htmlEl[langAttrIndex-1:langAttrIndex] != ":" {
			regularLangAttrIsHandled = true
		} else {
			xmlLangAttributeIsHandled = true
		}

		if endOfAttr == len(htmlEl) {
			break
		}

		htmlEl = htmlEl[endOfAttr+1:]
		langAttrIndex = strings.Index(htmlEl, langAttribute)
	}

	if !regularLangAttrIsHandled {
		newHtmlEl.WriteString(" ")
		newHtmlEl.WriteString(langAttribute)
		newHtmlEl.WriteString("\"")
		newHtmlEl.WriteString(lang)
		newHtmlEl.WriteString("\"")
	}

	if !xmlLangAttributeIsHandled {
		newHtmlEl.WriteString(" ")
		newHtmlEl.WriteString(xmlLangAttribute)
		newHtmlEl.WriteString("\"")
		newHtmlEl.WriteString(lang)
		newHtmlEl.WriteString("\"")
	}

	return text[0:htmlOpenStart] + newHtmlEl.String() + text[htmlOpenEnd:]
}
