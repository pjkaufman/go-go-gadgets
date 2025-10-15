package potentiallyfixableissue

import (
	"regexp"
)

const PageBrakeEl = `<hr class="blankSpace" />`

var emptyParagraphsOrDivs = regexp.MustCompile(`(?m)^([\r\t\f\v ]*?<(p|div)[^\n>]*?>)[\r\t\f\v ]*?(</(p|div)>)`)

func GetPotentialPageBreaks(fileContent string) (map[string]string, error) {
	var subMatches = emptyParagraphsOrDivs.FindAllStringSubmatch(fileContent, -1)
	var originalToSuggested = make(map[string]string, len(subMatches))
	if len(subMatches) == 0 {
		return originalToSuggested, nil
	}

	for _, groups := range subMatches {
		originalToSuggested[groups[0]] = PageBrakeEl
	}

	return originalToSuggested, nil
}
