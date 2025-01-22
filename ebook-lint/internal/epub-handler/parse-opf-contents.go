package epubhandler

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type EpubInfo struct {
	HtmlFiles   map[string]struct{}
	ImagesFiles map[string]struct{}
	CssFiles    map[string]struct{}
	OtherFiles  map[string]struct{}
	NcxFile     string
	NavFile     string
	TocFile     string
	Version     int
}

var (
	ErrNoPackageInfo   = errors.New("no package info found for the epub - please verify that the opf has a version in it")
	ErrNoItemEls       = errors.New("no manifest items found for the epub - please verify that the opf has items in it")
	ErrNoManifest      = errors.New("no manifest found")
	ErrNoEndOfManifest = errors.New("manifest is incorrectly formatted since it has no closing manifest element")
)

func ParseOpfFile(text string) (EpubInfo, error) {
	var epubInfo = EpubInfo{
		HtmlFiles:   make(map[string]struct{}),
		ImagesFiles: make(map[string]struct{}),
		OtherFiles:  make(map[string]struct{}),
		CssFiles:    make(map[string]struct{}),
	}

	opfInfo, err := GetOpfXml(text)
	if err != nil {
		return epubInfo, err
	}

	epubInfo.Version, err = versionTextToInt(opfInfo.Version)
	if err != nil {
		return epubInfo, err
	}

	if opfInfo.Manifest == nil {
		return epubInfo, ErrNoManifest
	}

	var filePath string
	for _, manifestItem := range opfInfo.Manifest.Items {
		filePath, err = hrefToFile(manifestItem.Href)
		if err != nil {
			return epubInfo, fmt.Errorf("failed to convert manifest href %q to file path: %w", manifestItem.Href, err)
		}

		if epubInfo.Version == 3 && manifestItem.Properties != nil && strings.Contains(*manifestItem.Properties, "nav") {
			epubInfo.NavFile = filePath
		}

		if strings.Contains(manifestItem.MediaType, "xhtml") {
			epubInfo.HtmlFiles[filePath] = struct{}{}
		} else if strings.Contains(manifestItem.MediaType, "image") {
			epubInfo.ImagesFiles[filePath] = struct{}{}
		} else if strings.Contains(manifestItem.MediaType, "css") {
			epubInfo.CssFiles[filePath] = struct{}{}
		} else {
			if strings.HasSuffix(filePath, ".ncx") {
				epubInfo.NcxFile = filePath
			}

			epubInfo.OtherFiles[filePath] = struct{}{}
		}
	}

	if len(opfInfo.Manifest.Items) == 0 {
		return epubInfo, ErrNoItemEls
	}

	if opfInfo.Guide != nil {
		for _, guideReference := range opfInfo.Guide.References {
			if guideReference.Type == "toc" {
				epubInfo.TocFile, err = hrefToFile(guideReference.Href)
				if err != nil {
					return epubInfo, fmt.Errorf("failed to convert toc href %q to file path: %w", guideReference.Href, err)
				}

				break
			}
		}
	}

	return epubInfo, nil
}

func hrefToFile(href string) (string, error) {
	var poundIndex = strings.Index(href, "#")
	if poundIndex == -1 {
		return url.QueryUnescape(href)
	}

	return url.QueryUnescape(href[0:poundIndex])
}

func versionTextToInt(versionText string) (int, error) {
	versionText = strings.TrimSpace(versionText)
	if versionText == "" {
		return 0, ErrNoPackageInfo
	}

	var periodIndex = strings.Index(versionText, ".")
	if periodIndex != -1 {
		versionText = versionText[0:periodIndex]
	}

	version, err := strconv.Atoi(versionText)
	if err != nil {
		return 0, fmt.Errorf(`failed to convert version text %q to an integer: %w`, versionText, err)
	}

	return version, nil
}
