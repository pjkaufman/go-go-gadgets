package wikipedia

import (
	"strings"
)

func ParseDateFromTd(tdHtml string) string {
	var endOfRow = strings.Index(tdHtml, tableDataEndingElIndicator)
	if endOfRow != -1 {
		tdHtml = tdHtml[:endOfRow]
	}

	var digitalVersionIndex = strings.Index(strings.ToLower(tdHtml), "(digital")
	if digitalVersionIndex != -1 {
		tdHtml = tdHtml[:digitalVersionIndex]

		// get values up until the last >
		for i := len(tdHtml) - 1; i >= 0; i-- {
			if tdHtml[i] == '>' {
				tdHtml = tdHtml[i+1:]
				break
			}
		}
	}

	if strings.HasPrefix(tdHtml, "<") {
		tdHtml = tdHtml[strings.Index(tdHtml, ">")+1:]
	}

	var firstOpeningHtmlIndicator = strings.Index(tdHtml, "<")
	if firstOpeningHtmlIndicator != -1 {
		tdHtml = tdHtml[:firstOpeningHtmlIndicator]
	}

	tdHtml = strings.TrimSpace(tdHtml)
	if tdHtml == "—" || tdHtml == "TBA" || strings.Contains(strings.ToLower(tdHtml), "(physical") || strings.Contains(strings.ToLower(tdHtml), "(print") {
		return ""
	}

	return tdHtml
}
