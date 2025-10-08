package cmd

import (
	"archive/zip"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"slices"

	"github.com/MakeNowJust/heredoc"
	epubcheck "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check"
	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/linter"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

const (
	invalidIdPrefix         = "Error while parsing file: value of attribute \""
	invalidAttribute        = "Error while parsing file: attribute \""
	missingUniqueIdentifier = "The unique-identifier \""
	emptyMetadataProperty   = "Error while parsing file: character content of element \""
	invalidPlayOrder        = "Error while parsing file: identical playOrder values for navPoint/navTarget/pageTarget that do not refer to same target"
	duplicateIdPrefix       = "Error while parsing file: Duplicate \""
	invalidBlockquote       = "Error while parsing file: element \"blockquote\" incomplete;"
	jnovelsFile             = "jnovels.xhtml"
	jnovelsImage            = "1.png"
)

var (
	validationIssuesFilePath string
	removeJNovelInfo         bool
	ErrIssueFileArgEmpty     = errors.New("issue-file must have a non-whitespace value")
)

// autoFixValidationCmd represents the auto fix validation command
var autoFixValidationCmd = &cobra.Command{
	Use:   "fix-validation",
	Short: "Reads in the output of EPUBCheck and fixes as many issues as are able to be fixed without the user making any changes.",
	Long: heredoc.Doc(`Uses the provided epub and EPUBCheck output file to fix auto fixable auto fix issues. Here is a list of all of the error codes that are currently handled:
	- OPF-014: add scripted to the list of values in the properties attribute on the manifest item
	- OPF-015: remove scripted to the list of values in the properties attribute on the manifest item
	- NCX-001: fix discrepancy in identifier between the OPF and NCX files
	- OPF-030: add the unique identifier id to the first dc:identifier element that does not have an id already
	- RSC-005: seems to be a catch all error id, but the following are handled around it
	  - Update ids/attributes to have valid xml ids that conform to the xml and epub spec by removing colons and any other invalid characters with an underscore
	    and starting the value with an underscore instead of a number if it currently is started by a number
	  - Move attribute properties to their own meta elements that refine the element they were on to fix incorrect scheme declarations or other prefixes
	  - Remove empty elements that should not be empty but are empty which is typically an identifier or description that has 0 content in it
		- Update duplicate ids to no longer be duplicates
		- Add paragraph tags inside of blockquote elements that were not able to be parsed and either were a self-closing element, just text, or a span tag
	- RSC-012: try to fix broken links by removing the id link in the href attribute
	`),
	Example: heredoc.Doc(`
		epub-lint fix-validation -f test.epub --issue-file epubCheckOutput.txt
		will read in the contents of the file and try to fix any of the fixable
		validation issues

		epub-lint fix-validation -f test.epub --issue-file epubCheckOutput.txt --cleanup-jnovels
		will read in the contents of the file and try to fix any of the fixable
		validation issues as well as remove any jnovels specific files
	`),
	Run: func(cmd *cobra.Command, args []string) {
		err := ValidateAutoFixValidationFlags(epubFile, validationIssuesFilePath)
		if err != nil {
			logger.WriteError(err.Error())
		}

		err = filehandler.FileArgExists(epubFile, "file")
		if err != nil {
			logger.WriteError(err.Error())
		}

		err = filehandler.FileArgExists(validationIssuesFilePath, "issue-file")
		if err != nil {
			logger.WriteError(err.Error())
		}

		logger.WriteInfo("Starting epub validation fixes...")

		validationOutput, err := filehandler.ReadInFileContents(validationIssuesFilePath)
		if err != nil {
			logger.WriteError(err.Error())
		}

		validationIssues, err := epubcheck.ParseEPUBCheckOutput(validationOutput)
		if err != nil {
			logger.WriteError(err.Error())
		}

		sort.Slice(validationIssues, func(i, j int) bool {
			msgI := validationIssues[i]
			msgJ := validationIssues[j]

			// Prioritize delete-required messages
			if strings.HasPrefix(msgI.Message, emptyMetadataProperty) && !strings.HasPrefix(msgJ.Message, emptyMetadataProperty) {
				return true
			}

			if !strings.HasPrefix(msgI.Message, emptyMetadataProperty) && strings.HasPrefix(msgJ.Message, emptyMetadataProperty) {
				return false
			}

			// Compare by path ascending
			if msgI.FilePath != msgJ.FilePath {
				return msgI.FilePath < msgJ.FilePath
			}

			if msgI.Location == nil && msgJ.Location == nil {
				return true
			}

			// If paths are the same, compare by line descending
			if msgI.Location.Line != msgJ.Location.Line {
				return msgI.Location.Line > msgJ.Location.Line
			}
			// If lines are the same, compare by column descending
			return msgI.Location.Column > msgJ.Location.Column
		})

		var elementNameToNumber = make(map[string]int)

		err = epubhandler.UpdateEpub(epubFile, func(zipFiles map[string]*zip.File, w *zip.Writer, epubInfo epubhandler.EpubInfo, opfFolder string) ([]string, error) {
			var (
				opfFilename = epubInfo.OpfFile
				opfFile     = zipFiles[opfFilename]
			)

			opfFileContents, err := filehandler.ReadInZipFileContents(opfFile)
			if err != nil {
				return nil, err
			}

			var ncxFilename = filepath.Join(opfFolder, epubInfo.NcxFile)
			ncxFileContents, err := filehandler.ReadInZipFileContents(zipFiles[ncxFilename])
			if err != nil {
				return nil, err
			}

			var (
				nameToUpdatedContents = map[string]string{
					ncxFilename: ncxFileContents,
					opfFilename: opfFileContents,
				}
				handledFiles []string
			)
			for i := 0; i < len(validationIssues); i++ {
				message := validationIssues[i]

				switch message.Code {
				case "OPF-014":
					nameToUpdatedContents[opfFilename], err = linter.AddScriptedToManifest(nameToUpdatedContents[opfFilename], strings.TrimLeft(message.FilePath, opfFolder+"/"))

					if err != nil {
						return nil, err
					}
				case "OPF-015":
					nameToUpdatedContents[opfFilename], err = linter.RemoveScriptedFromManifest(nameToUpdatedContents[opfFilename], strings.TrimLeft(message.FilePath, opfFolder+"/"))

					if err != nil {
						return nil, err
					}
				case "NCX-001":
					nameToUpdatedContents[opfFilename], err = linter.FixIdentifierDiscrepancy(nameToUpdatedContents[opfFilename], nameToUpdatedContents[ncxFilename])

					if err != nil {
						return nil, err
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

						// for now we will just fix the values in the opf and ncx files and we will handle the other cases separately
						// when that is encountered since it requires keeping track of which files have already been modified
						// and which ones have not been modified yet
						if strings.HasSuffix(message.FilePath, ".opf") {
							nameToUpdatedContents[opfFilename] = linter.FixXmlIdValue(nameToUpdatedContents[opfFilename], message.Location.Line, attribute)
						} else if strings.HasSuffix(message.FilePath, ".ncx") {
							nameToUpdatedContents[ncxFilename] = linter.FixXmlIdValue(nameToUpdatedContents[ncxFilename], message.Location.Line, attribute)
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
							nameToUpdatedContents[opfFilename], err = linter.FixManifestAttribute(nameToUpdatedContents[opfFilename], attribute, message.Location.Line-1, elementNameToNumber)
							if err != nil {
								return nil, err
							}

							incrementLineNumbers(message.Location.Line, message.FilePath, validationIssues)
						}
					} else if strings.HasPrefix(message.Message, emptyMetadataProperty) {
						startIndex := strings.Index(message.Message, emptyMetadataProperty)
						if startIndex == -1 {
							continue
						}
						startIndex += len(emptyMetadataProperty)
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
							nameToUpdatedContents[opfFilename], deletedLine, err = linter.RemoveEmptyOpfElements(elementName, message.Location.Line-1, nameToUpdatedContents[opfFilename])
							if err != nil {
								return nil, err
							}

							if deletedLine {
								validationIssues = decrementLineNumbersAndRemoveLineReferences(message.Location.Line, message.FilePath, validationIssues)
								oneDeleted = true
							}
						}

						if oneDeleted {
							i--
						}
					} else if message.Message == invalidPlayOrder {
						nameToUpdatedContents[ncxFilename] = linter.FixPlayOrder(nameToUpdatedContents[ncxFilename])
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

						fileContents, ok := nameToUpdatedContents[message.FilePath]
						if !ok {
							zipFile, ok := zipFiles[message.FilePath]
							if !ok {
								return nil, fmt.Errorf("failed to find %q in the epub", message.FilePath)
							}

							fileContents, err = filehandler.ReadInZipFileContents(zipFile)
							if err != nil {
								return nil, err
							}
						}

						fileContents, charactersAdded := linter.UpdateDuplicateIds(fileContents, id)
						nameToUpdatedContents[message.FilePath] = fileContents

						if charactersAdded > 0 {
							updateLineColumnPosition(message.Location.Line, message.Location.Column, charactersAdded, message.FilePath, validationIssues)
						}
					} else if strings.HasPrefix(message.Message, invalidBlockquote) {
						fileContents, ok := nameToUpdatedContents[message.FilePath]
						if !ok {
							zipFile, ok := zipFiles[message.FilePath]
							if !ok {
								return nil, fmt.Errorf("failed to find %q in the epub", message.FilePath)
							}

							fileContents, err = filehandler.ReadInZipFileContents(zipFile)
							if err != nil {
								return nil, err
							}
						}

						fileContents, charactersAdded := linter.FixFailedBlockquoteParsing(message.Location.Line, message.Location.Column, fileContents)
						nameToUpdatedContents[message.FilePath] = fileContents

						if charactersAdded > 0 {
							updateLineColumnPosition(message.Location.Line, message.Location.Column, charactersAdded, message.FilePath, validationIssues)
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

					nameToUpdatedContents[opfFilename], err = linter.FixMissingUniqueIdentifierId(nameToUpdatedContents[opfFilename], message.Message[startIndex:startIndex+endIndex])
					if err != nil {
						return nil, err
					}
				case "RSC-012":
					if strings.HasSuffix(message.FilePath, ".opf") {
						nameToUpdatedContents[opfFilename] = linter.RemoveLinkId(nameToUpdatedContents[opfFilename], message.Location.Line-1, message.Location.Column-1)
					} else if strings.HasSuffix(message.FilePath, ".ncx") {
						nameToUpdatedContents[ncxFilename] = linter.RemoveLinkId(nameToUpdatedContents[ncxFilename], message.Location.Line-1, message.Location.Column-1)
					} else {
						if fileContents, ok := nameToUpdatedContents[message.FilePath]; ok {
							nameToUpdatedContents[message.FilePath] = linter.RemoveLinkId(fileContents, message.Location.Line-1, message.Location.Column-1)
						} else {
							zipFile, ok := zipFiles[message.FilePath]
							if !ok {
								return nil, fmt.Errorf("failed to find %q in the epub", message.FilePath)
							}

							fileContents, err := filehandler.ReadInZipFileContents(zipFile)
							if err != nil {
								return nil, err
							}

							nameToUpdatedContents[message.FilePath] = linter.RemoveLinkId(fileContents, message.Location.Line-1, message.Location.Column-1)
						}
					}
				}
			}

			if removeJNovelInfo {
				// remove the jnovels file and the png associated with it
				for file := range zipFiles {
					if strings.HasSuffix(file, jnovelsFile) || strings.HasSuffix(file, jnovelsImage) {
						handledFiles = append(handledFiles, file)
					}
				}

				// remove the associated files from the opf manifest and spine
				nameToUpdatedContents[opfFilename], err = linter.RemoveFileFromOpf(nameToUpdatedContents[opfFilename], jnovelsFile)
				if err != nil {
					return nil, err
				}

				nameToUpdatedContents[opfFilename], err = linter.RemoveFileFromOpf(nameToUpdatedContents[opfFilename], jnovelsImage)
				if err != nil {
					return nil, err
				}
			}

			for filename, updatedContents := range nameToUpdatedContents {
				if removeJNovelInfo && (strings.HasSuffix(filename, jnovelsFile) || strings.HasSuffix(filename, jnovelsImage)) {
					continue
				}

				err = filehandler.WriteZipCompressedString(w, filename, updatedContents)
				if err != nil {
					return nil, err
				}

				handledFiles = append(handledFiles, filename)
			}

			return handledFiles, nil
		})
		if err != nil {
			logger.WriteErrorf("failed to fix validation issues in %q: %s", epubFile, err)
		}

		logger.WriteInfo("Finished fixing epub validation issues.")
	},
}

func init() {
	rootCmd.AddCommand(autoFixValidationCmd)

	autoFixValidationCmd.Flags().BoolVarP(&removeJNovelInfo, "cleanup-jnovels", "", false, "whether or not to remove JNovels info if it is present")
	autoFixValidationCmd.Flags().StringVarP(&validationIssuesFilePath, "issue-file", "", "", "the path to the file with the validation issues")
	err := autoFixValidationCmd.MarkFlagRequired("issue-file")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"issue-file\" as required on validation fix command: %v\n", err)
	}

	autoFixValidationCmd.Flags().StringVarP(&epubFile, "file", "f", "", "the epub file to replace strings in")
	err = autoFixValidationCmd.MarkFlagRequired("file")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"file\" as required on validation fix command: %v\n", err)
	}

	err = autoFixValidationCmd.MarkFlagFilename("file", "epub")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"file\" as looking for specific file types on validation fix command: %v\n", err)
	}
}

