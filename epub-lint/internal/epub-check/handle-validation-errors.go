package epubcheck

// import (
// 	"fmt"
// 	"strings"

// 	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
// 	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
// )

// func test() {
// 	var (
// 		nameToUpdatedContents = map[string]string{
// 			ncxFilename: ncxFileContents,
// 			opfFilename: opfFileContents,
// 		}
// 		handledFiles []string
// 	)
// 	for i := 0; i < len(validationIssues); i++ {
// 		message := validationIssues[i]

// 		switch message.Code {
// 		case "OPF-014":
// 			nameToUpdatedContents[opfFilename], err = rulefixes.AddScriptedToManifest(nameToUpdatedContents[opfFilename], strings.TrimLeft(message.FilePath, opfFolder+"/"))

// 			if err != nil {
// 				return nil, err
// 			}
// 		case "OPF-015":
// 			nameToUpdatedContents[opfFilename], err = rulefixes.RemoveScriptedFromManifest(nameToUpdatedContents[opfFilename], strings.TrimLeft(message.FilePath, opfFolder+"/"))

// 			if err != nil {
// 				return nil, err
// 			}
// 		case "NCX-001":
// 			nameToUpdatedContents[opfFilename], err = rulefixes.FixIdentifierDiscrepancy(nameToUpdatedContents[opfFilename], nameToUpdatedContents[ncxFilename])

// 			if err != nil {
// 				return nil, err
// 			}
// 		case "RSC-005":
// 			if strings.HasPrefix(message.Message, invalidIdPrefix) {
// 				startIndex := strings.Index(message.Message, invalidIdPrefix)
// 				if startIndex == -1 {
// 					continue
// 				}
// 				startIndex += len(invalidIdPrefix)
// 				endIndex := strings.Index(message.Message[startIndex:], `"`)
// 				if endIndex == -1 {
// 					continue
// 				}

// 				attribute := message.Message[startIndex : startIndex+endIndex]

// 				// for now we will just fix the values in the opf and ncx files and we will handle the other cases separately
// 				// when that is encountered since it requires keeping track of which files have already been modified
// 				// and which ones have not been modified yet
// 				if strings.HasSuffix(message.FilePath, ".opf") {
// 					nameToUpdatedContents[opfFilename] = rulefixes.FixXmlIdValue(nameToUpdatedContents[opfFilename], message.Location.Line, attribute)
// 				} else if strings.HasSuffix(message.FilePath, ".ncx") {
// 					nameToUpdatedContents[ncxFilename] = rulefixes.FixXmlIdValue(nameToUpdatedContents[ncxFilename], message.Location.Line, attribute)
// 				}
// 			} else if strings.HasPrefix(message.Message, invalidAttribute) {
// 				startIndex := strings.Index(message.Message, invalidAttribute)
// 				if startIndex == -1 {
// 					continue
// 				}
// 				startIndex += len(invalidAttribute)
// 				endIndex := strings.Index(message.Message[startIndex:], `"`)
// 				if endIndex == -1 {
// 					continue
// 				}

// 				attribute := message.Message[startIndex : startIndex+endIndex]

// 				// for now we will just fix the values in the opf file and we will handle the other cases separately
// 				// when that is encountered since it requires keeping track of which files have already been modified
// 				// and which ones have not been modified yet
// 				if strings.HasSuffix(message.FilePath, ".opf") {
// 					nameToUpdatedContents[opfFilename], err = rulefixes.FixManifestAttribute(nameToUpdatedContents[opfFilename], attribute, message.Location.Line-1, elementNameToNumber)
// 					if err != nil {
// 						return nil, err
// 					}

// 					incrementLineNumbers(message.Location.Line, message.FilePath, validationIssues)
// 				}
// 			} else if strings.HasPrefix(message.Message, emptyMetadataProperty) {
// 				startIndex := strings.Index(message.Message, emptyMetadataProperty)
// 				if startIndex == -1 {
// 					continue
// 				}
// 				startIndex += len(emptyMetadataProperty)
// 				endIndex := strings.Index(message.Message[startIndex:], `"`)
// 				if endIndex == -1 {
// 					continue
// 				}

// 				elementName := message.Message[startIndex : startIndex+endIndex]

// 				var deletedLine, oneDeleted bool
// 				// for now we will just fix the values in the opf file and we will handle the other cases separately
// 				// when that is encountered since it requires keeping track of which files have already been modified
// 				// and which ones have not been modified yet
// 				if strings.HasSuffix(message.FilePath, ".opf") {
// 					nameToUpdatedContents[opfFilename], deletedLine, err = rulefixes.RemoveEmptyOpfElements(elementName, message.Location.Line-1, nameToUpdatedContents[opfFilename])
// 					if err != nil {
// 						return nil, err
// 					}

