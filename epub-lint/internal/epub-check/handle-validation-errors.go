package epubcheck

import (
	"strings"

	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
)

func HandleValidationErrors(opfFolder, ncxFilename, opfFilename string, nameToUpdatedContents map[string]string, validationErrors *ValidationErrors, getContentByFileName func(string) (string, error)) error {
	var (
		err                         error
		fileContent, ncxFileContent string
		elementNameToNumber         = make(map[string]int)
	)
	for i := 0; i < len(validationErrors.ValidationIssues); i++ {
		message := validationErrors.ValidationIssues[i]

		switch message.Code {
		case "OPF-014":
			property, foundPropertyName := getFirstQuotedValue(message.Message, -1)
			if !foundPropertyName {
				continue
			}

			fileContent, err = getContentByFileName(opfFilename)
			if err != nil {
				return err
			}

			nameToUpdatedContents[opfFilename], err = rulefixes.AddPropertyToManifest(fileContent, strings.TrimLeft(message.FilePath, opfFolder+"/"), property)
			if err != nil {
				return err
			}
		case "OPF-015":
			property, foundPropertyName := getFirstQuotedValue(message.Message, -1)
			if !foundPropertyName {
				continue
			}

			fileContent, err = getContentByFileName(opfFilename)
			if err != nil {
				return err
			}

			nameToUpdatedContents[opfFilename], err = rulefixes.RemovePropertyFromManifest(fileContent, strings.TrimLeft(message.FilePath, opfFolder+"/"), property)
			if err != nil {
				return err
			}
		case "OPF-074":
			fileContent, err = getContentByFileName(opfFilename)
			if err != nil {
				return err
			}

			nameToUpdatedContents[opfFilename], err = rulefixes.RemoveDuplicateManifestEntry(message.Location.Line, message.Location.Column, fileContent)
			if err != nil {
				return err
			}
		case "NCX-001":
			fileContent, err = getContentByFileName(opfFilename)
			if err != nil {
				return err
			}

			ncxFileContent, err = getContentByFileName(ncxFilename)
			if err != nil {
				return err
			}

			nameToUpdatedContents[opfFilename], err = rulefixes.FixIdentifierDiscrepancy(fileContent, ncxFileContent)
			if err != nil {
				return err
			}
		case "RSC-005":
			if strings.HasPrefix(message.Message, invalidIdPrefix) {
				attribute, foundAttributeName := getFirstQuotedValue(message.Message, len(invalidIdPrefix))
				if !foundAttributeName {
					continue
				}

				fileContent, err = getContentByFileName(message.FilePath)
				if err != nil {
					return err
				}

				nameToUpdatedContents[message.FilePath] = rulefixes.FixXmlIdValue(fileContent, message.Location.Line, attribute)
			} else if strings.HasPrefix(message.Message, invalidAttribute) {
				attribute, foundAttributeName := getFirstQuotedValue(message.Message, len(invalidAttribute))
				if !foundAttributeName {
					continue
				}

				// for now we will just fix the values in the opf file and we will handle the other cases separately
				// when that is encountered since it requires keeping track of which files have already been modified
				// and which ones have not been modified yet
				if strings.HasSuffix(message.FilePath, ".opf") {
					fileContent, err = getContentByFileName(opfFilename)
					if err != nil {
						return err
					}

					nameToUpdatedContents[opfFilename], err = rulefixes.FixManifestAttribute(fileContent, attribute, message.Location.Line-1, elementNameToNumber)
					if err != nil {
						return err
					}
				}
			} else if strings.HasPrefix(message.Message, EmptyMetadataProperty) {
				elementName, foundElementName := getFirstQuotedValue(message.Message, len(EmptyMetadataProperty))
				if !foundElementName {
					continue
				}

				var deletedLine bool
				// for now we will just fix the values in the opf file and we will handle the other cases separately
				// when that is encountered since it requires keeping track of which files have already been modified
				// and which ones have not been modified yet
				if strings.HasSuffix(message.FilePath, ".opf") {
					fileContent, err = getContentByFileName(opfFilename)
					if err != nil {
						return err
					}

					nameToUpdatedContents[opfFilename], deletedLine, err = rulefixes.RemoveEmptyOpfElements(elementName, message.Location.Line-1, fileContent)
					if err != nil {
						return err
					}

					if deletedLine {
						validationErrors.DecrementLineNumbersAndRemoveLineReferences(message.Location.Line, message.FilePath)
						i--
					}
				}
			} else if message.Message == invalidPlayOrder {
				fileContent, err = getContentByFileName(ncxFilename)
				if err != nil {
					return err
				}

				nameToUpdatedContents[ncxFilename] = rulefixes.FixPlayOrder(fileContent)
			} else if strings.HasPrefix(message.Message, duplicateIdPrefix) {
				id, foundId := getFirstQuotedValue(message.Message, len(duplicateIdPrefix))
				if !foundId {
					continue
				}

				fileContent, err = getContentByFileName(message.FilePath)
				if err != nil {
					return err
				}

				nameToUpdatedContents[message.FilePath] = rulefixes.UpdateDuplicateIds(fileContent, id)
			} else if strings.HasPrefix(message.Message, invalidBlockquote) {
				fileContent, err = getContentByFileName(message.FilePath)
				if err != nil {
					return err
				}

				nameToUpdatedContents[message.FilePath] = rulefixes.FixFailedBlockquoteParsing(message.Location.Line, message.Location.Column, fileContent)
			} else if message.Message == missingImgAlt {
				fileContent, err = getContentByFileName(message.FilePath)
				if err != nil {
					return err
				}

				nameToUpdatedContents[message.FilePath] = rulefixes.FixMissingImageAlt(message.Location.Line, message.Location.Column, fileContent)
			} else if strings.HasPrefix(message.Message, unexpectedSectionEl) {
				fileContent, err = getContentByFileName(message.FilePath)
				if err != nil {
					return err
				}

				nameToUpdatedContents[message.FilePath] = rulefixes.FixSectionElementUnexpected(message.Location.Line, message.Location.Column, fileContent)
			}
		case "OPF-030":
			id, foundId := getFirstQuotedValue(message.Message, len(missingUniqueIdentifier))
			if !foundId {
				continue
			}

			fileContent, err = getContentByFileName(opfFilename)
			if err != nil {
				return err
			}

			nameToUpdatedContents[opfFilename], err = rulefixes.FixMissingUniqueIdentifierId(fileContent, id)
			if err != nil {
				return err
			}
		case "RSC-012":
			fileContent, err = getContentByFileName(message.FilePath)
			if err != nil {
				return err
			}

			nameToUpdatedContents[message.FilePath] = rulefixes.RemoveLinkId(fileContent, message.Location.Line-1, message.Location.Column-1)
		}
	}

	return nil
}

// getFirstQuotedValue takes in a message and a potential start index
// if start index is -1 then it will find the first double quote itself
func getFirstQuotedValue(message string, startIndex int) (string, bool) {
	if startIndex == -1 {
		startIndex = strings.Index(message, `"`)
		if startIndex == -1 {
			return "", false
		}

		startIndex++
	}

	endIndex := strings.Index(message[startIndex:], `"`)
	if endIndex == -1 {
		return "", false
	}

	quotedValue := message[startIndex : startIndex+endIndex]
	// there is a situation where EPUBCheck returns null as the id when it does not exist
	if quotedValue == "null" {
		quotedValue = ""
	}

	return quotedValue, true
}
