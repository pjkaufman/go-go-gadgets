package epubhandler

import (
	"fmt"
	"slices"
	"strings"
)

const (
	spineStartTag = "<spine"
	spineEndTag   = "</spine>"
)

var ErrNoSpine = fmt.Errorf("spine tag not found in OPF contents")

func RemoveIdFromSpine(opfContents, fileId string) (string, error) {
	startIndex, endIndex, spineContent, err := getSpineContents(opfContents)
	if err != nil {
		return "", err
	}

	lines := strings.Split(spineContent, "\n")
	var idRef = fmt.Sprintf(`idref="%s"`, fileId)

	for i, line := range lines {
		if strings.Contains(line, idRef) {
			// start iterating all idref els
			var (
				found          = false
				lineSubset     = line
				startOfItemref int
			)
			for startOfItemref != -1 {
				startOfItemref = strings.Index(lineSubset, itemRefStartTag)
				if startOfItemref == -1 {
					break
				}

				// for now we will assume the itemrefs are self-closing
				var endOfItemref = strings.Index(lineSubset, "/>")
				if endOfItemref == -1 {
					return "", fmt.Errorf("failed to parse itemref out of line contents %q due to missing %q", lineSubset, "/>")
				}

				var (
					endOfItemRef = endOfItemref + 2
					itemrefEl    = lineSubset[startOfItemref:endOfItemRef]
				)
				if strings.Contains(itemrefEl, idRef) {
					line = strings.Replace(line, itemrefEl, "", 1)

					if strings.TrimSpace(line) == "" {
						lines = slices.Delete(lines, i, i+1)
					} else {
						lines[i] = line
					}

					found = true
					break
				}

				lineSubset = lineSubset[endOfItemRef:]
			}

			if found {
				break
			}
		}
	}

	updatedSpineContent := strings.Join(lines, "\n")
	updatedOpfContents := opfContents[:startIndex+len(spineStartTag)] + updatedSpineContent + opfContents[endIndex:]

	return updatedOpfContents, nil
}

func getSpineContents(opfContents string) (int, int, string, error) {
	startIndex := strings.Index(opfContents, spineStartTag)
	endIndex := strings.Index(opfContents, spineEndTag)

	if startIndex == -1 || endIndex == -1 {
		return 0, 0, "", ErrNoSpine
	}

	return startIndex, endIndex, opfContents[startIndex+len(spineStartTag) : endIndex], nil
}
