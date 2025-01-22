package epubhandler

import (
	"encoding/xml"
	"fmt"
)

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
	Properties *string  `xml:"properties,attr"`
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

func GetOpfXml(opfContents string) (*Package, error) {
	var opfInfo Package
	err := xml.Unmarshal([]byte(opfContents), &opfInfo)
	if err != nil {
		return nil, fmt.Errorf(ErrorParsingXmlMessageStart+"%v", err)
	}

	return &opfInfo, nil
}
