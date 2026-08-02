package converter

import (
	"fmt"
	"strings"
	"time"
)

func BuildHtmlCover(coverMd, bookType, extraStyleCss string, currentTime time.Time) string {
	if extraStyleCss != "" {
		extraStyleCss = "; " + extraStyleCss
	}

	coverMd = strings.Replace(coverMd, "{{DATE_GENERATED}}", currentTime.Format("Jan 2006"), 1)
	coverMd = strings.Replace(coverMd, "{{TYPE}}", bookType, 1)

	coverHtml := mdToHTML([]byte(coverMd))
	// TODO: allow specifying left for the text alignment for the book type
	coverHtml = fmt.Sprintf("<div style=\"text-align: center%s\">\n%s</div>\n", extraStyleCss, coverHtml)
	coverHtml = strings.ReplaceAll(coverHtml, "&amp;nbsp;", "&nbsp;")

	return strings.ReplaceAll(coverHtml, "\n\n", "\n")
}
