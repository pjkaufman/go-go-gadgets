package rulefixes

import "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"

func FixMissingImageAlt(line, column int, contents string) (edit positions.TextEdit) {
	if line < 1 {
		return
	}

	// column is the index of the `>` in `/>`
	offset := positions.GetPositionOffset(contents, line, column)
	if offset == -1 {
		return
	}

	var emptyAlt = "alt=\"\""
	if contents[offset-3] != ' ' && contents[offset-3] != '\t' {
		emptyAlt = " " + emptyAlt
	}

	insertStartPos := positions.IndexToPosition(contents, offset-2)
	edit.Range.Start = insertStartPos
	edit.Range.End = insertStartPos
	edit.NewText = emptyAlt

	return
}
