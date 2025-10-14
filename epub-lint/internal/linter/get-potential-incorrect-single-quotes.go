package linter

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// commonContractions that either start or end with a single quote
var commonContractions = map[string]struct{}{
	"'bout": {}, "'cause": {}, "'cept": {}, "'em": {}, "'gainst": {}, "'neath": {}, "ol'": {},
	"'round": {}, "'s": {}, "shan'": {}, "'thout": {}, "'til": {}, "'tis": {}, "'twas": {},
	"'tween": {}, "'twere": {}, "a'ight": {}, "ain't": {}, "amn't": {}, "'n'": {}, "aren't": {},
	"can't": {}, "cap'n": {}, "c'mon": {}, "could've": {}, "couldn't": {}, "couldn't've": {},
	"daren't": {}, "daresn't": {}, "dasn't": {}, "didn't": {}, "doesn't": {}, "don't": {},
	"d'ye": {}, "d'ya": {}, "e'en": {}, "e'er": {}, "everybody's": {}, "everyone's": {},
	"everything's": {}, "fo'c'sle": {}, "g'day": {}, "giv'n": {}, "gi'z": {}, "gon't": {},
	"hadn't": {}, "had've": {}, "hasn't": {}, "haven't": {}, "he'd": {}, "he'd'nt've": {},
	"he'll": {}, "yesn't": {}, "he's": {}, "here's": {}, "how'd": {}, "how'll": {}, "how're": {},
	"how's": {}, "I'd": {}, "I'd've": {}, "I'd'nt": {}, "I'd'nt've": {}, "If'n": {}, "I'll": {},
	"I'm": {}, "I'm'onna": {}, "I'm'o": {}, "I'm'na": {}, "I've": {}, "isn't": {}, "it'd": {},
	"it'll": {}, "it's": {}, "let's": {}, "loven't": {}, "ma'am": {}, "mayn't": {}, "may've": {},
	"mightn't": {}, "might've": {}, "mine's": {}, "mustn't": {}, "mustn't've": {}, "must've": {},
	"needn't": {}, "ne'er": {}, "nothing's": {}, "o'clock": {}, "o'er": {}, "ought've": {},
	"oughtn't": {}, "oughtn't've": {}, "shalln't": {}, "shan't": {}, "she'd": {}, "she'll": {},
	"she's": {}, "she'd'nt've": {}, "should've": {}, "shouldn't": {}, "shouldn't've": {}, "somebody's": {},
	"someone's": {}, "something's": {}, "so're": {}, "so's": {}, "so've": {}, "that'll": {},
	"that're": {}, "that's": {}, "that'd": {}, "there'd": {}, "there'll": {}, "there're": {},
	"there's": {}, "these're": {}, "these've": {}, "they'd": {}, "they'd've": {}, "they'll": {},
	"they're": {}, "they've": {}, "this's": {}, "those're": {}, "those've": {}, "to've": {},
	"w'all": {}, "w'at": {}, "wasn't": {}, "we'd": {}, "we'd've": {}, "we'll": {}, "we're": {},
	"we've": {}, "weren't": {}, "what'd": {}, "what'll": {}, "what're": {}, "what's": {}, "what've": {},
	"when'd": {}, "when's": {}, "where'd": {}, "where'll": {}, "where're": {}, "where's": {},
	"where've": {}, "which'd": {}, "which'll": {}, "which're": {}, "which's": {}, "which've": {},
	"who'd": {}, "who'd've": {}, "who'll": {}, "who're": {}, "who's": {}, "who've": {},
	"why'd": {}, "why'dja": {}, "why're": {}, "why's": {}, "willn't": {}, "won't": {},
	"would've": {}, "wouldn't": {}, "wouldn't've": {}, "y'ain't": {}, "y'all": {}, "y'all'd've": {},
	"y'all'dn't've": {}, "y'all're": {}, "y'all'ren't": {}, "y'at": {}, "yes'm": {}, "y'ever": {},
	"y'know": {}, "you'd": {}, "you'dn't've": {}, "you'll": {}, "you're": {}, "you've": {},
	"mornin'": {},
}

var paragraphsWithSingleQuotes = regexp.MustCompile(`(?m)^([\r\t\f\v ]*?<p[^\n>]*?>)([^\n]*?'[^\n]*?)(</p>)`)

