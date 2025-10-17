package linter

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

var translatorsNotesRegex = regexp.MustCompile(`(?i)(TL Note:|Translator's Note:|Note:|T/N:)([^<]*?)(</(p|span|div)>|<br/>)`)

func GetTranslatorsNotes(text, fileName, noteFileName string, startingNoteNumber int) (string, []string, int) {
	var (
		indices = translatorsNotesRegex.FindAllStringSubmatchIndex(text, -1)
		tlNotes = make([]string, len(indices))
	)

	// we go in reverse so we do not mess up the actual contents of the file and positions for future values
	slices.Reverse(indices)

	startingNoteNumber += len(indices)
	noteNum := startingNoteNumber
	for i, loc := range indices {
		var (
			refId      = fmt.Sprintf("note_ref_%d", noteNum)
			noteId     = fmt.Sprintf("tl_note_%d", noteNum)
			noteAnchor = fmt.Sprintf(`<a id="%s" href="%s#%s"><sup>%d</sup></a>`, refId, noteFileName, noteId, noteNum)
			noteText   = strings.TrimSpace(strings.TrimSpace(text[loc[4]:loc[5]]))
		)

		var startIndex = loc[2]
		for startIndex > 0 && text[startIndex-1] != '>' {
			startIndex--
		}

		var wholeText = strings.TrimSpace(text[startIndex:loc[5]])
		if strings.HasPrefix(wholeText, "(") && strings.HasSuffix(wholeText, ")") {
			noteText = noteText[0 : len(noteText)-1]
		}

		tlNotes[i] = fmt.Sprintf(`<li id="%s">%s<br/><a href="%s#%s">Back to Reference</a></li>`+"\n", noteId, noteText, fileName, refId)

		text = text[0:startIndex] + noteAnchor + text[loc[5]:]
		noteNum--
	}

	slices.Reverse(tlNotes)

	return text, tlNotes, startingNoteNumber
}
