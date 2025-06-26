package sevenseasentertainment

import (
	"slices"
	"time"

	"github.com/gocolly/colly/v2"
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
			logger.WriteErrorf("failed to get content body: %s\n", err)
		}

		volumeContent = append(volumeContent, contentHtml)
	})

	// These are old versions of how this was handled. I am going to just use the site info directly.
	// var url = googlecache.BuildCacheURL(baseURL + seriesPath + seriesSlug + "/")
	// url, err := internetarchive.GetLatestPageSnapshot(baseURL+seriesPath+seriesSlug, verbose)
	// if err != nil {
	// 	logger.WriteErrorf("failed call to internet archive to get latest page snapshot for %q: %s\n", baseURL+seriesPath+seriesSlug, err)
	// }
	var url = BaseURL + seriesPath + seriesSlug + "/"
	err = c.Visit(url)
	if err != nil {
		logger.WriteErrorf("failed call to internet archive for %q: %s\n", url, err)
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
