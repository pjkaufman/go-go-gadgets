package converter

import (
	"fmt"
	"regexp"
	"strings"
)

type MdFileInfo struct {
	FilePath     string
	FileName     string
	FileContents string
}

var h1Regex = regexp.MustCompile(`<h1[^\n>]+id="([a-zA-Z-\d]+)"[^\n>]*>`)

func BuildHtmlSongs(mdInfo []MdFileInfo) (string, []string, error) {
	var html = strings.Builder{}
	var headerIdMap = make(map[string]int, len(mdInfo))
	var headerIds = make([]string, len(mdInfo))
	var h1Match []string
	var headerId string
	for i, mdData := range mdInfo {
		fileContentInHtml, err := ConvertMdToHtmlSong(mdData.FilePath, mdData.FileContents)
		if err != nil {
			return "", nil, err
		}

		h1Match = h1Regex.FindStringSubmatch(fileContentInHtml)
		if len(h1Match) > 0 {
			headerId = h1Match[1]
			if num, ok := headerIdMap[headerId]; ok {
				num++
				headerIdMap[headerId] = num

				var newHeaderId = fmt.Sprintf(`%s-%d`, headerId, num)
				fileContentInHtml = strings.Replace(fileContentInHtml, "id=\""+headerId+"\"", "id=\""+newHeaderId+"\"", 1)
				headerId = newHeaderId
			}

			headerIdMap[headerId] = 1
			headerIds[i] = headerId
		} else {
			return "", nil, fmt.Errorf(`no heading found for file "%s"`, mdData.FilePath)
		}

		html.WriteString(fileContentInHtml + "\n")
	}

	return html.String(), headerIds, nil
}
