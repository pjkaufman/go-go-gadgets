package linter

import (
	"fmt"
	"strings"
	"unicode"
)

func FixIdentifierDiscrepancy(opfContents, ncxContents string) (string, error) {
	// Extract the unique identifier from the NCX
	ncxIdentifier, err := getNcxIdentifier(ncxContents)
	if err != nil {
		return "", err
	}

	// Extract the unique identifier from the OPF
	opfIdentifierEl, opfIdentifier, opfIdentifierID := getOpfIdentifier(opfContents)

	// Scenario 1: No unique identifier in OPF, but present in NCX
	if opfIdentifier == "" && ncxIdentifier != "" {
		opfContents = addOpfIdentifier(opfContents, ncxIdentifier)
		return opfContents, nil
	}

	// Scenario 2: Different unique identifier in OPF and NCX and the NCX identifier is not already present
	if opfIdentifier != "" && ncxIdentifier != "" && opfIdentifier != ncxIdentifier && !strings.Contains(opfContents, ">"+ncxIdentifier) {
		opfContents = addOpfIdentifierAndUpdateExistingOne(opfIdentifierEl, opfContents, opfIdentifierID, ncxIdentifier)
		return opfContents, nil
	}

	// Scenario 3: Different unique identifier in OPF and NCX where the OPF has the identifier from the NCX, but it is not the identifier specified in the OPF
	if opfIdentifier != "" && ncxIdentifier != "" && opfIdentifier != ncxIdentifier && strings.Contains(opfContents, ">"+ncxIdentifier) {
		opfContents = moveOpfIdentifierID(opfContents, opfIdentifier, ncxIdentifier, opfIdentifierID, opfIdentifierEl)
		return opfContents, nil
	}

	return opfContents, nil
}

// getNcxIdentifier extracts the unique identifier from the NCX content.
func getNcxIdentifier(ncxContents string) (string, error) {
	startTag := `<meta name="dtb:uid" content="`
	startIndex := strings.Index(ncxContents, startTag)
	if startIndex == -1 {
		return "", fmt.Errorf("unique identifier not found in NCX")
	}
	startIndex += len(startTag)
	endIndex := strings.Index(ncxContents[startIndex:], `"`)
	if endIndex == -1 {
		return "", fmt.Errorf("unique identifier not found in NCX")
	}
	identifier := ncxContents[startIndex : startIndex+endIndex]

	return identifier, nil
}

// getOpfIdentifier extracts the unique identifier from the OPF content.
func getOpfIdentifier(opfContents string) (string, string, string) {
	// Attempt to find the unique-identifier attribute value
	uniqueIdAttr := `unique-identifier="`
	uniqueIdStart := strings.Index(opfContents, uniqueIdAttr)
	var uniqueId string
	if uniqueIdStart != -1 {
		uniqueIdStart += len(uniqueIdAttr)
		uniqueIdEnd := strings.Index(opfContents[uniqueIdStart:], `"`)
		if uniqueIdEnd != -1 {
			uniqueId = opfContents[uniqueIdStart : uniqueIdStart+uniqueIdEnd]
		}
	}

	// If uniqueId is found, try to find the corresponding dc:identifier element
	if uniqueId != "" {
		idTag := fmt.Sprintf(`id="%s"`, uniqueId)
		idStart := strings.Index(opfContents, idTag)
		if idStart != -1 {
			// Find the start of the line containing the id
			lineStart := strings.LastIndex(opfContents[:idStart], "\n") + 1
			lineEnd := strings.Index(opfContents[idStart:], "\n")
			if lineEnd == -1 {
				lineEnd = len(opfContents)
			} else {
				lineEnd += idStart
			}
			fullLine := opfContents[lineStart:lineEnd]

			// Extract the identifier value within the line
			identifierTag := `>`
			identifierStart := strings.Index(fullLine, identifierTag)
			if identifierStart != -1 {
				identifierStart += len(identifierTag)
				identifierEnd := strings.Index(fullLine[identifierStart:], `<`)
				if identifierEnd != -1 {
					identifier := fullLine[identifierStart : identifierStart+identifierEnd]
					return fullLine, identifier, uniqueId
				}
			}
		}
	}

	// Fallback to the first dc:identifier element if uniqueId or matching dc:identifier is not found
	firstIdTag := `<dc:identifier`
	firstIdStart := strings.Index(opfContents, firstIdTag)
	if firstIdStart != -1 {
		firstIdStart += len(firstIdTag)
		firstIdEnd := strings.Index(opfContents[firstIdStart:], `</dc:identifier>`)
		if firstIdEnd != -1 {
			fullLine := opfContents[firstIdStart : firstIdStart+firstIdEnd]

			identifierTag := `>`
			identifierStart := strings.Index(fullLine, identifierTag)
			if identifierStart != -1 {
				identifierStart += len(identifierTag)
				identifierEnd := strings.Index(fullLine[identifierStart:], `<`)
				if identifierEnd != -1 {
					identifier := fullLine[identifierStart : identifierStart+identifierEnd]
					return fullLine, identifier, ""
				}
			}
		}
	}

	return "", "", ""
}

