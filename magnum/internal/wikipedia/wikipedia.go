package wikipedia

import (
	"github.com/gocolly/colly/v2"
	sitehandler "github.com/pjkaufman/go-go-gadgets/magnum/internal/site-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/crawler"
)

type Wikipedia struct {
	options  sitehandler.SiteHandlerOptions
	scrapper *colly.Collector
	api      *WikipediaApi
}

func NewWikipediaHandler(options sitehandler.SiteHandlerOptions) sitehandler.SiteHandler {
	return &Wikipedia{
		options:  options,
		scrapper: crawler.CreateNewCollyCrawler(options.UserAgent, options.Verbose, options.AllowedDomains),
		api:      NewWikipediaApi(options.BaseURL, options.UserAgent, options.Verbose, options.ApiPath),
	}
}
