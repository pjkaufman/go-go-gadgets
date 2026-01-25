package rulefixes

import (
	"fmt"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
)

const (
	metadataStartTag = "<metadata"
	metadataEndTag   = "</metadata>"
)

var ErrNoMetadata = fmt.Errorf("metadata tag not found in OPF contents")

func FixMissingUniqueIdentifierId(opfContents string, id string) (positions.TextEdit, error) {
	var edit positions.TextEdit
	startIndex, _, manifestContent, err := getMetadataContents(opfContents)
	if err != nil {
		return edit, err
	}

	var (
		remainingManifestContent                 = manifestContent
		identifierStartIndex, identifierEndIndex int
		insertIndex                              = startIndex
	)

	for identifierStartIndex != -1 {
		identifierStartIndex = strings.Index(remainingManifestContent, "<dc:identifier")
		if identifierStartIndex == -1 {
			return edit, nil
		}

		insertIndex += identifierStartIndex
		identifierEndIndex = strings.Index(remainingManifestContent[identifierStartIndex:], ">")
		if identifierEndIndex == -1 {
			remainingManifestContent = remainingManifestContent[identifierStartIndex:]
			continue // something is wrong, so ignore this element
		}

		insertIndex += identifierEndIndex
		identifierOpeningEl := remainingManifestContent[identifierStartIndex : identifierStartIndex+identifierEndIndex]
		if strings.Contains(identifierOpeningEl, "id=") {
			remainingManifestContent = remainingManifestContent[identifierStartIndex+identifierEndIndex:]
			continue
		}

		insertIdPos := positions.IndexToPosition(opfContents, insertIndex)
		edit.Range.Start = insertIdPos
		edit.Range.End = insertIdPos
		edit.NewText = ` id="` + id + `"`
		return edit, nil
	}

	return edit, nil
}

func getMetadataContents(opfContents string) (int, int, string, error) {
	startIndex := strings.Index(opfContents, metadataStartTag)
	endIndex := strings.Index(opfContents, metadataEndTag)

	if startIndex == -1 || endIndex == -1 {
		return 0, 0, "", ErrNoMetadata
	}

	return startIndex, endIndex, opfContents[startIndex+len(metadataStartTag) : endIndex], nil
}
