package wikipedia

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	sitehandler "github.com/pjkaufman/go-go-gadgets/magnum/internal/site-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

type WikipediaApi struct {
	BaseURL      string
	UserAgent    string
	Verbose      bool
	BuildApiPath sitehandler.ApiPathBuilder
}

func NewWikipediaApi(baseURL, userAgent string, verbose bool, buildApiPath sitehandler.ApiPathBuilder) *WikipediaApi {
	return &WikipediaApi{
		BaseURL:      baseURL,
		UserAgent:    userAgent,
		Verbose:      verbose,
		BuildApiPath: buildApiPath,
	}
}

func (wa *WikipediaApi) GetSectionInfo(pageTitle string) (*WikipediaSectionInfo, error) {
	url := wa.BuildApiPath(wa.BaseURL, apiPath, pageTitle)
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
