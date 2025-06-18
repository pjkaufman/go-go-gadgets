package vizmedia

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gocolly/colly/v2"
	sitehandler "github.com/pjkaufman/go-go-gadgets/magnum/internal/site-handler"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/slug"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

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

	volumes, err := v.getListOfVolumesWithInfoV2(fullVolumeLink, seriesName)
	if err != nil {
		return nil, -1, err
	}

	slices.Reverse(volumes)

	return volumes, len(volumes), nil
}

// TODO: see about converting this into something attached to the VizMedia struct
func (v *VizMedia) getListOfVolumesWithInfoV2(fullVolumeLink, seriesName string) ([]*sitehandler.VolumeInfo, error) {
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

		// TODO: swap this to be on VizMedia
		releaseDate := getVolumeReleaseDate(v.scrapper.Clone(), volumeReleasePage)

		volumes = append(volumes, &sitehandler.VolumeInfo{
			Name:        name,
			ReleaseDate: &releaseDate,
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
