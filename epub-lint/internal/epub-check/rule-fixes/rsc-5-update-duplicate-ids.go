package rulefixes

import (
	"fmt"
	"sort"
	"strings"
)

// UpdateDuplicateIds finds and renames duplicate IDs in file contents.
func UpdateDuplicateIds(contents, id string) (edits []TextEdit) {
	var indexes = getAllIndexesInStringForLastCharOfSubstring(contents, "id=\""+id+"\"")
	indexes = append(indexes, getAllIndexesInStringForLastCharOfSubstring(contents, "id='"+id+"'")...)

	if len(indexes) <= 1 {
		return
	}

	sort.Ints(indexes)

	for i := 1; i < len(indexes); i++ {
		start := indexToPosition(contents, indexes[i])
		edits = append(edits, TextEdit{
			Range: Range{
				Start: start,
				End:   start,
			},
			NewText: fmt.Sprintf("_%d", i+1),
		})
	}

	return
}

func getAllIndexesInStringForLastCharOfSubstring(content, substring string) []int {
	var indexes []int
	offset := 0

	for {
		index := strings.Index(content[offset:], substring)
		if index == -1 {
			break
		}

		index += len(substring) - 1 // should get us to the last character before the single or double quote

		indexes = append(indexes, offset+index)
		offset += index + 1
	}

	return indexes
}
