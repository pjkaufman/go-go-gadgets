package converter

import (
	"regexp"
	"sort"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
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

		mdData.PrimaryPageNumbers = pageNumbers
		mdData.Header = getHeaderText(mdData.FileContents)

		for _, pageNumber := range pageNumbers {
			pageInfo = append(pageInfo, filterData{
				mdInfo:     &mdData,
				pageNumber: pageNumber,
			})
		}
	}

	c := collate.New(language.English)
	sort.Slice(pageInfo, func(i, j int) bool {
		if pageInfo[i].pageNumber != pageInfo[j].pageNumber {
			return pageInfo[i].pageNumber < pageInfo[j].pageNumber
		}

		return c.CompareString(pageInfo[i].mdInfo.FileName, pageInfo[j].mdInfo.FileName) < 0
	})

	var newMdInfo = make([]MdFileInfo, len(pageInfo))
	for i, pageData := range pageInfo {
		newMdInfo[i] = *pageData.mdInfo
	}

	return newMdInfo, nil
}

func SortSongs(mdInfo []MdFileInfo) []MdFileInfo {
	var pageInfo = make([]filterData, 0, len(mdInfo))
	for _, mdData := range mdInfo {
		if len(mdData.PrimaryPageNumbers) == 0 {
			continue
		}

		for _, pageNumber := range mdData.PrimaryPageNumbers {
			pageInfo = append(pageInfo, filterData{
				mdInfo:     &mdData,
				pageNumber: pageNumber,
			})
		}
	}

	c := collate.New(language.English)
	sort.Slice(pageInfo, func(i, j int) bool {
		if pageInfo[i].pageNumber != pageInfo[j].pageNumber {
			return pageInfo[i].pageNumber < pageInfo[j].pageNumber
		}

		return c.CompareString(pageInfo[i].mdInfo.FileName, pageInfo[j].mdInfo.FileName) < 0
	})

	var newMdInfo = make([]MdFileInfo, len(pageInfo))
	for i, pageData := range pageInfo {
		newMdInfo[i] = *pageData.mdInfo
	}

	return newMdInfo
}
