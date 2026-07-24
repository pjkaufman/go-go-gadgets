package linter

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

// these values are lowercased because that makes the checks later on more performant since we don't need
// to lowercase them
var noteIndicators = []string{"tl note:", "translator's note:", "t/n:", "tn:", "tln:", "tl:", "author's note:", "note:", "ed:"}

type noteMatch struct {
	start   int
	end     int
	content string
	isEmpty bool
}

func GetTranslatorsNotes(text, fileName, noteFileName string, startingNoteNumber int) (string, []string, int, error) {
	matches, err := findNotesWithXML(text)
	if err != nil {
		return "", []string{}, 0, fmt.Errorf("file %q had issues determining translator's notes: %w", fileName, err)
	}

	if len(matches) == 0 {
		return text, []string{}, startingNoteNumber, nil
	}

	slices.Reverse(matches)

	realMatches := matches[:0]

	for _, match := range matches {
		if match.isEmpty {
			text = text[:match.start] + text[match.end:]
			continue
		}

		realMatches = append(realMatches, match)
	}

	matches = realMatches

	if len(matches) == 0 {
		return text, []string{}, startingNoteNumber, nil
	}

	var tlNotes = make([]string, len(matches))

	startingNoteNumber += len(matches)
	noteNum := startingNoteNumber

	for i, match := range matches {
		refId := fmt.Sprintf("note_ref_%d", noteNum)
		noteId := fmt.Sprintf("tl_note_%d", noteNum)
		noteAnchor := fmt.Sprintf(`<a id=%q href="%s#%s"><sup>%d</sup></a>`, refId, noteFileName, noteId, noteNum)

		tlNotes[i] = fmt.Sprintf(`<li id=%q>%s<br/><a href="%s#%s">Back to Reference</a></li>`+"\n",
			noteId, match.content, fileName, refId)

		text = text[:match.start] + noteAnchor + text[match.end:]
		noteNum--
	}

	slices.Reverse(tlNotes)
	return text, tlNotes, startingNoteNumber, nil
}

func findNotesWithXML(text string) ([]noteMatch, error) {
	var matches []noteMatch

	decoder := xml.NewDecoder(strings.NewReader(text))
	decoder.Strict = false

	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to decode xml content: %w", err)
		}

		switch start := token.(type) {
		case xml.StartElement:
			if start.Name.Local != "p" &&
				start.Name.Local != "div" &&
				start.Name.Local != "span" {
				continue
			}

			// The opening tag has already been consumed by the XML decoder.
			// InputOffset points after the token, so use the decoder offset
			// to determine the element boundary.
			elementStart := findElementStart(text, int(decoder.InputOffset()))
			startPos := findElementContentStart(
				text,
				elementStart,
			)

			innerContent, textOnlyContent, endPos, encounteredPTag := getInnerContent(decoder, text)

			if encounteredPTag {
				innerContent, textOnlyContent, endPos, _ = getInnerContent(decoder, text)
			}

			indicator, tlNotePos := translatorNoteIndicatorPosInfo(innerContent)
			if tlNotePos == -1 {
				continue
			}

			if tlNotePos != 0 {
				r, _ := utf8.DecodeLastRuneInString(innerContent[:tlNotePos])

				if unicode.IsLetter(r) {
					continue
				}
			}

			var (
				elementEnd = findFullElementEnd(text, endPos)
				match      = extractNoteContent(
					indicator,
					innerContent,
					strings.TrimSpace(textOnlyContent),
					tlNotePos,
					startPos,
					endPos,
					elementStart,
					elementEnd,
				)
			)
			matches = append(matches, match)
		}
	}

	return matches, nil
}

func getInnerContent(decoder *xml.Decoder, text string) (string, string, int, bool) {
	var (
		textOnly strings.Builder
		depth    = 1
	)

	// The decoder is positioned immediately after the opening element.
	contentStart := int(decoder.InputOffset())

	for depth > 0 {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			continue
		}

		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == "p" {
				return "", "", 0, true
			}

			depth++

		case xml.EndElement:
			depth--

			if depth == 0 {
				contentEnd := findElementEnd(text, int(decoder.InputOffset()))

				rawContent := text[contentStart:contentEnd]

				return rawContent, textOnly.String(), contentEnd, false
			}

		case xml.CharData:
			textOnly.Write(t)
		}
	}

	return "", textOnly.String(), int(decoder.InputOffset()), false
}

func translatorNoteIndicatorPosInfo(text string) (string, int) {
	var (
		lowerText = strings.ToLower(text)
		pos       int
	)
	for _, indicator := range noteIndicators {
		pos = strings.Index(lowerText, indicator)
		if pos != -1 {
			return indicator, pos
		}
	}

	return "", -1
}

