package potentiallyfixableissue

import "regexp"

var squareBracketConversationRegex = regexp.MustCompile(`(<p[^\n>]*?>[\r\t\f\v ]*?(<a[^>]*?></a>[\r\t\f\v ]*?)?)\[([^\n]*?)\]([\r\t\f\v ]*?</p>)`)

func GetPotentialSquareBracketConversationInstances(fileContent string) (map[string]string, error) {
	var subMatches = squareBracketConversationRegex.FindAllStringSubmatch(fileContent, -1)
	var originalToSuggested = make(map[string]string, len(subMatches))
	if len(subMatches) == 0 {
		return originalToSuggested, nil
	}

	for _, groups := range subMatches {
		originalToSuggested[groups[0]] = groups[1] + `"` + groups[3] + `"` + groups[4]
	}

	return originalToSuggested, nil
}
