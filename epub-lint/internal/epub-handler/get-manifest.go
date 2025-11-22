package epubhandler

import "strings"

const (
	ManifestStartTag = "<manifest>"
	ManifestEndTag   = "</manifest>"
)

func GetManifestContents(opfContents string) (int, int, string, error) {
	startIndex := strings.Index(opfContents, ManifestStartTag)
	endIndex := strings.Index(opfContents, ManifestEndTag)

	if startIndex == -1 || endIndex == -1 {
		return 0, 0, "", ErrNoManifest
	}

	return startIndex, endIndex, opfContents[startIndex+len(ManifestStartTag) : endIndex], nil
}
