package tableofcontents

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/converter"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

// BuildMultipleBookAlphabeticalListItems
// Note: that alphabetical order will not be perfect given the discrepancy between some of the names in the digital vs. book versions
func BuildMultipleBookAlphabeticalListItems(_ []string, headerInfo []converter.MdFileInfo) (string, error) {
	if len(headerInfo) == 0 {
		return "", nil
	}

	// TODO: should some scenarios be removed from the following?
	for i, headerData := range headerInfo {
		if headerData.AlternateTitle == "" {
			continue
		}

		var (
			primaryPageCount   = len(headerData.PrimaryPageNumbers)
			secondaryPageCount = len(headerData.SecondaryPageNumbers)
			hasPrimaryPages    = primaryPageCount > 0
			hasSecondaryPages  = secondaryPageCount > 0
		)
		if hasPrimaryPages && !hasSecondaryPages {
			switch primaryPageCount {
			case 2: // likely wrong, but not sure I can do much about it right now...
				headerData.Header = headerData.AlternateTitle
				headerData.AlternateTitle = ""
				headerInfo = append(headerInfo, headerData)

				continue
			case 1:
				headerData.PrimaryPageNumbers = append(headerData.PrimaryPageNumbers, headerData.PrimaryPageNumbers...)
				headerData.Header = headerData.AlternateTitle
				headerData.AlternateTitle = ""
				headerInfo = append(headerInfo, headerData)

				continue
			}
		} else if hasSecondaryPages && !hasPrimaryPages {
			switch secondaryPageCount {
			case 2:
				headerData.HasBeenDuplicatedForSecondaryOnly = true
				headerInfo[i] = headerData

				headerData.Header = headerData.AlternateTitle
				headerData.AlternateTitle = ""
				headerInfo = append(headerInfo, headerData)

				continue
			case 1:
				headerData.HasBeenDuplicatedForSecondaryOnly = true
				headerData.SecondaryPageNumbers = append(headerData.SecondaryPageNumbers, headerData.SecondaryPageNumbers...)
				headerInfo[i] = headerData

				headerData.Header = headerData.AlternateTitle
				headerData.AlternateTitle = ""
				headerInfo = append(headerInfo, headerData)

				continue
			}
		} else if primaryPageCount == 1 && secondaryPageCount == 1 {
			var secondaryPages = headerData.SecondaryPageNumbers

			headerData.SecondaryPageNumbers = []int{}
			headerInfo[i] = headerData

			headerData.PrimaryPageNumbers = headerData.SecondaryPageNumbers
			headerData.SecondaryPageNumbers = secondaryPages
			headerData.Header = headerData.AlternateTitle
			headerData.AlternateTitle = ""
			headerInfo = append(headerInfo, headerData)

			continue
		}

		return "", fmt.Errorf("Encountered an unhandled situation for creating book ToC entries: title %q; alternate title %q; primary page count %d; secondary page count %d\n", headerData.Header, headerData.AlternateTitle, primaryPageCount, secondaryPageCount)
	}

	c := collate.New(language.English)
	sort.Slice(headerInfo, func(i, j int) bool {
		if headerInfo[i].Header != headerInfo[j].Header {
			return c.CompareString(headerInfo[i].Header, headerInfo[j].Header) < 0
		}

		return c.CompareString(headerInfo[i].FileName, headerInfo[j].FileName) < 0
	})

	var (
		primaryPageNumberIndex   = make(map[string]int)
		secondaryPageNumberIndex = make(map[string]int)
		listItems                = strings.Builder{}
		startingLetter           = ""
	)
	fmt.Fprintln(&listItems, `<div class="divider">&#9836;</div>`)
	for _, headerData := range headerInfo {
		if startingLetter == "" {
			startingLetter = headerData.Header[0:1]
		} else if startingLetter != headerData.Header[0:1] {
			fmt.Fprintln(&listItems, `<div class="divider">&#9836;</div>`)
			startingLetter = headerData.Header[0:1]
		}

		addToCEntry(&listItems, headerData.FileName, headerData.Header, "%d", headerData.PrimaryPageNumbers, primaryPageNumberIndex)
		if _, ok := primaryPageNumberIndex[headerData.FileName]; !ok {
			// we can have scenarios where we have an alternate and regular title for a book that is only a secondary page number, so if we have two tiles for the same song
			// we need to not iterate over all secondary numbers
			if headerData.HasBeenDuplicatedForSecondaryOnly {
				addToCEntry(&listItems, headerData.FileName, headerData.Header, "(%d)", headerData.SecondaryPageNumbers, secondaryPageNumberIndex)
			} else {
				for range headerData.SecondaryPageNumbers {
					addToCEntry(&listItems, headerData.FileName, headerData.Header, "(%d)", headerData.SecondaryPageNumbers, secondaryPageNumberIndex)
				}
			}
		}
	}

	return listItems.String(), nil
}

func addToCEntry(tocItems *strings.Builder, fileName, header, pageFormat string, pageNumbers []int, pageNumberIndex map[string]int) {
	if len(pageNumbers) == 0 {
		return
	}

	var pageNumber int
	if val, ok := pageNumberIndex[fileName]; ok {
		pageNumber = pageNumbers[val]
	} else {
		pageNumber = pageNumbers[0]
		pageNumberIndex[fileName] = 1
	}

	fmt.Fprintf(tocItems, `<li><span class="name">%s</span><span class="page">%s</span></li>`+"\n", header, fmt.Sprintf(pageFormat, pageNumber))
}
