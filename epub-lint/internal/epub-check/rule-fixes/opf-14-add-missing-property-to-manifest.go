package rulefixes

import (
	"errors"
	"path"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
)

var ErrNoManifest = errors.New("manifest tag not found in OPF contents")

func AddPropertyToManifest(opfContents string, opfPath string, fileName string, property string) (positions.TextEdit, error) {
	var edit positions.TextEdit

	startIndex, _, manifestContent, err := epubhandler.GetManifestContents(opfContents)
	if err != nil {
		return edit, err
	}

	targetPath := normalizeEPUBPath(fileName)
	opfDir := path.Dir(normalizeEPUBPath(opfPath))

	var (
		element        string
		startOfElement int
		endOfElement   int
	)

	for searchStart := 0; searchStart < len(manifestContent); {
		itemStart := strings.Index(manifestContent[searchStart:], "<item")
		if itemStart == -1 {
			break
		}
		itemStart += searchStart

		itemEnd := strings.Index(manifestContent[itemStart:], "/>")
		if itemEnd == -1 {
			break
		}
		itemEnd += itemStart + 2

		candidate := manifestContent[itemStart:itemEnd]

		href, _, _, err := epubhandler.GetAttributeValue(candidate, "href")
		if err == nil {
			resolvedHref := path.Join(opfDir, href)

			if normalizeEPUBPath(resolvedHref) == targetPath {
				element = candidate
				startOfElement = itemStart
				endOfElement = itemEnd - 2
				break
			}
		}

		searchStart = itemEnd
	}

	if element == "" {
		return edit, nil
	}

	_, propertiesStart, _, err := epubhandler.GetAttributeValue(element, "properties")
	if err == nil {
		insertPropertiesPos := positions.IndexToPosition(
			opfContents,
			startIndex+startOfElement+propertiesStart,
		)

		newText := property
		if element[propertiesStart] != '"' {
			newText += " "
		}

		return positions.TextEdit{
			Range: positions.Range{
				Start: insertPropertiesPos,
				End:   insertPropertiesPos,
			},
			NewText: newText,
		}, nil
	}

	insertPropertiesPos := positions.IndexToPosition(
		opfContents,
		startIndex+endOfElement,
	)

	return positions.TextEdit{
		Range: positions.Range{
			Start: insertPropertiesPos,
			End:   insertPropertiesPos,
		},
		NewText: ` properties="` + property + `"`,
	}, nil
}

func normalizeEPUBPath(p string) string {
	p = strings.ReplaceAll(p, `\`, "/")
	p = strings.TrimPrefix(p, "/")

	return path.Clean(p)
}
