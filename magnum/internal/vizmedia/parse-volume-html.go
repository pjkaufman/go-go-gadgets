package vizmedia

import (
	"fmt"
	"strings"
)

const (
	volumeTypeIndicator         = `<a class="color-mid-gray hover-red">`
	parseVolumeHtmlEndIndicator = "</"
	volumeLinkIndicator         = `<a class="color-off-black hover-red" href="/manga-books`
)

func ParseVolumeHtml(html, seriesName string, volume int) (string, string, bool, error) {
	volumeTypeStartIndex := strings.Index(html, volumeTypeIndicator)
	if volumeTypeStartIndex == -1 {
		return "", "", false, fmt.Errorf("failed to get the start of the type of literature on sale for volume %d of series %q", volume, seriesName)
	}

	volumeTypeStartIndex += len(volumeTypeIndicator)

	volumeTypeEndIndex := strings.Index(html[volumeTypeStartIndex:], parseVolumeHtmlEndIndicator)
	if volumeTypeEndIndex == -1 {
		return "", "", false, fmt.Errorf("failed to get the end of the type of literature on sale for volume %d of series %q", volume, seriesName)
	}

	volumeTypeEndIndex += volumeTypeStartIndex

	var volumeType = strings.TrimSpace(html[volumeTypeStartIndex:volumeTypeEndIndex])
	if volumeType != "Manga" {
		return "", "", false, nil
	}

	var volumeLinkStartIndex = strings.Index(html, volumeLinkIndicator)
	if volumeLinkStartIndex == -1 {
		return "", "", false, fmt.Errorf("failed to get the start of the link of literature on sale for volume %d of series %q", volume, seriesName)
	}

	volumeLinkStartIndex += len(volumeLinkIndicator)

	volumeLinkEndIndex := strings.Index(html[volumeLinkStartIndex:], "\"")
	if volumeLinkEndIndex == -1 {
		return "", "", false, fmt.Errorf("failed to get the end of the link of literature on sale for volume %d of series %q", volume, seriesName)
	}

	volumeLinkEndIndex += volumeLinkStartIndex

	var volumeLink = html[volumeLinkStartIndex:volumeLinkEndIndex]
	var volumeNameStartIndex = strings.Index(html[volumeLinkEndIndex:], ">")
	if volumeNameStartIndex == -1 {
		return "", "", false, fmt.Errorf("failed to get the start of the name of literature on sale for volume %d of series %q", volume, seriesName)
	}

	volumeNameStartIndex += volumeLinkEndIndex + 1

	volumeNameEndIndex := strings.Index(html[volumeNameStartIndex:], parseVolumeHtmlEndIndicator)
	if volumeNameEndIndex == -1 {
		return "", "", false, fmt.Errorf("failed to get the end of the name of literature on sale for volume %d of series %q", volume, seriesName)
	}

	volumeNameEndIndex += volumeNameStartIndex

	var volumeName = strings.TrimSpace(html[volumeNameStartIndex:volumeNameEndIndex])
	if volumeName == "" {
		return "", "", false, fmt.Errorf("the volume name is empty for volume %d of series %q", volume, seriesName)
	}

	if strings.Contains(strings.ToLower(volumeName), "box set") {
		return "", "", false, nil
	}

	return volumeName, "/manga-books" + volumeLink, !strings.Contains(html, "Pre-Order"), nil
}
