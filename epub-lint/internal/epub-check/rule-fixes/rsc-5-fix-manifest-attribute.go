package rulefixes

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
)

func FixManifestAttribute(opfContents, attribute string, lineNum int, elementNameToNumber map[string]int) ([]positions.TextEdit, error) {
	var edits []positions.TextEdit
	lineNum--
	lines := strings.Split(opfContents, "\n")
	if lineNum < 0 || lineNum >= len(lines) {
		return edits, errors.New("line number out of range")
	}

	// Find the target line
	line := lines[lineNum]
	if !strings.Contains(line, attribute) {
		return edits, errors.New("attribute not found on the specified line")
	}

	// Determine the element name
	elementStart := strings.Index(line, "<dc:")
	if elementStart == -1 {
		return edits, nil
	}
	elementEnd := strings.Index(line[elementStart:], ">")
	if elementEnd == -1 {
		return edits, errors.New("malformed element")
	}
	elementEnd += elementStart
	element := line[elementStart : elementEnd+1]

	// Determine the id
	id, _, _, err := epubhandler.GetAttributeValue(line, "id")
	if err != nil { // we will assume that any parsing error means no id for now, we can amend this if that is not the case
		var (
			elementName = strings.TrimSuffix(strings.TrimPrefix(element, "<dc:"), ">")
			num         = "1"
		)

		elementName = elementName[0:strings.Index(elementName, " ")]

		if val, ok := elementNameToNumber[elementName]; ok {
			num = strconv.Itoa(val)
			elementNameToNumber[elementName] += 1
		} else {
			elementNameToNumber[elementName] = 2
		}

		id = elementName + num
		insertIdPos := positions.Position{
			Line:   lineNum + 1,
			Column: positions.GetColumnForLine(line, elementStart+len(element)-1),
		}
		edits = append(edits, positions.TextEdit{
			Range: positions.Range{
				Start: insertIdPos,
				End:   insertIdPos,
			},
			NewText: fmt.Sprintf(` id="%s"`, id),
		})
	}

	attrValue, attrStart, attrEnd, err := epubhandler.GetAttributeValue(line, attribute)
	if err != nil {
		return edits, err
	}

	// Remove the attribute from the line
	edits = append(edits, positions.TextEdit{
		Range: positions.Range{
			Start: positions.Position{
				Line:   lineNum + 1,
				Column: positions.GetColumnForLine(line, attrStart-len(attribute)-3), // account for "=", quote, and attribute name
			},
			End: positions.Position{
				Line:   lineNum + 1,
				Column: positions.GetColumnForLine(line, attrEnd+1),
			},
		},
	})

	// Create the meta tag
	metaTag := getLeadingWhitespace(line) + fmt.Sprintf(`<meta refines="#%s" property="%s">%s</meta>`, id, attribute[strings.Index(attribute, ":")+1:], attrValue) + "\n"

	newTagInsertPos := positions.Position{
		Line:   lineNum + 2,
		Column: 1,
	}
	edits = append(edits, positions.TextEdit{
		Range: positions.Range{
			Start: newTagInsertPos,
			End:   newTagInsertPos,
		},
		NewText: metaTag,
	})

	return edits, nil
}
