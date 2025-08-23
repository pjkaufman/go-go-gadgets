package linter

import "regexp"

var (
	subordinateClausesRegex        = regexp.MustCompile(`(?m)^[\r\t\f\v ]*?<p[^>]*?>[^\n]*?(?i:although|because|while)[^.!?\n]*?, (?:but|thus|therefore|furthermore|however)[^\n]*?</p>`)
	subordinateClausesReplaceRegex = regexp.MustCompile(`(.*?,)\s*?(?:but|thus|therefore|furthermore|however)(.*)`)
)

func GetPotentiallyLackingSubordinateClauseInstances(fileContent string) (map[string]string, error) {
	matches := subordinateClausesRegex.FindAllString(fileContent, -1)
	originalToSuggested := make(map[string]string, len(matches))
	if len(matches) == 0 {
		return originalToSuggested, nil
	}

	for _, match := range matches {
		suggestion := subordinateClausesReplaceRegex.ReplaceAllString(match, "$1$2")
		originalToSuggested[match] = suggestion
	}

	return originalToSuggested, nil
}
