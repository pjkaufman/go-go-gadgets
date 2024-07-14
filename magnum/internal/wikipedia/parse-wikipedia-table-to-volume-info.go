package wikipedia

import (
	"fmt"
	"strings"
	"time"

	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

func ParseWikipediaTableToVolumeInfo(namePrefix, tableHtml string) []VolumeInfo {
	var rows = volumeRowHeaderRegex.FindAllStringSubmatch(tableHtml, -1)
	if len(rows) == 0 {
		logger.WriteError("failed to find table row info for: " + namePrefix)
	}

	var volumeInfo = []VolumeInfo{}
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
			logger.WriteErrorf("failed to parse rows for %q: %s\n", namePrefix, err)
		}

		if !hasValidAmountOfColumns {
			logger.WriteWarnf("skipped rows for %q since it did not have the expected amount of rows\n", namePrefix)
			return volumeInfo
		}
		var date *time.Time
		if releaseDateString != "" {
			tempDate, err := time.Parse(releaseDateFormat, releaseDateString)
			if err != nil {
				logger.WriteErrorf("failed to parse %q to a date time value: %v\n", releaseDateString, err)
			}

			date = &tempDate
		}

		volumeInfo = append(volumeInfo, VolumeInfo{
			Name:        fmt.Sprintf("%s Vol. %s", namePrefix, strings.TrimSpace(rowSubmatches[1])),
			ReleaseDate: date,
		})

		rowHtml = rowHtml[endOfRow:]
	}

	return volumeInfo
}

func getEnglishReleaseDateFromRow(rowHtml string) (string, bool, error) {
	numTds, actualColumns, err := GetColumnCountFromTr(rowHtml)
	if err != nil {
		return "", false, err
	}

	expectedDateColumn, ok := columnAmountToExpectedColumn[actualColumns]
	if !ok || expectedDateColumn > numTds {
		return "", false, nil
	}

	var releaseDateColumn = rowHtml
	for i := 0; i < expectedDateColumn; i++ {
		releaseDateColumn = releaseDateColumn[strings.Index(releaseDateColumn, tableDataStartingElIndicator)+4:]
	}

	return ParseDateFromTd(releaseDateColumn), true, nil
}
