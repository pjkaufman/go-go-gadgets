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
		result             = make([]rune, 0, len(runes))
		insideDoubleQuotes = false
		doubleQuoteCount   = 0
		singleQuoteCount   = 0 // Only counts non-contraction single quotes
		updateMade         = false
	)

	for i := 0; i < len(runes); i++ {
		currentRune := runes[i]

		if currentRune == '"' {
			insideDoubleQuotes = !insideDoubleQuotes
			doubleQuoteCount++
			result = append(result, currentRune)
		} else if currentRune == '\'' {
			// Check if it's a contraction (surrounded by letters)
			isPrevLetter := i > 0 && unicode.IsLetter(runes[i-1])
			isNextLetter := i < len(runes)-1 && unicode.IsLetter(runes[i+1])
			isContraction := isPrevLetter && isNextLetter

			if !isContraction {
				singleQuoteCount++
			}

			// If it's not a contraction and not inside double quotes, convert to double quote
			if !isContraction && !insideDoubleQuotes {
				result = append(result, '"')
				updateMade = true
			} else {
				result = append(result, '\'')
			}
		} else {
			result = append(result, currentRune)
		}
	}

	// Validate quote counts
	if doubleQuoteCount%2 != 0 {
		return "", false, fmt.Errorf("unmatched double quotes: found %d double quotes", doubleQuoteCount)
	}

	if singleQuoteCount%2 != 0 {
		return "", false, fmt.Errorf("unmatched single quotes: found %d non-contraction single quotes", singleQuoteCount)
	}

	return string(result), updateMade, nil
}
