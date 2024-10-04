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

	var (
		attributeEnd  = htmlOpenStart + htmlOpenEnd
		htmlEl        = text[htmlOpenStart:attributeEnd]
		newHtmlEl     strings.Builder
		langAttrIndex = strings.Index(htmlEl, langAttribute)
	)
	if langAttrIndex == -1 {
		return text[:attributeEnd] + " " + langAttribute + "\"" + lang + "\" " + xmlLangAttribute + "\"" + lang + "\"" + text[attributeEnd:]
	}

	var (
		regularLangAttrIsHandled, xmlLangAttributeIsHandled bool
		startOfAttr, endOfAttr                              int
		langVal, endingIndicator, char                      string
	)

	handleLangAttribute := func() {
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

		if strings.TrimSpace(langVal) == "" {
			newHtmlEl.WriteString(htmlEl[langAttrIndex : startOfAttr+1])
			newHtmlEl.WriteString(lang)
			newHtmlEl.WriteString(endingIndicator)
		} else {
			newHtmlEl.WriteString(htmlEl[langAttrIndex : endOfAttr+1])
		}
	}

	newHtmlEl.WriteString(htmlEl[:langAttrIndex])
	if htmlEl[langAttrIndex-1:langAttrIndex] == ":" {
		xmlLangAttributeIsHandled = true
	} else {
		regularLangAttrIsHandled = true
	}

	handleLangAttribute()

	langAttrIndex = strings.Index(htmlEl[endOfAttr:], langAttribute)
	if langAttrIndex != -1 {
		langAttrIndex += endOfAttr

		newHtmlEl.WriteString(htmlEl[endOfAttr+1 : langAttrIndex])
		if htmlEl[langAttrIndex-1:langAttrIndex] == ":" {
			xmlLangAttributeIsHandled = true
		} else {
			regularLangAttrIsHandled = true
		}

		handleLangAttribute()
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

	if endOfAttr < len(htmlEl)-1 {
		newHtmlEl.WriteString(htmlEl[endOfAttr+1:])
	}

	return text[0:htmlOpenStart] + newHtmlEl.String() + text[attributeEnd:]
}
