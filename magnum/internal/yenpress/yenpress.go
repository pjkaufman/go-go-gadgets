package yenpress

import (
	"github.com/gocolly/colly/v2"
	sitehandler "github.com/pjkaufman/go-go-gadgets/magnum/internal/site-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/crawler"
)

type YenPress struct {
	options  sitehandler.SiteHandlerOptions
	scrapper *colly.Collector
}

func NewYenPressHandler(options sitehandler.SiteHandlerOptions) sitehandler.SiteHandler {
	return &YenPress{
		options:  options,
		scrapper: crawler.CreateNewCollyCrawler(options.Verbose),
	}
}
