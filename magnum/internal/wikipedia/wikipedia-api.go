package wikipedia

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

type WikipediaSectionInfo struct {
	Parse PageSectionInfo `json:"parse"`
}

type PageSectionInfo struct {
	Title    string        `json:"title"`
	PageId   int64         `json:"pageid"`
	Sections []SectionInfo `json:"sections"`
}

type SectionInfo struct {
	TocLevel   int    `json:"toclevel"`
	Level      string `json:"level"`
	Heading    string `json:"line"`
	Number     string `json:"number"`
	Index      string `json:"index"`
	FromTitle  string `json:"fromtitle"`
	ByteOffset int32  `json:"byteoffset"`
	Anchor     string `json:"anchor"`
	LinkAnchor string `json:"linkAnchor"`
}

func getSectionInfo(userAgent, pageTitle string) *WikipediaSectionInfo {
	var url = fmt.Sprintf("%s%s?action=parse&prop=sections&page=%s&format=json", BaseURL, apiPath, pageTitle)
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.WriteErrorf("failed to build http request for section info for %q: %s\n", url, err)
	}
	request.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(request)
	if err != nil {
		logger.WriteErrorf("failed to get section info for %q: %s\n", url, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.WriteErrorf("failed to get section info body for %q: %s\n", url, err)
	}

	var sectionInfo = &WikipediaSectionInfo{}
	err = json.Unmarshal(body, sectionInfo)
	if err != nil {
		logger.WriteErrorf("failed to unmarshal section info for %q: %s\n", url, err)
	}

	return sectionInfo
}
