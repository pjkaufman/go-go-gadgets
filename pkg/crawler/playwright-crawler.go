package crawler

import (
	"fmt"

	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/playwright-community/playwright-go"
)

func CreateNewPlaywrightCrawler() (*playwright.Playwright, playwright.Browser, playwright.Page) {
	pw, err := playwright.Run()
	if err != nil {
		logger.WriteError(fmt.Sprintf("could not start playwright: %v", err))
	}

	browser, err := pw.Chromium.Launch()
	if err != nil {
		logger.WriteError(fmt.Sprintf("could not launch browser: %v", err))
	}

	page, err := browser.NewPage()
	if err != nil {
		logger.WriteError(fmt.Sprintf("could not create page: %v", err))
	}

	return pw, browser, page
}

func ClosePlaywrightCrawler(pw *playwright.Playwright, browser playwright.Browser) {
	err := browser.Close()
	if err != nil {
		logger.WriteError(fmt.Sprintf("could not close browser: %v", err))
	}

	if pw == nil {
		return
	}

	err = pw.Stop()
	if err != nil {
		logger.WriteError(fmt.Sprintf("could not stop Playwright: %v", err))
	}
}
