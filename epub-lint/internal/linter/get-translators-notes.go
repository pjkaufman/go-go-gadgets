package linter

import (
	"encoding/xml"
	"fmt"
	"io"
	"slices"
	"strings"
	"unicode"
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
		noteAnchor := fmt.Sprintf(`<a id="%s" href="%s#%s"><sup>%d</sup></a>`, refId, noteFileName, noteId, noteNum)

		tlNotes[i] = fmt.Sprintf(`<li id="%s">%s<br/><a href="%s#%s">Back to Reference</a></li>`+"\n",
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
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		switch start := token.(type) {
		case xml.StartElement:
			// Only care about p, div, span
			if start.Name.Local != "p" && start.Name.Local != "div" && start.Name.Local != "span" {
				continue
			}

			// Get inner content
			innerContent, textOnlyContent, encounteredPTag := getInnerContent(decoder)
			// to avoid odd nesting scenarios, only handle the direct parent when possible
			if encounteredPTag { // may want to revise this as there could be text in the div prior to the p tag, but this may handle things
				innerContent, textOnlyContent, _ = getInnerContent(decoder)
			}

			indicator, tlNotePos := translatorNoteIndicatorPosInfo(innerContent)
			if tlNotePos == -1 {
				continue
			}

			// Find positions in original text
			var (
				startPos = strings.Index(text, innerContent)
				endPos   = startPos + len(innerContent)
			)
			if startPos == -1 {
				return nil, fmt.Errorf("attempting to find translator's note text %q failed. This likely means that the source text has html entities. Please convert them to the corresponding character and then try again.", innerContent)
			}

			// check to make sure that the character prior to the tl note is not a letter to help filter out false positives for things like ed
			if tlNotePos != 0 && unicode.IsLetter(rune(innerContent[tlNotePos-1])) {
				continue
			}

			matches = append(matches, extractNoteContent(indicator, innerContent, strings.TrimSpace(textOnlyContent), tlNotePos, startPos, endPos))
		}
	}

	return matches, nil
}

func getInnerContent(decoder *xml.Decoder) (string, string, bool) {
	var (
		content  strings.Builder
		textOnly strings.Builder
		depth    = 1
	)

	for depth > 0 {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == "p" {
				return "", "", true
			}

			depth++
			// Write the tag back
			content.WriteString("<" + t.Name.Local)
			for _, attr := range t.Attr {
				content.WriteString(" " + attr.Name.Local + "=\"" + attr.Value + "\"")
			}
			content.WriteString(">")
		case xml.EndElement:
			depth--
			if depth > 0 {
				content.WriteString("</" + t.Name.Local + ">")
			}
		case xml.CharData:
			content.Write(t)
			textOnly.Write(t)
		}
	}

	return content.String(), textOnly.String(), false
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
		if innerElContent[startOfNote] == ' ' {
			startOfNote++
		}

		match.Content = strings.TrimSpace(innerElContent[:indicatorPos] + innerElContent[startOfNote:])

		return
	}

	beforeIndicator := innerElContent[:indicatorPos]
	afterIndicator := innerElContent[startOfNote:]

	// Has opening paren?
	if strings.Contains(beforeIndicator, "(") {
		var (
			isInOpeningParen bool
			char             rune
			priorChars       = []rune(beforeIndicator)
		)
		for i := len(priorChars) - 1; i >= 0; i-- {
			char = priorChars[i]
			if char == '(' {
				isInOpeningParen = true
				match.Start += i
				break
			}

			if !unicode.IsSpace(char) {
				break
			}
		}

		if isInOpeningParen {
			var (
				openCount  = 1
				closeCount = 0
			)
			for i, ch := range afterIndicator {
				switch ch {
				case ')':
					closeCount++

					if closeCount >= openCount {
						match.End = startPos + startOfNote + i + 1
						match.Content = strings.TrimSpace(afterIndicator[:i])

						return
					}
				case '(':
					openCount++
				}
			}

			match.Content = strings.TrimSpace(afterIndicator)

			return
		}
	}

	match.Start += indicatorPos

	// No paren - until next tag
	endIdx := strings.Index(afterIndicator, "<")
	if endIdx == -1 {
		match.End = startPos + startOfNote + endIdx
		match.Content = strings.TrimSpace(afterIndicator)

		return
	}

	match.Content = strings.TrimSpace(afterIndicator[:endIdx])

	return
}
