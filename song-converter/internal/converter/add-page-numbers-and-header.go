package converter

import (
	"strconv"
	"strings"
)

func AddPageNumbersAndHeader(mdInfo []MdFileInfo, bodyLocation string, otherTocLocation string) error {
	var mdData *MdFileInfo
	for i := range mdInfo {
		mdData = &(mdInfo[i])

		var metadata SongMetadata
		_, err := parseFrontmatter(mdData.FilePath, mdData.FileContents, &metadata)
		if err != nil {
			return err
		}

		pageNumbers := getPageNumbers(bodyLocation, metadata.BookLocation)
		if len(pageNumbers) != 0 {
			mdData.PrimaryPageNumbers = pageNumbers
		} else {
			pageNumbers = getPageNumbers(otherTocLocation, metadata.BookLocation)
			if len(pageNumbers) != 0 {
				mdData.SecondaryPageNumbers = pageNumbers
			}
		}

		if len(pageNumbers) == 0 {
			mdData.PrimaryPageNumbers = pageNumbers
		}

		mdData.Header = getHeaderText(mdData.FileContents)
	}

	return nil
}

func getPageNumbers(location, locations string) []int {
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
