package linter

import (
	"fmt"
	"regexp"
	"unicode"

	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

var paragraphsWithSingleQuotes = regexp.MustCompile(`(?m)^([\r\t\f\v ]*?<p[^\n>]*?>)([^\n]*?'[^\n]*?)(</p>)`)

func GetPotentialIncorrectSingleQuotes(fileContent string) map[string]string {
	var subMatches = paragraphsWithSingleQuotes.FindAllStringSubmatch(fileContent, -1)
	var originalToSuggested = make(map[string]string, 0)
	if len(subMatches) == 0 {
		return originalToSuggested
	}

	for _, groups := range subMatches {
		replacedSingleQuoteString, updateMade, err := convertQuotes(groups[2])
		if err != nil {
			logger.WriteErrorf("Failed to convert single quotes to double as needed on string %q: %s", groups[0], err)
		}

		if updateMade {
			originalToSuggested[groups[0]] = groups[1] + replacedSingleQuoteString + groups[3]
		}
	}

	return originalToSuggested
}

func convertQuotes(input string) (string, bool, error) {
	var (
		runes              = []rune(input)
		insideDoubleQuotes = false
		doubleQuoteCount   = 0
		singleQuoteCount   = 0 // Only counts non-possesive and non-contraction single quotes
		updateMade         = false
	)

	for i := 0; i < len(runes); i++ {
		currentRune := runes[i]

		if currentRune == '"' {
			insideDoubleQuotes = !insideDoubleQuotes
			doubleQuoteCount++
		} else if currentRune == '\'' {
			// Check if it's a contraction (surrounded by letters)
			isPrevLetter := i > 0 && unicode.IsLetter(runes[i-1])
			isNextLetter := i < len(runes)-1 && unicode.IsLetter(runes[i+1])
			isContraction := isPrevLetter && isNextLetter

			isPrevS := i > 0 && runes[i-1] == 's'
			isNextS := i < len(runes)-1 && runes[i+1] == 's'
			isPrevWord := i > 0 && unicode.IsLetter(runes[i-1])
			// we will assume that no possesives show up inside a single quote as that gets hairy and is not valid
			isPossessive := (isPrevS || (isPrevWord && isNextS)) && singleQuoteCount%2 == 0

			if !isContraction && !isPossessive {
				singleQuoteCount++
			}

			// If it's not a contraction, not a possessive, and not inside double quotes, convert to double quote
			if !isContraction && !isPossessive && !insideDoubleQuotes {
				runes[i] = '"'
				updateMade = true
			}
		}
	}

	if doubleQuoteCount%2 != 0 {
		return "", false, fmt.Errorf("unmatched double quotes: found %d double quotes", doubleQuoteCount)
	}

	if singleQuoteCount%2 != 0 {
		return "", false, fmt.Errorf("unmatched single quotes: found %d non-contraction single quotes", singleQuoteCount)
	}

	return string(runes), updateMade, nil
}
