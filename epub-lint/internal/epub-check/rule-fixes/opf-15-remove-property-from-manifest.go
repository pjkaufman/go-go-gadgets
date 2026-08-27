package rulefixes

import (
	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
)

func RemovePropertyFromManifest(opfContents, fileName, property string) (positions.TextEdit, error) {
	var edit positions.TextEdit

	startIndex, _, manifestContent, err := epubhandler.GetManifestContents(opfContents)
	if err != nil {
		return edit, err
	}

	targetPath := normalizeEPUBPath(fileName)
	item, found := findManifestItem(manifestContent, targetPath)
	if !found {
		return edit, nil
	}

	edit, found = propertyAttributeEdit(opfContents, item.element, startIndex+item.startOfElement, property)
	if !found {
		return positions.TextEdit{}, nil
	}

	return edit, nil
}
