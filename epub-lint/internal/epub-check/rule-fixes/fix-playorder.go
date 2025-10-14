package rulefixes

import (
	"fmt"
	"regexp"
	"strings"
)

var navPointRegex = regexp.MustCompile(`(?i)(<navPoint[^>]*)(playOrder="\d*")?`)

func FixPlayOrder(fileContents string) string {
	playOrder := 1

	updatedContent := navPointRegex.ReplaceAllStringFunc(fileContents, func(match string) string {
		var (
			playOrderStart = strings.Index(match, "playOrder=")
			result         string
		)
		if playOrderStart != -1 {
			result = fmt.Sprintf("%s playOrder=\"%d\"", match[:playOrderStart-1], playOrder)
		} else {
			result = fmt.Sprintf("%s playOrder=\"%d\"", match, playOrder)
		}

		playOrder++

		return result
	})

	return updatedContent
}
