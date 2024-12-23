package linter

import "regexp"

var (
	althoughButRegex = regexp.MustCompile(`(?m)^[\r\t\f\v ]*?<p[^>]*?>.*?(?i:although|because|while)[^.!?]*?, (?:but|thus|therefore|furthermore|however).*?</p>`)
	replaceRegex     = regexp.MustCompile(`(.*?,)\s*?(?:but|thus|therefore|furthermore|however)(.*)`)
)

func GetPotentiallyLackingSubordinateClauseInstances(fileContent string) map[string]string {
	matches := althoughButRegex.FindAllString(fileContent, -1)
	originalToSuggested := make(map[string]string, len(matches))
	if len(matches) == 0 {
		return originalToSuggested
	}

	for _, match := range matches {
		suggestion := replaceRegex.ReplaceAllString(match, "$1$2")
		originalToSuggested[match] = suggestion
	}

	return originalToSuggested
}
