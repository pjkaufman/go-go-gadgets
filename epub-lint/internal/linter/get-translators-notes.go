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

func GetTranslatorsNotes(text, fileName, noteFileName string, startingNoteNumber int) (string, []string, int, error) {
	matches, err := findNotesWithXML(text)
	if err != nil {
		return "", []string{}, 0, fmt.Errorf("file %q had issues determining translator's notes: %w", fileName, err)
	}

	if len(matches) == 0 {
		return text, []string{}, startingNoteNumber, nil
	}

	var tlNotes = make([]string, len(matches))
	slices.Reverse(matches)

	startingNoteNumber += len(matches)
	noteNum := startingNoteNumber

	for i, match := range matches {
		refId := fmt.Sprintf("note_ref_%d", noteNum)
		noteId := fmt.Sprintf("tl_note_%d", noteNum)
		noteAnchor := fmt.Sprintf(`<a id=%q href="%s#%s"><sup>%d</sup></a>`, refId, noteFileName, noteId, noteNum)

		tlNotes[i] = fmt.Sprintf(`<li id=%q>%s<br/><a href="%s#%s">Back to Reference</a></li>`+"\n",
			noteId, match.Content, fileName, refId)

		text = text[:match.Start] + noteAnchor + text[match.End:]
		noteNum--
	}

	slices.Reverse(tlNotes)
	return text, tlNotes, startingNoteNumber, nil
}

type noteMatch struct {
	Start   int
	End     int
	Content string
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
			continue
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

			matches = append(matches, extractNoteContent(
				indicator,
				innerContent,
				strings.TrimSpace(textOnlyContent),
				tlNotePos,
				startPos,
				endPos,
			))
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

func extractNoteContent(indicator, innerElContent, textOnlyContent string, indicatorPos, startPos, endPos int) (match noteMatch) {
	var (
		startOfNote     = indicatorPos + len(indicator)
		startOfTextNote = strings.Index(strings.ToLower(textOnlyContent), indicator)
	)

	match.Start = startPos
	match.End = endPos

	// If indicator at start, return all
	if startOfTextNote == 0 {
		if len(innerElContent) <= startOfNote { // no actual content is present
			return
		}

		if innerElContent[startOfNote] == ' ' {
			startOfNote++
		}

		match.Content = strings.TrimSpace(innerElContent[:indicatorPos] + innerElContent[startOfNote:])

		return
	}

	beforeIndicator := innerElContent[:indicatorPos]
	afterIndicator := innerElContent[startOfNote:]

	// Has opening paren?
	var updated bool
	match, updated = updateNoteForOpeningChar(match, beforeIndicator, afterIndicator, '(', ')', startPos, startOfNote)
	if updated {
		return
	}

	// Has opening square bracket?
	match, updated = updateNoteForOpeningChar(match, beforeIndicator, afterIndicator, '[', ']', startPos, startOfNote)
	if updated {
		return
	}

	match.Start += indicatorPos

	// No paren - until next tag
	before, _, ok := strings.Cut(afterIndicator, "<")
	if !ok {
		match.Content = strings.TrimSpace(afterIndicator)

		return
	}

	match.Content = strings.TrimSpace(before)

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
			match.Start += i
			break
		}

		if !unicode.IsSpace(char) {
			break
		}
	}

	if isInOpeningChar {
		var (
			openCount  = 1
			closeCount = 0
		)
		for i, ch := range afterIndicator {
			switch ch {
			case closingChar:
				closeCount++

				if closeCount >= openCount {
					match.End = startPos + startOfNote + i + 1
					match.Content = strings.TrimSpace(afterIndicator[:i])

					return match, true
				}
			case openingChar:
				openCount++
			}
		}

		match.Content = strings.TrimSpace(afterIndicator)

		return match, true
	}

	return match, false
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
