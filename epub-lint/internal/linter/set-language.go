package linter

import (
	"strings"
	"unicode"
)

const (
	openingHTMLTag = "<html"
)

func EnsureLanguageIsSet(text, lang string) string {
	htmlStart := strings.Index(text, openingHTMLTag)
	if htmlStart == -1 {
		return text
	}

	tagEnd := strings.IndexByte(text[htmlStart:], '>')
	if tagEnd == -1 {
		return text
	}
	tagEnd += htmlStart

	tag := text[htmlStart:tagEnd]

	var out strings.Builder
	out.Grow(len(tag) + 32)

	out.WriteString("<html")

	i := len("<html")

	foundLang := false
	foundXMLLang := false

	for i < len(tag) {
		// Copy whitespace.
		start := i
		for i < len(tag) && unicode.IsSpace(rune(tag[i])) {
			i++
		}
		out.WriteString(tag[start:i])

		if i >= len(tag) {
			break
		}

		// Parse attribute name.
		nameStart := i
		for i < len(tag) &&
			tag[i] != '=' &&
			!unicode.IsSpace(rune(tag[i])) {
			i++
		}

		name := tag[nameStart:i]

		// Copy spaces before '='.
		wsStart := i
		for i < len(tag) && unicode.IsSpace(rune(tag[i])) {
			i++
		}
		out.WriteString(tag[nameStart:wsStart])

		// Boolean attribute.
		if i >= len(tag) || tag[i] != '=' {
			continue
		}

		out.WriteByte('=')
		i++

		// Spaces after '='.
		wsStart = i
		for i < len(tag) && unicode.IsSpace(rune(tag[i])) {
			i++
		}
		out.WriteString(tag[wsStart:i])

		if i >= len(tag) {
			break
		}

		quote := tag[i]
		if quote != '"' && quote != '\'' {
			// Invalid HTML; just copy remainder.
			out.WriteString(tag[i:])
			break
		}

		valueStart := i + 1
		valueEnd := strings.IndexByte(tag[valueStart:], quote)
		if valueEnd == -1 {
			out.WriteString(tag[i:])
			break
		}
		valueEnd += valueStart

		switch name {
		case "lang":
			foundLang = true
			out.WriteByte(quote)
			if strings.TrimSpace(tag[valueStart:valueEnd]) == "" {
				out.WriteString(lang)
			} else {
				out.WriteString(tag[valueStart:valueEnd])
			}
			out.WriteByte(quote)

		case "xml:lang":
			foundXMLLang = true
			out.WriteByte(quote)
			if strings.TrimSpace(tag[valueStart:valueEnd]) == "" {
				out.WriteString(lang)
			} else {
				out.WriteString(tag[valueStart:valueEnd])
			}
			out.WriteByte(quote)

		default:
			out.WriteByte(quote)
			out.WriteString(tag[valueStart:valueEnd])
			out.WriteByte(quote)
		}

		i = valueEnd + 1
	}

	if !foundLang {
		out.WriteString(` lang="`)
		out.WriteString(lang)
		out.WriteByte('"')
	}

	if !foundXMLLang {
		out.WriteString(` xml:lang="`)
		out.WriteString(lang)
		out.WriteByte('"')
	}

	return text[:htmlStart] + out.String() + text[tagEnd:]
}
