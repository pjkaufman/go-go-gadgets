package jnovelclub

import (
	"github.com/gocolly/colly/v2"
	sitehandler "github.com/pjkaufman/go-go-gadgets/magnum/internal/site-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/crawler"
)

type JNovelClub struct {
	options  sitehandler.SiteHandlerOptions
	scrapper *colly.Collector
}

func NewJNovelClubHandler(options sitehandler.SiteHandlerOptions) sitehandler.SiteHandler {
	return &JNovelClub{
		options:  options,
		scrapper: crawler.CreateNewCollyCrawler(options.UserAgent, options.Verbose, options.AllowedDomains),
	}
}
