package sevenseasentertainment

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	sitehandler "github.com/pjkaufman/go-go-gadgets/magnum/internal/site-handler"
)

var volumeNameRegex = regexp.MustCompile(`<h3>([^<]+)</h3>`)
var earlyDigitalAccessRegex = regexp.MustCompile(`<b>Early Digital:</b> ([a-zA-Z]+ \d{2}, \d{4})`)
var releaseDateRegex = regexp.MustCompile(`<b>Release Date</b>: ([a-zA-Z]+ \d{2}, \d{4})`)

func ParseVolumeInfo(series, contentHtml string, volume int) (*sitehandler.VolumeInfo, error) {
	// get name from the anchor in the h3
	var firstHeading = volumeNameRegex.FindStringSubmatch(contentHtml)
	if len(firstHeading) < 2 {
		return nil, fmt.Errorf(`failed to get the name of volume %d for series %q`, volume, series)
	}

	var heading = firstHeading[1]
	var lowercaseHeading = strings.ToLower(heading)
	if strings.Contains(lowercaseHeading, "(audiobook)") || strings.Contains(lowercaseHeading, "hardcover") || strings.Contains(lowercaseHeading, "hardback") {
		return nil, nil
	}

	// get early digital release if present
	var earlyDigitalAccessDateInfo = earlyDigitalAccessRegex.FindStringSubmatch(contentHtml)
	var releaseDateString string
	if len(earlyDigitalAccessDateInfo) > 1 {
		releaseDateString = earlyDigitalAccessDateInfo[1]
	}

	// if not present get release date
	if releaseDateString == "" {
		var releaseDateInfo = releaseDateRegex.FindStringSubmatch(contentHtml)
		if len(releaseDateInfo) > 1 {
			releaseDateString = releaseDateInfo[1]
		}
	}

	var releaseDate *time.Time
	if releaseDateString != "" {
		tempDate, err := time.Parse(releaseDateFormat, releaseDateString)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %q to a date time value: %v", releaseDateString, err)
		}

		releaseDate = &tempDate
	}

	return &sitehandler.VolumeInfo{
		Name:        heading,
		ReleaseDate: releaseDate,
	}, nil
}
