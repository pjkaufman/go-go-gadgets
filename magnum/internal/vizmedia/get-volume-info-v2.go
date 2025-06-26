package vizmedia

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	sitehandler "github.com/pjkaufman/go-go-gadgets/magnum/internal/site-handler"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/slug"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

var alreadyReleasedDate = time.Now().Add(-1 * 24 * time.Hour)

func (v *VizMedia) GetVolumeInfo(seriesName string, options sitehandler.ScrapingOptions) ([]*sitehandler.VolumeInfo, int, error) {
	var fullVolumeLink string
	v.scrapper.OnHTML("#section1 > div > div.clearfix.mar-t-md.mar-b-lg > div > a", func(e *colly.HTMLElement) {
		var link = e.Attr("href")
		if strings.TrimSpace(link) != "" {
			fullVolumeLink = link
		}
	})

	var seriesSlug string
	if options.SlugOverride != nil {
		seriesSlug = *options.SlugOverride
	} else {
		seriesSlug = slug.GetSeriesSlugFromName(seriesName)
	}

	var seriesURL = v.options.BaseURL + seriesSlug
	err := v.scrapper.Visit(seriesURL)
	if err != nil {
		return nil, -1, fmt.Errorf("failed call to viz media series page: %w", err)
	}

	if strings.TrimSpace(fullVolumeLink) == "" {
		return nil, -1, fmt.Errorf("failed to get the list of volumes link for %q", seriesName)
	}

	// reset colly before moving to the next page to avoid the potential blowback of having
	// the OnHtml logic and other logic running again
	v.scrapper = v.scrapper.Clone()

	volumes, err := v.getListOfVolumesWithInfo(fullVolumeLink, seriesName)
	if err != nil {
		return nil, -1, err
	}

	slices.Reverse(volumes)

	return volumes, len(volumes), nil
}

func (v *VizMedia) getListOfVolumesWithInfo(fullVolumeLink, seriesName string) ([]*sitehandler.VolumeInfo, error) {
	var volumes = []*sitehandler.VolumeInfo{}

	var volumeNum = 1
	v.scrapper.OnHTML("body > div.bg-off-white.overflow-hide > section > section.row.mar-t-lg.mar-t-xl--lg.mar-last-row > div > article", func(e *colly.HTMLElement) {
		var html, err = e.DOM.Html()
		if err != nil {
			logger.WriteErrorf("failed to get the html for the volume info for %q\n", fullVolumeLink)
		}

		// TODO: update this to only care about Manga (see UT/validator for a sample of one with non-manga)
		name, volumeReleasePage, isReleased, err := ParseVolumeHtml(html, seriesName, volumeNum)
		if err != nil {
			logger.WriteError(err.Error())
		}

		volumeNum++

		// to minimize the amount of API calls we need to make, we will only get the actual release dates for unreleased
		if isReleased {
			volumes = append(volumes, &sitehandler.VolumeInfo{
				Name:        name,
				ReleaseDate: &alreadyReleasedDate,
			})

			return
		}

		volumes = append(volumes, &sitehandler.VolumeInfo{
			Name:        name,
			ReleaseDate: v.getVolumeReleaseDate(volumeReleasePage),
		})
	})

	var mangaVolumesLink = v.options.BaseURL + fullVolumeLink
	err := v.scrapper.Visit(mangaVolumesLink)
	if err != nil {
		return nil, fmt.Errorf("failed call to viz media volumes page: %w", err)
	}

	slices.Reverse(volumes)

	return volumes, nil
}

func (v *VizMedia) getVolumeReleaseDate(volumeReleasePage string) *time.Time {
	var releaseDate time.Time
	v.scrapper.OnHTML("#product_row > div.row.pad-b-xl > div.g-6--lg.type-sm.type-rg--md.line-caption > div:nth-child(1) > div.o_release-date.mar-b-md", func(e *colly.HTMLElement) {
		var text = e.DOM.Text()

		text = strings.TrimSpace(strings.Replace(text, "Release", "", 1))
		tempDate, err := time.Parse(releaseDateFormat, text)
		if err != nil {
			logger.WriteErrorf("failed to parse %q to a date time value: %v\n", text, err)
		}

		releaseDate = tempDate
	})

	var mangaVolumesLink = v.options.BaseURL + volumeReleasePage
	err := v.scrapper.Visit(mangaVolumesLink)
	if err != nil {
		logger.WriteErrorf("failed call to viz media volume release page: %s\n", err)
	}

	return &releaseDate
}