// addOpfIdentifier adds a unique identifier to the OPF content.
func addOpfIdentifier(opfContents, identifier string) string {
	var identifierTag = fmt.Sprintf(`<dc:identifier id="pub-id">%s</dc:identifier>`, identifier)

	metadataEndTag := `</metadata>`
	return strings.Replace(opfContents, metadataEndTag, identifierTag+"\n"+metadataEndTag, 1)
}

// addOpfIdentifierAndUpdateExistingOne replaces the unique identifier in the OPF content.
func addOpfIdentifierAndUpdateExistingOne(oldIdentifierEl, opfContents, identifierID, newIdentifier string) string {
	var (
		idAttribute            = fmt.Sprintf(` id="%s"`, identifierID)
		updatedOldIdentifierEl = strings.Replace(oldIdentifierEl, idAttribute, "", 1)
		format                 strings.Builder
	)

	format.WriteString("\n" + getLeadingWhitespace(oldIdentifierEl))
	format.WriteString("<dc:identifier")
	if identifierID != "" {
		format.WriteString(idAttribute)
	}

	format.WriteString(">")
	format.WriteString(newIdentifier)
	format.WriteString("</dc:identifier>")

	return strings.Replace(opfContents, oldIdentifierEl, updatedOldIdentifierEl+format.String(), 1)
}

// moveOpfIdentifierID moves the identifier's id from the current identifier in the OPF to the other identifier in the OPF that matches the NCX.
// moveOpfIdentifierID updates the identifier line, adding or replacing the id attribute,
// and removes the id attribute from the old identifier element.
func moveOpfIdentifierID(opfContents, opfIdentifier, ncxIdentifier, uniqueId, oldIdentifierEl string) string {
	// Find the line containing the ncxIdentifier
	ncxIdentifierLineStart := strings.Index(opfContents, ncxIdentifier)
	if ncxIdentifierLineStart == -1 {
		return opfContents // ncxIdentifier not found, return the content unchanged
	}

	lineStart := strings.LastIndex(opfContents[:ncxIdentifierLineStart], "\n") + 1
	lineEnd := strings.Index(opfContents[ncxIdentifierLineStart:], "\n")
	if lineEnd == -1 {
		lineEnd = len(opfContents)
	} else {
		lineEnd += ncxIdentifierLineStart
	}

	var (
		line                   = opfContents[lineStart:lineEnd]
		idAttribute            = fmt.Sprintf(` id="%s"`, uniqueId)
		updatedOldIdentifierEl = strings.Replace(oldIdentifierEl, idAttribute, "", 1)
	)

	opfContents = strings.Replace(opfContents, oldIdentifierEl, updatedOldIdentifierEl, 1)

	// Check if the line already has an id attribute
	idAttr := ` id="`
	idStart := strings.Index(line, idAttr)
	if idStart == -1 {
		// No id attribute, add it
		newLine := strings.Replace(line, ">"+ncxIdentifier, fmt.Sprintf(` id="%s">%s`, uniqueId, ncxIdentifier), 1)
		opfContents = strings.Replace(opfContents, line, newLine, 1)
	} else {
		// Replace the existing id attribute value with the uniqueId
		idEnd := strings.Index(line[idStart+len(idAttr):], `"`) + idStart + len(idAttr)
		newLine := line[:idStart+len(idAttr)] + uniqueId + line[idEnd:]
		opfContents = strings.Replace(opfContents, line, newLine, 1)
	}

	return opfContents
}

// getLeadingWhitespace returns the leading whitespace from the input string.
func getLeadingWhitespace(input string) string {
	// Initialize a variable to store the leading whitespace
	var leadingWhitespace strings.Builder

	// Iterate over the string to find leading whitespace characters
	for _, char := range input {
		if unicode.IsSpace(char) {
			leadingWhitespace.WriteRune(char)
		} else {
			break
		}
	}

	return leadingWhitespace.String()
}
