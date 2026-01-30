package rulefixes

import (
	"strings"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
)

func FixIrregularDoctype(contents, expectedDoctype string) (edit positions.TextEdit) {
	doctypeStart := strings.Index(contents, "<!DOCTYPE")
	if doctypeStart == -1 {
		return // no doctype tag
	}

	doctypeEnd := strings.Index(contents[doctypeStart:], ">")
	if doctypeEnd == -1 {
		return // failed to find the end of the doctype for some reason
	}

	edit.NewText = expectedDoctype
	edit.Range.Start = positions.IndexToPosition(contents, doctypeStart)
	edit.Range.End = positions.IndexToPosition(contents, doctypeStart+doctypeEnd+1)

	return
}
