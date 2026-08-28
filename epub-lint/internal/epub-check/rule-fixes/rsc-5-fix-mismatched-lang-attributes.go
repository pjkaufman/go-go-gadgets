package rulefixes

import (
	"strings"
	"unicode"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
)

func FixMismatchedLangAttributes(line, column int, contents, language string) (edits []positions.TextEdit) {
	if line < 1 || language == "" {
		return
	}

	offset := positions.GetPositionOffset(contents, line, column)
	if offset == -1 || offset > len(contents) {
		return
	}

	openStart := strings.LastIndex(contents[:offset], "<")
	if openStart == -1 {
		return
	}

	tag := contents[openStart:offset]
	if !strings.HasSuffix(tag, ">") {
		return
	}

	langStart, langEnd, foundLang := findAttributeValue(tag, "lang")
	xmlLangStart, xmlLangEnd, foundXMLLang := findAttributeValue(tag, "xml:lang")

	if !foundLang || !foundXMLLang {
		return
	}

	edits = append(edits,
		positions.TextEdit{
			Range: positions.Range{
				Start: positions.IndexToPosition(contents, openStart+langStart),
				End:   positions.IndexToPosition(contents, openStart+langEnd),
			},
			NewText: language,
		},
		positions.TextEdit{
			Range: positions.Range{
				Start: positions.IndexToPosition(contents, openStart+xmlLangStart),
				End:   positions.IndexToPosition(contents, openStart+xmlLangEnd),
			},
			NewText: language,
		},
	)

	return
}

func findAttributeValue(tag, attributeName string) (start, end int, found bool) {
	i := strings.IndexByte(tag, '<')
	if i == -1 {
		return
	}
	i++

	// Skip the element name.
	for i < len(tag) &&
		!unicode.IsSpace(rune(tag[i])) &&
		tag[i] != '>' {
		i++
	}

	for i < len(tag) {
		for i < len(tag) && unicode.IsSpace(rune(tag[i])) {
			i++
		}

		if i >= len(tag) || tag[i] == '>' {
			break
		}

		nameStart := i
		for i < len(tag) &&
			!unicode.IsSpace(rune(tag[i])) &&
			tag[i] != '=' &&
			tag[i] != '>' {
			i++
		}

		name := tag[nameStart:i]

		for i < len(tag) && unicode.IsSpace(rune(tag[i])) {
			i++
		}

		if i >= len(tag) || tag[i] != '=' {
			continue
		}
		i++

		for i < len(tag) && unicode.IsSpace(rune(tag[i])) {
			i++
		}

		if i >= len(tag) {
			break
		}

		quote := tag[i]
		if quote != '"' && quote != '\'' {
			return
		}

		valueStart := i + 1
		valueEnd := strings.IndexByte(tag[valueStart:], quote)
		if valueEnd == -1 {
			return
		}
		valueEnd += valueStart

		if strings.EqualFold(name, attributeName) {
			return valueStart, valueEnd, true
		}

		i = valueEnd + 1
	}

	return
}
