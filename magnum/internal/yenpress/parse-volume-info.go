package yenpress

import (
	"fmt"
	"html"
	"strings"
)

const nameElIndicator = `<p class="paragraph"><span>`

func ParseVolumeInfo(series, contentHtml string) (*VolumeInfo, error) {
	var nameStart = strings.Index(contentHtml, nameElIndicator)
	if nameStart == -1 {
		return nil, fmt.Errorf("failed to find the start of a volume name for a volume in series %q with html content %q", series, contentHtml)
	}

	nameStart += len(nameElIndicator)
	var nameEnd = strings.Index(contentHtml[nameStart:], "</span>")
	if nameEnd == -1 {
		return nil, fmt.Errorf("failed to find the end of a volume name for a volume in series %q with html content %q", series, contentHtml)
	}

	var name = html.UnescapeString(contentHtml[nameStart : nameStart+nameEnd])

	var lowercaseName = strings.ToLower(name)
	if strings.Contains(lowercaseName, "collector's edition") || strings.Contains(lowercaseName, "omnibus edition") {
		return nil, nil
	}

	var relativeLinkStart = strings.Index(contentHtml, `href="`)
	if relativeLinkStart == -1 {
		return nil, fmt.Errorf("failed to find the start of the relative link for a volume for series %q with the given html %q", series, contentHtml)
	}

	relativeLinkStart += 6

	var relativeLinkEnd = strings.Index(contentHtml[relativeLinkStart:], `"`)
	if relativeLinkEnd == -1 {
		return nil, fmt.Errorf("failed to find the end of the relative link for a volume for series %q with the given html %q", series, contentHtml)
	}

	var relativeLink = contentHtml[relativeLinkStart : relativeLinkStart+relativeLinkEnd]

	return &VolumeInfo{
		RelativeLink: relativeLink,
		Name:         name,
	}, nil
}
