package config

import (
	"fmt"
	"strings"
)

type ReleaseInfo struct {
	Name        string `json:"name"`
	ReleaseDate string `json:"release_date"`
}

type SeriesInfo struct {
	Name                           string        `json:"name,omitempty"`
	TotalVolumes                   int           `json:"total_volumes,omitempty"`
	LatestVolume                   string        `json:"latest_volume,omitempty"`
	UnreleasedVolumes              []ReleaseInfo `json:"unreleased_volumes,omitempty"`
	SlugOverride                   *string       `json:"slug_override,omitempty"`
	Type                           SeriesType    `json:"type,omitempty"`
	Publisher                      PublisherType `json:"publisher,omitempty"`
	Status                         SeriesStatus  `json:"status,omitempty"`
	WikipediaTablesToParseOverride *int          `json:"tables_to_parse_override,omitempty"`
}

type Config struct {
	Series []SeriesInfo `json:"series"`
}

var WikipediaTablesToParseOverrideWarningMsg = fmt.Sprintf("wikipedia tables to parse override is only valid on the publisher %s or %s", OnePeaceBooks, HanashiMedia)

func (c *Config) HasSeries(name string) bool {
	for _, series := range c.Series {
		if strings.EqualFold(name, series.Name) {
			return true
		}
	}

	return false
}

func (c *Config) RemoveSeriesIfExists(name string) bool {
	var newSeries []SeriesInfo
	for _, series := range c.Series {
		if !strings.EqualFold(name, series.Name) {
			newSeries = append(newSeries, series)
		}
	}

	var changeMade = len(newSeries) != len(c.Series)
	if changeMade {
		c.Series = newSeries
	}

	return changeMade
}

func (c *Config) AddSeries(series SeriesInfo, wikipediaTablesToParseOverride int) string {
	var warning string
	if wikipediaTablesToParseOverride > 0 {
		if series.Publisher == OnePeaceBooks || series.Publisher == HanashiMedia {
			series.WikipediaTablesToParseOverride = &wikipediaTablesToParseOverride
		} else {
			warning = WikipediaTablesToParseOverrideWarningMsg
		}
	}

	c.Series = append(c.Series, series)

	return warning
}
