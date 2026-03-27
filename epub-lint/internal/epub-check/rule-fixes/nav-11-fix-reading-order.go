package rulefixes

import (
	"fmt"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
)

type navItemPosInfo struct {
	fullPath                         string
	anchorTag                        string
	anchorStartIndex, anchorEndIndex int
}

func FixReadingOrder(spineOrder []string, navContents, navPath, opfFolder string) (edits []positions.TextEdit) {
	const tocEpubType = `epub:type="toc"`
	epubTypeIndex := strings.Index(navContents, tocEpubType)
	if epubTypeIndex == -1 {
		return
	}

	var (
		startOfEl = strings.LastIndex(navContents[:epubTypeIndex], "<")
		endOfEl   = strings.Index(navContents[epubTypeIndex:], ">")
	)
	if startOfEl == -1 || endOfEl == -1 {
		return
	}

	endOfEl += 1 + epubTypeIndex

	elNameEndIndex := strings.Index(navContents[startOfEl:], " ")
	if elNameEndIndex == -1 {
		return
	}

	elNameEndIndex += startOfEl
	elName := navContents[startOfEl+1 : elNameEndIndex]

	endOfElIndex := strings.Index(navContents[endOfEl:], fmt.Sprintf("</%s>", elName))
	if endOfElIndex == -1 {
		return
	}

	endOfElIndex += endOfEl

	var (
		startOfContent      = endOfEl
		remainingTocContent = navContents[endOfEl:endOfElIndex]
		navItemInfo         = make([]navItemPosInfo, 0, len(spineOrder))
		navItemInfoSorted   []navItemPosInfo
		sortWeights         = make(map[string]int, len(spineOrder))
	)

	for i, path := range spineOrder {
		var fullPath = path
		if opfFolder != "" && opfFolder != "." {
			fullPath = filepath.Join(opfFolder, fullPath)
		}

		sortWeights[fullPath] = i
	}

	const (
		hrefAttribute  = `href="`
		anchorTagStart = "<a "
		anchorTagEnd   = "</a>"
	)

	var (
		anchorStartIndex, anchorEndIndex int
	)
	for anchorStartIndex != -1 {
		anchorStartIndex = strings.Index(remainingTocContent, anchorTagStart)
		if anchorStartIndex == -1 {
			break
		}

		anchorStartIndex += len(anchorTagStart)

		anchorEndIndex = strings.Index(remainingTocContent[anchorStartIndex:], anchorTagEnd)
		if anchorEndIndex == -1 {
			remainingTocContent = remainingTocContent[anchorStartIndex:]
			startOfContent += anchorStartIndex
			continue
		}

		anchorEndIndex += anchorStartIndex + len(anchorTagEnd)

		var (
			anchorTag = remainingTocContent[anchorStartIndex:anchorEndIndex]
			hrefIndex = strings.Index(anchorTag, hrefAttribute)
		)
		if hrefIndex == -1 {
			remainingTocContent = remainingTocContent[anchorEndIndex:]
			startOfContent += anchorEndIndex
			continue
		}

		hrefIndex += len(hrefAttribute)
		endOfHref := strings.Index(anchorTag[hrefIndex:], `"`)
		if endOfHref == -1 {
			remainingTocContent = remainingTocContent[anchorEndIndex:]
			startOfContent += anchorEndIndex
			continue
		}

		endOfHref += hrefIndex

		var (
			filePath       = anchorTag[hrefIndex:endOfHref]
			referenceIndex = strings.Index(filePath, "#")
		)

		if referenceIndex != -1 {
			filePath = filePath[:referenceIndex]
		}

		var fullPath = filepath.Join(navPath, filePath)

		navItemInfo = append(navItemInfo, navItemPosInfo{
			fullPath:         fullPath,
			anchorTag:        anchorTag,
			anchorStartIndex: startOfContent + anchorStartIndex,
			anchorEndIndex:   startOfContent + anchorEndIndex,
		})

		navItemInfoSorted = append(navItemInfoSorted, navItemPosInfo{
			fullPath:         fullPath,
			anchorTag:        anchorTag,
			anchorStartIndex: startOfContent + anchorStartIndex,
			anchorEndIndex:   startOfContent + anchorEndIndex,
		})

		remainingTocContent = remainingTocContent[anchorEndIndex:]
		startOfContent += anchorEndIndex
	}

	slices.Reverse(navItemInfo)

	sort.Slice(navItemInfoSorted, func(i, j int) bool {
		return sortWeights[navItemInfoSorted[i].fullPath] > sortWeights[navItemInfoSorted[j].fullPath]
	})

	for i, navInfo := range navItemInfoSorted {
		if navInfo.fullPath == navItemInfo[i].fullPath {
			continue
		}

		edits = append(edits, positions.TextEdit{
			Range: positions.Range{
				Start: positions.IndexToPosition(navContents, navItemInfo[i].anchorStartIndex),
				End:   positions.IndexToPosition(navContents, navItemInfo[i].anchorEndIndex),
			},
			NewText: navInfo.anchorTag,
		})
	}

	return
}
