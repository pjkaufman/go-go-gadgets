package vizmedia

import (
	"fmt"
	"regexp"
	"strings"
)

var volumeNameAndRedirectLinkRegex = regexp.MustCompile(`<a[^>]*href=['"](/manga-books[^"']*)['"]>([^<\n]+)</a>`)

func ParseVolumeHtml(html, seriesName string, volume int) (string, string, bool, error) {
	var nameAndLinkInfo = volumeNameAndRedirectLinkRegex.FindStringSubmatch(html)
	if len(nameAndLinkInfo) < 3 {
		return "", "", false, fmt.Errorf(`failed to get the name and or redirect link of volume %d for series %q`, volume, seriesName)
	}

	return nameAndLinkInfo[2], nameAndLinkInfo[1], !strings.Contains(html, "Pre-Order"), nil
}
