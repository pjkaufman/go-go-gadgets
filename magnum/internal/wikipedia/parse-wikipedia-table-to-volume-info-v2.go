package wikipedia

import (
	"fmt"
	"strings"
	"time"

	sitehandler "github.com/pjkaufman/go-go-gadgets/magnum/internal/site-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

func ParseWikipediaTableToVolumeInfoV2(namePrefix, tableHtml string) ([]*sitehandler.VolumeInfo, error) {
	var rows = volumeRowHeaderRegex.FindAllStringSubmatch(tableHtml, -1)
	if len(rows) == 0 {
		return nil, fmt.Errorf("failed to find table row info for: %s", namePrefix)
	}

	var volumeInfo = []*sitehandler.VolumeInfo{}
	var rowHtml = tableHtml
	var startOfRow, endOfRow int
	var releaseDateString string
	var hasValidAmountOfColumns bool
	var err error
	for _, rowSubmatches := range rows {
		startOfRow = strings.Index(rowHtml, rowSubmatches[0])
		rowHtml = rowHtml[startOfRow:]
		endOfRow = strings.Index(rowHtml, wikiTableRowEnd)

		releaseDateString, hasValidAmountOfColumns, err = getEnglishReleaseDateFromRow(rowHtml[:endOfRow])
		if err != nil {
			return nil, fmt.Errorf("failed to parse rows for %q: %w", namePrefix, err)
		}

		if !hasValidAmountOfColumns {
			logger.WriteWarnf("skipped rows for %q since it did not have the expected amount of rows", namePrefix)
			return volumeInfo, nil
		}
		var date *time.Time
		if releaseDateString != "" {
			tempDate, err := time.Parse(releaseDateFormat, releaseDateString)
			if err != nil {
				return nil, fmt.Errorf("failed to parse %q to a date time value: %w", releaseDateString, err)
			}

			date = &tempDate
		}

		volumeInfo = append(volumeInfo, &sitehandler.VolumeInfo{
			Name:        fmt.Sprintf("%s Vol. %s", namePrefix, strings.TrimSpace(rowSubmatches[1])),
			ReleaseDate: date,
		})

		rowHtml = rowHtml[endOfRow:]
	}

	return volumeInfo, nil
}
