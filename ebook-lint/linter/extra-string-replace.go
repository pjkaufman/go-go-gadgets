package linter

import "strings"

func ExtraStringReplace(text string, extraFindAndReplaces map[string]string, numHits map[string]int) string {
	var newText = text

	var stringsToReplace []string = make([]string, 2*len(extraFindAndReplaces))
	var i = 0;
	for search, replace := range extraFindAndReplaces {
		if hits, ok := numHits[search]; ok {
			numHits[search] = hits + strings.Count(newText, search)
		} else {
			numHits[search] = strings.Count(newText, search)
		}

		stringsToReplace[2*i] = search
		stringsToReplace[2*i+1] = replace
		i++
	}

	var replacer = strings.NewReplacer(stringsToReplace...)
	newText = replacer.Replace(newText)

	return newText
}
