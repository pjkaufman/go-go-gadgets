package linter

import (
	"fmt"
	"strings"
)

const (
	metadataStartTag = "<metadata"
	metadataEndTag   = "</metadata>"
)

var ErrNoMetadata = fmt.Errorf("metadata tag not found in OPF contents")

func FixMissingUniqueIdentifierId(opfContents string, id string) (string, error) {
	startIndex, endIndex, manifestContent, err := getMetadataContents(opfContents)
	if err != nil {
		return "", err
	}

	lines := strings.Split(manifestContent, "\n")

	for i, line := range lines {
		if strings.Contains(line, "<dc:identifier") {
			if strings.Contains(line, "id=") {
				continue
			}

			closeTagIndex := strings.Index(line, ">")
			if closeTagIndex == -1 {
				continue
			}

			lines[i] = line[:closeTagIndex] + ` id="` + id + `"` + line[closeTagIndex:]
			break
		}
	}

	updatedManifestContent := strings.Join(lines, "\n")
	updatedOpfContents := opfContents[:startIndex+len(metadataStartTag)] + updatedManifestContent + opfContents[endIndex:]

	return updatedOpfContents, nil
}

func getMetadataContents(opfContents string) (int, int, string, error) {
	startIndex := strings.Index(opfContents, metadataStartTag)
	endIndex := strings.Index(opfContents, metadataEndTag)

	if startIndex == -1 || endIndex == -1 {
		return 0, 0, "", ErrNoMetadata
	}

	return startIndex, endIndex, opfContents[startIndex+len(metadataStartTag) : endIndex], nil
}
