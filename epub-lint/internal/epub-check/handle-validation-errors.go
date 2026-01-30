package epubcheck

import (
	"strings"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
)

func HandleValidationErrors(opfFolder, ncxFilename, opfFilename string, nameToUpdatedContents map[string]string, basenameToFilePaths map[string][]string, validationErrors *ValidationErrors, getContentByFileName func(string) (string, error)) error {
	var (
		err                         error
		fileContent, ncxFileContent string
		elementNameToNumber         = make(map[string]int)
		fileToChanges               = make(map[string]positions.TextDocumentEdit)
	)
	for i := 0; i < len(validationErrors.ValidationIssues); i++ {
		var (
			fileUpdated string
			edits       []positions.TextEdit
			message     = validationErrors.ValidationIssues[i]
		)

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

			// apply immediately since it is has little chance of conflict with other rules, but can be triggered
			// for the same line multiple times and previous edits need to be present when the next attempt to make
			// an edit to the line is present
			if !update.IsEmpty() {
				fileContent, err = positions.ApplyEdits(opfFilename, fileContent, []positions.TextEdit{update})
				if err != nil {
					return err
				}

				nameToUpdatedContents[opfFilename] = fileContent
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

			// apply immediately since it is has little chance of conflict with other rules, but can be triggered
			// for the same line multiple times and previous edits need to be present when the next attempt to make
			// an edit to the line is present
			if !update.IsEmpty() {
				fileContent, err = positions.ApplyEdits(opfFilename, fileContent, []positions.TextEdit{update})
				if err != nil {
					return err
				}

				nameToUpdatedContents[opfFilename] = fileContent
			}
		case "OPF-074":
			fileContent, err = getContentByFileName(opfFilename)
			if err != nil {
				return err
			}

			fileUpdated = opfFilename
			edits, err = rulefixes.RemoveDuplicateManifestEntry(message.Location.Line, message.Location.Column, fileContent)
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

			fileUpdated = opfFilename
			edits, err = rulefixes.FixIdentifierDiscrepancy(fileContent, ncxFileContent)
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

				update := rulefixes.FixXmlIdValue(fileContent, message.Location.Line, attribute)
				if !update.IsEmpty() {
					fileUpdated = message.FilePath
					edits = append(edits, update)
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

					fileUpdated = opfFilename
					edits, err = rulefixes.FixManifestAttribute(fileContent, attribute, message.Location.Line, elementNameToNumber)
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

					var update positions.TextEdit
					update, deletedLine, err = rulefixes.RemoveEmptyOpfElements(elementName, message.Location.Line, fileContent)
					if err != nil {
						return err
					}

					if !update.IsEmpty() {
						fileUpdated = opfFilename
						edits = append(edits, update)
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

				fileUpdated = ncxFilename
				edits = rulefixes.FixPlayOrder(fileContent)
			} else if strings.HasPrefix(message.Message, duplicateIdPrefix) {
				id, foundId := getFirstQuotedValue(message.Message, len(duplicateIdPrefix))
				if !foundId {
					continue
				}

				fileContent, err = getContentByFileName(message.FilePath)
				if err != nil {
					return err
				}

				fileUpdated = message.FilePath
				edits = rulefixes.UpdateDuplicateIds(fileContent, id)
			} else if strings.HasPrefix(message.Message, invalidBlockquote) {
				fileContent, err = getContentByFileName(message.FilePath)
				if err != nil {
					return err
				}

				fileUpdated = message.FilePath
				edits = rulefixes.FixFailedBlockquoteParsing(message.Location.Line, message.Location.Column, fileContent)
			} else if message.Message == missingImgAlt {
				fileContent, err = getContentByFileName(message.FilePath)
				if err != nil {
					return err
				}

				update := rulefixes.FixMissingImageAlt(message.Location.Line, message.Location.Column, fileContent)
				if !update.IsEmpty() {
					fileUpdated = message.FilePath
					edits = append(edits, update)
				}
			} else if strings.HasPrefix(message.Message, unexpectedSectionEl) {
				fileContent, err = getContentByFileName(message.FilePath)
				if err != nil {
					return err
				}

				fileUpdated = message.FilePath
				edits = rulefixes.FixSectionElementUnexpected(message.Location.Line, message.Location.Column, fileContent)
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
				fileUpdated = opfFilename
				edits = append(edits, update)
			}
		case "RSC-007":
			// we will not update the NCX files since they deal with the nav and that can get pretty tricky in the end
			// as it could remove a TOC entry you want to keep
			if strings.HasSuffix(message.FilePath, ".ncx") {
				continue
			}

			resource, foundId := getFirstQuotedValue(message.Message, len("Referenced resource "))
			if !foundId {
				continue
			}

			fileContent, err = getContentByFileName(message.FilePath)
			if err != nil {
				return err
			}

			var update positions.TextEdit
			update, err = rulefixes.FixFileNotFound(fileContent, resource, message.Message, message.Location.Line, message.Location.Column, basenameToFilePaths)
			if err != nil {
				return err
			}

			if !update.IsEmpty() {
				fileUpdated = message.FilePath
				edits = append(edits, update)
			}

			// TODO: handle the scenario where a JS file is removed in the reference by checking if there is a js file currently referenced and if not create a message for removing the scripted tag

		case "RSC-012":
			fileContent, err = getContentByFileName(message.FilePath)
			if err != nil {
				return err
			}

			update := rulefixes.RemoveLinkId(fileContent, message.Location.Line, message.Location.Column)
			if !update.IsEmpty() {
				fileUpdated = message.FilePath
				edits = append(edits, update)
			}
		case "HTM-004":
			fileContent, err = getContentByFileName(message.FilePath)
			if err != nil {
				return err
			}

			const expectedStart = `expected "`
			startOfExpected := strings.Index(message.Message, expectedStart)
			if startOfExpected == -1 {
				continue
			}

			startOfExpected += len(expectedStart)

			// the expected doctype does not seem to be properly escaped in the string provided,
			// so we will just go ahead and grab the last double quote and that should be the correct one
			endOfExpected := strings.LastIndex(message.Message[startOfExpected:], `"`)
			if endOfExpected == -1 {
				continue
			}
			endOfExpected += startOfExpected

			update := rulefixes.FixIrregularDoctype(fileContent, message.Message[startOfExpected:endOfExpected])
			if !update.IsEmpty() {
				fileUpdated = message.FilePath
				edits = append(edits, update)
			}

		}

		if len(edits) != 0 {
			if existingUpdates, ok := fileToChanges[fileUpdated]; ok {
				existingUpdates.Edits = append(existingUpdates.Edits, edits...)
				fileToChanges[fileUpdated] = existingUpdates
			} else {
				fileToChanges[fileUpdated] = positions.TextDocumentEdit{
					FilePath: fileUpdated,
					Edits:    edits,
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
