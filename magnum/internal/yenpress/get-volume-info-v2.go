package yenpress

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	sitehandler "github.com/pjkaufman/go-go-gadgets/magnum/internal/site-handler"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/slug"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

func (y *YenPress) GetVolumeInfo(seriesName string, options sitehandler.ScrapingOptions) ([]*sitehandler.VolumeInfo, int, error) {
	var volumes = []*VolumeInfo{}

	y.scrapper.OnHTML("#volumes-list > div > div > div.inline_block", func(e *colly.HTMLElement) {
		contentHtml, err := e.DOM.Html()
		if err != nil {
			logger.WriteErrorf("failed to get content body: %s\n", err)
		}

		volumeInfo, err := ParseVolumeInfo(seriesName, contentHtml)
		if err != nil {
			logger.WriteError(err.Error())
		}

		if volumeInfo != nil {
			volumes = append(volumes, volumeInfo)
		}
	})

	var numVolumes int = -1
	y.scrapper.OnHTML("body > div > div:nth-child(4) > div > section.content-heading.fade-in-container > div > h1 > sup", func(e *colly.HTMLElement) {
		if strings.TrimSpace(e.Text) != "" {
			val, err := strconv.Atoi(e.Text)
			if err == nil {
				numVolumes = val
			}
		}
	})

	var seriesSlug string
	if options.SlugOverride != nil {
		seriesSlug = *options.SlugOverride
	} else {
		seriesSlug = slug.GetSeriesSlugFromName(seriesName)
	}

	var seriesURL = y.options.BaseURL + seriesPath + seriesSlug
	err := y.scrapper.Visit(seriesURL)
	if err != nil {
		return nil, -1, fmt.Errorf("failed call to Yen Press: %w", err)
	}

	var today = time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	var volumeData []*sitehandler.VolumeInfo
	for _, info := range volumes {
		releaseDate, err := y.getReleaseDateInfo(info)
		if err != nil {
			return nil, -1, err
		}

		if releaseDate != nil {
			if releaseDate.Before(today) {
				break
			} else {
				volumeData = append(volumeData, &sitehandler.VolumeInfo{
					Name:        info.Name,
					ReleaseDate: releaseDate,
				})
			}
		}
	}

	return volumeData, numVolumes, nil
}

func (y *YenPress) getReleaseDateInfo(info *VolumeInfo) (*time.Time, error) {
	if info == nil {
		if y.options.Verbose {
			logger.WriteInfo("no volume info provided...")
		}

		return nil, nil
	}

	var releaseDate string
	y.scrapper.OnHTML("div.books-page.series-page > section.book-details.wrapper-1410.prel.fade-in-container > div.detail.active > div.detail-info.fade-el > div:nth-child(3) > div:nth-child(1) > p", func(e *colly.HTMLElement) {
		releaseDate = e.Text
	})

	var volumeURL = y.options.BaseURL + info.RelativeLink
	err := y.scrapper.Visit(volumeURL)
	if err != nil {
		return nil, fmt.Errorf("failed call to Yen Press: %w", err)
	}

	if releaseDate == "" {
		if y.options.Verbose {
			logger.WriteInfof("no release date found on the page: %q\n", volumeURL)
		}

		return nil, nil
	}

	date, err := time.Parse(releaseDateFormat, releaseDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %q to a date time value: %w", releaseDate, err)
	}

	return &date, nil
}
