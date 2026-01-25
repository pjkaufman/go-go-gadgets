package epubcheck

import (
	"strings"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
)

func HandleValidationErrors(opfFolder, ncxFilename, opfFilename string, nameToUpdatedContents map[string]string, validationErrors *ValidationErrors, getContentByFileName func(string) (string, error)) error {
	var (
		err                         error
		fileContent, ncxFileContent string
		elementNameToNumber         = make(map[string]int)
		fileToChanges               = make(map[string]positions.TextDocumentEdit)
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

			var update positions.TextEdit
			update, err = rulefixes.AddPropertyToManifest(fileContent, strings.TrimLeft(message.FilePath, opfFolder+"/"), property)
			if err != nil {
				return err
			}

			if !update.IsEmpty() {
				if existingUpdates, ok := fileToChanges[opfFilename]; ok {
					existingUpdates.Edits = append(existingUpdates.Edits, update)
					fileToChanges[opfFilename] = existingUpdates
				} else {
					fileToChanges[opfFilename] = positions.TextDocumentEdit{
						FilePath: opfFilename,
						Edits:    []positions.TextEdit{update},
					}
				}
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

			var update positions.TextEdit
			update, err = rulefixes.RemovePropertyFromManifest(fileContent, strings.TrimLeft(message.FilePath, opfFolder+"/"), property)
			if err != nil {
				return err
			}

			if !update.IsEmpty() {
				if existingUpdates, ok := fileToChanges[opfFilename]; ok {
					existingUpdates.Edits = append(existingUpdates.Edits, update)
					fileToChanges[opfFilename] = existingUpdates
				} else {
					fileToChanges[opfFilename] = positions.TextDocumentEdit{
						FilePath: opfFilename,
						Edits:    []positions.TextEdit{update},
					}
				}
			}
		case "OPF-074":
			fileContent, err = getContentByFileName(opfFilename)
			if err != nil {
				return err
			}

			var updates []positions.TextEdit
			updates, err = rulefixes.RemoveDuplicateManifestEntry(message.Location.Line, message.Location.Column, fileContent)
			if err != nil {
				return err
			}

			if len(updates) != 0 {
				if existingUpdates, ok := fileToChanges[message.FilePath]; ok {
					existingUpdates.Edits = append(existingUpdates.Edits, updates...)
					fileToChanges[message.FilePath] = existingUpdates
				} else {
					fileToChanges[message.FilePath] = positions.TextDocumentEdit{
						FilePath: message.FilePath,
						Edits:    updates,
					}
				}
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

			var updates []positions.TextEdit
			updates, err = rulefixes.FixIdentifierDiscrepancy(fileContent, ncxFileContent)
			if err != nil {
				return err
			}

			if len(updates) != 0 {
				if existingUpdates, ok := fileToChanges[opfFilename]; ok {
					existingUpdates.Edits = append(existingUpdates.Edits, updates...)
					fileToChanges[opfFilename] = existingUpdates
				} else {
					fileToChanges[opfFilename] = positions.TextDocumentEdit{
						FilePath: opfFilename,
						Edits:    updates,
					}
				}
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

				update := rulefixes.FixXmlIdValue(fileContent, message.Location.Line, attribute)
				if !update.IsEmpty() {
					if existingUpdates, ok := fileToChanges[message.FilePath]; ok {
						existingUpdates.Edits = append(existingUpdates.Edits, update)
						fileToChanges[message.FilePath] = existingUpdates
					} else {
						fileToChanges[message.FilePath] = positions.TextDocumentEdit{
							FilePath: message.FilePath,
							Edits:    []positions.TextEdit{update},
						}
					}
				}
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

					var updates []positions.TextEdit
					updates, err = rulefixes.FixManifestAttribute(fileContent, attribute, message.Location.Line, elementNameToNumber)
					if err != nil {
						return err
					}

					if len(updates) != 0 {
						if existingUpdates, ok := fileToChanges[message.FilePath]; ok {
							existingUpdates.Edits = append(existingUpdates.Edits, updates...)
							fileToChanges[message.FilePath] = existingUpdates
						} else {
							fileToChanges[message.FilePath] = positions.TextDocumentEdit{
								FilePath: message.FilePath,
								Edits:    updates,
							}
						}
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

					var update positions.TextEdit
					update, deletedLine, err = rulefixes.RemoveEmptyOpfElements(elementName, message.Location.Line, fileContent)
					if err != nil {
						return err
					}

					if !update.IsEmpty() {
						if existingUpdates, ok := fileToChanges[opfFilename]; ok {
							existingUpdates.Edits = append(existingUpdates.Edits, update)
							fileToChanges[opfFilename] = existingUpdates
						} else {
							fileToChanges[opfFilename] = positions.TextDocumentEdit{
								FilePath: opfFilename,
								Edits:    []positions.TextEdit{update},
							}
						}
					}

					if deletedLine {
						validationErrors.RemoveLineReferences(message.Location.Line, message.FilePath)
						i--
					}
				}
			} else if message.Message == invalidPlayOrder {
				fileContent, err = getContentByFileName(ncxFilename)
				if err != nil {
					return err
				}

				updates := rulefixes.FixPlayOrder(fileContent)
				if len(updates) != 0 {
					if existingUpdates, ok := fileToChanges[message.FilePath]; ok {
						existingUpdates.Edits = append(existingUpdates.Edits, updates...)
						fileToChanges[message.FilePath] = existingUpdates
					} else {
						fileToChanges[message.FilePath] = positions.TextDocumentEdit{
							FilePath: message.FilePath,
							Edits:    updates,
						}
					}
				}
			} else if strings.HasPrefix(message.Message, duplicateIdPrefix) {
				id, foundId := getFirstQuotedValue(message.Message, len(duplicateIdPrefix))
				if !foundId {
					continue
				}

				fileContent, err = getContentByFileName(message.FilePath)
				if err != nil {
					return err
				}

				updates := rulefixes.UpdateDuplicateIds(fileContent, id)
				if len(updates) != 0 {
					if existingUpdates, ok := fileToChanges[message.FilePath]; ok {
						existingUpdates.Edits = append(existingUpdates.Edits, updates...)
						fileToChanges[message.FilePath] = existingUpdates
					} else {
						fileToChanges[message.FilePath] = positions.TextDocumentEdit{
							FilePath: message.FilePath,
							Edits:    updates,
						}
					}
				}
			} else if strings.HasPrefix(message.Message, invalidBlockquote) {
				fileContent, err = getContentByFileName(message.FilePath)
				if err != nil {
					return err
				}

				updates := rulefixes.FixFailedBlockquoteParsing(message.Location.Line, message.Location.Column, fileContent)
				if len(updates) != 0 {
					if existingUpdates, ok := fileToChanges[message.FilePath]; ok {
						existingUpdates.Edits = append(existingUpdates.Edits, updates...)
						fileToChanges[message.FilePath] = existingUpdates
					} else {
						fileToChanges[message.FilePath] = positions.TextDocumentEdit{
							FilePath: message.FilePath,
							Edits:    updates,
						}
					}
				}
			} else if message.Message == missingImgAlt {
				fileContent, err = getContentByFileName(message.FilePath)
				if err != nil {
					return err
				}

				update := rulefixes.FixMissingImageAlt(message.Location.Line, message.Location.Column, fileContent)
				if !update.IsEmpty() {
					if existingUpdates, ok := fileToChanges[message.FilePath]; ok {
						existingUpdates.Edits = append(existingUpdates.Edits, update)
						fileToChanges[message.FilePath] = existingUpdates
					} else {
						fileToChanges[message.FilePath] = positions.TextDocumentEdit{
							FilePath: message.FilePath,
							Edits:    []positions.TextEdit{update},
						}
					}
				}
			} else if strings.HasPrefix(message.Message, unexpectedSectionEl) {
				fileContent, err = getContentByFileName(message.FilePath)
				if err != nil {
					return err
				}

				updates := rulefixes.FixSectionElementUnexpected(message.Location.Line, message.Location.Column, fileContent)
				if len(updates) != 0 {
					if existingUpdates, ok := fileToChanges[message.FilePath]; ok {
						existingUpdates.Edits = append(existingUpdates.Edits, updates...)
						fileToChanges[message.FilePath] = existingUpdates
					} else {
						fileToChanges[message.FilePath] = positions.TextDocumentEdit{
							FilePath: message.FilePath,
							Edits:    updates,
						}
					}
				}
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

			var update positions.TextEdit
			update, err = rulefixes.FixMissingUniqueIdentifierId(fileContent, id)
			if err != nil {
				return err
			}

			if !update.IsEmpty() {
				if existingUpdates, ok := fileToChanges[message.FilePath]; ok {
					existingUpdates.Edits = append(existingUpdates.Edits, update)
					fileToChanges[message.FilePath] = existingUpdates
				} else {
					fileToChanges[message.FilePath] = positions.TextDocumentEdit{
						FilePath: message.FilePath,
						Edits:    []positions.TextEdit{update},
					}
				}
			}
		case "RSC-012":
			fileContent, err = getContentByFileName(message.FilePath)
			if err != nil {
				return err
			}

			update := rulefixes.RemoveLinkId(fileContent, message.Location.Line, message.Location.Column)
			if !update.IsEmpty() {
				if existingUpdates, ok := fileToChanges[message.FilePath]; ok {
					existingUpdates.Edits = append(existingUpdates.Edits, update)
					fileToChanges[message.FilePath] = existingUpdates
				} else {
					fileToChanges[message.FilePath] = positions.TextDocumentEdit{
						FilePath: message.FilePath,
						Edits:    []positions.TextEdit{update},
					}
				}
			}
		}
	}

	var updatedContents string
	for filePath, documentEdit := range fileToChanges {
		updatedContents, err = getContentByFileName(filePath)
		if err != nil {
			return err
		}

		updatedContents, err = positions.ApplyEdits(filePath, updatedContents, documentEdit.Edits)
		if err != nil {
			return err
		}

		nameToUpdatedContents[filePath] = updatedContents
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
