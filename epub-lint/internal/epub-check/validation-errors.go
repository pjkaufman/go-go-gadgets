package epubcheck

import (
	"slices"
	"sort"
	"strings"
)

type ValidationErrors struct {
	ValidationIssues []ValidationError
}

type ValidationError struct {
	Code     string
	FilePath string
	Location *Position
	Message  string
}

type Position struct {
	Line   int
	Column int
}

func (ve *ValidationErrors) DecrementLineNumbersAndRemoveLineReferences(lineNum int, path string) {
	for i := 0; i < len(ve.ValidationIssues); i++ {
		if ve.ValidationIssues[i].Location != nil {
			if ve.ValidationIssues[i].FilePath == path {
				if ve.ValidationIssues[i].Location.Line == lineNum {
					ve.ValidationIssues = slices.Delete(ve.ValidationIssues, i, i+1)
					i--
				} else if ve.ValidationIssues[i].Location.Line > lineNum {
					ve.ValidationIssues[i].Location.Line--
				}
			}
		}
	}
}

// Sort sorts the ValidationIssues in the following order:
// 1. Deleted line fixes
// 2. Path Ascending
// 3. Line descending (nil positions will be considered to be after the last line in the file)
// 4. Column descending
func (ve *ValidationErrors) Sort() {
	sort.Slice(ve.ValidationIssues, func(i, j int) bool {
		msgI := ve.ValidationIssues[i]
		msgJ := ve.ValidationIssues[j]

		// Prioritize delete-required messages
		if strings.HasPrefix(msgI.Message, EmptyMetadataProperty) && !strings.HasPrefix(msgJ.Message, EmptyMetadataProperty) {
			return true
		}

		if !strings.HasPrefix(msgI.Message, EmptyMetadataProperty) && strings.HasPrefix(msgJ.Message, EmptyMetadataProperty) {
			return false
		}

		// Compare by path ascending
		if msgI.FilePath != msgJ.FilePath {
			return msgI.FilePath < msgJ.FilePath
		}

		if msgI.Location == nil && msgJ.Location == nil {
			return true
		}

		if msgI.Location == nil {
			return false
		} else if msgJ.Location == nil {
			return true
		}

		// If paths are the same, compare by line descending
		if msgI.Location.Line != msgJ.Location.Line {
			return msgI.Location.Line > msgJ.Location.Line
		}
		// If lines are the same, compare by column descending
		return msgI.Location.Column > msgJ.Location.Column
	})
}
