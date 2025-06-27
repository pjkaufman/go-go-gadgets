package crawler

import (
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

func CreateNewCollyCrawler(userAgent string, verbose bool, allowedDomains []string) *colly.Collector {
	c := colly.NewCollector()

	var agent = strings.TrimSpace(userAgent)
	if agent != "" {
		c.UserAgent = agent
	}

	c.AllowedDomains = allowedDomains

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
