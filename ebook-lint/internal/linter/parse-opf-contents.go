package linter

import (
	"encoding/xml"
	"fmt"
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

type Package struct {
	XMLName  xml.Name  `xml:"package"`
	Version  string    `xml:"version,attr"`
	Manifest *Manifest `xml:"manifest"`
	Guide    *Guide    `xml:"guide"`
}

type Manifest struct {
	XMLName xml.Name        `xml:"manifest"`
	Items   []*ManifestItem `xml:"item"`
}

type ManifestItem struct {
	XMLName    xml.Name `xml:"item"`
	Href       string   `xml:"href,attr"`
	MediaType  string   `xml:"media-type,attr"`
	Properties string   `xml:"properties,attr"`
}

type Guide struct {
	XMLName    xml.Name          `xml:"guide"`
	References []*GuideReference `xml:"reference"`
}

type GuideReference struct {
	XMLName xml.Name `xml:"reference"`
	Href    string   `xml:"href,attr"`
	Type    string   `xml:"type,attr"`
}

const ErrorParsingXmlMessageStart = "error parsing xml: "

var (
	ErrNoPackageInfo   = fmt.Errorf("no package info found for the epub - please verify that the opf has a version in it")
	ErrNoItemEls       = fmt.Errorf("no manifest items found for the epub - please verify that the opf has items in it")
	ErrNoManifest      = fmt.Errorf("no manifest found")
	ErrNoEndOfManifest = fmt.Errorf("manifest is incorrectly formatted since it has no closing manifest element")
)

func ParseOpfFile(text string) (EpubInfo, error) {
	var epubInfo = EpubInfo{
		HtmlFiles:   make(map[string]struct{}),
		ImagesFiles: make(map[string]struct{}),
		OtherFiles:  make(map[string]struct{}),
		CssFiles:    make(map[string]struct{}),
	}

	var opfInfo Package
	err := xml.Unmarshal([]byte(text), &opfInfo)
	if err != nil {
		return epubInfo, fmt.Errorf(ErrorParsingXmlMessageStart+"%v", err)
	}

	epubInfo.Version, err = versionTextToInt(opfInfo.Version)
	if err != nil {
		return epubInfo, err
	}

	if opfInfo.Manifest == nil {
		return epubInfo, ErrNoManifest
	}

	for _, manifestItem := range opfInfo.Manifest.Items {
		var filePath = hrefToFile(manifestItem.Href)
		if epubInfo.Version == 3 && strings.Contains(manifestItem.Properties, "nav") {
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
				epubInfo.TocFile = hrefToFile(guideReference.Href)
				break
			}
		}
	}

	return epubInfo, nil
}

func hrefToFile(href string) string {
	var poundIndex = strings.Index(href, "#")
	if poundIndex == -1 {
		return href
	}

	return href[0:poundIndex]
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
