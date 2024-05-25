package vizmedia

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/slug"
	"github.com/pjkaufman/go-go-gadgets/pkg/crawler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

type VolumeInfo struct {
	Name        string
	ReleaseDate time.Time
}

var alreadyReleasedDate = time.Now().Add(-1 * 24 * time.Hour)

func GetVolumeInfo(seriesName string, slugOverride *string, verbose bool) []*VolumeInfo {
	c := crawler.CreateNewCollyCrawler(verbose)

	var fullVolumeLink string
	c.OnHTML("#section1 > div > div.clearfix.mar-t-md.mar-b-lg > div > a", func(e *colly.HTMLElement) {
		var link = e.Attr("href")
		if strings.TrimSpace(link) != "" {
			fullVolumeLink = link
		}
	})

	var seriesSlug string
	if slugOverride != nil {
		seriesSlug = *slugOverride
	} else {
		seriesSlug = slug.GetSeriesSlugFromName(seriesName)
	}

	var seriesURL = baseURL + seriesSlug
	err := c.Visit(seriesURL)
	if err != nil {
		logger.WriteError(fmt.Sprintf("failed call to viz media series page: %s", err))
	}

	if strings.TrimSpace(fullVolumeLink) == "" {
		logger.WriteError(fmt.Sprintf(`failed to get the list of volumes link for %q`, seriesName))
	}

	return getListOfVolumesWithInfo(c.Clone(), fullVolumeLink, seriesName)
}

func getListOfVolumesWithInfo(c *colly.Collector, fullVolumeLink, seriesName string) []*VolumeInfo {
	var volumes = []*VolumeInfo{}
	if c == nil {
		return volumes
	}

	var volumeNum = 1
	c.OnHTML("body > div.bg-off-white.overflow-hide > section > section.row.mar-t-lg.mar-t-xl--lg.mar-last-row > div > article", func(e *colly.HTMLElement) {
		var html, err = e.DOM.Html()
		if err != nil {
			logger.WriteError(fmt.Sprintf(`failed to get the html for the volume info for %q`, fullVolumeLink))
		}

		name, volumeReleasePage, isReleased, err := ParseVolumeHtml(html, seriesName, volumeNum)
		if err != nil {
			logger.WriteError(err.Error())
		}

		volumeNum++

		// to minimize the amount of API calls we need to make, we will only get the actual release dates for unreleased
		if isReleased {
			volumes = append(volumes, &VolumeInfo{
				Name:        name,
				ReleaseDate: alreadyReleasedDate,
			})

			return
		}

		volumes = append(volumes, &VolumeInfo{
			Name:        name,
			ReleaseDate: getVolumeReleaseDate(c.Clone(), volumeReleasePage),
		})
	})

	var mangaVolumesLink = baseURL + fullVolumeLink
	err := c.Visit(mangaVolumesLink)
	if err != nil {
		logger.WriteError(fmt.Sprintf("failed call to viz media volumes page: %s", err))
	}

	slices.Reverse(volumes)

	return volumes
}

func getVolumeReleaseDate(c *colly.Collector, volumeReleasePage string) time.Time {
	var releaseDate time.Time
	c.OnHTML("#product_row > div.row.pad-b-xl > div.g-6--lg.type-sm.type-rg--md.line-caption > div:nth-child(1) > div.o_release-date.mar-b-md", func(e *colly.HTMLElement) {
		var text = e.DOM.Text()

		text = strings.TrimSpace(strings.Replace(text, "Release", "", 1))
		tempDate, err := time.Parse(releaseDateFormat, text)
		if err != nil {
			logger.WriteError(fmt.Sprintf("failed to parse %q to a date time value: %v", text, err))
		}

		releaseDate = tempDate
	})

	var mangaVolumesLink = baseURL + volumeReleasePage
	err := c.Visit(mangaVolumesLink)
	if err != nil {
		logger.WriteError(fmt.Sprintf("failed call to viz media volume release page: %s", err))
	}

	return releaseDate
}
