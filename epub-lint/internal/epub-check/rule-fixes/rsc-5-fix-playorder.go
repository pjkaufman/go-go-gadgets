package rulefixes

import (
	"fmt"
	"regexp"
	"strings"
)

var navPointRegex = regexp.MustCompile(`(?i)(<navPoint[^>]*)`)

func FixPlayOrder(fileContents string) (edits []TextEdit) {
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
			playOrderAttr  = "playOrder="
			match          = fileContents[navPointIndex[0]:navPointIndex[1]]
			playOrderStart = strings.Index(match, "playOrder=")
		)
		if playOrderStart == -1 {
			insertStartPos := indexToPosition(fileContents, navPointIndex[1])
			edits = append(edits, TextEdit{
				Range: Range{
					Start: insertStartPos,
					End:   insertStartPos,
				},
				NewText: fmt.Sprintf(" playOrder=\"%d\"", i+1),
			})
		} else {
			playOrderStart += len(playOrderAttr) + 1 // include whichever quote

			var (
				playOrderQuote    = string(match[playOrderStart-1])
				playOrderEnd      = strings.Index(match[playOrderStart:], playOrderQuote)
				expectedPlayOrder = fmt.Sprint(i + 1)
			)

			if playOrderEnd == -1 {
				// something went wrong, we are not able to handle this scenario...
				continue
			}

			if match[playOrderStart:playOrderStart+playOrderEnd] != expectedPlayOrder {
				insertStartPos := indexToPosition(fileContents, navPointIndex[0]+playOrderStart)
				insertEndPos := indexToPosition(fileContents, navPointIndex[0]+playOrderStart+playOrderEnd)
				edits = append(edits, TextEdit{
					Range: Range{
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
