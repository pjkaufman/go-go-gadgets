package converter

import (
	"fmt"
	"strings"
)

type MdFileInfo struct {
	FilePath     string
	FileName     string
	FileContents string
}

const (
	h1Indicator        = "<h1"
	endingTagIndicator = ">"
)

func BuildHtmlSongs(mdInfo []MdFileInfo) (string, []string, error) {
	var (
		html                                                   = strings.Builder{}
		headerIdMap                                            = make(map[string]int, len(mdInfo))
		headerIds                                              = make([]string, len(mdInfo))
		headerId, h1OpeningTag                                 string
		firstH1IndexStart, firstH1IndexEnd, h1IdStart, h1IdEnd int
	)
	for i, mdData := range mdInfo {
		fileContentInHtml, err := ConvertMdToHtmlSong(mdData.FilePath, mdData.FileContents)
		if err != nil {
			return "", nil, err
		}

		firstH1IndexStart = strings.Index(fileContentInHtml, h1Indicator)
		if firstH1IndexStart != -1 {
			firstH1IndexEnd = strings.Index(fileContentInHtml[firstH1IndexStart:], endingTagIndicator)
			if firstH1IndexEnd == -1 {
				return "", nil, fmt.Errorf("no h1 heading found for file %q", mdData.FilePath)
			}

			h1OpeningTag = fileContentInHtml[firstH1IndexStart : firstH1IndexStart+firstH1IndexEnd]
			h1IdStart = strings.Index(h1OpeningTag, "id=\"")
			if h1IdStart == -1 {
				return "", nil, fmt.Errorf("no h1 heading id found for file %q", mdData.FilePath)
			}

			h1IdEnd = strings.Index(h1OpeningTag[h1IdStart+4:], "\"")
			if h1IdStart == -1 {
				return "", nil, fmt.Errorf("no h1 heading id found for file %q", mdData.FilePath)
			}

			headerId = h1OpeningTag[h1IdStart+4 : h1IdStart+4+h1IdEnd]

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
			return "", nil, fmt.Errorf("no h1 heading found for file %q", mdData.FilePath)
		}

		html.WriteString(fileContentInHtml + "\n")
	}

	return html.String(), headerIds, nil
}
