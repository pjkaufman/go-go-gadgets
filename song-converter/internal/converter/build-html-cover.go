package converter

import (
	"fmt"
	"strings"
	"time"
)

func BuildHtmlCover(coverMd, bookType string, currentTime time.Time) string {
	coverMd = strings.Replace(coverMd, "{{DATE_GENERATED}}", currentTime.Format("Jan 2006"), 1)
	coverMd = strings.Replace(coverMd, "{{TYPE}}", bookType, 1)

	coverHtml := mdToHTML([]byte(coverMd))
	coverHtml = fmt.Sprintf("<div style=\"text-align: center\">\n%s</div>\n", coverHtml)

	return strings.ReplaceAll(coverHtml, "\n\n", "\n")
}
