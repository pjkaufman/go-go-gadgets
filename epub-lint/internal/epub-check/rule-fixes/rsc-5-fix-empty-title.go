package rulefixes

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
)

// FixEmptyTitle updates the title element's contents to be the first heading's text
// or the first paragraph's text if there are no headings
func FixEmptyTitle(line, column int, contents string) (edit positions.TextEdit) {
	offset := positions.GetPositionOffset(contents, line, column)
	if offset == -1 {
		return
	}

	title := firstHeaderOrParagraphText(contents)
	if strings.TrimSpace(title) == "" {
		return
	}

	var (
		newText       = title
		startPosition = positions.Position{
			Line:   line,
			Column: column,
		}
		endPosition = startPosition
	)
	if contents[offset-2] == '/' { // we are dealing with a self-closing title
		startIndex := strings.LastIndex(contents[:offset], "<")
		if startIndex == -1 { // something went wrong, so we will skip this one
			return
		}

		startPosition = positions.IndexToPosition(contents, startIndex)
		newText = fmt.Sprintf("<title>%s</title>", title)
	}

	edit.Range.Start = startPosition
	edit.Range.End = endPosition
	edit.NewText = newText

	return
}

type textNode struct {
	Text string `xml:",chardata"`
}

func firstHeaderOrParagraphText(contents string) string {
	var (
		isFirstParagraph   = true
		firstParagraphText string
		decoder            = xml.NewDecoder(strings.NewReader(contents))
	)

	for {
		tok, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return firstParagraphText
			}

			return ""
		}

		switch se := tok.(type) {
		case xml.StartElement:
			if se.Name.Local == "h1" || se.Name.Local == "h2" ||
				se.Name.Local == "h3" || se.Name.Local == "h4" ||
				se.Name.Local == "h5" || se.Name.Local == "h6" {
				var h textNode
				decoder.DecodeElement(&h, &se)

				return strings.TrimSpace(h.Text)
			} else if se.Name.Local == "p" && isFirstParagraph {
				var h textNode
				decoder.DecodeElement(&h, &se)

				firstParagraphText = strings.TrimSpace(h.Text)
				isFirstParagraph = false
			}
		}
	}
}
