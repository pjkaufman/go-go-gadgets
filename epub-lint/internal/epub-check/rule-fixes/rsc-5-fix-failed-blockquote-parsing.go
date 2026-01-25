package rulefixes

import (
	"strings"
	"unicode"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
)

// FixFailedBlockquoteParsing takes a line and column (both 1-based), and
// file contents, and returns updated contents
// It inserts </p> before the ending blockquote element at the given position,
// and <p> after the opening blockquote element corresponding to that blockquote.
func FixFailedBlockquoteParsing(line, column int, contents string) (edits []positions.TextEdit) {
	if line < 1 {
		return
	}

	offset := positions.GetPositionOffset(contents, line, column)
	if offset == -1 {
		return
	}

	closeTag := "</blockquote>"
	openTag := "<blockquote"
	pOpen := "<p>"
	pClose := "</p>"

	if offset < len(closeTag) {
		return
	}

	closeIdx := offset - len(closeTag)

	// Scan backwards for previous non-whitespace chars before </blockquote>
	beforeClose := contents[:closeIdx]
	trimmed := trimRightSpace(beforeClose)
	if !strings.HasSuffix(trimmed, "</span>") &&
		!strings.HasSuffix(trimmed, "/>") &&
		trimmed[len(trimmed)-1] == '>' {
		return
	}

	// Find the opening <blockquote before closeIdx
	openIdx := strings.LastIndex(contents[:closeIdx], openTag)
	if openIdx == -1 {
		return
	}
	openEnd := strings.Index(contents[openIdx:], ">")
	if openEnd == -1 || openIdx+openEnd > closeIdx {
		return
	}
	openEndIdx := openIdx + openEnd

	// Insert <p> after opening tag and </p> before closing tag
	var (
		insertStartTagPos = positions.IndexToPosition(contents, openEndIdx+1)
		insertEndTagPos   = positions.IndexToPosition(contents, closeIdx)
	)
	edits = append(edits, positions.TextEdit{
		Range: positions.Range{
			Start: insertStartTagPos,
			End:   insertStartTagPos,
		},
		NewText: pOpen,
	}, positions.TextEdit{
		Range: positions.Range{
			Start: insertEndTagPos,
			End:   insertEndTagPos,
		},
		NewText: pClose,
	})

	return
}

func trimRightSpace(s string) string {
	i := len(s)
	for i > 0 && unicode.IsSpace(rune(s[i-1])) {
		i--
	}

	return s[:i]
}
