package hanashimedia

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	sitehandler "github.com/pjkaufman/go-go-gadgets/magnum/internal/site-handler"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/slug"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

type SeriesRequest struct {
	At string `json:"at"`
}

func (hm *HanashiMedia) GetVolumeInfo(seriesName string, options sitehandler.ScrapingOptions) ([]*sitehandler.VolumeInfo, int, error) {
	var seriesSlug string
	if options.SlugOverride != nil {
		seriesSlug = *options.SlugOverride
	} else {
		seriesSlug = slug.GetSeriesSlugFromName(seriesName)
	}

	seriesData, err := getSeriesData(hm.options.BaseURL, seriesSlug, hm.options.UserAgent, hm.options.Verbose)
	if err != nil {
		return nil, 0, err
	}

	sort.Slice(seriesData.Data.Entries, func(i, j int) bool {
		return seriesData.Data.Entries[i].Order > seriesData.Data.Entries[j].Order
	})

	var volumes = make([]*sitehandler.VolumeInfo, len(seriesData.Data.Entries))
	for i, volumeData := range seriesData.Data.Entries {
		volumes[i] = &sitehandler.VolumeInfo{
			Name:        volumeData.Ebook.Title,
			ReleaseDate: &volumeData.Ebook.Release,
		}
	}

	return volumes, len(volumes), nil
}

func getSeriesData(baseUrl, series, userAgent string, verbose bool) (*JSONVolumeInfo, error) {
	if verbose {
		logger.WriteInfof("Fetching series info for %q\n", series)
	}

	payload, err := json.Marshal(SeriesRequest{
		At: series,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to encode series title to JSON series request: %w", err)
	}

	if verbose {
		logger.WriteInfof("payload: %q\n", payload)
	}

	if strings.HasSuffix(baseUrl, "/") {
		baseUrl = strings.TrimRight(baseUrl, "/")
	}

	var url = fmt.Sprintf("%s/%s", baseUrl, seriesPath)
	if verbose {
		logger.WriteInfof("Calling out to %q\n", url)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request to Hanashi Media: %w", err)
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make http request to Hanashi Media: %w", err)
	}
	defer resp.Body.Close()

	var (
		body, _    = io.ReadAll(resp.Body)
		seriesInfo JSONVolumeInfo
	)
	err = json.Unmarshal(body, &seriesInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to convert Hanashi Media response body to json: %q: %w", body, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("response from Hanashi Media indicates that the request failed: %s: %q", resp.Status, body)
	}

	return &seriesInfo, nil
}
