package linter

import (
	"fmt"
	"regexp"
	"strings"
)

const SectionBreakEl = `<hr class="character" />`

func GetPotentialSectionBreaks(fileContent, sectionBreakIndicator string) map[string]string {
	var contextBreakRegex = regexp.MustCompile(fmt.Sprintf(`(\n[ \t]*<p[^\n>]*>([^\n])*)%s(([^\n]*)</p>)`, regexp.QuoteMeta(sectionBreakIndicator)))

	var subMatches = contextBreakRegex.FindAllStringSubmatch(fileContent, -1)
	var originalToSuggested = make(map[string]string, len(subMatches))
	if len(subMatches) == 0 {
		return originalToSuggested
	}

	for _, groups := range subMatches {
		var replaceValue = "\n" + SectionBreakEl
		if strings.TrimSpace(groups[2]) != "" || strings.TrimSpace(groups[4]) != "" {
			replaceValue = groups[1] + SectionBreakEl + groups[3]
		}

		originalToSuggested[groups[0]] = replaceValue
	}

	return originalToSuggested
}
