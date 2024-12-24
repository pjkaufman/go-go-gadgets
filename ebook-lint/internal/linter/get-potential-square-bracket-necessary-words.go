package linter

import (
	"regexp"
	"strings"
)

var (
	squareBracketNecessaryWordRegex = regexp.MustCompile(`(<p[^\n>]*?>[^\n[]*?)\[([^\n[\]]*?)\]([^\n]*?)(</p>)`)
	squareBracketContentRegex       = regexp.MustCompile(`\[([^\n\]]*?)\]`)
)

func GetPotentialSquareBracketNecessaryWords(fileContent string) map[string]string {
	var subMatches = squareBracketNecessaryWordRegex.FindAllStringSubmatch(fileContent, -1)
	var originalToSuggested = make(map[string]string, len(subMatches))
	if len(subMatches) == 0 {
		return originalToSuggested
	}

	for _, groups := range subMatches {
		// we need to skip the match if there is no content other than whitespace in group 3 and
		// group 1 is just the opening paragraph tag and whitespace
		if strings.TrimSpace(groups[3]) == "" && strings.HasSuffix(strings.TrimSpace(groups[1]), ">") {
			continue
		}

		var replaceValue = groups[1]

		var squareBracketMatches = squareBracketContentRegex.FindAllStringSubmatch(groups[0], -1)
		if len(squareBracketMatches) == 0 {
			replaceValue += groups[2] + groups[3]
		} else {
			var textReplacement = "[" + groups[2] + "]" + groups[3]
			for _, parenGroups := range squareBracketMatches {
				textReplacement = strings.Replace(textReplacement, parenGroups[0], parenGroups[1], 1)
			}

			replaceValue += textReplacement
		}

		replaceValue += groups[4]

		originalToSuggested[groups[0]] = replaceValue
	}

	return originalToSuggested
}
