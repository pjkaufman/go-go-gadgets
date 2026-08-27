package rulefixes

import (
	"errors"
	"net/url"
	"path"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
)

var ErrNoManifest = errors.New("manifest tag not found in OPF contents")

func AddPropertyToManifest(opfContents, fileName, property string) (positions.TextEdit, error) {
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

	_, propertiesStart, _, err :=
		epubhandler.GetAttributeValue(item.element, "properties")

	if err == nil {
		insertPropertiesPos := positions.IndexToPosition(
			opfContents,
			startIndex+item.startOfElement+propertiesStart,
		)

		newText := property
		if item.element[propertiesStart] != '"' {
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
		startIndex+item.endOfElement,
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

type manifestItem struct {
	element        string
	startOfElement int
	endOfElement   int // position immediately before "/>"
}

// findManifestItem finds an <item> whose href resolves to targetPath.
//
// All offsets returned by this function are offsets into manifestContent,
// which is the original, undecoded XML source.
func findManifestItem(manifestContent, targetPath string) (manifestItem, bool) {
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

		element := manifestContent[itemStart:itemEnd]

		href, _, _, err := epubhandler.GetAttributeValue(element, "href")
		if err == nil {
			// href is a URI, so decode it before resolving the EPUB path.
			if decodedHref, err := url.PathUnescape(href); err == nil {
				href = decodedHref
			}

			if normalizeEPUBPath(href) == targetPath {
				return manifestItem{
					element:        element,
					startOfElement: itemStart,
					endOfElement:   itemEnd - 2, // immediately before "/>"
				}, true
			}
		}

		searchStart = itemEnd
	}

	return manifestItem{}, false
}

// propertyAttributeEdit returns the edit required to remove property from
// the "properties" attribute of item.
//
// elementStart is the absolute offset of the element in opfContents.
func propertyAttributeEdit(opfContents string, element string, elementStart int, property string) (positions.TextEdit, bool) {
	properties, valueStart, attributeEnd, err :=
		epubhandler.GetAttributeValue(element, "properties")
	if err != nil || valueStart == -1 {
		return positions.TextEdit{}, false
	}

	valueStart += elementStart

	// Treat properties as whitespace-separated tokens.
	propertyIndex := -1
	for start := 0; start < len(properties); {
		for start < len(properties) && isXMLWhitespace(properties[start]) {
			start++
		}
		if start >= len(properties) {
			break
		}

		end := start
		for end < len(properties) && !isXMLWhitespace(properties[end]) {
			end++
		}

		if properties[start:end] == property {
			propertyIndex = start
			break
		}

		start = end
	}

	if propertyIndex == -1 {
		return positions.TextEdit{}, false
	}

	var start, end int

	switch {
	case strings.TrimSpace(properties) == property:
		// Remove the entire properties attribute, including the preceding
		// space:
		//
		//   ... media-type="..." properties="scripted"/>
		//                         ^^^^^^^^^^^^^^^^^^^
		start = elementStart + valueStart - elementStart - len("properties") - 3
		end = elementStart + attributeEnd + 1

	case propertyIndex == 0:
		// Remove the property and the following space.
		start = valueStart + propertyIndex
		end = start + len(property) + 1

	default:
		// Remove the preceding space and the property.
		start = valueStart + propertyIndex - 1
		end = valueStart + propertyIndex + len(property)
	}

	return positions.TextEdit{
		Range: positions.Range{
			Start: positions.IndexToPosition(opfContents, start),
			End:   positions.IndexToPosition(opfContents, end),
		},
	}, true
}

func isXMLWhitespace(b byte) bool {
	switch b {
	case ' ', '\t', '\r', '\n':
		return true
	default:
		return false
	}
}
