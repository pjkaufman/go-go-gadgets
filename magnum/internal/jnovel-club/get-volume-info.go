package jnovelclub

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gocolly/colly/v2"
	sitehandler "github.com/pjkaufman/go-go-gadgets/magnum/internal/site-handler"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/slug"
)

func (j *JNovelClub) GetVolumeInfo(seriesName string, options sitehandler.ScrapingOptions) ([]*sitehandler.VolumeInfo, int, error) {
	var seriesSlug string
	if options.SlugOverride != nil {
		seriesSlug = *options.SlugOverride
	} else {
		seriesSlug = slug.GetSeriesSlugFromName(seriesName)
	}

	var firstErr error
	var jsonVolumeInfo JSONVolumeInfo
	j.scrapper.OnHTML("#__NEXT_DATA__", func(e *colly.HTMLElement) {
		err := json.Unmarshal([]byte(e.Text), &jsonVolumeInfo)
		if err != nil {
			firstErr = fmt.Errorf("failed to deserialize json %q to volume info: %w", e.Text, err)
			e.Request.Abort()

			return
		}
	})

	var seriesURL = j.options.BaseURL + seriesPath + seriesSlug
	err := j.scrapper.Visit(seriesURL)
	if err != nil {
		return nil, -1, fmt.Errorf("failed call to JNovel Club for %q: %w", seriesURL, err)
	}

	if firstErr != nil {
		return nil, -1, firstErr
	}

	var numVolumes = len(jsonVolumeInfo.Props.PageProps.Aggregate.Volumes)
	var volumes = make([]*sitehandler.VolumeInfo, numVolumes)
	for i, volume := range jsonVolumeInfo.Props.PageProps.Aggregate.Volumes {
		// no release data is present, but this should not happen
		if volume.Volume.Publishing.Seconds == "" {
			return nil, -1, fmt.Errorf("failed to get volume info properly for series %q as there is no publishing data", seriesName)
		}

		secondsFromEpoch, err := strconv.Atoi(volume.Volume.Publishing.Seconds)
		if err != nil {
			return nil, -1, fmt.Errorf("failed to parse out seconds for volume %q: %w", volume.Volume.Title, err)
		}

		var releaseDate = time.Unix(int64(secondsFromEpoch), int64(volume.Volume.Publishing.Nanos))
		volumes[numVolumes-i-1] = &sitehandler.VolumeInfo{
			Name:        volume.Volume.Title,
			ReleaseDate: &releaseDate,
		}
	}

	return volumes, len(volumes), nil
}
