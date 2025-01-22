package linter

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func FixIdentifierDiscrepancy(opfContents, ncxContents string) (string, error) {
	// Extract the unique identifier from the NCX
	ncxIdentifier, ncxScheme, err := getNcxIdentifier(ncxContents)
	if err != nil {
		return "", err
	}

	// Extract the unique identifier from the OPF
	opfIdentifier, opfIdentifierID, _, err := getOpfIdentifier(opfContents)
	if err != nil {
		return "", err
	}

	// Scenario 1: No unique identifier in OPF, but present in NCX
	if opfIdentifier == "" && ncxIdentifier != "" {
		opfContents = addOpfIdentifier(opfContents, ncxIdentifier, ncxScheme)
		return opfContents, nil
	}

	// Scenario 2: Different unique identifier in OPF and NCX
	if opfIdentifier != "" && ncxIdentifier != "" && opfIdentifier != ncxIdentifier {
		opfContents = replaceOpfIdentifier(opfContents, opfIdentifierID, ncxIdentifier, ncxScheme)
		return opfContents, nil
	}

	// Scenario 3: Different unique identifier in OPF and NCX where the OPF has the identifier from the NCX, but it is not the identifier specified in the OPF
	if opfIdentifier != "" && ncxIdentifier != "" && opfIdentifierID != ncxIdentifier {
		opfContents = moveOpfIdentifierID(opfContents, opfIdentifierID, ncxIdentifier)
		return opfContents, nil
	}

	return opfContents, nil
}

// getNcxIdentifier extracts the unique identifier and scheme from the NCX content.
func getNcxIdentifier(ncxContents string) (string, string, error) {
	startTag := `<meta name="dtb:uid" content="`
	startIndex := strings.Index(ncxContents, startTag)
	if startIndex == -1 {
		return "", "", fmt.Errorf("unique identifier not found in NCX")
	}
	startIndex += len(startTag)
	endIndex := strings.Index(ncxContents[startIndex:], `"`)
	if endIndex == -1 {
		return "", "", fmt.Errorf("unique identifier not found in NCX")
	}
	identifier := ncxContents[startIndex : startIndex+endIndex]

	// Determine the scheme based on the identifier format
	scheme := ""
	if IsValidISBN(identifier) {
		scheme = "ISBN"
	} else if _, err := uuid.Parse(identifier); err == nil {
		scheme = "UUID"
	}

	return identifier, scheme, nil
}

// getOpfIdentifier extracts the unique identifier and scheme from the OPF content.
func getOpfIdentifier(opfContents string) (string, string, string, error) {
	startTag := `<dc:identifier id="`
	startIndex := strings.Index(opfContents, startTag)
	if startIndex == -1 {
		return "", "", "", nil // Return nil error if unique identifier is not found in OPF
	}
	startIndex += len(startTag)
	endIndex := strings.Index(opfContents[startIndex:], `"`)
	if endIndex == -1 {
		return "", "", "", fmt.Errorf("unique identifier not found in OPF")
	}
	id := opfContents[startIndex : startIndex+endIndex]

	contentTag := `opf:scheme="`
	contentStartIndex := strings.Index(opfContents, contentTag)
	if contentStartIndex == -1 {
		return "", "", "", fmt.Errorf("unique identifier content not found in OPF")
	}
	contentStartIndex += len(contentTag)
	contentEndIndex := strings.Index(opfContents[contentStartIndex:], `"`)
	if contentEndIndex == -1 {
		return "", "", "", fmt.Errorf("unique identifier content not found in OPF")
	}
	scheme := opfContents[contentStartIndex : contentStartIndex+contentEndIndex]

	identifierTag := fmt.Sprintf(`opf:scheme="%s">`, scheme)
	identifierStartIndex := strings.Index(opfContents, identifierTag)
	if identifierStartIndex == -1 {
		return "", "", "", fmt.Errorf("unique identifier content not found in OPF")
	}
	identifierStartIndex += len(identifierTag)
	identifierEndIndex := strings.Index(opfContents[identifierStartIndex:], `<`)
	if identifierEndIndex == -1 {
		return "", "", "", fmt.Errorf("unique identifier content not found in OPF")
	}
	identifier := opfContents[identifierStartIndex : identifierStartIndex+identifierEndIndex]

	return identifier, id, scheme, nil
}

// addOpfIdentifier adds a unique identifier to the OPF content.
func addOpfIdentifier(opfContents, identifier, scheme string) string {
	var identifierTag string
	if scheme == "" {
		identifierTag = fmt.Sprintf(`<dc:identifier id="pub-id">%s</dc:identifier>`, identifier)
	} else {
		identifierTag = fmt.Sprintf(`<dc:identifier id="pub-id" opf:scheme="%s">%s</dc:identifier>`, scheme, identifier)
	}

	metadataEndTag := `</metadata>`
	return strings.Replace(opfContents, metadataEndTag, identifierTag+"\n"+metadataEndTag, 1)
}

// replaceOpfIdentifier replaces the unique identifier in the OPF content.
func replaceOpfIdentifier(opfContents, identifierID, newIdentifier, scheme string) string {
	oldIdentifierTag := fmt.Sprintf(`<dc:identifier id="%s" opf:scheme="`, identifierID)
	oldIdentifierStart := strings.Index(opfContents, oldIdentifierTag)
	if oldIdentifierStart == -1 {
		return opfContents
	}
	oldIdentifierStart += len(oldIdentifierTag)
	oldIdentifierEnd := strings.Index(opfContents[oldIdentifierStart:], `">`) + oldIdentifierStart + 2
	oldIdentifierContentEnd := strings.Index(opfContents[oldIdentifierEnd:], `<`) + oldIdentifierEnd

	oldIdentifier := opfContents[oldIdentifierStart:oldIdentifierEnd]
	oldIdentifierContent := opfContents[oldIdentifierEnd:oldIdentifierContentEnd]

	newIdentifierTag := fmt.Sprintf(`%s%s">%s`, oldIdentifier[:len(oldIdentifier)-2], scheme, newIdentifier)
	return strings.Replace(opfContents, oldIdentifier+oldIdentifierContent, newIdentifierTag, 1)
}

// moveOpfIdentifierID moves the identifier's id from the current identifier in the OPF to the other identifier in the OPF that matches the NCX.
func moveOpfIdentifierID(opfContents, oldIdentifierID, newIdentifierID string) string {
	oldIdentifierTag := fmt.Sprintf(` id="%s"`, oldIdentifierID)
	newIdentifierTag := fmt.Sprintf(` id="%s"`, newIdentifierID)
	return strings.Replace(opfContents, oldIdentifierTag, newIdentifierTag, 1)
}
