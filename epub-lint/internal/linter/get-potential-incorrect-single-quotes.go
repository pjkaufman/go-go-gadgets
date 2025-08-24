package linter

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var contractions = map[string]struct{}{
	"a'ight": {}, "ain't": {}, "amn't": {}, "aren't": {}, "'bout": {}, "can't": {}, "cap'n": {},
	"'cause": {}, "'cept": {}, "c'mon": {}, "could've": {}, "couldn't": {}, "couldn't've": {},
	"daren't": {}, "daresn't": {}, "dasn't": {}, "didn't": {}, "doesn't": {}, "don't": {},
	"d'ye": {}, "d'ya": {}, "e'en": {}, "e'er": {}, "'em": {}, "everybody's": {}, "everyone's": {},
	"everything's": {}, "fo'c'sle": {}, "'gainst": {}, "g'day": {}, "giv'n": {}, "gi'z": {},
	"gon't": {}, "hadn't": {}, "had've": {}, "hasn't": {}, "haven't": {}, "he'd": {}, "he'd'nt've": {},
	"he'll": {}, "yesn't": {}, "he's": {}, "here's": {}, "how'd": {}, "how'll": {}, "how're": {},
	"how's": {}, "i'd": {}, "i'd've": {}, "i'd'nt": {}, "i'd'nt've": {}, "if'n": {}, "i'll": {},
	"i'm": {}, "i'm'onna": {}, "i'm'o": {}, "i'm'na": {}, "i've": {}, "isn't": {}, "it'd": {},
	"it'll": {}, "it's": {}, "let's": {}, "loven't": {}, "ma'am": {}, "mayn't": {}, "may've": {},
	"mightn't": {}, "might've": {}, "mine's": {}, "mustn't": {}, "mustn't've": {}, "must've": {},
	"'neath": {}, "needn't": {}, "ne'er": {}, "nothing's": {}, "o'clock": {}, "o'er": {}, "ol'": {},
	"ought've": {}, "oughtn't": {}, "oughtn't've": {}, "'round": {}, "'s": {}, "shalln't": {},
	"shan'": {}, "shan't": {}, "she'd": {}, "she'll": {}, "she's": {}, "she'd'nt've": {}, "should've": {},
	"shouldn't": {}, "shouldn't've": {}, "somebody's": {}, "someone's": {}, "something's": {},
	"so're": {}, "so's": {}, "so've": {}, "that'll": {}, "that're": {}, "that's": {}, "that'd": {},
	"there'd": {}, "there'll": {}, "there're": {}, "there's": {}, "these're": {}, "these've": {},
	"they'd": {}, "they'd've": {}, "they'll": {}, "they're": {}, "they've": {}, "this's": {},
	"those're": {}, "those've": {}, "'thout": {}, "'til": {}, "'tis": {}, "to've": {}, "'twas": {},
	"'tween": {}, "'twere": {}, "w'all": {}, "w'at": {}, "wasn't": {}, "we'd": {}, "we'd've": {},
	"we'll": {}, "we're": {}, "we've": {}, "weren't": {}, "what'd": {}, "what'll": {}, "what're": {},
	"what's": {}, "what've": {}, "when'd": {}, "when's": {}, "where'd": {}, "where'll": {},
	"where're": {}, "where's": {}, "where've": {}, "which'd": {}, "which'll": {}, "which're": {},
	"which's": {}, "which've": {}, "who'd": {}, "who'd've": {}, "who'll": {}, "who're": {},
	"who's": {}, "who've": {}, "why'd": {}, "why'dja": {}, "why're": {}, "why's": {}, "willn't": {},
	"won't": {}, "would've": {}, "wouldn't": {}, "wouldn't've": {}, "y'ain't": {}, "y'all": {},
	"y'all'd've": {}, "y'all'dn't've": {}, "y'all're": {}, "y'all'ren't": {}, "y'at": {},
	"yes'm": {}, "y'ever": {}, "y'know": {}, "you'd": {}, "you'dn't've": {}, "you'll": {},
	"you're": {}, "you've": {},
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
		runes                             = []rune(input)
		insideDoubleQuotes                = false
		doubleQuoteCount                  = 0
		singleQuoteCount                  = 0 // Only counts non-possesive, non-contraction, and non-plural or omission digit single quotes
		updateMade                        = false
		checkForContractionAndGetNewStart = func(startIndex int) int {
			var start = startIndex
			for start > 0 && (unicode.IsLetter(runes[start-1]) || runes[start-1] == '\'') {
				start--
			}

			var end = startIndex
			for end < len(runes)-1 && (unicode.IsLetter(runes[end+1]) || runes[end+1] == '\'') {
				end++
			}

			if _, ok := contractions[strings.ToLower(string(runes[start:end+1]))]; !ok {
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

			return end + 1
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
			isPrevWord := i > 0 && unicode.IsLetter(runes[i-1])

			// is a plural, possesive, or omitted number scenario
			isDigitScenarios := (isPrevDigit && isNextS) || (!isPrevWord && isNextDigit)
			// we will assume that no possesives show up inside a single quote as that gets hairy and is not valid
			isPossessive := (isPrevS || (isPrevWord && isNextS)) && singleQuoteCount%2 == 0

			if isPossessive || isDigitScenarios {
				continue
			}

			i = checkForContractionAndGetNewStart(i)
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
