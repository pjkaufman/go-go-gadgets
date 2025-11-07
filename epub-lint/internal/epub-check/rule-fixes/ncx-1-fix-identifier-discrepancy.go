package rulefixes

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

	var (
		indexOfEndTag    = strings.Index(opfContents, metadataEndTag)
		textUntilEndTag  = opfContents[:indexOfEndTag]
		hasNcxIdentifier = strings.Contains(opfContents, ">"+ncxIdentifier+"<")
	)

	if opfIdentifier == "" && ncxIdentifier != "" {
		// Scenario 1: No unique identifier in OPF, but an identifier el exists and matches the NCX id
		if hasNcxIdentifier {
			opfContents = moveOrSetOpfIdentifierID(opfContents, ncxIdentifier, opfIdentifierID, "")

			return opfContents, nil
		}

		// Scenario 2: No unique identifier in OPF, but it is present in NCX
		var (
			previousNewLineIndex = strings.LastIndex(textUntilEndTag, "\n")
			previousNewLine      string
		)
		if previousNewLineIndex == -1 {
			previousNewLine = textUntilEndTag
		} else {
			previousNewLine = opfContents[previousNewLineIndex+1 : indexOfEndTag]
		}

		opfContents = addOpfIdentifier(opfContents, ncxIdentifier, opfIdentifierID, previousNewLine)
		return opfContents, nil
	}

	// Scenario 3: Different unique identifier in OPF and NCX and the NCX identifier is not already present
	if opfIdentifier != "" && ncxIdentifier != "" && opfIdentifier != ncxIdentifier && !hasNcxIdentifier {
		opfContents = addOpfIdentifierAndUpdateExistingOne(opfIdentifierEl, opfContents, opfIdentifierID, ncxIdentifier)

		return opfContents, nil
	}

	// Scenario 4: Different unique identifier in OPF and NCX where the OPF has the identifier from the NCX, but it is not the identifier specified in the OPF
	if opfIdentifier != "" && ncxIdentifier != "" && opfIdentifier != ncxIdentifier && hasNcxIdentifier {
		opfContents = moveOrSetOpfIdentifierID(opfContents, ncxIdentifier, opfIdentifierID, opfIdentifierEl)
		return opfContents, nil
	}

	return opfContents, nil
}

// getNcxIdentifier extracts the unique identifier from the NCX content.
func getNcxIdentifier(ncxContents string) (string, error) {
	startTag := `name="dtb:uid"`
	startIndex := strings.Index(ncxContents, startTag)
	if startIndex == -1 {
		return "", fmt.Errorf("unique identifier not found in NCX")
	}

	// Find the line containing the startTag
	lineStart := strings.LastIndex(ncxContents[:startIndex], "\n") + 1
	lineEnd := strings.Index(ncxContents[startIndex:], "\n")
	if lineEnd == -1 {
		lineEnd = len(ncxContents)
	} else {
		lineEnd += startIndex
	}
	line := ncxContents[lineStart:lineEnd]

	// Extract the content attribute value
	contentAttr := `content="`
	contentStart := strings.Index(line, contentAttr)
	if contentStart == -1 {
		return "", fmt.Errorf("content attribute not found in the line")
	}
	contentStart += len(contentAttr)
	contentEnd := strings.Index(line[contentStart:], `"`)
	if contentEnd == -1 {
		return "", fmt.Errorf("content attribute value not found in the line")
	}
	content := line[contentStart : contentStart+contentEnd]

	return content, nil
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

	return "", "", uniqueId
}

// addOpfIdentifier adds a unique identifier to the OPF content.
func addOpfIdentifier(opfContents, identifier, identifierID, metadataEndElPriorToEl string) string {
	if identifierID == "" {
		identifierID = "pub-id"
	}
	var (
		identifierTag = fmt.Sprintf(`<dc:identifier id="%s">%s</dc:identifier>`, identifierID, identifier)
		endingEl      = metadataEndTag
	)
	if strings.TrimSpace(metadataEndElPriorToEl) == "" {
		// Assuming the metadata tag was on its own line, double the space
		// behind the identifierTag since that should make the tag align
		// with the others make sure the manifest tag has the same indentation as it did
		if metadataEndElPriorToEl == "" {
			identifierTag = "\t" + identifierTag
		} else {
			endingEl = metadataEndElPriorToEl + endingEl
			identifierTag = metadataEndElPriorToEl + identifierTag
		}
	} else {
		var currentLineWhitespace = getLeadingWhitespace(metadataEndElPriorToEl)
		identifierTag = currentLineWhitespace + identifierTag
		endingEl = getMetadataWhitespaceForNewLine(currentLineWhitespace) + endingEl
	}

	return strings.Replace(opfContents, metadataEndTag, identifierTag+"\n"+endingEl, 1)
}

// addOpfIdentifierAndUpdateExistingOne replaces the unique identifier in the OPF content.
func addOpfIdentifierAndUpdateExistingOne(oldIdentifierEl, opfContents, identifierID, newIdentifier string) string {
	var (
		idAttribute                    = fmt.Sprintf(` id="%s"`, identifierID)
		updatedOldIdentifierEl         = strings.Replace(oldIdentifierEl, idAttribute, "", 1)
		format                         strings.Builder
		oldIdentifierLeadingWhitespace = getLeadingWhitespace(oldIdentifierEl)
	)

	format.WriteString("\n" + oldIdentifierLeadingWhitespace)
	format.WriteString("<dc:identifier")
	if identifierID != "" {
		format.WriteString(idAttribute)
	}

	format.WriteString(">")
	format.WriteString(newIdentifier)
	format.WriteString("</dc:identifier>")

	if strings.Contains(updatedOldIdentifierEl, metadataEndTag) {
		updatedOldIdentifierEl = strings.Replace(updatedOldIdentifierEl, metadataEndTag, "", 1)

		format.WriteString("\n")
		format.WriteString(getMetadataWhitespaceForNewLine(oldIdentifierLeadingWhitespace) + metadataEndTag)
	}

	return strings.Replace(opfContents, oldIdentifierEl, updatedOldIdentifierEl+format.String(), 1)
}

// moveOrSetOpfIdentifierID moves the identifier's id from the current identifier in the OPF to the other identifier in the OPF that matches the NCX
// assuming that the oldIdentifierEl is not an empty string. If it is an empty string, it just sets the identifier's id.
func moveOrSetOpfIdentifierID(opfContents, ncxIdentifier, uniqueId, oldIdentifierEl string) string {
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

	if oldIdentifierEl != "" {
		var (
			idAttribute            = fmt.Sprintf(` id="%s"`, uniqueId)
			updatedOldIdentifierEl = strings.Replace(oldIdentifierEl, idAttribute, "", 1)
		)

		opfContents = strings.Replace(opfContents, oldIdentifierEl, updatedOldIdentifierEl, 1)
	}

	var line = opfContents[lineStart:lineEnd]

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

// getMetadataWhitespaceForNewLine determines what the whitespace should be when an ending metadata el
// was on the same line as some other element
func getMetadataWhitespaceForNewLine(currentLineWhitespace string) string {
	var potentialWhitespace = currentLineWhitespace[len(currentLineWhitespace)/2:]
	// if an element was indented two spaces we will assume that there was no indentation
	// instead of a single space of indentation
	if potentialWhitespace == " " {
		return ""
	}

	return potentialWhitespace
}
