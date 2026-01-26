package positions

import (
	"fmt"
	"sort"
)

func ApplyEdits(filePath, content string, edits []TextEdit) (string, error) {
	sort.Slice(edits, func(i, j int) bool {
		if edits[i].Range.Start.Line != edits[j].Range.Start.Line {
			return edits[i].Range.Start.Line > edits[j].Range.Start.Line
		}

		return edits[i].Range.Start.Column > edits[j].Range.Start.Column
	})

	for _, e := range edits {
		if e.IsEmpty() {
			continue
		}

		startOffset := GetPositionOffset(content, e.Range.Start.Line, e.Range.Start.Column)
		endOffset := GetPositionOffset(content, e.Range.End.Line, e.Range.End.Column)
		if startOffset < 0 || endOffset < startOffset {
			return "", fmt.Errorf("failed to update %q due to invalid range of %d to %d", filePath, startOffset, endOffset)
		}

		content = content[:startOffset] + e.NewText + content[endOffset:]
	}

	return content, nil
}
