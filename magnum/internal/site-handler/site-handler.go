package sitehandler

import "time"

type VolumeInfo struct {
	Name        string
	ReleaseDate *time.Time
}

type SiteHandlerOptions struct {
	BaseURL           string
	ReleaseDateFormat string
	Verbose           bool
}

type ScrapingOptions struct {
	SlugOverride  *string
	TablesToParse *int   // for Wikipedia
	UserAgent     string // for Wikipedia
}

type SiteHandler interface {
	GetVolumeInfo(seriesName string, options ScrapingOptions) ([]*VolumeInfo, int, error)
}
