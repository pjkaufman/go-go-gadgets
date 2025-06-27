package sitehandler

import "time"

type VolumeInfo struct {
	Name        string
	ReleaseDate *time.Time
}

type ApiPathBuilder func(string, string, string) string

type SiteHandlerOptions struct {
	BaseURL           string
	ReleaseDateFormat string
	Verbose           bool
	UserAgent         string
	BuildApiPath      ApiPathBuilder // for Wikipedia
	AllowedDomains    []string
}

type ScrapingOptions struct {
	SlugOverride          *string
	TablesToParseOverride *int // for Wikipedia
}

type SiteHandler interface {
	GetVolumeInfo(seriesName string, options ScrapingOptions) ([]*VolumeInfo, int, error)
}
