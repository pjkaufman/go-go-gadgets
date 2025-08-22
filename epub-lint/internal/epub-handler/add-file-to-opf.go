package cmdhandler

import (
	"fmt"
	"strings"
)

func AddFileToOpf(text, filename, id, mediaType string) string {
	itemEntry := fmt.Sprintf(`<item id="%s" href="%s" media-type="%s"/>`, id, filename, mediaType)
	itemrefEntry := fmt.Sprintf(`<itemref idref="%s"/>`, id)

	manifestClose := "</manifest>"
	manifestIndex := strings.Index(text, manifestClose)
	if manifestIndex != -1 {
		text = text[:manifestIndex] + "  " + itemEntry + "\n" + text[manifestIndex:]
	}

	spineClose := "</spine>"
	spineIndex := strings.Index(text, spineClose)
	if spineIndex != -1 {
		text = text[:spineIndex] + "  " + itemrefEntry + "\n" + text[spineIndex:]
	}

	return text
}
