package sevenseasentertainment

import (
	"fmt"
	"slices"
	"time"

	"github.com/gocolly/colly/v2"
	googlecache "github.com/pjkaufman/go-go-gadgets/magnum/internal/google-cache"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/slug"
	"github.com/pjkaufman/go-go-gadgets/pkg/crawler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

type VolumeInfo struct {
	Name        string
	ReleaseDate *time.Time
}

func GetVolumeInfo(seriesName string, slugOverride *string, verbose bool) []VolumeInfo {
	var seriesSlug string
	if slugOverride != nil {
		seriesSlug = *slugOverride
	} else {
		seriesSlug = slug.GetSeriesSlugFromName(seriesName)
	}

	c := crawler.CreateNewCollyCrawler(verbose)

	var err error

	var volumeContent = []string{}
	c.OnHTML(".series-volume", func(e *colly.HTMLElement) {
		contentHtml, err := e.DOM.Html()
		if err != nil {
			logger.WriteError(fmt.Sprintf("failed to get content body: %s", err))
		}

		volumeContent = append(volumeContent, contentHtml)
	})

	var url = googlecache.BuildCacheURL(baseURL + seriesPath + seriesSlug + "/")
	err = c.Visit(url)
	if err != nil {
		logger.WriteError(fmt.Sprintf("failed call to google cache for %q: %s", url, err))
	}

	var volumeInfo = []VolumeInfo{}
	var index = 1
	for _, contentHtml := range volumeContent {
		var tempVolumeInfo, err = ParseVolumeInfo(seriesName, contentHtml, index)
		if err != nil {
			logger.WriteError(err.Error())
		}

		if tempVolumeInfo != nil {
			volumeInfo = append(volumeInfo, *tempVolumeInfo)
			index++
		}
	}

	slices.Reverse(volumeInfo)

	return volumeInfo
}
