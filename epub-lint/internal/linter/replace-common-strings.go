package linter

import (
	"bytes"
	"io"
	"strings"

	"golang.org/x/net/html"
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
		Search:   "`‘`,",
		Replace:  "'",
		Rational: "Replace smart single quotes with straight single quotes",
	},
	{
		Search:   "`’`,",
		Replace:  "'",
		Rational: "Replace smart single quotes with straight single quotes",
	},
	{
		Search:   "...",
		Replace:  "…",
		Rational: "Proper ellipses should be used instead of 3 periods as it keeps things clean and consistent",
	},
	{
		Search:   ". . .",
		Replace:  "…",
		Rational: "Proper ellipses should be used instead of 3 periods with spaces between them as it keeps things clean and consistent",
	},
}

func CommonStringReplace(text string) string {
	// Replace multiple spaces in a row between words with a single space since this can cause issues with replace strings
	var newText = replaceTwoPlusSpacesBetweenWords(text)

	var stringsToReplace = make([]string, 2*len(commonReplaceWords))
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
			newText.WriteString(text[0:index])
			newText.WriteString("—")
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
			newText.WriteString(text[0:startWhitespace])
			newText.WriteString(" ")
		}

		text = text[endingWhitespace+1:]
		index = strings.Index(text, doubleSpace)
	}

	newText.WriteString(text)

	return newText.String()
}

// ReplaceTextNodesInXHTML applies CommonStringReplace only to text nodes in the input
// XHTML/HTML bytes. It streams tokens, preserves the raw bytes of all non-text tokens,
// and writes replaced text directly (no HTML escaping). Replacements are skipped inside
// tags listed in skipTags.
func ReplaceTextNodesInXHTML(input []byte) ([]byte, error) {
	z := html.NewTokenizer(bytes.NewReader(input))
	var out bytes.Buffer

	skipTags := map[string]bool{
		"script": true,
		"style":  true,
		"code":   true,
		"pre":    true,
	}

	var tagStack []string
	push := func(name string) { tagStack = append(tagStack, name) }
	pop := func() {
		if len(tagStack) > 0 {
			tagStack = tagStack[:len(tagStack)-1]
		}
	}
	inSkip := func() bool {
		if len(tagStack) == 0 {
			return false
		}
		return skipTags[strings.ToLower(tagStack[len(tagStack)-1])]
	}

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return out.Bytes(), nil
			}
			return nil, z.Err()

		case html.TextToken:
			if inSkip() {
				// preserve raw token bytes so we keep original entity representation
				out.Write(z.Raw())
			} else {
				// Apply replacements to decoded text and write it directly (no escaping)
				replaced := CommonStringReplace(string(z.Text()))
				out.WriteString(replaced)
			}

		case html.StartTagToken:
			t := z.Token()
			out.Write(z.Raw())
			push(t.Data)

		case html.SelfClosingTagToken:
			out.Write(z.Raw())

		case html.EndTagToken:
			out.Write(z.Raw())
			pop()

		default:
			// comments, doctype, etc.
			out.Write(z.Raw())
		}
	}
}
