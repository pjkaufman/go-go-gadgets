package wikipedia

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/pjkaufman/go-go-gadgets/pkg/crawler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

type VolumeInfo struct {
	Name        string
	ReleaseDate *time.Time
}

var wikiTableRegex = regexp.MustCompile(`<table[^>]*class="wikitable"[^>]*>`)
var volumeRowHeaderRegex = regexp.MustCompile(`<th[^>]*scope="row"[^>]*>([^<]*)</th>`)
var columnAmountToExpectedColumn = map[int]int{
	4: 3,
	5: 4,
	6: 4,
}

const (
	tableStart                   = `<table`
	tableEnd                     = `</table>`
	wikiTableRowEnd              = `</tr>`
	tableDataStartingElIndicator = "<td"
	tableDataEndingElIndicator   = "</td"
)

func GetVolumeInfo(userAgent, title string, slugOverride *string, tablesToParseOverride *int, verbose bool) []VolumeInfo {
	var seriesSlug string
	if slugOverride != nil {
		seriesSlug = *slugOverride
	} else {
		seriesSlug = convertTitleToSlug(title)
	}

	sections := getSectionInfo(userAgent, seriesSlug)
	var lnSection SectionInfo
	var sectionAfterLn SectionInfo
	var subSectionTiles []string
	for _, section := range sections.Parse.Sections {
		if lnSection.Anchor != "" {
			if section.Level <= lnSection.Level {
				sectionAfterLn = section
				break
			} else {
				var heading = section.Heading
				var htmlElEndIndicatorIndex = strings.Index(heading, ">")
				if htmlElEndIndicatorIndex != -1 {
					heading = heading[htmlElEndIndicatorIndex+1:]
					heading = heading[:strings.Index(heading, "<")]
				}

				subSectionTiles = append(subSectionTiles, heading)
			}

			continue
		}

		if strings.HasPrefix(strings.ToLower(section.Heading), "light novel") {
			lnSection = section
		}
	}

	if lnSection.Heading == "" {
		logger.WriteError("failed to get light novel section")
	}

	c := crawler.CreateNewCollyCrawler(verbose)

	var err error

	var contentHtml string
	c.OnHTML("#content > div.vector-page-toolbar", func(e *colly.HTMLElement) {
		var content = e.DOM.Parent()
		contentHtml, err = content.Html()

		if err != nil {
			logger.WriteError(fmt.Sprintf("failed to get content body: %s", err))
		}
	})

	var lnHeadingHtml string
	var startIndexOfLnSection int
	c.OnHTML("#"+lnSection.Anchor, func(e *colly.HTMLElement) {
		var parents = e.DOM.Parent()
		lnHeadingHtml, err = parents.Html()
		if err != nil {
			logger.WriteError(fmt.Sprintf("failed to get content body: %s", err))
		}

		startIndexOfLnSection = strings.Index(contentHtml, lnHeadingHtml)
		if startIndexOfLnSection == -1 {
			logger.WriteError(fmt.Sprintf("failed to find light novel section: %s", err))
		}
	})

	var lnSectionAfterLnHtml string
	var endIndexOfLnSection int = -1
	if sectionAfterLn.Heading != "" {
		c.OnHTML("#"+sectionAfterLn.Anchor, func(e *colly.HTMLElement) {
			var parents = e.DOM.Parent()
			lnSectionAfterLnHtml, err = parents.Html()
			if err != nil {
				logger.WriteError(fmt.Sprintf("failed to get section after light novel section: %s", err))
			}

			endIndexOfLnSection = strings.Index(contentHtml, lnSectionAfterLnHtml)
			if endIndexOfLnSection == -1 {
				logger.WriteError(fmt.Sprintf("failed to find section after light novel section: %s", err))
			}
		})
	}

	var url = baseURL + wikiPath + seriesSlug
	err = c.Visit(url)
	if err != nil {
		logger.WriteError(fmt.Sprintf("failed call to wikipedia for \"%s\": %s", url, err))
	}

	var lnSectionHtml string
	if endIndexOfLnSection != -1 {
		lnSectionHtml = contentHtml[startIndexOfLnSection:endIndexOfLnSection]
	} else {
		lnSectionHtml = contentHtml[startIndexOfLnSection:]
	}

	if len(subSectionTiles) == 0 {
		subSectionTiles = []string{title}
	}

	var numTables = strings.Count(lnSectionHtml, "wikitable")
	if numTables == 0 {
		logger.WriteError(fmt.Sprintf("could not find tables for light novel section: %s", err))
	} else if len(subSectionTiles)+1 == numTables {
		subSectionTiles = append([]string{title}, subSectionTiles...)
	} else if len(subSectionTiles) != numTables {
		logger.WriteError(fmt.Sprintf("number of tables does not match number of table title prefixes for \"%s\": %d vs. %d", title, len(subSectionTiles), numTables))
	}

	var volumeInfo = []VolumeInfo{}
	for i, subSectionTitle := range subSectionTiles {
		if tablesToParseOverride != nil && *tablesToParseOverride < i+1 {
			break
		}

		tableHtml, stop := GetNextTableAndItsEndPosition(lnSectionHtml)
		volumeInfo = append(volumeInfo, ParseWikipediaTableToVolumeInfo(subSectionTitle, tableHtml)...)
		lnSectionHtml = lnSectionHtml[stop:]
	}

	slices.Reverse(volumeInfo)

	return volumeInfo
}

func convertTitleToSlug(title string) string {
	return strings.ReplaceAll(title, " ", "_")
}
