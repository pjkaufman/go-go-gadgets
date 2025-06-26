package vizmedia

import (
	"github.com/gocolly/colly/v2"
	sitehandler "github.com/pjkaufman/go-go-gadgets/magnum/internal/site-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/crawler"
)

type VizMedia struct {
	options  sitehandler.SiteHandlerOptions
	scrapper *colly.Collector
}

func NewVizMediaHandler(options sitehandler.SiteHandlerOptions) sitehandler.SiteHandler {
	return &VizMedia{
		options:  options,
		scrapper: crawler.CreateNewCollyCrawler(options.UserAgent, options.Verbose, options.AllowedDomains),
	}
}
