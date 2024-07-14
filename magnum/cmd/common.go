package cmd

import (
	"fmt"
	"time"

	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

var (
	verbose bool
)

const (
	defaultReleaseDate = "TBA"
	releaseDateFormat  = "January 2, 2006"
	userAgent          = "Magnum/1.0"
)

func getUnreleasedVolumeDisplayText(unreleasedVol, releaseDate string) string {
	if releaseDate == defaultReleaseDate {
		return fmt.Sprintf("%q release has not been announced yet", unreleasedVol)
	}

	return fmt.Sprintf("%q releases on %s", unreleasedVol, releaseDate)
}

func unreleasedDateIsBeforeDate(releaseDate string, date time.Time) bool {
	if releaseDate == defaultReleaseDate {
		return false
	}

	release, err := time.Parse(releaseDateFormat, releaseDate)
	if err != nil {
		logger.WriteErrorf("failed to convert date %q to date time: %s\n", releaseDate, err)
	}

	return release.Before(date)
}
