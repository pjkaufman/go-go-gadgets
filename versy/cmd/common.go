package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

const (
	bibleGatewayBaseUrl = "https://www.biblegateway.com/"
	searchPath          = "passage/"
	searchParamFormat   = "?search=%s&version=%s"
	userAgent           = "Versy/1.0"
)

var (
	errFailedToGetVerseOfDay                = errors.New("failed to get the verse of the day as it is not present")
	errFailedToGetVerseReferenceForLanguage = errors.New("failed to get the verse reference in the proper language as it is not present")
	allowedDomains                          = []string{"www.biblegateway.com"}
)

func getVerse(reference, version string, scrapper *colly.Collector) (verse string, err error) {
	var referenceInLanguage string
	scrapper.OnHTML(".bcv > .dropdown-display > .dropdown-display-text", func(e *colly.HTMLElement) {
		referenceInLanguage = strings.TrimSpace(e.Text)
	})

	var (
		verseContent string
		firstErr     error
	)
	scrapper.OnHTML(".passage-content", func(e *colly.HTMLElement) {
		var verseHtml string
		verseHtml, firstErr = e.DOM.Html()
		if firstErr != nil {
			e.Request.Abort()

			return
		}

		verseContent = htmlToPlaintext(verseHtml)
	})

	var url = bibleGatewayBaseUrl + searchPath + fmt.Sprintf(searchParamFormat, url.QueryEscape(reference), version)
	err = scrapper.Visit(url)
	if err != nil {
		err = fmt.Errorf("failed call to Bible Gateway for %q: %w", url, err)

		return
	}

	if firstErr != nil {
		err = fmt.Errorf("failed to get verse content for %s %s: %w", reference, version, firstErr)

		return
	}

	if referenceInLanguage == "" {
		err = errFailedToGetVerseReferenceForLanguage
		return
	}

	referenceInLanguage = strings.Replace(referenceInLanguage, "-", " - ", 1)

	return fmt.Sprintf(`%s %q (%s)`, referenceInLanguage, verseContent, version), nil
}

func getVerseOfTheDayReference(scrapper *colly.Collector) (reference string, err error) {
	scrapper.OnHTML(".passage-box > .verse-bar > a", func(e *colly.HTMLElement) {
		reference = strings.TrimSpace(e.Attr("title"))
	})

	var url = bibleGatewayBaseUrl
	err = scrapper.Visit(url)
	if err != nil {
		err = fmt.Errorf("failed call to Bible Gateway for %q: %w", url, err)

		return
	}

	if reference == "" {
		err = errFailedToGetVerseOfDay
		return
	}

	return reference, nil
}

var (
	supRegex     = regexp.MustCompile(`<sup[^>]*>.*?</sup>`)
	headingRegex = regexp.MustCompile(`<h[1-6][^>]*>.*?</h[1-6]>`)
	spaceRegex   = regexp.MustCompile(`\s+`)
	tagRegex     = regexp.MustCompile(`<[^>]*>`)
)

func htmlToPlaintext(htmlContent string) string {
	// fmt.Println("Test: " + htmlContent)
	// Remove all sup elements and their content
	result := supRegex.ReplaceAllString(htmlContent, "")

	// fmt.Println("Test: " + result)

	// Remove all heading elements (h1-h6) and their content
	result = headingRegex.ReplaceAllString(result, "")

	// Remove any remaining HTML tags
	result = tagRegex.ReplaceAllString(result, "")

	// Remove HTML encoded non-breaking spaces
	result = strings.ReplaceAll(result, "&nbsp;", "")
	result = strings.ReplaceAll(result, "&#160;", "")
	result = strings.ReplaceAll(result, "\u00a0", " ")

	var startIndexOfExtraneousText = strings.Index(result, "Read full chapter")
	if startIndexOfExtraneousText != -1 {
		result = result[:startIndexOfExtraneousText]
	}

	// convert Spanish quotes into English quotes
	result = strings.ReplaceAll(result, "«", "")
	result = strings.ReplaceAll(result, "»", "")

	// Clean up extra whitespace
	result = strings.TrimSpace(result)
	result = spaceRegex.ReplaceAllString(result, " ")

	return result
}

func getAndDisplayBothVerses(reference, version1, version2 string, scrapper *colly.Collector) {
	firstVerse, err := getVerse(reference, version1, scrapper)
	if err != nil {
		logger.WriteError(err.Error())
	}

	logger.WriteInfo(firstVerse)
	logger.WriteInfo("")

	secondVerse, err := getVerse(reference, version2, scrapper.Clone())
	if err != nil {
		logger.WriteError(err.Error())
	}

	logger.WriteInfo(secondVerse)

}
