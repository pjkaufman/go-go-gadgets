package jnovelclub

import (
	"encoding/json"
	"strconv"
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

func GetVolumeInfo(seriesName string, slugOverride *string, verbose bool) []VolumeInfo {
	var seriesSlug string
	if slugOverride != nil {
		seriesSlug = *slugOverride
	} else {
		seriesSlug = slug.GetSeriesSlugFromName(seriesName)
	}

	c := crawler.CreateNewCollyCrawler(verbose)

	var jsonVolumeInfo JSONVolumeInfo
	c.OnHTML("#__NEXT_DATA__", func(e *colly.HTMLElement) {
		err := json.Unmarshal([]byte(e.Text), &jsonVolumeInfo)
		if err != nil {
			logger.WriteErrorf("failed to deserialize json to volume info: %s\n", err)
		}
	})

	var seriesURL = BaseURL + seriesPath + seriesSlug
	err := c.Visit(seriesURL)
	if err != nil {
		logger.WriteErrorf("failed call to JNovel Club for %q: %s\n", seriesURL, err)
	}

	var numVolumes = len(jsonVolumeInfo.Props.PageProps.Aggregate.Volumes)
	var volumes = make([]VolumeInfo, numVolumes)
	for i, volume := range jsonVolumeInfo.Props.PageProps.Aggregate.Volumes {
		// no release data is present, but this should not happen
		if volume.Volume.Publishing.Seconds == "" {
			logger.WriteErrorf("failed to get volume info properly for series %q as there is no publishing data\n", seriesName)
		}

		secondsFromEpoch, err := strconv.Atoi(volume.Volume.Publishing.Seconds)
		if err != nil {
			logger.WriteErrorf("failed to parse out seconds for volume %q: %s\n", volume.Volume.Title, err)
		}

		volumes[numVolumes-i-1] = VolumeInfo{
			Name:        volume.Volume.Title,
			ReleaseDate: time.Unix(int64(secondsFromEpoch), int64(volume.Volume.Publishing.Nanos)),
		}
	}

	return volumes
}
