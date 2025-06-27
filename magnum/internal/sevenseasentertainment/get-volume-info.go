package sevenseasentertainment

import (
	"fmt"
	"slices"

	"github.com/gocolly/colly/v2"
	sitehandler "github.com/pjkaufman/go-go-gadgets/magnum/internal/site-handler"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/slug"
)

func (s *SevenSeasEntertainment) GetVolumeInfo(seriesName string, options sitehandler.ScrapingOptions) ([]*sitehandler.VolumeInfo, int, error) {
	var seriesSlug string
	if options.SlugOverride != nil {
		seriesSlug = *options.SlugOverride
	} else {
		seriesSlug = slug.GetSeriesSlugFromName(seriesName)
	}

	var (
		volumeContent = []string{}
		firstErr      error
	)
	s.scrapper.OnHTML(".series-volume", func(e *colly.HTMLElement) {
		contentHtml, err := e.DOM.Html()
		if err != nil {
			firstErr = fmt.Errorf("failed to get content body: %w", err)
			e.Request.Abort()

			return
		}

		volumeContent = append(volumeContent, contentHtml)
	})

	var url = s.options.BaseURL + seriesPath + seriesSlug + "/"
	err := s.scrapper.Visit(url)
	if err != nil {
		return nil, -1, fmt.Errorf("failed call to Seven Seas Entertainment for %q: %w", url, err)
	}

	if firstErr != nil {
		return nil, -1, firstErr
	}

	var volumeInfo = []*sitehandler.VolumeInfo{}
	var index = 1
	for _, contentHtml := range volumeContent {
		var tempVolumeInfo, err = ParseVolumeInfo(seriesName, contentHtml, index)
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
