package sevenseasentertainment

import (
	"fmt"
	"slices"

	"github.com/gocolly/colly/v2"
	sitehandler "github.com/pjkaufman/go-go-gadgets/magnum/internal/site-handler"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/slug"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

func (s *SevenSeasEntertainment) GetVolumeInfo(seriesName string, options sitehandler.ScrapingOptions) ([]*sitehandler.VolumeInfo, int, error) {
	var seriesSlug string
	if options.SlugOverride != nil {
		seriesSlug = *options.SlugOverride
	} else {
		seriesSlug = slug.GetSeriesSlugFromName(seriesName)
	}

	var volumeContent = []string{}
	s.scrapper.OnHTML(".series-volume", func(e *colly.HTMLElement) {
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

	var url = s.options.BaseURL + seriesPath + seriesSlug + "/"
	err := s.scrapper.Visit(url)
	if err != nil {
		return nil, -1, fmt.Errorf("failed call to Seven Seas Entertainment for %q: %w", url, err)
	}

	var volumeInfo = []*sitehandler.VolumeInfo{}
	var index = 1
	for _, contentHtml := range volumeContent {
		var tempVolumeInfo, err = ParseVolumeInfoV2(seriesName, contentHtml, index)
		if err != nil {
			return nil, -1, err
		}

		if tempVolumeInfo != nil {
			volumeInfo = append(volumeInfo, tempVolumeInfo)
			index++
		}
	}

	slices.Reverse(volumeInfo)

	return volumeInfo, len(volumeInfo), nil
}
