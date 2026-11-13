package potentiallyfixableissue

import (
	"regexp"
	"unicode"
)

var paragraphsWithPotentialNestedQuotes = regexp.MustCompile(`(?m)^([\r\t\f\v ]*?<p[^\n>]*?>)([^\n]*)(</p>)`)

// GetPotentialIncorrectNestedDoubleQuotes finds paragraphs where double
// quotes appear nested and suggests using single quotes for the inner quotes.
// For example, "He said, "hello" to her." becomes
// "He said, 'hello' to her.".
func GetPotentialIncorrectNestedDoubleQuotes(fileContent string) (map[string]string, error) {
	matches := paragraphsWithPotentialNestedQuotes.FindAllStringSubmatch(fileContent, -1)
	suggestions := make(map[string]string)

	for _, groups := range matches {
		replaced, changed := convertNestedDoubleQuotes(groups[2])
		if changed {
			suggestions[groups[0]] = groups[1] + replaced + groups[3]
		}
	}

	return suggestions, nil
}

// convertNestedDoubleQuotes converts the inner quote marks in runs of at
// least four double quotes. A run is broken by whitespace-only text between
// adjacent quotes, so separate quotes such as "What?" "Why?" are ignored.
// Quote characters in HTML tags and attributes are not considered content.
func convertNestedDoubleQuotes(input string) (string, bool) {
	runes := []rune(input)
	quotePositions := make([]int, 0, 8)
	insideTag := false

	for i, r := range runes {
		switch r {
		case '<':
			insideTag = true
		case '>':
			insideTag = false
		case '"':
			if !insideTag {
				quotePositions = append(quotePositions, i)
			}
		}
	}

	if len(quotePositions) < 4 {
		return input, false
	}

	changed := false
	for i := 0; i < len(quotePositions); {
		end := i
		for end+1 < len(quotePositions) && containsNonWhitespace(runes[quotePositions[end]+1:quotePositions[end+1]]) {
			end++
		}

		if end-i+1 >= 4 {
			for inner := i + 1; inner < end; inner++ {
				runes[quotePositions[inner]] = '\''
			}
			changed = true
		}
		i = end + 1
	}

	if !changed {
		return input, false
	}
	return string(runes), true
}

func containsNonWhitespace(runes []rune) bool {
	for _, r := range runes {
		if !unicode.IsSpace(r) {
			return true
		}
	}
	return false
}
