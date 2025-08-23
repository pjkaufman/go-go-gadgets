package epubhandler

import (
	"fmt"
	"strings"
)

func AddFileToNcx(text, filePath, title, id string) string {
	navMapClose := "</navMap>"
	navMapIdx := strings.Index(text, navMapClose)
	if navMapIdx == -1 {
		return text
	}

	var playOrder = strings.Count(text, "<navPoint") + 1

	navPoint := fmt.Sprintf(
		`  <navPoint id="%s" playOrder="%d">
    <navLabel>
      <text>%s</text>
    </navLabel>
    <content src="%s"/>
  </navPoint>
`, id, playOrder, title, filePath)

	result := text[:navMapIdx] + navPoint + text[navMapIdx:]
	return result
}
