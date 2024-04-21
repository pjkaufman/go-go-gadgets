package converter

import (
	"fmt"
	"strings"
)

func BuildHtmlCover(coverMd string) string {
	coverHtml := mdToHTML([]byte(coverMd))
	coverHtml = fmt.Sprintf("<div style=\"text-align: center\">\n%s</div>\n", coverHtml)

	return strings.ReplaceAll(coverHtml, "\n\n", "\n")
}
