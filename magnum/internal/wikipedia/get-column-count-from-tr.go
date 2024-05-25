package wikipedia

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var colspanRegex = regexp.MustCompile(`colspan=['"](\d+)['"]`)

// GetColumnCountFromTr takes in a table row and returns the amount of actual td elements and
// then the amount of tds taking into account colspan values
func GetColumnCountFromTr(rowHtml string) (int, int, error) {
	var numberOfTableDataEls = strings.Count(rowHtml, `<td`)
	var actualColumnNum = numberOfTableDataEls
	colspanInfo := colspanRegex.FindAllStringSubmatch(rowHtml, -1)
	for _, colspan := range colspanInfo {
		if len(colspan) < 2 {
			return 0, 0, fmt.Errorf(`failed to get the colspan info from row %q as colspan regex match was "%v"`, rowHtml, colspan)
		}

		cols, err := strconv.Atoi(colspan[1])
		if err != nil {
			return 0, 0, fmt.Errorf(`failed to get the colspan column number from row "%s: as colspan match was %q`, rowHtml, colspan[1])
		}

		actualColumnNum += cols - 1
	}

	return numberOfTableDataEls, actualColumnNum, nil
}
