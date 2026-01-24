package rulefixes

import (
	"fmt"
	"strings"
)

func FixManifestAttribute(opfContents, attribute string, lineNum int, elementNameToNumber map[string]int) ([]TextEdit, error) {
	var edits []TextEdit
	lineNum--
	lines := strings.Split(opfContents, "\n")
	if lineNum < 0 || lineNum >= len(lines) {
		return edits, fmt.Errorf("line number out of range")
	}

	// Find the target line
	line := lines[lineNum]
	if !strings.Contains(line, attribute) {
		return edits, fmt.Errorf("attribute not found on the specified line")
	}

	// Determine the element name
	elementStart := strings.Index(line, "<dc:")
	if elementStart == -1 {
		return edits, nil
	}
	elementEnd := strings.Index(line[elementStart:], ">")
	if elementEnd == -1 {
		return edits, fmt.Errorf("malformed element")
	}
	elementEnd += elementStart
	element := line[elementStart : elementEnd+1]

	// Determine the id
	idAttr := ` id="`
	idStart := strings.Index(line, idAttr)
	var id string
	if idStart != -1 {
		idStart += len(idAttr)
		idEnd := strings.Index(line[idStart:], `"`)
		if idEnd == -1 {
			return edits, fmt.Errorf("malformed id attribute")
		}
		id = line[idStart : idStart+idEnd]
	} else {
		var (
			elementName = strings.TrimSuffix(strings.TrimPrefix(element, "<dc:"), ">")
			num         = "1"
		)

		elementName = elementName[0:strings.Index(elementName, " ")]

		if val, ok := elementNameToNumber[elementName]; ok {
			num = fmt.Sprint(val)
			elementNameToNumber[elementName] += 1
		} else {
			elementNameToNumber[elementName] = 2
		}

		id = elementName + num
		insertIdPos := Position{
			Line:   lineNum + 1,
			Column: getColumnForLine(line, elementStart+len(element)-1),
		}
		edits = append(edits, TextEdit{
			Range: Range{
				Start: insertIdPos,
				End:   insertIdPos,
			},
			NewText: fmt.Sprintf(` id="%s"`, id),
		})
	}

	// Parse out the value of the attribute
	attrStart := strings.Index(line, attribute+`="`)
	if attrStart == -1 {
		return edits, fmt.Errorf("attribute not found")
	}

	attrValueStart := attrStart + len(attribute) + 2
	attrEnd := strings.Index(line[attrValueStart:], `"`)
	if attrEnd == -1 {
		return edits, fmt.Errorf("malformed attribute value")
	}
	attrValue := line[attrValueStart : attrValueStart+attrEnd]

	// Remove the attribute from the line
	edits = append(edits, TextEdit{
		Range: Range{
			Start: Position{
				Line:   lineNum + 1,
				Column: getColumnForLine(line, attrStart-1),
			},
			End: Position{
				Line:   lineNum + 1,
				Column: getColumnForLine(line, attrValueStart+attrEnd+1),
			},
		},
	})

	// Create the meta tag
	metaTag := fmt.Sprintf(`<meta refines="#%s" property="%s">%s</meta>`, id, attribute[strings.Index(attribute, ":")+1:], attrValue) + "\n" + getLeadingWhitespace(line)

	newTagInsertPos := Position{
		Line:   lineNum + 2,
		Column: 1,
	}
	edits = append(edits, TextEdit{
		Range: Range{
			Start: newTagInsertPos,
			End:   newTagInsertPos,
		},
		NewText: metaTag,
	})

	return edits, nil
}