func ValidateAutoFixValidationFlags(epubPath, validationIssuesPath string) error {
	err := validateCommonEpubFlags(epubPath)
	if err != nil {
		return err
	}

	if strings.TrimSpace(validationIssuesPath) == "" {
		return ErrIssueFileArgEmpty
	}

	return nil
}

func decrementLineNumbersAndRemoveLineReferences(lineNum int, path string, validationIssues []epubcheck.ValidationError) []epubcheck.ValidationError {
	for i := 0; i < len(validationIssues); i++ {
		if validationIssues[i].Location != nil {
			if validationIssues[i].FilePath == path {
				if validationIssues[i].Location.Line == lineNum {
					validationIssues = slices.Delete(validationIssues, i, i+1)
					i--
				} else if validationIssues[i].Location.Line > lineNum {
					validationIssues[i].Location.Line--
				}
			}
		}
	}

	return validationIssues
}

func incrementLineNumbers(lineNum int, path string, validationIssues []epubcheck.ValidationError) {
	for i := range validationIssues {
		if validationIssues[i].Location != nil {
			if validationIssues[i].FilePath == path && validationIssues[i].Location.Line > lineNum {
				validationIssues[i].Location.Line++
			}
		}
	}
}

func updateLineColumnPosition(lineNum, column, offset int, path string, validationIssues []epubcheck.ValidationError) {
	for i := range validationIssues {
		if validationIssues[i].Location != nil {
			if validationIssues[i].FilePath == path && validationIssues[i].Location.Line == lineNum && validationIssues[i].Location.Column > column {
				validationIssues[i].Location.Line += offset
			}
		}
	}
}
