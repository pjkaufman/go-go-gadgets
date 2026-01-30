package rulefixes

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
)

func FixFileNotFound(contents, referencedResource, currentFile string, line, column int, basenameToFilePaths map[string][]string) (positions.TextEdit, error) {
	var edit positions.TextEdit
	// the line, column combination should be the end of "/>" for self-closing elements and ">" on the opening el for other elements
	offset := positions.GetPositionOffset(contents, line, column)
	if offset == -1 {
		return edit, nil
	}

	var startOfEl = strings.LastIndex(contents[:offset], "<")
	if startOfEl == -1 {
		return edit, nil
	}

	var (
		startingEl         = contents[startOfEl:offset]
		attributeIndicator = `src="`
		// it is possible for it to be a single quote, but we will cross that bridge when we get to it
		attributeIndex = strings.Index(startingEl, attributeIndicator)
	)
	if attributeIndex == -1 {
		attributeIndicator = `href="`
		attributeIndex = strings.Index(startingEl, attributeIndicator)

		if attributeIndex == -1 { // this should never happen, but maybe it will
			return edit, nil
		}
	}

	attributeIndex += len(attributeIndicator)

	endOfAttributeIndex := strings.Index(startingEl[attributeIndex:], `"`)
	if endOfAttributeIndex == -1 {
		return edit, nil
	}

	// TODO: is it possible for the attribute value to be a URL?
	// I guess it is, but it needs more research
	var (
		basename      = filepath.Base(referencedResource)
		possibleFiles = basenameToFilePaths[basename]
	)

	if len(possibleFiles) != 0 {
		// determine relative path between current and the first possible file
		// which may want to swap to be the best match down the road...
		relativePath, err := filepath.Rel(filepath.Dir(currentFile), filepath.Dir(possibleFiles[0]))
		if err != nil {
			return edit, fmt.Errorf("failed to determine the relative file path for %q referenced in %q: %w", possibleFiles[0], currentFile, err)
		}

		edit.Range.Start = positions.IndexToPosition(contents, startOfEl+attributeIndex)
		edit.Range.End = positions.IndexToPosition(contents, startOfEl+attributeIndex+endOfAttributeIndex)
		edit.NewText = filepath.Join(relativePath, basename)

		return edit, nil
	}

	var endOfElName = strings.Index(startingEl, " ")
	if endOfElName == -1 {
		// I don't think this can happen since it would have to have a src or href present...
		return edit, nil
	}

	var (
		isSelfClosing = contents[offset-2] == '/'
		endOfEl       = offset
	)

	if !isSelfClosing {
		var closingEl = fmt.Sprintf("</%s>", startingEl[1:endOfElName])

		endOfEl = strings.Index(contents[offset:], closingEl)
		if endOfEl == -1 {
			return edit, nil
		}

		endOfEl += len(closingEl) + offset
	}

	// check if line will be empty
	var (
		startOfLine = strings.LastIndex(contents[:startOfEl], "\n") + 1
		endOfLine   = strings.Index(contents[endOfEl:], "\n")
	)

	if endOfLine == -1 {
		endOfLine = len(contents) - 1
	} else {
		endOfLine += endOfEl
	}

	var (
		lineContents = strings.Replace(contents[startOfLine:endOfLine], contents[startOfEl:endOfEl], "", 1)
		beginReplace = startOfEl
		endReplace   = endOfEl
		removeLine   bool
	)
	if strings.TrimSpace(lineContents) == "" {
		beginReplace = startOfLine
		endReplace = endOfLine
		removeLine = true
	}

	edit.Range.Start = positions.IndexToPosition(contents, beginReplace)
	edit.Range.End = positions.IndexToPosition(contents, endReplace)
	if removeLine {
		edit.Range.End.Column = 1
		edit.Range.End.Line += 1
	}

	return edit, nil
}
