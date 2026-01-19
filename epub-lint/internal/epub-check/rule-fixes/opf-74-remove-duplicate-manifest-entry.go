package rulefixes

import (
	"strings"

	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
)

func RemoveDuplicateManifestEntry(line, column int, opfContents string) (string, error) {
	offset := getColumnOffset(opfContents, line, column)
	if offset == -1 {
		return opfContents, nil
	}

	openItemTag := "<item"
	openIdx := strings.LastIndex(opfContents[:offset], openItemTag)
	if openIdx == -1 {
		return opfContents, nil
	}

	duplicateEl := opfContents[openIdx:offset]
	opfContents = opfContents[:openIdx] + opfContents[offset:]

	fileId := epubhandler.ExtractID(duplicateEl)
	if fileId == "" {
		return opfContents, nil
	}

	return epubhandler.RemoveIdFromSpine(opfContents, fileId)
}
