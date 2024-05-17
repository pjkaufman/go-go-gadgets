package linter

import "regexp"

var squareBracketConversationRegex = regexp.MustCompile(`(<p[^\n>]*>\s*)\[([^\n]*)\](\s*</p>)`)

func GetPotentialSquareBracketConversationInstances(fileContent string) map[string]string {
	var subMatches = squareBracketConversationRegex.FindAllStringSubmatch(fileContent, -1)
	var originalToSuggested = make(map[string]string, len(subMatches))
	if len(subMatches) == 0 {
		return originalToSuggested
	}

	for _, groups := range subMatches {
		originalToSuggested[groups[0]] = groups[1] + `"` + groups[2] + `"` + groups[3]
	}

	return originalToSuggested
}
