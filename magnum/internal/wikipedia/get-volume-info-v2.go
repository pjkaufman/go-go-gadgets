package wikipedia

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gocolly/colly/v2"
	sitehandler "github.com/pjkaufman/go-go-gadgets/magnum/internal/site-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

func (w *Wikipedia) GetVolumeInfo(seriesName string, options sitehandler.ScrapingOptions) ([]*sitehandler.VolumeInfo, int, error) {
	var seriesSlug string
	if options.SlugOverride != nil {
		seriesSlug = *options.SlugOverride
	} else {
		seriesSlug = convertTitleToSlug(seriesName)
	}

	sections, err := w.api.GetSectionInfo(seriesSlug)
	if err != nil {
		return nil, -1, err
	}

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
		return nil, -1, fmt.Errorf("failed to get light novel section")
	}

	var contentHtml string
	w.scrapper.OnHTML("#content > div.vector-page-toolbar", func(e *colly.HTMLElement) {
		var content = e.DOM.Parent()
		contentHtml, err = content.Html()

		if err != nil {
			logger.WriteErrorf("failed to get content body: %s\n", err)
		}
	})

	var lnHeadingHtml string
	var startIndexOfLnSection int
	w.scrapper.OnHTML("#"+lnSection.Anchor, func(e *colly.HTMLElement) {
		var parents = e.DOM.Parent()
		lnHeadingHtml, err = parents.Html()
		if err != nil {
			logger.WriteErrorf("failed to get content body: %s\n", err)
		}

		startIndexOfLnSection = strings.Index(contentHtml, lnHeadingHtml)
		if startIndexOfLnSection == -1 {
			logger.WriteErrorf("failed to find light novel section: %s\n", err)
		}
	})

	var lnSectionAfterLnHtml string
	var endIndexOfLnSection int = -1
	if sectionAfterLn.Heading != "" {
		w.scrapper.OnHTML("#"+sectionAfterLn.Anchor, func(e *colly.HTMLElement) {
			var parents = e.DOM.Parent()
			lnSectionAfterLnHtml, err = parents.Html()
			if err != nil {
				logger.WriteErrorf("failed to get section after light novel section: %s\n", err)
			}

			endIndexOfLnSection = strings.Index(contentHtml, lnSectionAfterLnHtml)
			if endIndexOfLnSection == -1 {
				logger.WriteErrorf("failed to find section after light novel section: %s\n", err)
			}
		})
	}

	var url = w.options.BaseURL + wikiPath + seriesSlug
	err = w.scrapper.Visit(url)
	if err != nil {
		return nil, -1, fmt.Errorf("failed call to wikipedia for %q: %w", url, err)
	}

	var lnSectionHtml string
	if endIndexOfLnSection != -1 {
		lnSectionHtml = contentHtml[startIndexOfLnSection:endIndexOfLnSection]
	} else {
		lnSectionHtml = contentHtml[startIndexOfLnSection:]
	}

	if len(subSectionTiles) == 0 {
		subSectionTiles = []string{seriesName}
	}

	var numTables = strings.Count(lnSectionHtml, "wikitable")
	if numTables == 0 {
		return nil, -1, fmt.Errorf("could not find tables for light novel section: %w", err)
	} else if len(subSectionTiles)+1 == numTables {
		subSectionTiles = append([]string{seriesName}, subSectionTiles...)
	} else if len(subSectionTiles) != numTables {
		return nil, -1, fmt.Errorf("number of tables does not match number of table title prefixes for %q: %d vs. %d", seriesName, len(subSectionTiles), numTables)
	}

	var volumeInfo = []*sitehandler.VolumeInfo{}
	for i, subSectionTitle := range subSectionTiles {
		if options.TablesToParseOverride != nil && *options.TablesToParseOverride < i+1 {
			break
		}

		tableHtml, stop := GetNextTableAndItsEndPosition(lnSectionHtml)
		volumes, err := ParseWikipediaTableToVolumeInfoV2(subSectionTitle, tableHtml)
		if err != nil {
			return nil, -1, err
		}

		volumeInfo = append(volumeInfo, volumes...)
		lnSectionHtml = lnSectionHtml[stop:]
	}

	slices.Reverse(volumeInfo)

	return volumeInfo, len(volumeInfo), nil
}
