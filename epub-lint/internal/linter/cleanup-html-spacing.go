package linter

import (
	"regexp"
	"unicode"
)

var (
	blankLineRegex                                = regexp.MustCompile(`([ \t]*\n){2,}`)
	paragraphElRemoveEndingInternalSpacingRegex   = regexp.MustCompile(`\s*\n\s*<\/p>`)
	paragraphElRemoveStartingInternalSpacingRegex = regexp.MustCompile(`(<p[^\n>]*>)\s*\n\s*`)
)

func CleanupHtmlSpacing(text string) string {
	// general whitespace
	text = removeStartingSpacing(text)
	text = removeEndingSpacing(text)
	text = blankLineRegex.ReplaceAllString(text, "\n") // remove blank lines

	// html whitespace
	text = paragraphElRemoveEndingInternalSpacingRegex.ReplaceAllString(text, "</p>")
	text = paragraphElRemoveStartingInternalSpacingRegex.ReplaceAllString(text, "$1")

	return text
}

func removeStartingSpacing(text string) string {
	for i := 0; i < len(text); i++ {
		if !unicode.IsSpace(rune(text[i])) {
			return text[i:]
		}
	}

	return text
}

func removeEndingSpacing(text string) string {
	for i := len(text) - 1; i >= 0; i-- {
		if !unicode.IsSpace(rune(text[i])) {
			return text[0:i+1] + "\n"
		}
	}

	return text
}
