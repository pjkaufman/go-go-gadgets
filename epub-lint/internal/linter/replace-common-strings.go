package linter

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
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
		Search:   "‘",
		Replace:  "'",
		Rational: "Replace smart single quotes with straight single quotes",
	},
	{
		Search:   "’",
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
	{
		Search:   "--",
		Replace:  "—",
		Rational: "An em dash should be used instead of two consecutive regular dashes",
	},
}

// CommonStringReplace applies common string replacements to character data
// in XHTML/XML input.
//
// The original source is preserved byte-for-byte except for character data
// that is intentionally modified by applyStringReplacements.
//
// Whitespace-only character data is preserved exactly. Text inside script,
// style, code, and pre elements is also preserved exactly.
func CommonStringReplace(text string) (string, error) {
	input := []byte(text)
	decoder := xml.NewDecoder(bytes.NewReader(input))

	// EPUB XHTML should generally be valid XML. Keep this false if the
	// existing linter intentionally needs to tolerate malformed input.
	decoder.Strict = false

	var out bytes.Buffer
	out.Grow(len(input))

	var stringsToReplace = make([]string, 2*len(commonReplaceWords))
	for i, replaceWord := range commonReplaceWords {
		stringsToReplace[2*i] = replaceWord.Search
		stringsToReplace[2*i+1] = replaceWord.Replace
	}

	replacer := strings.NewReplacer(stringsToReplace...)

	skipTags := map[string]bool{
		"script": true,
		"style":  true,
		"code":   true,
		"pre":    true,
	}

	var skipStack []string
	var tokenStart int

	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			if tokenStart < len(input) {
				out.Write(input[tokenStart:])
			}

			return out.String(), nil
		}

		if err != nil {
			return "", fmt.Errorf("failed to decode XML content: %w", err)
		}

		tokenEnd := int(decoder.InputOffset())

		switch token := token.(type) {
		case xml.CharData:
			raw := input[tokenStart:tokenEnd]

			// Preserve whitespace-only character data exactly as it appeared
			// in the source. This includes indentation and whitespace between
			// elements, e.g. the two spaces before <p>.
			if strings.TrimSpace(string(token)) == "" {
				out.Write(raw)
				break
			}

			// Preserve character data inside skip elements.
			if len(skipStack) > 0 {
				out.Write(raw)
				break
			}

			// Apply replacements to the decoded character data.
			//
			// If nothing changed, use the original bytes. This ensures that
			// entities and other source representation remain untouched when
			// there is no actual replacement.
			replaced := applyStringReplacements(string(token), replacer)

			if replaced == string(token) {
				out.Write(raw)
			} else {
				out.WriteString(replaced)
			}

		case xml.StartElement:
			// Copy the exact original opening tag.
			out.Write(input[tokenStart:tokenEnd])

			name := strings.ToLower(token.Name.Local)
			if skipTags[name] {
				skipStack = append(skipStack, name)
			}

		case xml.EndElement:
			// Copy the exact original closing tag.
			out.Write(input[tokenStart:tokenEnd])

			name := strings.ToLower(token.Name.Local)

			if len(skipStack) > 0 && skipStack[len(skipStack)-1] == name {
				skipStack = skipStack[:len(skipStack)-1]
			}

		default:
			// Comments, directives, processing instructions, CDATA-related
			// tokens, etc. are copied exactly as they appeared in the input.
			out.Write(input[tokenStart:tokenEnd])
		}

		tokenStart = tokenEnd
	}
}

func applyStringReplacements(text string, replacer *strings.Replacer) string {
	// Replace multiple spaces in a row between words with a single space
	// since this can cause issues with replace strings.
	var newText = replaceTwoPlusSpacesBetweenWords(text)

	return replacer.Replace(newText)
}

func replaceTwoPlusSpacesBetweenWords(text string) string {
	var index = strings.Index(text, doubleSpace)
	if index == -1 {
		return text
	}

	var newText strings.Builder
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

		if startWhitespace > 0 &&
			(text[startWhitespace-1] == '\n' ||
				text[startWhitespace-1] == '\t') {
			newText.WriteString(text[0 : index+2])
		} else if endingWhitespace+1 < len(text) &&
			(text[endingWhitespace+1] == '<' ||
				text[endingWhitespace+1] == '\n') {
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
