package rulefixes

import (
	"strings"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
)

// FixFailedBlockquoteParsing takes a line and column (both 1-based), and
// file contents, and returns updated contents
// It inserts </div> before the ending blockquote element at the given position,
// and <div> after the opening blockquote element corresponding to that blockquote.
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
	divOpen := "<div>"
	divClose := "</div>"

	if offset < len(closeTag) {
		return
	}

	closeIdx := offset - len(closeTag)

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
	// skip blockquotes with blockquotes inside of them, but otherwise proceed with inserting a div
	if strings.Contains(contents[openEndIdx:closeIdx], closeTag) {
		return
	}

	// Insert <div> after opening tag and </div> before closing tag
	var (
		insertStartTagPos = positions.IndexToPosition(contents, openEndIdx+1)
		insertEndTagPos   = positions.IndexToPosition(contents, closeIdx)
	)
	edits = append(edits, positions.TextEdit{
		Range: positions.Range{
			Start: insertStartTagPos,
			End:   insertStartTagPos,
		},
		NewText: divOpen,
	}, positions.TextEdit{
		Range: positions.Range{
			Start: insertEndTagPos,
			End:   insertEndTagPos,
		},
		NewText: divClose,
	})

	return
}
