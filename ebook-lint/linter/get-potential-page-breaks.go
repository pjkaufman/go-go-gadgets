package linter

import (
	"regexp"
)

const PageBrakeEl = `<hr class="blankSpace" />`

var emptyParagraphsOrDivs = regexp.MustCompile(`(\n[ \t]*<(p|div)[^\n>]*>)[ \t]*(</(p|div)>)`)

func GetPotentialPageBreaks(fileContent string) map[string]string {
	var subMatches = emptyParagraphsOrDivs.FindAllStringSubmatch(fileContent, -1)
	var originalToSuggested = make(map[string]string, len(subMatches))
	if len(subMatches) == 0 {
		return originalToSuggested
	}

	for _, groups := range subMatches {
		originalToSuggested[groups[0]] = "\n" + PageBrakeEl
	}

	return originalToSuggested
}