// 					if deletedLine {
// 						validationIssues = decrementLineNumbersAndRemoveLineReferences(message.Location.Line, message.FilePath, validationIssues)
// 						oneDeleted = true
// 					}
// 				}

// 				if oneDeleted {
// 					i--
// 				}
// 			} else if message.Message == invalidPlayOrder {
// 				nameToUpdatedContents[ncxFilename] = rulefixes.FixPlayOrder(nameToUpdatedContents[ncxFilename])
// 			} else if strings.HasPrefix(message.Message, duplicateIdPrefix) {
// 				startIndex := strings.Index(message.Message, duplicateIdPrefix)
// 				if startIndex == -1 {
// 					continue
// 				}
// 				startIndex += len(duplicateIdPrefix)
// 				endIndex := strings.Index(message.Message[startIndex:], `"`)
// 				if endIndex == -1 {
// 					continue
// 				}

// 				id := message.Message[startIndex : startIndex+endIndex]

// 				fileContents, ok := nameToUpdatedContents[message.FilePath]
// 				if !ok {
// 					zipFile, ok := zipFiles[message.FilePath]
// 					if !ok {
// 						return nil, fmt.Errorf("failed to find %q in the epub", message.FilePath)
// 					}

// 					fileContents, err = filehandler.ReadInZipFileContents(zipFile)
// 					if err != nil {
// 						return nil, err
// 					}
// 				}

// 				fileContents, charactersAdded := rulefixes.UpdateDuplicateIds(fileContents, id)
// 				nameToUpdatedContents[message.FilePath] = fileContents

// 				if charactersAdded > 0 {
// 					updateLineColumnPosition(message.Location.Line, message.Location.Column, charactersAdded, message.FilePath, validationIssues)
// 				}
// 			} else if strings.HasPrefix(message.Message, invalidBlockquote) {
// 				fileContents, ok := nameToUpdatedContents[message.FilePath]
// 				if !ok {
// 					zipFile, ok := zipFiles[message.FilePath]
// 					if !ok {
// 						return nil, fmt.Errorf("failed to find %q in the epub", message.FilePath)
// 					}

// 					fileContents, err = filehandler.ReadInZipFileContents(zipFile)
// 					if err != nil {
// 						return nil, err
// 					}
// 				}

// 				fileContents, charactersAdded := rulefixes.FixFailedBlockquoteParsing(message.Location.Line, message.Location.Column, fileContents)
// 				nameToUpdatedContents[message.FilePath] = fileContents

// 				if charactersAdded > 0 {
// 					updateLineColumnPosition(message.Location.Line, message.Location.Column, charactersAdded, message.FilePath, validationIssues)
// 				}
// 			}
// 		case "OPF-030":
// 			startIndex := strings.Index(message.Message, missingUniqueIdentifier)
// 			if startIndex == -1 {
// 				continue
// 			}
// 			startIndex += len(missingUniqueIdentifier)
// 			endIndex := strings.Index(message.Message[startIndex:], `"`)
// 			if endIndex == -1 {
// 				continue
// 			}

// 			nameToUpdatedContents[opfFilename], err = rulefixes.FixMissingUniqueIdentifierId(nameToUpdatedContents[opfFilename], message.Message[startIndex:startIndex+endIndex])
// 			if err != nil {
// 				return nil, err
// 			}
// 		case "RSC-012":
// 			if strings.HasSuffix(message.FilePath, ".opf") {
// 				nameToUpdatedContents[opfFilename] = rulefixes.RemoveLinkId(nameToUpdatedContents[opfFilename], message.Location.Line-1, message.Location.Column-1)
// 			} else if strings.HasSuffix(message.FilePath, ".ncx") {
// 				nameToUpdatedContents[ncxFilename] = rulefixes.RemoveLinkId(nameToUpdatedContents[ncxFilename], message.Location.Line-1, message.Location.Column-1)
// 			} else {
// 				if fileContents, ok := nameToUpdatedContents[message.FilePath]; ok {
// 					nameToUpdatedContents[message.FilePath] = rulefixes.RemoveLinkId(fileContents, message.Location.Line-1, message.Location.Column-1)
// 				} else {
// 					zipFile, ok := zipFiles[message.FilePath]
// 					if !ok {
// 						return nil, fmt.Errorf("failed to find %q in the epub", message.FilePath)
// 					}

// 					fileContents, err := filehandler.ReadInZipFileContents(zipFile)
// 					if err != nil {
// 						return nil, err
// 					}

// 					nameToUpdatedContents[message.FilePath] = rulefixes.RemoveLinkId(fileContents, message.Location.Line-1, message.Location.Column-1)
// 				}
// 			}
// 		}
// 	}
// }
