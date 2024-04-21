package linter

import (
	"strings"
)

type ReplaceWords struct {
	Search   string
	Replace  string
	Rational string
}

const (
	emIndicator = "--"
	doubleSpace = "  "
)

var commonReplaceWords = []ReplaceWords{
	{
		Search:   "Sneaked",
		Replace:  "Snuck",
		Rational: "Use snuck instead of sneaked as it is the more commonly used version of the word nowadays",
	},
	{
		Search:   "sneaked",
		Replace:  "snuck",
		Rational: "Use snuck instead of sneaked as it is the more commonly used version of the word nowadays",
	},
	{
		Search:   "“",
		Replace:  "\"",
		Rational: "Replace smart double quotes with straight double quotes",
	},
	{
		Search:   "”",
		Replace:  "\"",
		Rational: "Replace smart double quotes with straight double quotes",
	},
	{
		Search:   "”",
		Replace:  "\"",
		Rational: "Replace smart double quotes with straight double quotes",
	},
	{
		Search:   `‘`,
		Replace:  "'",
		Rational: "Replace smart single quotes with straight single quotes",
	},
	{
		Search:   `’`,
		Replace:  "'",
		Rational: "Replace smart single quotes with straight single quotes",
	},
	{
		Search:   "...",
		Replace:  "…",
		Rational: "Proper ellipses should be used where possible as it keeps things clean and consistentReplace smart single quotes with straight single quotes",
	},
}

func CommonStringReplace(text string) string {

	// Replace multiple spaces in a row between words with a single space since this can cause issues with replace strings
	var newText = replaceTwoPlusSpacesBetweenWords(text)

	var stringsToReplace []string = make([]string, 2*len(commonReplaceWords))
	for i, replaceWord := range commonReplaceWords {
		stringsToReplace[2*i] = replaceWord.Search
		stringsToReplace[2*i+1] = replaceWord.Replace
	}

	var replacer = strings.NewReplacer(stringsToReplace...)
	newText = replacer.Replace(newText)

	return replaceDoubleDashesWithEmDashes(newText)
}

func replaceDoubleDashesWithEmDashes(text string) string {
	var index = strings.Index(text, emIndicator)
	if index == -1 {
		return text
	}

	var newText = strings.Builder{}
	for index != -1 {
		if index > 0 && text[index-1] == '!' {
			newText.WriteString(text[0 : index+2])
		} else if index+2 < len(text) && text[index+2] == '>' {
			newText.WriteString(text[0 : index+2])
		} else {
			newText.WriteString(text[0:index] + "—")
		}

		text = text[index+2:]
		index = strings.Index(text, emIndicator)
	}

	newText.WriteString(text)

	return newText.String()
}

func replaceTwoPlusSpacesBetweenWords(text string) string {
	var index = strings.Index(text, doubleSpace)
	if index == -1 {
		return text
	}

	var newText = strings.Builder{}
	var endingWhitespace, startWhitespace int
	for index != -1 {
		startWhitespace = index
		endingWhitespace = index + 1
		for startWhitespace > 0 && text[startWhitespace-1] == ' ' {
			startWhitespace--
		}

		for endingWhitespace+1 < len(text) && text[endingWhitespace+1] == ' ' {
			endingWhitespace++
		}

		if startWhitespace > 0 && (text[startWhitespace-1] == '\n' || text[startWhitespace-1] == '\t') {
			newText.WriteString(text[0 : index+2])
		} else if endingWhitespace+1 < len(text) && (text[endingWhitespace+1] == '<' || text[endingWhitespace+1] == '\n') {
			newText.WriteString(text[0 : index+2])
		} else {
			newText.WriteString(text[0:startWhitespace] + " ")
		}

		text = text[endingWhitespace+1:]
		index = strings.Index(text, doubleSpace)
	}

	newText.WriteString(text)

	return newText.String()
}
