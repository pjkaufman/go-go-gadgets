package tableofcontents

import (
	"fmt"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/converter"
)

func BuildDigitalListItems(headerIds []string, _ []converter.MdFileInfo) (string, error) {
	if len(headerIds) == 0 {
		return "", nil
	}

	var listItems = strings.Builder{}
	for _, headerId := range headerIds {
		fmt.Fprintf(&listItems, `<li><a href="#%s"></a></li>`+"\n", headerId)
	}

	return listItems.String(), nil
}
