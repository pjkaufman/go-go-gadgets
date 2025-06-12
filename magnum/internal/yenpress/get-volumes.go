package yenpress

import (
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/slug"
	"github.com/pjkaufman/go-go-gadgets/pkg/crawler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

type VolumeInfo struct {
	Name         string
	RelativeLink string
}

func GetVolumes(seriesName string, slugOverride *string, verbose bool) ([]*VolumeInfo, int) {
	c := crawler.CreateNewCollyCrawler(verbose)

	var volumes = []*VolumeInfo{}

	c.OnHTML("#volumes-list > div > div > div.inline_block", func(e *colly.HTMLElement) {
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
	c.OnHTML("body > div > div:nth-child(4) > div > section.content-heading.fade-in-container > div > h1 > sup", func(e *colly.HTMLElement) {
		if strings.TrimSpace(e.Text) != "" {
			val, err := strconv.Atoi(e.Text)
			if err == nil {
				numVolumes = val
			}
		}
	})

	var seriesSlug string
	if slugOverride != nil {
		seriesSlug = *slugOverride
	} else {
		seriesSlug = slug.GetSeriesSlugFromName(seriesName)
	}

	var seriesURL = baseURL + seriesPath + seriesSlug
	err := c.Visit(seriesURL)
	if err != nil {
		logger.WriteErrorf("failed call to yen press: %s\n", err)
		return nil, 0
	}

	return volumes, numVolumes
}
