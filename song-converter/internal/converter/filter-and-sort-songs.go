package converter

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type filterData struct {
	mdInfo     *MdFileInfo
	pageNumber int
}

var h1Regex = regexp.MustCompile("\n#\\s+(.+)\n")

func FilterAndSortSongs(mdInfo []MdFileInfo, location string) ([]MdFileInfo, error) {
	var pageInfo = make([]filterData, 0, len(mdInfo))
	for _, mdData := range mdInfo {
		var metadata SongMetadata
		_, err := parseFrontmatter(mdData.FilePath, mdData.FileContents, &metadata)
		if err != nil {
			return nil, err
		}

		pageNumbers := getPageNumbers(location, metadata.BookLocation)
		if len(pageNumbers) == 0 {
			continue
		}

		mdData.PageNumbers = pageNumbers
		mdData.Header = getHeaderText(mdData.FileContents)

		for _, pageNumber := range pageNumbers {
			pageInfo = append(pageInfo, filterData{
				mdInfo:     &mdData,
				pageNumber: pageNumber,
			})
		}
	}

	sort.Slice(pageInfo, func(i, j int) bool {
		if pageInfo[i].pageNumber != pageInfo[j].pageNumber {
			return pageInfo[i].pageNumber < pageInfo[j].pageNumber
		}

		return pageInfo[i].mdInfo.FileName < pageInfo[j].mdInfo.FileName
	})

	var newMdInfo = make([]MdFileInfo, len(pageInfo))
	for i, pageData := range pageInfo {
		newMdInfo[i] = *pageData.mdInfo
	}

	return newMdInfo, nil
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
