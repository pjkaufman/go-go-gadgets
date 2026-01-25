package rulefixes

import (
	"fmt"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
)

var ErrNoManifest = fmt.Errorf("manifest tag not found in OPF contents")

func AddPropertyToManifest(opfContents, fileName, property string) (positions.TextEdit, error) {
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
		newText                string
	)

	var insertPropertiesPos positions.Position
	if startOfPropertiesIndex != -1 {
		var startOfAttributeValue = startOfPropertiesIndex + len(propertiesAttr)
		insertPropertiesPos = positions.IndexToPosition(opfContents, startIndex+startOfElement+startOfAttributeValue)

		if element[startOfAttributeValue] == '"' {
			newText = property
		} else {
			newText = property + " "
		}
	} else {
		insertPropertiesPos = positions.IndexToPosition(opfContents, startIndex+startOfHref+endOfElement)
		newText = ` properties="` + property + `"`
	}

	edit = positions.TextEdit{
		Range: positions.Range{
			Start: insertPropertiesPos,
			End:   insertPropertiesPos,
		},
		NewText: newText,
	}

	return edit, nil
}
