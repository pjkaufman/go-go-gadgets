package hanashimedia

import (
	sitehandler "github.com/pjkaufman/go-go-gadgets/magnum/internal/site-handler"
)

type HanashiMedia struct {
	options sitehandler.SiteHandlerOptions
}

func NewHanashiMediaHandler(options sitehandler.SiteHandlerOptions) sitehandler.SiteHandler {
	return &HanashiMedia{
		options: options,
	}
}