func extractNoteContent(indicator, innerElContent, textOnlyContent string, indicatorPos, startPos, endPos, openingElPos, closingElPos int) (match noteMatch) {
	var (
		startOfNote     = indicatorPos + len(indicator)
		startOfTextNote = strings.Index(strings.ToLower(textOnlyContent), indicator)
	)

	match.start = startPos
	match.end = endPos

	// If indicator at start, return all
	if startOfTextNote == 0 {
		if len(innerElContent) <= startOfNote ||
			strings.TrimSpace(innerElContent[startOfNote:]) == "" {
			match.start = openingElPos
			match.end = closingElPos
			match.isEmpty = true

			return
		}

		if innerElContent[startOfNote] == ' ' {
			startOfNote++
		}

		match.content = strings.TrimSpace(innerElContent[:indicatorPos] + innerElContent[startOfNote:])

		return
	}

	var (
		beforeIndicator = innerElContent[:indicatorPos]
		afterIndicator  = innerElContent[startOfNote:]
	)

	// Has opening paren?
	var updated bool
	match, updated = updateNoteForOpeningChar(match, beforeIndicator, afterIndicator, '(', ')', startPos, startOfNote)
	if updated {
		match = normalizeEmptyNote(match, innerElContent, startPos, openingElPos, closingElPos)

		return
	}

	// Has opening square bracket?
	match, updated = updateNoteForOpeningChar(match, beforeIndicator, afterIndicator, '[', ']', startPos, startOfNote)
	if updated {
		match = normalizeEmptyNote(match, innerElContent, startPos, openingElPos, closingElPos)

		return
	}

	match.start += indicatorPos

	// No paren - until next tag
	before, _, ok := strings.Cut(afterIndicator, "<")
	if !ok {
		match.content = strings.TrimSpace(afterIndicator)
	} else {
		match.content = strings.TrimSpace(before)
	}

	if match.content == "" {
		match.isEmpty = true
	}

	return
}

func updateNoteForOpeningChar(match noteMatch, beforeIndicator, afterIndicator string, openingChar, closingChar rune, startPos, startOfNote int) (noteMatch, bool) {
	if !strings.Contains(beforeIndicator, string(openingChar)) {
		return match, false
	}

	var (
		isInOpeningChar bool
		char            rune
		size            int
	)
	for i := len(beforeIndicator); i > 0; {
		char, size = utf8.DecodeLastRuneInString(beforeIndicator[:i])
		i -= size

		if char == openingChar {
			isInOpeningChar = true
			match.start += i
			break
		}

		if !unicode.IsSpace(char) {
			break
		}
	}

	if isInOpeningChar {
		var depth = 1
		for i, ch := range afterIndicator {
			switch ch {
			case closingChar:
				depth--

				if depth <= 0 {
					match.end = startPos + startOfNote + i + 1
					match.content = strings.TrimSpace(afterIndicator[:i])

					return match, true
				}
			case openingChar:
				depth++
			}
		}

		match.content = strings.TrimSpace(afterIndicator)

		return match, true
	}

	return match, false
}

// when the content of a note is empty we need to determine if we are removing the whole element because removing
// the empty translator's note will leave an empty element if it is the entirety of the actual element
func normalizeEmptyNote(match noteMatch, innerElContent string, startPos, openingElPos, closingElPos int) noteMatch {
	if strings.TrimSpace(match.content) != "" {
		return match
	}

	outside := strings.TrimSpace(
		stripTags(
			innerElContent[:match.start-startPos] +
				innerElContent[match.end-startPos:],
		),
	)

	if outside == "" {
		match.start = openingElPos
		match.end = closingElPos
	}

	match.isEmpty = true
	match.content = ""

	return match
}

func findElementStart(text string, offset int) int {
	if offset > len(text) {
		offset = len(text)
	}

	for i := offset - 1; i >= 0; i-- {
		if text[i] == '<' {
			return i
		}
	}

	return offset
}

func findElementContentStart(text string, elementStart int) int {
	end := strings.IndexByte(text[elementStart:], '>')

	if end == -1 {
		return elementStart
	}

	return elementStart + end + 1
}

func findElementEnd(text string, offset int) int {
	if offset > len(text) {
		offset = len(text)
	}

	return strings.LastIndex(text[:offset], "</")
}

func findFullElementEnd(text string, offset int) int {
	if offset > len(text) {
		offset = len(text)
	}

	end := strings.Index(text[offset:], ">")
	if end == -1 {
		return offset
	}

	return offset + end + 1
}

func stripTags(s string) string {
	var b strings.Builder
	inTag := false

	for _, r := range s {
		switch r {
		case '<':
			inTag = true
		case '>':
			inTag = false
		default:
			if !inTag {
				b.WriteRune(r)
			}
		}
	}

	return b.String()
}
