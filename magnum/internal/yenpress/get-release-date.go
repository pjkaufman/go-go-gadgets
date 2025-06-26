package yenpress

import (
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/pjkaufman/go-go-gadgets/pkg/crawler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

func GetReleaseDateInfo(info *VolumeInfo, verbose bool) *time.Time {
	if info == nil {
		if verbose {
			logger.WriteInfo("no volume info provided...")
		}

		return nil
	}

	c := crawler.CreateNewCollyCrawler(verbose)

	var releaseDate string
	c.OnHTML("div.books-page.series-page > section.book-details.wrapper-1410.prel.fade-in-container > div.detail.active > div.detail-info.fade-el > div:nth-child(3) > div:nth-child(1) > p", func(e *colly.HTMLElement) {
		releaseDate = e.Text
	})

	var volumeURL = BaseURL + info.RelativeLink
	err := c.Visit(volumeURL)
	if err != nil {
		logger.WriteErrorf("failed call to Yen Press: %s\n", err)
	}

	if releaseDate == "" {
		if verbose {
			logger.WriteInfof("no release date found on the page: %q\n", volumeURL)
		}

		return nil
	}

	date, err := time.Parse(releaseDateFormat, releaseDate)
	if err != nil {
		logger.WriteErrorf("failed to parse %q to a date time value: %v\n", releaseDate, err)
	}

	return &date
}
