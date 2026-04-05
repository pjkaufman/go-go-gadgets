package rulefixes

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
)

var navPointRegex = regexp.MustCompile(`(?i)(<navPoint[^>]*)`)

func FixPlayOrder(fileContents string) (edits []positions.TextEdit) {
	navPoints := navPointRegex.FindAllStringIndex(fileContents, -1)
	if len(navPoints) == 0 {
		return
	}

	for i, navPointIndex := range navPoints {
		if len(navPointIndex) != 2 {
			continue
		}

		// for some reason the indexes for the groups are not present beyond the first, so I am going to just
		// go ahead and string parse the nav point element
		var (
			match                                        = fileContents[navPointIndex[0]:navPointIndex[1]]
			playOrder, playOrderStart, playOrderEnd, err = epubhandler.GetAttributeValue(match, "playOrder")
		)
		if err != nil { // it is either malformed or not present, so we shall add it...
			insertStartPos := positions.IndexToPosition(fileContents, navPointIndex[1])
			edits = append(edits, positions.TextEdit{
				Range: positions.Range{
					Start: insertStartPos,
					End:   insertStartPos,
				},
				NewText: fmt.Sprintf(" playOrder=\"%d\"", i+1),
			})
		} else {
			var expectedPlayOrder = strconv.Itoa(i + 1)

			if playOrder != expectedPlayOrder {
				insertStartPos := positions.IndexToPosition(fileContents, navPointIndex[0]+playOrderStart)
				insertEndPos := positions.IndexToPosition(fileContents, navPointIndex[0]+playOrderEnd)
				edits = append(edits, positions.TextEdit{
					Range: positions.Range{
						Start: insertStartPos,
						End:   insertEndPos,
					},
					NewText: expectedPlayOrder,
				})
			}
		}
	}

	return
}
