package jnovelclub

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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
	j.scrapper.OnHTML("script", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "publishing") {
			// this is a bit brittle, but seems to get the job done.
			// it parses out the JSON object, then unquotes the logic since it is now JS instead of JSON
			var jsonText = e.Text[strings.Index(e.Text, "{"):]
			jsonText, _ = strings.CutSuffix(jsonText, "]\\n\"])")

			var err error
			jsonText, err = strconv.Unquote("\"" + jsonText + "\"")
			if err != nil {
				firstErr = fmt.Errorf("failed to unquote JS to get it into proper JSON %q to volume info: %w", jsonText, err)
				e.Request.Abort()

				return
			}

			err = json.Unmarshal([]byte(jsonText), &jsonVolumeInfo)
			if err != nil {
				firstErr = fmt.Errorf("failed to deserialize json %q to volume info: %w", jsonText, err)
				e.Request.Abort()

				return
			}
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

	var numVolumes = len(jsonVolumeInfo.Volumes)
	var volumes = make([]*sitehandler.VolumeInfo, numVolumes)
	for i, volume := range jsonVolumeInfo.Volumes {
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
