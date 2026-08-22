package converter

import (
	"fmt"
	"strings"
	"time"
)

func BuildHtmlCover(coverMd, bookType, titleAlignment, extraStyleCss string, currentTime time.Time) string {
	if extraStyleCss != "" {
		extraStyleCss = "; " + extraStyleCss
	}

	coverMd = strings.Replace(coverMd, "{{DATE_GENERATED}}", currentTime.Format("Jan 2006"), 1)
	coverMd = strings.Replace(coverMd, "{{TYPE}}", bookType, 1)

	coverHtml := mdToHTML([]byte(coverMd))
	coverHtml = fmt.Sprintf("<div style=\"text-align: %s%s\">\n%s</div>\n", titleAlignment, extraStyleCss, coverHtml)
	coverHtml = strings.ReplaceAll(coverHtml, "&amp;nbsp;", "&nbsp;")
	coverHtml = strings.ReplaceAll(coverHtml, `id="songs"`, `id="songs-2"`) // causes issues if there is a songs id in the header

	return strings.ReplaceAll(coverHtml, "\n\n", "\n")
}
