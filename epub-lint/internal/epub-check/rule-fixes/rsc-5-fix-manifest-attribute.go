package rulefixes

import (
	"fmt"
	"strings"
)

// TODO: swap to lsp update method...
func FixManifestAttribute(opfContents, attribute string, lineNum int, elementNameToNumber map[string]int) (string, error) {
	lines := strings.Split(opfContents, "\n")
	if lineNum < 0 || lineNum >= len(lines) {
		return opfContents, fmt.Errorf("line number out of range")
	}

	// Find the target line
	line := lines[lineNum]
	if !strings.Contains(line, attribute) {
		return opfContents, fmt.Errorf("attribute not found on the specified line")
	}

	// Determine the element name
	elementStart := strings.Index(line, "<dc:")
	if elementStart == -1 {
		return opfContents, nil
	}
	elementEnd := strings.Index(line[elementStart:], ">")
	if elementEnd == -1 {
		return opfContents, fmt.Errorf("malformed element")
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
			return opfContents, fmt.Errorf("malformed id attribute")
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
		line = strings.Replace(line, element, fmt.Sprintf(`<%s id="%s">`, element[1:len(element)-1], id), 1)
	}

	// Parse out the value of the attribute
	attrStart := strings.Index(line, attribute+`="`)
	if attrStart == -1 {
		return opfContents, fmt.Errorf("attribute not found")
	}
	attrStart += len(attribute) + 2
	attrEnd := strings.Index(line[attrStart:], `"`)
	if attrEnd == -1 {
		return opfContents, fmt.Errorf("malformed attribute value")
	}
	attrValue := line[attrStart : attrStart+attrEnd]

	// Remove the attribute from the line
	line = strings.Replace(line, fmt.Sprintf(` %s="%s"`, attribute, attrValue), "", 1)

	// Create the meta tag
	metaTag := "\n" + getLeadingWhitespace(line) + fmt.Sprintf(`<meta refines="#%s" property="%s">%s</meta>`, id, attribute[strings.Index(attribute, ":")+1:], attrValue)

	lines[lineNum] = line + metaTag

	// Join the lines back together
	updatedOpfContents := strings.Join(lines, "\n")

	return updatedOpfContents, nil
}
