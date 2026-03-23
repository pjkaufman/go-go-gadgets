package linter

import (
	"encoding/xml"
	"fmt"
	"io"
	"slices"
	"strings"
)

// these values are lowercased because that makes the checks later on more performant since we don't need
// to lowercase them
var noteIndicators = []string{"tl note:", "translator's note:", "t/n:", "author's note:", "note:"}

func GetTranslatorsNotes(text, fileName, noteFileName string, startingNoteNumber int) (string, []string, int) {
	matches := findNotesWithXML(text)
	if len(matches) == 0 {
		return text, []string{}, startingNoteNumber
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
	return text, tlNotes, startingNoteNumber
}

type noteMatch struct {
	Start   int
	End     int
	Content string
}

func findNotesWithXML(text string) []noteMatch {
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
			innerContent, textOnlyContent := getInnerContent(decoder)
			indicator, tlNotePos := translatorNoteIndicatorPosInfo(innerContent)
			if tlNotePos == -1 {
				continue
			}

			// Find positions in original text
			var (
				startPos = strings.Index(text, innerContent)
				endPos   = startPos + len(innerContent)
			)

			matches = append(matches, noteMatch{
				Start:   startPos,
				End:     endPos,
				Content: strings.TrimSpace(extractNoteContent(indicator, innerContent, strings.TrimSpace(textOnlyContent), tlNotePos)),
			})
		}
	}

	return matches
}

func getInnerContent(decoder *xml.Decoder) (string, string) {
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

	return content.String(), textOnly.String()
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

func extractNoteContent(indicator, innerElContent, textOnlyContent string, indicatorPos int) string {
	var (
		startOfNote     = indicatorPos + len(indicator)
		startOfTextNote = strings.Index(strings.ToLower(textOnlyContent), indicator)
	)

	// If indicator at start, return all
	if startOfTextNote == 0 {
		if innerElContent[startOfNote] == ' ' {
			startOfNote++
		}

		return innerElContent[:indicatorPos] + innerElContent[startOfNote:]
	}

	beforeIndicator := innerElContent[:indicatorPos]
	afterIndicator := innerElContent[startOfNote:]

	// Has opening paren?
	if strings.Contains(beforeIndicator, "(") {
		openCount := strings.Count(beforeIndicator, "(") - strings.Count(beforeIndicator, ")")
		return extractUntilBalancedParens(afterIndicator, openCount)
	}

	// No paren - until next tag
	endIdx := strings.Index(afterIndicator, "<")
	if endIdx == -1 {
		return strings.TrimSpace(afterIndicator)
	}

	return strings.TrimSpace(afterIndicator[:endIdx])
}

func extractUntilBalancedParens(s string, openCount int) string {
	closeCount := 0
	for i, ch := range s {
		if ch == ')' {
			closeCount++

			if closeCount >= openCount {
				return strings.TrimSpace(s[:i])
			}
		}
	}

	return strings.TrimSpace(s)
}
