package converter

import (
	"fmt"
	"strings"
)

type MdFileInfo struct {
	FilePath     string
	FileName     string
	FileContents string
	// book generation properties
	Header      string
	PageNumbers []int
}

func BuildHtmlSongs(mdInfo []MdFileInfo, songType SongGenerationType) (string, []string, error) {
	html := strings.Builder{}
	html.Grow(estimateCapacity(mdInfo)) // Pre-allocate capacity

	var (
		headerIdMap     = make(map[string]int, len(mdInfo))
		pageNumberIndex = make(map[string]int)
		headerIds       = make([]string, len(mdInfo))
		pageNumber      int
		isLastOnPage    bool
	)

	for i, mdData := range mdInfo {
		if songType == Book {
			if val, ok := pageNumberIndex[mdData.FileName]; ok {
				pageNumber = mdData.PageNumbers[val]
			} else {
				pageNumber = mdData.PageNumbers[0]
				pageNumberIndex[mdData.FileName] = 1
			}

			var nextMdData *MdFileInfo
			if i+1 < len(mdInfo) {
				nextMdData = &mdInfo[i+1]
			}

			nextPageNumber := 0
			if nextMdData != nil {
				if val, ok := pageNumberIndex[nextMdData.FileName]; ok {
					nextPageNumber = nextMdData.PageNumbers[val]
				} else {
					nextPageNumber = mdData.PageNumbers[0]
				}

			}

			// this is not needed for the last page, so we will use this for all pages except that one
			isLastOnPage = nextMdData != nil && pageNumber != nextPageNumber
		}

		fileContentInHtml, err := ConvertMdToHtmlSong(mdData.FilePath, mdData.FileContents, songType, isLastOnPage)
		if err != nil {
			return "", nil, err
		}

		updatedContent, headerId, err := extractAndUpdateH1Id(fileContentInHtml, headerIdMap)
		if err != nil {
			return "", nil, fmt.Errorf("error processing file %q: %w", mdData.FilePath, err)
		}

		headerIds[i] = headerId
		html.WriteString(updatedContent)
		html.WriteByte('\n')
	}

	return html.String(), headerIds, nil
}

func estimateCapacity(mdInfo []MdFileInfo) int {
	totalSize := 0
	for _, info := range mdInfo {
		totalSize += len(info.FileContents) * 2 // Rough estimate, HTML might be larger than Markdown
	}
	return totalSize
}

func extractAndUpdateH1Id(content string, headerIdMap map[string]int) (string, string, error) {
	h1Start := strings.Index(content, "<h1")
	if h1Start == -1 {
		return "", "", fmt.Errorf("no h1 heading found")
	}

	idStart := strings.Index(content[h1Start:], "id=\"")
	if idStart == -1 {
		return "", "", fmt.Errorf("no h1 heading id found")
	}
	idStart += h1Start + 4

	idEnd := strings.IndexByte(content[idStart:], '"')
	if idEnd == -1 {
		return "", "", fmt.Errorf("malformed h1 heading id")
	}
	idEnd += idStart

	headerId := content[idStart:idEnd]
	newHeaderId := headerId

	if num, ok := headerIdMap[headerId]; ok {
		num++
		headerIdMap[headerId] = num
		newHeaderId = fmt.Sprintf("%s-%d", headerId, num)

		content = content[:idStart] + newHeaderId + content[idEnd:]
	}

	headerIdMap[newHeaderId] = 1
	return content, newHeaderId, nil
}
