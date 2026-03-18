package rulefixes

import (
	"fmt"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
)

// AddMissingTitle adds the title element to the header contents
// with the text being the first header's text or the first paragraph's
// text if there are no headings
func AddMissingTitle(line, column int, contents string) (edit positions.TextEdit) {
	offset := positions.GetPositionOffset(contents, line, column)
	if offset == -1 {
		return
	}

	title := firstHeaderOrParagraphText(contents)
	if strings.TrimSpace(title) == "" {
		return
	}

	var (
		newText       = fmt.Sprintf("<title>%s</title>", title)
		startPosition = positions.Position{
			Line:   line,
			Column: column,
		}
		endPosition = startPosition
	)
	if contents[offset-2] == '/' { // we are dealing with a self-closing head
		startIndex := strings.LastIndex(contents[:offset], "<")
		if startIndex == -1 { // something went wrong, so we will skip this one
			return
		}

		startPosition = positions.IndexToPosition(contents, startIndex)
		newText = fmt.Sprintf(`<head>%s</head>`, newText)
	} else {
		nextLineStart := strings.Index(contents[offset:], "\n")
		if nextLineStart != -1 {
			nextLineStart += offset + 1

			nextLineEnd := strings.Index(contents[nextLineStart:], "\n")
			var nextLine string
			if nextLineEnd == -1 {
				nextLine = contents[nextLineStart:]
			} else {
				nextLine = contents[nextLineStart : nextLineStart+nextLineEnd]
			}

			indent := getLeadingWhitespace(nextLine)
			newText = "\n" + indent + newText
		}
	}

	edit.Range.Start = startPosition
	edit.Range.End = endPosition
	edit.NewText = newText

	return
}
