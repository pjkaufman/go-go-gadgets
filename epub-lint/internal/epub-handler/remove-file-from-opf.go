package epubhandler

import (
	"fmt"
	"slices"
	"strings"
)

const (
	spineStartTag   = "<spine"
	spineEndTag     = "</spine>"
	hrefAttribute   = "href="
	itemStartTag    = "<item"
	itemRefStartTag = "<itemref"
)

var ErrNoSpine = fmt.Errorf("spine tag not found in OPF contents")

func RemoveFileFromOpf(opfContents, fileName string) (string, error) {
	startIndex, endIndex, manifestContent, err := GetManifestContents(opfContents)
	if err != nil {
		return "", err
	}

	lines := strings.Split(manifestContent, "\n")
	var (
		fileIDs []string
	)

	for i, line := range lines {
		// find each item in the line and check if the href is the one you are looking for
		// if it is, remove it from that line. If the line is now blank, remove it.
		var (
			lineSubset  = line
			startOfItem int
		)
		for startOfItem != -1 {
			startOfItem = strings.Index(lineSubset, itemStartTag)
			if startOfItem == -1 {
				break
			}

			// for now we will assume the items are self-closing
			var endOfItem = strings.Index(lineSubset, "/>")
			if endOfItem == -1 {
				return "", fmt.Errorf("failed to parse item out of line contents %q due to missing %q", lineSubset, "/>")
			}

			var (
				itemEndIndex = endOfItem + 2
				itemEl       = lineSubset[startOfItem:itemEndIndex]
			)
			// to make sure we are only operating on the href
			var startOfHref = strings.Index(itemEl, hrefAttribute)
			if startOfHref == -1 {
				lineSubset = lineSubset[itemEndIndex:]
				continue
			}

			startOfHref += len(hrefAttribute)
			var hrefQuote = itemEl[startOfHref : startOfHref+1]
			startOfHref++

			var endOfHrefIndex = strings.Index(itemEl[startOfHref:], hrefQuote)
			if endOfHrefIndex == -1 {
				lineSubset = lineSubset[itemEndIndex:]
				continue // something went wrong, so we have to ignore the el...
			}

			var hrefContent = itemEl[startOfHref : startOfHref+endOfHrefIndex]
			if !strings.HasSuffix(hrefContent, fileName) {
				lineSubset = lineSubset[endOfItem:] // start over with any other items on this line
				continue
			}

			hrefContent = strings.TrimSuffix(hrefContent, fileName)

			// check for a false positive by checking that previous char is not a slash or a quote
			var previousChar rune
			if len(hrefContent) == 0 {
				previousChar = rune(hrefQuote[0])
			} else {
				previousChar = rune(hrefContent[len(hrefContent)-1])
			}
			if previousChar != '\'' && previousChar != '"' && previousChar != '\\' && previousChar != '/' {
				lineSubset = lineSubset[endOfItem:]
				continue
			}

			var fileID = extractID(itemEl)
			if fileID != "" {
				fileIDs = append(fileIDs, fileID)
			}

			line = strings.Replace(line, itemEl, "", 1)

			if strings.TrimSpace(line) == "" {
				lines = slices.Delete(lines, i, i+1)
				break
			}

			lines[i] = line
			lineSubset = lineSubset[endOfItem:]
		}
	}

	updatedManifestContent := strings.Join(lines, "\n")
	updatedOpfContents := opfContents[:startIndex+len(ManifestStartTag)] + updatedManifestContent + opfContents[endIndex:]

	if len(fileIDs) == 0 {
		return updatedOpfContents, nil
	}

	startIndex, endIndex, spineContent, err := getSpineContents(updatedOpfContents)
	if err != nil {
		return "", err
	}

	lines = strings.Split(spineContent, "\n")
	for _, fileID := range fileIDs {
		var idRef = fmt.Sprintf(`idref="%s"`, fileID)

		for i, line := range lines {
			if strings.Contains(line, idRef) {
				// start iterating all idref els
				var (
					found       = false
					lineSubset  = line
					startOfItem int
				)
				for startOfItem != -1 {
					startOfItem = strings.Index(lineSubset, itemRefStartTag)
					if startOfItem == -1 {
						break
					}

					// for now we will assume the itemrefs are self-closing
					var endOfItem = strings.Index(lineSubset, "/>")
					if endOfItem == -1 {
						return "", fmt.Errorf("failed to parse itemref out of line contents %q due to missing %q", lineSubset, "/>")
					}

					var (
						endOfItemRef = endOfItem + 2
						itemrefEl    = lineSubset[startOfItem:endOfItemRef]
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
	}

	updatedSpineContent := strings.Join(lines, "\n")
	updatedOpfContents = updatedOpfContents[:startIndex+len(spineStartTag)] + updatedSpineContent + updatedOpfContents[endIndex:]

	return updatedOpfContents, nil
}

func extractID(line string) string {
	const idAttr = `id="`
	startIndex := strings.Index(line, idAttr)
	if startIndex == -1 {
		return ""
	}

	startIndex += len(idAttr)
	endIndex := strings.Index(line[startIndex:], `"`)
	if endIndex == -1 {
		return ""
	}

	return line[startIndex : startIndex+endIndex]
}

func getSpineContents(opfContents string) (int, int, string, error) {
	startIndex := strings.Index(opfContents, spineStartTag)
	endIndex := strings.Index(opfContents, spineEndTag)

	if startIndex == -1 || endIndex == -1 {
		return 0, 0, "", ErrNoSpine
	}

	return startIndex, endIndex, opfContents[startIndex+len(spineStartTag) : endIndex], nil
}
