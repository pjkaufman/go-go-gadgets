package sevenseasentertainment

import (
	"github.com/gocolly/colly/v2"
	sitehandler "github.com/pjkaufman/go-go-gadgets/magnum/internal/site-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/crawler"
)

type SevenSeasEntertainment struct {
	options  sitehandler.SiteHandlerOptions
	scrapper *colly.Collector
}

func NewSevenSeasEntertainmentHandler(options sitehandler.SiteHandlerOptions) sitehandler.SiteHandler {
	return &SevenSeasEntertainment{
		options:  options,
		scrapper: crawler.CreateNewCollyCrawler(options.Verbose),
	}
}
