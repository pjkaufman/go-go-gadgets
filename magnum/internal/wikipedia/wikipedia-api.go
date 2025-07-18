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

type WikipediaApi struct {
	BaseURL   string
	UserAgent string
	Verbose   bool
	ApiPath   string
}

func NewWikipediaApi(baseURL, userAgent string, verbose bool, apiPath string) *WikipediaApi {
	return &WikipediaApi{
		BaseURL:   baseURL,
		UserAgent: userAgent,
		Verbose:   verbose,
		ApiPath:   apiPath,
	}
}

func (wa *WikipediaApi) GetSectionInfo(pageTitle string) (*WikipediaSectionInfo, error) {
	url := GetWikipediaAPIUrl(wa.BaseURL, wa.ApiPath, pageTitle)

	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build http request for section info for %q: %w", url, err)
	}
	request.Header.Set("User-Agent", wa.UserAgent)
	if wa.Verbose {
		logger.WriteInfof("calling out to %q to get the section info for %q", url, pageTitle)
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get section info for %q: %w", url, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to get section info body for %q: %w", url, err)
	}

	var sectionInfo = &WikipediaSectionInfo{}
	err = json.Unmarshal(body, sectionInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal section info for %q: %w", url, err)
	}

	return sectionInfo, nil
}

func GetWikipediaAPIUrl(baseURL, apiPath, pageTitle string) string {
	return fmt.Sprintf("%s%s?action=parse&prop=sections&page=%s&format=json", baseURL, apiPath, pageTitle)
}
