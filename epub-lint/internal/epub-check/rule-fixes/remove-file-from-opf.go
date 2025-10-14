package rulefixes

import (
	"fmt"
	"slices"
	"strings"
)

const (
	spineStartTag = "<spine"
	spineEndTag   = "</spine>"
	hrefAttribute = "href="
)

var ErrNoSpine = fmt.Errorf("spine tag not found in OPF contents")

func RemoveFileFromOpf(opfContents, fileName string) (string, error) {
	startIndex, endIndex, manifestContent, err := getManifestContents(opfContents)
	if err != nil {
		return "", err
	}

	lines := strings.Split(manifestContent, "\n")
	var (
		fileID string
	)

	for i, line := range lines {
		// to make sure we are only operating on the href
		var startOfHref = strings.Index(line, hrefAttribute)
		if startOfHref == -1 {
			continue
		}

		startOfHref += len(hrefAttribute)
		var hrefQuote = line[startOfHref : startOfHref+1]
		startOfHref++

		var endOfHrefIndex = strings.Index(line[startOfHref:], hrefQuote)
		if endOfHrefIndex == -1 {
			continue
		}

		var hrefContent = line[startOfHref : startOfHref+endOfHrefIndex]
		if !strings.HasSuffix(hrefContent, fileName) {
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
			continue
		}

		fileID = extractID(line)
		lines = slices.Delete(lines, i, i+1)
		break
	}

	updatedManifestContent := strings.Join(lines, "\n")
	updatedOpfContents := opfContents[:startIndex+len(manifestStartTag)] + updatedManifestContent + opfContents[endIndex:]

	if strings.TrimSpace(fileID) == "" {
		return updatedOpfContents, nil
	}

	startIndex, endIndex, spineContent, err := getSpineContents(updatedOpfContents)
	if err != nil {
		return "", err
	}

	lines = strings.Split(spineContent, "\n")
	for i, line := range lines {
		if strings.Contains(line, fmt.Sprintf(`idref="%s"`, fileID)) {
			lines = append(lines[:i], lines[i+1:]...)
			break
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
