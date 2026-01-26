package rulefixes

import (
	"fmt"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
)

func RemovePropertyFromManifest(opfContents, fileName, property string) (positions.TextEdit, error) {
	var edit positions.TextEdit
	startIndex, _, manifestContent, err := epubhandler.GetManifestContents(opfContents)
	if err != nil {
		return edit, err
	}

	var (
		href        = fmt.Sprintf(`href="%s"`, fileName)
		startOfHref = strings.Index(manifestContent, href)
	)
	if startOfHref == -1 {
		return edit, nil
	}

	startOfElement := strings.LastIndex(manifestContent[:startOfHref], "<item")
	if startOfElement == -1 {
		return edit, nil
	}

	endOfElement := strings.Index(manifestContent[startOfHref:], "/>")
	if endOfElement == -1 {
		return edit, nil
	}

	var (
		propertiesAttr         = `properties="`
		element                = manifestContent[startOfElement : startOfHref+endOfElement]
		startOfPropertiesIndex = strings.Index(element, propertiesAttr)
	)

	var removePropertiesStartPos, removePropertiesEndPos positions.Position
	if startOfPropertiesIndex != -1 {
		var startOfAttributeValue = startOfPropertiesIndex + len(propertiesAttr)
		if element[startOfAttributeValue] == '"' {
			return edit, nil
		}

		attributeEnd := strings.Index(element[startOfAttributeValue:], "\"")
		if attributeEnd == -1 {
			return edit, nil
		}

		var (
			attributeValue    = element[startOfAttributeValue : startOfAttributeValue+attributeEnd]
			startOfValueIndex = startIndex + startOfElement + startOfAttributeValue
		)
		if strings.TrimSpace(attributeValue) == property {
			// remove properties attribute and the preceding space
			removePropertiesStartPos = positions.IndexToPosition(opfContents, startIndex+startOfElement+startOfPropertiesIndex-1)
			removePropertiesEndPos = positions.IndexToPosition(opfContents, startIndex+startOfElement+startOfAttributeValue+attributeEnd+1)
		} else {
			propertyIndex := strings.Index(attributeValue, property)
			if propertyIndex == -1 {
				return edit, nil
			}

			if propertyIndex == 0 {
				removePropertiesStartPos = positions.IndexToPosition(opfContents, startOfValueIndex)
				removePropertiesEndPos = positions.IndexToPosition(opfContents, startOfValueIndex+len(property)+1) // remove the following space
			} else {
				removePropertiesStartPos = positions.IndexToPosition(opfContents, startOfValueIndex+propertyIndex-1) // remove the preceding space
				removePropertiesEndPos = positions.IndexToPosition(opfContents, startOfValueIndex+propertyIndex+len(property))
			}
		}
	}

	edit = positions.TextEdit{
		Range: positions.Range{
			Start: removePropertiesStartPos,
			End:   removePropertiesEndPos,
		},
	}

	return edit, nil
}
