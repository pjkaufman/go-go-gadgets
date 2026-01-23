package rulefixes

import (
	"strings"
)

func RemoveLinkId(contents string, lineToUpdate, startOfFragment int) (edit TextEdit) {
	if lineToUpdate < 1 {
		return
	}

	offset := getColumnOffset(contents, lineToUpdate, startOfFragment)
	if offset == -1 {
		return
	}

	// Work backwards from startOfFragment to find href or src attribute
	var (
		hrefAttributeIndicator = "href="
		srcAttributeIndicator  = "src="
		fragmentStart          = strings.LastIndex(contents[:startOfFragment], hrefAttributeIndicator)
		srcStart               = strings.LastIndex(contents[:startOfFragment], srcAttributeIndicator)
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

	edit = TextEdit{
		Range: Range{
			Start: Position{
				Line:   lineToUpdate,
				Column: getColumnFromIndex(contents, lineToUpdate, startAttr+1+idIndicatorStart+1),
			},
			End: Position{
				Line:   lineToUpdate,
				Column: getColumnFromIndex(contents, lineToUpdate, endOfFragment+1),
			},
		},
		NewText: "",
	}

	return
}
