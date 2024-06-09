package linter

import (
	"regexp"
)

// missing oxford comma regex based on https://stackoverflow.com/questions/30006666/capture-a-list-of-words-that-doesnt-contain-an-oxford-comma/30006707#30006707
var oxfordCommaRegex = regexp.MustCompile(`(\n[\t ]*<p[^\n>]*>[^\n]*)(\w+)((,\s*\w+)+)(\s+)(and|or)(\s+\w+)([^\n]*</p>)`)

func GetPotentialMissingOxfordCommas(fileContent string) map[string]string {
	var subMatches = oxfordCommaRegex.FindAllStringSubmatch(fileContent, -1)
	var originalToSuggested = make(map[string]string, len(subMatches))
	if len(subMatches) == 0 {
		return originalToSuggested
	}

	for _, groups := range subMatches {
		originalToSuggested[groups[0]] = groups[1] + groups[2] + groups[3] + "," + groups[5] + groups[6] + groups[7] + groups[8]
	}

	return originalToSuggested
}
