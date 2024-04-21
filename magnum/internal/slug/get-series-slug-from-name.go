package slug

import (
	"regexp"
	"strings"
)

var seriesInvalidSlugCharacters = regexp.MustCompile(`[\(\),:?!]`)

func GetSeriesSlugFromName(seriesName string) string {
	var slug = strings.ToLower(seriesInvalidSlugCharacters.ReplaceAllString(seriesName, ""))

	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "'", "-")
	slug = strings.ReplaceAll(slug, "--", "-")

	return slug
}
