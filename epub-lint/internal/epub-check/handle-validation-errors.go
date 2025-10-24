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
			startIndex := strings.Index(message.Message, `"`)
			if startIndex == -1 {
				continue
			}
			endIndex := strings.Index(message.Message[startIndex:], `"`)
			if endIndex == -1 {
				continue
			}

			property := message.Message[startIndex : startIndex+endIndex]

			fileContent, err = getContentByFileName(opfFilename)
			if err != nil {
				return err
			}

			nameToUpdatedContents[opfFilename], err = rulefixes.AddPropertyToManifest(fileContent, strings.TrimLeft(message.FilePath, opfFolder+"/"), property)
			if err != nil {
				return err
			}
		case "OPF-015":
			startIndex := strings.Index(message.Message, `"`)
			if startIndex == -1 {
				continue
			}
			endIndex := strings.Index(message.Message[startIndex:], `"`)
			if endIndex == -1 {
				continue
			}

			property := message.Message[startIndex : startIndex+endIndex]

			fileContent, err = getContentByFileName(opfFilename)
			if err != nil {
				return err
			}

			nameToUpdatedContents[opfFilename], err = rulefixes.RemovePropertyFromManifest(fileContent, strings.TrimLeft(message.FilePath, opfFolder+"/"), property)
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

			var (
				newLineAddedAt, amountOfNewLinesAdded int
			)
			nameToUpdatedContents[opfFilename], newLineAddedAt, amountOfNewLinesAdded, err = rulefixes.FixIdentifierDiscrepancy(fileContent, ncxFileContent)
			if err != nil {
				return err
			}

			if amountOfNewLinesAdded > 0 {
				validationErrors.IncrementLineNumbersBy(newLineAddedAt, amountOfNewLinesAdded, message.FilePath)
			}
		case "RSC-005":
			if strings.HasPrefix(message.Message, invalidIdPrefix) {
				startIndex := strings.Index(message.Message, invalidIdPrefix)
				if startIndex == -1 {
					continue
				}
				startIndex += len(invalidIdPrefix)
				endIndex := strings.Index(message.Message[startIndex:], `"`)
				if endIndex == -1 {
					continue
				}

				attribute := message.Message[startIndex : startIndex+endIndex]

				// TODO: update to not care about file other than where it gets the file contents

				// for now we will just fix the values in the opf and ncx files and we will handle the other cases separately
				// when that is encountered since it requires keeping track of which files have already been modified
				// and which ones have not been modified yet
				if strings.HasSuffix(message.FilePath, ".opf") || strings.HasSuffix(message.FilePath, ".ncx") {
					fileContent, err = getContentByFileName(message.FilePath)
					if err != nil {
						return err
					}

					nameToUpdatedContents[message.FilePath] = rulefixes.FixXmlIdValue(fileContent, message.Location.Line, attribute)
				}
			} else if strings.HasPrefix(message.Message, invalidAttribute) {
				startIndex := strings.Index(message.Message, invalidAttribute)
				if startIndex == -1 {
					continue
				}
				startIndex += len(invalidAttribute)
				endIndex := strings.Index(message.Message[startIndex:], `"`)
				if endIndex == -1 {
					continue
				}

				attribute := message.Message[startIndex : startIndex+endIndex]

				// for now we will just fix the values in the opf file and we will handle the other cases separately
				// when that is encountered since it requires keeping track of which files have already been modified
				// and which ones have not been modified yet
				if strings.HasSuffix(message.FilePath, ".opf") {
					nameToUpdatedContents[opfFilename], err = rulefixes.FixManifestAttribute(nameToUpdatedContents[opfFilename], attribute, message.Location.Line-1, elementNameToNumber)
					if err != nil {
						return err
					}

					validationErrors.IncrementLineNumbers(message.Location.Line, message.FilePath)
				}
			} else if strings.HasPrefix(message.Message, EmptyMetadataProperty) {
				startIndex := strings.Index(message.Message, EmptyMetadataProperty)
				if startIndex == -1 {
					continue
				}
				startIndex += len(EmptyMetadataProperty)
				endIndex := strings.Index(message.Message[startIndex:], `"`)
				if endIndex == -1 {
					continue
				}

				elementName := message.Message[startIndex : startIndex+endIndex]

				var deletedLine, oneDeleted bool
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
						oneDeleted = true
					}
				}

				if oneDeleted {
					i--
				}
			} else if message.Message == invalidPlayOrder {
				nameToUpdatedContents[ncxFilename] = rulefixes.FixPlayOrder(nameToUpdatedContents[ncxFilename])
			} else if strings.HasPrefix(message.Message, duplicateIdPrefix) {
				startIndex := strings.Index(message.Message, duplicateIdPrefix)
				if startIndex == -1 {
					continue
				}
				startIndex += len(duplicateIdPrefix)
				endIndex := strings.Index(message.Message[startIndex:], `"`)
				if endIndex == -1 {
					continue
				}

				id := message.Message[startIndex : startIndex+endIndex]

				fileContent, err = getContentByFileName(message.FilePath)
				if err != nil {
					return err
				}

				fileContent, charactersAdded := rulefixes.UpdateDuplicateIds(fileContent, id)
				nameToUpdatedContents[message.FilePath] = fileContent

				if charactersAdded > 0 {
					validationErrors.UpdateLineColumnPosition(message.Location.Line, message.Location.Column, charactersAdded, message.FilePath)
				}
			} else if strings.HasPrefix(message.Message, invalidBlockquote) {
				fileContent, err = getContentByFileName(message.FilePath)
				if err != nil {
					return err
				}

				fileContent, charactersAdded := rulefixes.FixFailedBlockquoteParsing(message.Location.Line, message.Location.Column, fileContent)
				nameToUpdatedContents[message.FilePath] = fileContent

				if charactersAdded > 0 {
					validationErrors.UpdateLineColumnPosition(message.Location.Line, message.Location.Column, charactersAdded, message.FilePath)
				}
			}
		case "OPF-030":
			startIndex := strings.Index(message.Message, missingUniqueIdentifier)
			if startIndex == -1 {
				continue
			}
			startIndex += len(missingUniqueIdentifier)
			endIndex := strings.Index(message.Message[startIndex:], `"`)
			if endIndex == -1 {
				continue
			}

			fileContent, err = getContentByFileName(opfFilename)
			if err != nil {
				return err
			}

			nameToUpdatedContents[opfFilename], err = rulefixes.FixMissingUniqueIdentifierId(fileContent, message.Message[startIndex:startIndex+endIndex])
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
