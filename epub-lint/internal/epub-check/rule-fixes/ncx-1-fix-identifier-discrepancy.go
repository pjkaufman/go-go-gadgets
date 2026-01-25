package rulefixes

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
)

func FixIdentifierDiscrepancy(opfContents, ncxContents string) ([]positions.TextEdit, error) {
	var edits []positions.TextEdit

	ncxIdentifier, err := getNcxIdentifier(ncxContents)
	if err != nil {
		return edits, err
	}

	opfIdentifierEl, opfIdentifier, opfIdentifierID, opfIdentifierIndex := getOpfIdentifier(opfContents)

	indexOfEndTag := strings.Index(opfContents, metadataEndTag)
	textUntilEndTag := opfContents[:indexOfEndTag]
	hasNcxIdentifier := strings.Contains(opfContents, ">"+ncxIdentifier+"<")

	// Scenario 1: OPF has no unique identifier but contains the NCX identifier
	if opfIdentifier == "" && ncxIdentifier != "" && hasNcxIdentifier {
		return moveOrSetOpfIdentifierID(opfContents, ncxIdentifier, opfIdentifierID, opfIdentifierEl, opfIdentifierIndex), nil
	}

	// Scenario 2: OPF has no unique identifier, NCX does
	if opfIdentifier == "" && ncxIdentifier != "" {
		prevNL := strings.LastIndex(textUntilEndTag, "\n")
		var prevLine string
		if prevNL == -1 {
			prevLine = textUntilEndTag
		} else {
			prevLine = opfContents[prevNL+1 : indexOfEndTag]
		}

		edit := addOpfIdentifier(opfContents, ncxIdentifier, opfIdentifierID, prevLine)
		return []positions.TextEdit{edit}, nil
	}

	// Scenario 3: OPF and NCX identifiers differ, NCX identifier not present in OPF
	if opfIdentifier != "" && ncxIdentifier != "" && opfIdentifier != ncxIdentifier && !hasNcxIdentifier {
		return addOpfIdentifierAndUpdateExistingOne(opfIdentifierEl, opfContents, opfIdentifierID, ncxIdentifier, opfIdentifierIndex), nil
	}

	// Scenario 4: OPF and NCX differ, but OPF already contains NCX identifier
	if opfIdentifier != "" && ncxIdentifier != "" && opfIdentifier != ncxIdentifier && hasNcxIdentifier {
		return moveOrSetOpfIdentifierID(opfContents, ncxIdentifier, opfIdentifierID, opfIdentifierEl, opfIdentifierIndex), nil
	}

	return edits, nil
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
func getOpfIdentifier(opfContents string) (string, string, string, int) {
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
					return fullLine, identifier, uniqueId, lineStart
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
					return fullLine, identifier, "", firstIdStart
				}
			}
		}
	}

	return "", "", uniqueId, -1
}

// addOpfIdentifier adds a unique identifier to the OPF content.
func addOpfIdentifier(opfContents, identifier, identifierID, metadataEndElPriorToEl string) positions.TextEdit {
	if identifierID == "" {
		identifierID = "pub-id"
	}

	identifierTag := fmt.Sprintf(`<dc:identifier id="%s">%s</dc:identifier>`, identifierID, identifier)

	var insertText string
	if strings.TrimSpace(metadataEndElPriorToEl) == "" {
		if metadataEndElPriorToEl == "" {
			insertText = "\t" + identifierTag + "\n"
		} else {
			insertText = metadataEndElPriorToEl + identifierTag + "\n" + metadataEndElPriorToEl
		}
	} else {
		ws := getLeadingWhitespace(metadataEndElPriorToEl)
		insertText = ws + identifierTag + "\n" + getMetadataWhitespaceForNewLine(ws)
	}

	idx := strings.Index(opfContents, metadataEndTag)
	pos := positions.IndexToPosition(opfContents, idx)

	return positions.TextEdit{
		Range: positions.Range{
			Start: pos,
			End:   pos,
		},
		NewText: insertText,
	}
}

