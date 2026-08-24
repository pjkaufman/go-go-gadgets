package tableofcontents

import (
	"fmt"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/converter"
)

func BuildBookPageOrderListItems(_ []string, headerInfo []converter.MdFileInfo) (string, error) {
	if len(headerInfo) == 0 {
		return "", nil
	}

	for i, headerData := range headerInfo {
		if len(headerData.PrimaryPageNumbers) == 0 {
			return "", fmt.Errorf("Unable to determine where a song goes due to no primary page numbers: title %q; alternate title %q\n", headerData.Header, headerData.AlternateTitle)
		}

		if len(headerData.SecondaryPageNumbers) != 0 {
			return "", fmt.Errorf("Songs for single book page order cannot have secondary page numbers: title %q; alternate title %q; secondary page numbers: %v\n", headerData.Header, headerData.AlternateTitle, headerData.SecondaryPageNumbers)
		}

		if len(headerData.PrimaryPageNumbers) == 1 {
			continue
		}

		var originalPrimaryPageNumbers = headerData.PrimaryPageNumbers
		headerData.PrimaryPageNumbers = []int{originalPrimaryPageNumbers[0]}
		headerInfo[i] = headerData
		for i = 1; i < len(originalPrimaryPageNumbers); i++ {
			headerData.PrimaryPageNumbers = []int{originalPrimaryPageNumbers[i]}
			headerInfo = append(headerInfo, headerData)
		}
	}

	headerInfo = converter.SortSongs(headerInfo)

	var (
		listItems = strings.Builder{}
	)
	for _, headerData := range headerInfo {
		fmt.Fprintf(&listItems, `<li><span class="name">%s</span><span class="page">%d</span></li>`+"\n", headerData.Header, headerData.PrimaryPageNumbers[0])
	}

	return listItems.String(), nil
}
