package converter

import (
	"strconv"
	"strings"
)

func AddPageNumbersAlternateTitleAndHeader(mdInfo []MdFileInfo, bodyLocation string, otherTocLocation string) error {
	var mdData *MdFileInfo
	for i := range mdInfo {
		mdData = &(mdInfo[i])

		var metadata SongMetadata
		_, err := parseFrontmatter(mdData.FilePath, mdData.FileContents, &metadata)
		if err != nil {
			return err
		}

		if metadata.SkipBook {
			continue
		}

		mdData.PrimaryPageNumbers = getPageNumbers(bodyLocation, metadata.BookLocation)
		mdData.SecondaryPageNumbers = getPageNumbers(otherTocLocation, metadata.BookLocation)

		mdData.Header = strings.TrimSpace(metadata.BookTitle)
		if mdData.Header == "" {
			mdData.Header = getHeaderText(mdData.FileContents)
		}

		mdData.AlternateTitle = strings.TrimSpace(metadata.AlternateTitle)
	}

	return nil
}

func getPageNumbers(location, locations string) []int {
	if strings.TrimSpace(location) == "" {
		return nil
	}

	var possibleLocations = strings.Split(strings.ReplaceAll(strings.ReplaceAll(locations, "(", ""), ")", ""), " ")

	var pageNumbers []int
	for _, possibleLocation := range possibleLocations {
		if pageNumberString, hasPrefix := strings.CutPrefix(possibleLocation, location); hasPrefix {
			pageNumber, err := strconv.ParseFloat(pageNumberString, 64)
			if err != nil {
				continue
			}

			pageNumbers = append(pageNumbers, int(pageNumber))
		}
	}

	return pageNumbers
}

func getHeaderText(content string) string {
	var m = h1Regex.FindStringSubmatch(content)
	if len(m) != 0 {
		return m[1]
	}

	return ""
}
