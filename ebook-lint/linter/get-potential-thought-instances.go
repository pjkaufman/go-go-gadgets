package linter

import (
	"regexp"
	"strings"
)

var thoughtParagraphs = regexp.MustCompile(`(<p[^\n>]*>[^\n(]*)\(([^\n()]*)\)([^\n]*)(</p>)`)
var parenthesesContent = regexp.MustCompile(`\(([^\n\)]*)\)`)

func GetPotentialThoughtInstances(fileContent string) map[string]string {
	var subMatches = thoughtParagraphs.FindAllStringSubmatch(fileContent, -1)
	var originalToSuggested = make(map[string]string, len(subMatches))
	if len(subMatches) == 0 {
		return originalToSuggested
	}

	for _, groups := range subMatches {
		var replaceValue = groups[1]

		var thoughtParenMatches = parenthesesContent.FindAllStringSubmatch(groups[0], -1)
		if len(thoughtParenMatches) == 0 {
			replaceValue += "<i>" + groups[2] + "</i>" + groups[3]
		} else {
			var textReplacement = "(" + groups[2] + ")" + groups[3]
			for _, parenGroups := range thoughtParenMatches {
				textReplacement = strings.Replace(textReplacement, parenGroups[0], "<i>"+parenGroups[1]+"</i>", 1)
			}

			replaceValue += textReplacement
		}

		replaceValue += groups[4]

		originalToSuggested[groups[0]] = replaceValue
	}

	return originalToSuggested
}