func GetPotentialIncorrectSingleQuotes(fileContent string) (map[string]string, error) {
	var subMatches = paragraphsWithSingleQuotes.FindAllStringSubmatch(fileContent, -1)
	var originalToSuggested = make(map[string]string, 0)
	if len(subMatches) == 0 {
		return originalToSuggested, nil
	}

	for _, groups := range subMatches {
		replacedSingleQuoteString, updateMade, err := convertQuotes(groups[2])
		if err != nil {
			return nil, fmt.Errorf("Failed to convert single quotes to double as needed on string %q: %s", groups[0], err)
		}

		if updateMade {
			originalToSuggested[groups[0]] = groups[1] + replacedSingleQuoteString + groups[3]
		}
	}

	return originalToSuggested, nil
}

func convertQuotes(input string) (string, bool, error) {
	var (
		runes                                     = []rune(input)
		insideDoubleQuotes                        = false
		doubleQuoteCount                          = 0
		singleQuoteCount                          = 0 // Only counts non-possesive, non-contraction, and non-plural or omission digit single quotes
		updateMade                                = false
		checkForSpecialContractionsAndGetNewStart = func(startIndex int) int {
			var start = startIndex
			for start > 0 && (unicode.IsLetter(runes[start-1]) || runes[start-1] == '\'') {
				start--
			}

			var end = startIndex
			for end < len(runes)-1 && (unicode.IsLetter(runes[end+1]) || runes[end+1] == '\'') {
				end++
			}

			if _, ok := commonContractions[strings.ToLower(string(runes[start:end+1]))]; !ok {
				var startsWithSingleQuote = runes[start] == '\''
				var endsWithSingleQuote = runes[end] == '\''
				// remove starting single quote and see if the string matches a common contraction
				if startsWithSingleQuote {
					_, ok = commonContractions[strings.ToLower(string(runes[start+1:end+1]))]
					if ok {
						singleQuoteCount++

						return end
					}
				}

				// remove ending single quote and see if the string matches a common contraction
				if endsWithSingleQuote {
					_, ok = commonContractions[strings.ToLower(string(runes[start:end]))]
					if ok {
						singleQuoteCount++

						return end
					}
				}

				// remove starting and ending single quote and see if the string matches a common contraction
				if startsWithSingleQuote && endsWithSingleQuote && start+1 <= end {
					_, ok = commonContractions[strings.ToLower(string(runes[start+1:end]))]
					if ok {
						singleQuoteCount += 2

						return end
					}
				}

				// for now, we will do this the less performant way
				for i := start; i <= end; i++ {
					if runes[i] == '\'' {
						singleQuoteCount++

						if !insideDoubleQuotes {
							runes[i] = '"'
							updateMade = true
						}
					}
				}
			}

			return end
		}
	)

	for i := 0; i < len(runes); i++ {
		currentRune := runes[i]

		if currentRune == '"' {
			insideDoubleQuotes = !insideDoubleQuotes
			doubleQuoteCount++
		} else if currentRune == '\'' {
			isPrevDigit := i > 0 && unicode.IsDigit(runes[i-1])
			isNextDigit := i < len(runes)-1 && unicode.IsDigit(runes[i+1])
			isPrevS := i > 0 && (runes[i-1] == 's' || runes[i-1] == 'S')
			isNextS := i < len(runes)-1 && (runes[i+1] == 's' || runes[i+1] == 'S')
			isPrevLetter := i > 0 && unicode.IsLetter(runes[i-1])
			isNextLetter := i < len(runes)-1 && unicode.IsLetter(runes[i+1])

			// is a plural, possesive, or omitted number scenario
			isDigitScenarios := (isPrevDigit && isNextS) || (!isPrevLetter && isNextDigit)
			// we will assume that no possesives show up inside a single quote as that gets hairy and is not valid
			isPossessive := (isPrevS || (isPrevLetter && isNextS)) && singleQuoteCount%2 == 0
			// handles many names that have single quotes in them as well as many contractions
			isBetweenLetters := isPrevLetter && isNextLetter

			fmt.Println(string(runes[i-1]), string(runes[i+1]))
			if isPossessive || isDigitScenarios || isBetweenLetters {
				continue
			}

			i = checkForSpecialContractionsAndGetNewStart(i)
		}
	}

	// Note: this will fail any time we get into measurements like 6'2" (6 foot and 2 inches)
	if doubleQuoteCount%2 != 0 {
		return "", false, fmt.Errorf("unmatched double quotes: found %d double quotes", doubleQuoteCount)
	}

	if singleQuoteCount%2 != 0 {
		return "", false, fmt.Errorf("unmatched single quotes: found %d non-contraction, non-possesive, non-plural or omission digit single quotes", singleQuoteCount)
	}

	return string(runes), updateMade, nil
}