// addOpfIdentifierAndUpdateExistingOne replaces the unique identifier in the OPF content.
func addOpfIdentifierAndUpdateExistingOne(oldIdentifierEl, opfContents, identifierID, newIdentifier string, oldIdentifierElIndex int) []positions.TextEdit {
	var edits []positions.TextEdit

	idAttr := fmt.Sprintf(` id="%s"`, identifierID)
	oldLeadingWS := getLeadingWhitespace(oldIdentifierEl)

	// Remove id="..." from old element (minimal deletion)
	if idx := strings.Index(oldIdentifierEl, idAttr); idx != -1 {
		edits = append(edits, positions.TextEdit{
			Range: positions.Range{
				Start: positions.IndexToPosition(opfContents, oldIdentifierElIndex+idx),
				End:   positions.IndexToPosition(opfContents, oldIdentifierElIndex+idx+len(idAttr)),
			},
			NewText: "",
		})
	}

	// Insert new identifier element after old one
	insertPoint := oldIdentifierElIndex + len(oldIdentifierEl)
	pos := positions.IndexToPosition(opfContents, insertPoint)

	newLine := "\n" + oldLeadingWS +
		fmt.Sprintf(`<dc:identifier id="%s">%s</dc:identifier>`, identifierID, newIdentifier)

	edits = append(edits, positions.TextEdit{
		Range: positions.Range{
			Start: pos,
			End:   pos,
		},
		NewText: newLine,
	})

	return edits
}

// moveOrSetOpfIdentifierID moves the identifier's id from the current identifier in the OPF to the other identifier in the OPF that matches the NCX
// assuming that the oldIdentifierEl is not an empty string. If it is an empty string, it just sets the identifier's id.
func moveOrSetOpfIdentifierID(opfContents, ncxIdentifier, uniqueId, oldIdentifierEl string, oldOpfIdentifierIndex int) []positions.TextEdit {
	// Find the line containing the ncxIdentifier
	ncxIdentifierLineStart := strings.Index(opfContents, ncxIdentifier)
	if ncxIdentifierLineStart == -1 {
		return nil // ncxIdentifier not found
	}

	var (
		edits     []positions.TextEdit
		lineStart = strings.LastIndex(opfContents[:ncxIdentifierLineStart], "\n") + 1
		lineEnd   = strings.Index(opfContents[ncxIdentifierLineStart:], "\n")
	)
	if lineEnd == -1 {
		lineEnd = len(opfContents)
	} else {
		lineEnd += ncxIdentifierLineStart
	}

	if oldIdentifierEl != "" {
		var (
			idAttribute            = fmt.Sprintf(` id="%s"`, uniqueId)
			indexOfOldIdentifierId = strings.Index(oldIdentifierEl, idAttribute)
		)

		if indexOfOldIdentifierId != -1 {
			edits = append(edits, positions.TextEdit{
				Range: positions.Range{
					Start: positions.IndexToPosition(opfContents, oldOpfIdentifierIndex+indexOfOldIdentifierId),
					End:   positions.IndexToPosition(opfContents, oldOpfIdentifierIndex+indexOfOldIdentifierId+len(idAttribute)),
				},
			})
		}
	}

	var line = opfContents[lineStart:lineEnd]

	// Check if the line already has an id attribute
	idAttr := ` id="`
	idStart := strings.Index(line, idAttr)
	if idStart == -1 {
		// No id attribute, add it
		var (
			ncxIdentifierIdInsertIndex = strings.Index(line, ">"+ncxIdentifier)
			// we are ignoring that the index could be -1 for now and will deal with it if we need to
			ncxInsertPos = positions.IndexToPosition(opfContents, lineStart+ncxIdentifierIdInsertIndex)
		)
		edits = append(edits, positions.TextEdit{
			Range: positions.Range{
				Start: ncxInsertPos,
				End:   ncxInsertPos,
			},
			NewText: fmt.Sprintf(` id="%s"`, uniqueId),
		})
	} else {
		// Replace the existing id attribute value with the uniqueId
		idEnd := strings.Index(line[idStart+len(idAttr):], `"`) + idStart + len(idAttr)
		// we are ignoring that the idEnd could be -1 for now and will deal with it if we need to
		edits = append(edits, positions.TextEdit{
			Range: positions.Range{
				Start: positions.IndexToPosition(opfContents, lineStart+idStart+len(idAttr)),
				End:   positions.IndexToPosition(opfContents, lineStart+idEnd),
			},
			NewText: uniqueId,
		})
	}

	return edits
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
