package jnovelclub

import (
	"fmt"
	"slices"
	"time"

	"github.com/pjkaufman/go-go-gadgets/magnum/internal/slug"
	"github.com/pjkaufman/go-go-gadgets/pkg/crawler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

type VolumeInfo struct {
	Name        string
	ReleaseDate time.Time
}

func GetVolumeInfo(seriesName string, slugOverride *string, verbose bool) []VolumeInfo {
	var volumes = []VolumeInfo{}

	// playwright is used here instead of colly since the release date info is loaded by JS
	// and so we need to wait for the page to load to scrape the page
	pw, browser, page := crawler.CreateNewPlaywrightCrawler()

	var seriesSlug string
	if slugOverride != nil {
		seriesSlug = *slugOverride
	} else {
		seriesSlug = slug.GetSeriesSlugFromName(seriesName)
	}

	var seriesURL = baseURL + seriesPath + seriesSlug
	_, err := page.Goto(seriesURL)
	if err != nil {
		logger.WriteError(fmt.Sprintf("failed to visit \"%s\": %v", seriesURL, err))
	}

	err = page.WaitForLoadState()
	if err != nil {
		logger.WriteError(fmt.Sprintf("failed to wait for the load state: %v", err))
	}

	volumeTitles, err := page.Locator("div.header > h2 > a").All()
	if err != nil {
		logger.WriteError(fmt.Sprintf("could not get entries: %v", err))
	}

	for _, entry := range volumeTitles {
		title, err := entry.TextContent()
		if err != nil {
			logger.WriteError(fmt.Sprintf("could not get text content for volume title: %v", err))
		}

		volumes = append(volumes, VolumeInfo{
			Name: title,
		})
	}

	entries, err := page.Locator("div.f1g4j9eh > div.f1k2es0r > div.f1ijq7jq > div.f1s7hfqq.label.f1n072s.color-digital > div.f1oyfch5.text").All()
	if err != nil {
		logger.WriteError(fmt.Sprintf("could not get entries: %v", err))
	}

	for i, entry := range entries {
		releaseDate, err := entry.TextContent()
		if err != nil {
			logger.WriteError(fmt.Sprintf("could not get text content for release date: %v", err))
		}

		date, err := time.Parse(releaseDateFormat, releaseDate)
		if err != nil {
			logger.WriteError(fmt.Sprintf("failed to parse \"%s\" to a date time value: %v", releaseDate, err))
		}

		if i >= len(volumes) {
			break
		}

		volumes[i].ReleaseDate = date
	}

	crawler.ClosePlaywrightCrawler(pw, browser)

	slices.Reverse(volumes)

	return volumes
}
