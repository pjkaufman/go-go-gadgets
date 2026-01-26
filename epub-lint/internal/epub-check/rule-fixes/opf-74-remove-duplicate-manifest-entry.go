package rulefixes

import (
	"strings"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
)

func RemoveDuplicateManifestEntry(line, column int, opfContents string) ([]positions.TextEdit, error) {
	offset := positions.GetPositionOffset(opfContents, line, column)
	if offset == -1 {
		return nil, nil
	}

	openItemTag := "<item"
	openIdx := strings.LastIndex(opfContents[:offset], openItemTag)
	if openIdx == -1 {
		return nil, nil
	}

	var edits []positions.TextEdit
	duplicateEl := opfContents[openIdx:offset]

	edits = append(edits, positions.TextEdit{
		Range: positions.Range{
			Start: positions.IndexToPosition(opfContents, openIdx),
			End: positions.Position{
				Line:   line,
				Column: column,
			},
		},
	})

	fileId := epubhandler.ExtractID(duplicateEl)
	if fileId == "" {
		return edits, nil
	}

	update, err := epubhandler.RemoveIdFromSpine(opfContents, fileId)
	if err != nil {
		return nil, err
	}

	return append(edits, update), nil
}
