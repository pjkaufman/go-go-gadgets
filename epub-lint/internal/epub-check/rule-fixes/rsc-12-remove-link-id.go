package rulefixes

import (
	"strings"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
)

func RemoveLinkId(contents string, lineToUpdate, startOfFragment int) (edit positions.TextEdit) {
	if lineToUpdate < 1 {
		return
	}

	offset := positions.GetPositionOffset(contents, lineToUpdate, startOfFragment)
	if offset == -1 {
		return
	}

	// Work backwards from startOfFragment to find href or src attribute
	var (
		hrefAttributeIndicator = "href="
		srcAttributeIndicator  = "src="
		fragmentStart          = strings.LastIndex(contents[:offset], hrefAttributeIndicator)
		srcStart               = strings.LastIndex(contents[:offset], srcAttributeIndicator)
		startAttr              int
	)
	if fragmentStart == -1 && srcStart == -1 {
		return
	} else if fragmentStart > srcStart {
		startAttr = fragmentStart + len(hrefAttributeIndicator)
	} else {
		startAttr = srcStart + len(srcAttributeIndicator)
	}

	var endingQuote = string(contents[startAttr])
	endOfFragment := startAttr + 1 + strings.Index(contents[startAttr+1:], endingQuote)
	if endOfFragment == -1 {
		return
	}

	fragment := contents[startAttr+1 : endOfFragment]
	if strings.Contains(fragment, "\"") || strings.Contains(fragment, "'") {
		return
	}

	idIndicatorStart := strings.Index(fragment, "#")
	if idIndicatorStart == -1 {
		return
	}

	edit = positions.TextEdit{
		Range: positions.Range{
			Start: positions.Position{
				Line:   lineToUpdate,
				Column: positions.GetColumnFromIndex(contents, lineToUpdate, startAttr+1+idIndicatorStart),
			},
			End: positions.Position{
				Line:   lineToUpdate,
				Column: positions.GetColumnFromIndex(contents, lineToUpdate, endOfFragment),
			},
		},
	}

	return
}
