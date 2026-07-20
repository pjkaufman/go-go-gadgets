package potentiallyfixableissue

import (
	"fmt"
	"regexp"
	"strings"
)

var thoughtParagraphs = regexp.MustCompile(`(<p[^\n>]*?>[^\n(]*?)\(([^\n()]*?)\)([^\n]*?)(</p>)`)
var parenthesesContent = regexp.MustCompile(`\(([^\n\)]*?)\)`)

func GetPotentialThoughtInstances(fileContent string) (map[string]string, error) {
	var subMatches = thoughtParagraphs.FindAllStringSubmatch(fileContent, -1)
	var originalToSuggested = make(map[string]string, len(subMatches))
	if len(subMatches) == 0 {
		return originalToSuggested, nil
	}

	for _, groups := range subMatches {
		var (
			replaceValue           = groups[1]
			italicsAlreadyIncluded bool
			priorValue             = strings.TrimSpace(groups[1])
			actualThoughtText      = strings.TrimSpace(groups[2])
			followingValue         = strings.TrimSpace(groups[3])
			thoughtParenMatches    = parenthesesContent.FindAllStringSubmatch(groups[0], -1)
		)
		italicsAlreadyIncluded = checkIfElStartsAndEndsGrouping(priorValue, actualThoughtText, followingValue, "<i>", "</i>") || checkIfElStartsAndEndsGrouping(priorValue, actualThoughtText, followingValue, "<em>", "</em>")
		if len(thoughtParenMatches) == 0 {
			if italicsAlreadyIncluded {
				replaceValue += groups[2] + groups[3]
			} else {
				replaceValue += "<i>" + groups[2] + "</i>" + groups[3]
			}
		} else {
			var (
				textReplacement   = "(" + groups[2] + ")" + groups[3]
				replacementFormat = "<i>%s</i>"
			)
			if italicsAlreadyIncluded {
				replacementFormat = "%s"
			}

			for _, parenGroups := range thoughtParenMatches {
				textReplacement = strings.Replace(textReplacement, parenGroups[0], fmt.Sprintf(replacementFormat, parenGroups[1]), 1)
			}

			replaceValue += textReplacement
		}

		replaceValue += groups[4]

		originalToSuggested[groups[0]] = replaceValue
	}

	return originalToSuggested, nil
}

func checkIfElStartsAndEndsGrouping(prior, actual, following, startEl, endEl string) bool {
	return (strings.HasSuffix(prior, startEl) || strings.HasPrefix(actual, startEl)) &&
		(strings.HasSuffix(actual, endEl) || strings.HasPrefix(following, endEl))
}
