package linter

import "regexp"

var althoughButRegex = regexp.MustCompile(`(\n[\t ]*<p[^\n>]*>[^\n]*[Aa]lthough[^\n.!?]*,)( but)([^\n]*</p>)`)

func GetPotentialAlthoughButInstances(fileContent string) map[string]string {
	var subMatches = althoughButRegex.FindAllStringSubmatch(fileContent, -1)
	var originalToSuggested = make(map[string]string, len(subMatches))
	if len(subMatches) == 0 {
		return originalToSuggested
	}

	for _, groups := range subMatches {
		originalToSuggested[groups[0]] = groups[1] + groups[3]
	}

	return originalToSuggested
}
