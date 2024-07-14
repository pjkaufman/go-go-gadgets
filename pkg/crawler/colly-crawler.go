package crawler

import (
	"github.com/gocolly/colly/v2"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

func CreateNewCollyCrawler(verbose bool) *colly.Collector {
	c := colly.NewCollector()

	if verbose {
		c.OnRequest(func(r *colly.Request) {
			logger.WriteInfof("Visiting: %v\n", r.URL)
		})

		c.OnResponse(func(r *colly.Response) {
			logger.WriteInfof("Page visited: %v\n", r.Request.URL)
		})

		c.OnScraped(func(r *colly.Response) {
			logger.WriteInfof("Finished visiting: %v\n", r.Request.URL)
		})
	}

	c.OnError(func(_ *colly.Response, err error) {
		logger.WriteErrorf("Something went wrong making an http call: %s\n", err)
	})

	return c
}
