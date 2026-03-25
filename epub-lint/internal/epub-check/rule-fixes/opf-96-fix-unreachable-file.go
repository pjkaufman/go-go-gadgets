package rulefixes

import (
	"fmt"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
)

func FixUnreachableFile(line, column int, contents string) (edit positions.TextEdit) {
	offset := positions.GetPositionOffset(contents, line, column)
	if offset == -1 {
		fmt.Println("Fail 1")
		return
	}

	const idAttribute = `id="`
	idAttributeIndex := strings.LastIndex(contents[:offset], idAttribute)
	if idAttributeIndex == -1 { // we should not hit this scenario, but it is possible single quotes were used instead...
		fmt.Println("Fail 2")
		return
	}

	idAttributeIndex += len(idAttribute)

	endOfIdIndex := strings.Index(contents[idAttributeIndex:], "\"")
	if endOfIdIndex == -1 {
		fmt.Println("Fail 3")
		return
	}

	endOfIdIndex += idAttributeIndex

	var (
		idref      = fmt.Sprintf(`idref="%s"`, contents[idAttributeIndex:endOfIdIndex])
		idRefIndex = strings.Index(contents[offset:], idref)
	)
	if idRefIndex == -1 {
		fmt.Println("No idref?")
		return
	}

	idRefIndex += offset

	var (
		startOfItemRef = strings.LastIndex(contents[:idRefIndex], "<itemref")
		endOfItemRef   = strings.Index(contents[idRefIndex:], "/>")
	)
	if startOfItemRef == -1 || endOfItemRef == -1 {
		fmt.Println("No itemref?")
		return
	}

	endOfItemRef += idRefIndex

	const linearAttribute = `linear="no"`
	var (
		itemRefEl            = contents[startOfItemRef:endOfItemRef]
		linearAttributeIndex = strings.Index(itemRefEl, linearAttribute)
	)
	if linearAttributeIndex == -1 {
		fmt.Println("No linear?")
		return
	}

	linearAttributeIndex += startOfItemRef
	endOfLinearAttribute := linearAttributeIndex + len(linearAttribute)
	if contents[linearAttributeIndex-1] == ' ' {
		linearAttributeIndex--
	}

	edit.Range.Start = positions.IndexToPosition(contents, linearAttributeIndex)
	edit.Range.End = positions.IndexToPosition(contents, endOfLinearAttribute)

	return
}
