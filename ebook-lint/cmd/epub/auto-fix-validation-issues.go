package epub

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"slices"

	"github.com/MakeNowJust/heredoc"
	epubhandler "github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/epub-handler"
	"github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/linter"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

const (
	invalidIdPrefix         = "Error while parsing file: value of attribute \""
	invalidAttribute        = "Error while parsing file: attribute \""
	missingUniqueIdentifier = "The unique-identifier \""
	emptyMetadataProperty   = "Error while parsing file: character content of element \""
	jnovelsFile             = "jnovels.xhtml"
	jnovelsImage            = "1.png"
)

var (
	validationIssuesFilePath string
	removeJNovelInfo         bool
	ErrIssueFileArgNonJson   = errors.New("issue-file must be a JSON file")
	ErrIssueFileArgEmpty     = errors.New("issue-file must have a non-whitespace value")
)

// autoFixValidationCmd represents the auto fix validation command
var autoFixValidationCmd = &cobra.Command{
	Use:   "fix-validation",
	Short: "Reads in the output of EPUBCheck and fixes as many issues as are able to be fixed without the user making any changes.",
	Long: heredoc.Doc(`Uses the provided epub and EPUBCheck JSON output file to fix auto fixable auto fix issues. Here is a list of all of the error codes that are currently handled:
	- OPF-014: add scripted to the list of values in the properties attribute on the manifest item
	- OPF-015: remove scripted to the list of values in the properties attribute on the manifest item
	- NCX-001: fix discrepancy in identifier between the OPF and NCX files
	- OPF-030: add the unique identifier id to the first dc:identifier element that does not have an id already
	- RSC-005: seems to be a catch all error id, but the following are handled around it
	  - Update ids/attributes to have valid xml ids that conform to the xml and epub spec by removing colons and any other invalid characters with an underscore
	    and starting the value with an underscore instead of a number if it currently is started by a number
	  - Move attribute properties to their own meta elements that refine the element they were on to fix incorrect scheme declarations or other prefixes
	  - Remove empty elements that should not be empty but are empty which is typically an identifier or description that has 0 content in it
	`),
	Example: heredoc.Doc(`
		ebook-lint epub fix-validation -f test.epub --issue-file epubCheckOutput.json
		will read in the contents of the JSON file and try to fix any of the fixable
		validation issues
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

		validationBytes, err := filehandler.ReadInBinaryFileContents(validationIssuesFilePath)
		if err != nil {
			logger.WriteError(err.Error())
		}

		var validationIssues EpubCheckInfo
		err = json.Unmarshal(validationBytes, &validationIssues)
		if err != nil {
			logger.WriteErrorf("failed to unmarshal validation issues: %s", err)
		}

		sort.Slice(validationIssues.Messages, func(i, j int) bool {
			msgI := validationIssues.Messages[i]
			msgJ := validationIssues.Messages[j]

			// Prioritize delete-required messages
			if strings.HasPrefix(msgI.Message, emptyMetadataProperty) && !strings.HasPrefix(msgJ.Message, emptyMetadataProperty) {
				return true
			}

			if !strings.HasPrefix(msgI.Message, emptyMetadataProperty) && strings.HasPrefix(msgJ.Message, emptyMetadataProperty) {
				return false
			}

			// Compare by path ascending
			if msgI.Locations[0].Path != msgJ.Locations[0].Path {
				return msgI.Locations[0].Path < msgJ.Locations[0].Path
			}
			// If paths are the same, compare by line descending
			if msgI.Locations[0].Line != msgJ.Locations[0].Line {
				return msgI.Locations[0].Line > msgJ.Locations[0].Line
			}
			// If lines are the same, compare by column descending
			return msgI.Locations[0].Column > msgJ.Locations[0].Column
		})

		for index := range validationIssues.Messages {
			if len(validationIssues.Messages[index].Locations) > 1 {
				sort.Slice(validationIssues.Messages[index].Locations, func(i, j int) bool {
					// Compare by line descending
					if validationIssues.Messages[index].Locations[i].Line != validationIssues.Messages[index].Locations[j].Line {
						return validationIssues.Messages[index].Locations[i].Line > validationIssues.Messages[index].Locations[j].Line
					}

					// If lines are the same, compare by column descending
					return validationIssues.Messages[index].Locations[i].Column > validationIssues.Messages[index].Locations[j].Column
				})
			}
		}

		var elementNameToNumber = make(map[string]int)

		err = epubhandler.UpdateEpub(epubFile, func(zipFiles map[string]*zip.File, w *zip.Writer, epubInfo epubhandler.EpubInfo, opfFolder string) ([]string, error) {
			var (
				opfFilename string
				opfFile     *zip.File
			)
			for filename, file := range zipFiles {
				if strings.HasSuffix(filename, "opf") {
					opfFilename = filename
					opfFile = file
					break
				}
			}

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
			for i := 0; i < len(validationIssues.Messages); i++ {
				message := validationIssues.Messages[i]

				switch message.ID {
				case "OPF-014":
					for _, location := range message.Locations {
						nameToUpdatedContents[opfFilename], err = linter.AddScriptedToManifest(nameToUpdatedContents[opfFilename], strings.TrimLeft(location.Path, opfFolder+"/"))

						if err != nil {
							return nil, err
						}
					}
				case "OPF-015":
					for _, location := range message.Locations {
						nameToUpdatedContents[opfFilename], err = linter.RemoveScriptedFromManifest(nameToUpdatedContents[opfFilename], strings.TrimLeft(location.Path, opfFolder+"/"))

						if err != nil {
							return nil, err
						}
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

						for _, location := range message.Locations {
							// for now we will just fix the values in the opf and ncx files and we will handle the other cases separately
							// when that is encountered since it requires keeping track of which files have already been modified
							// and which ones have not been modified yet
							if strings.HasSuffix(location.Path, ".opf") {
								nameToUpdatedContents[opfFilename] = linter.FixXmlIdValue(nameToUpdatedContents[opfFilename], location.Line, attribute)
							} else if strings.HasSuffix(location.Path, ".ncx") {
								nameToUpdatedContents[ncxFilename] = linter.FixXmlIdValue(nameToUpdatedContents[ncxFilename], location.Line, attribute)
							}
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

						for i := range message.Locations {
							location := message.Locations[i]
							// for now we will just fix the values in the opf file and we will handle the other cases separately
							// when that is encountered since it requires keeping track of which files have already been modified
							// and which ones have not been modified yet
							if strings.HasSuffix(location.Path, ".opf") {
								nameToUpdatedContents[opfFilename], err = linter.FixManifestAttribute(nameToUpdatedContents[opfFilename], attribute, location.Line-1, elementNameToNumber)
								if err != nil {
									return nil, err
								}

								incrementLineNumbers(location.Line, location.Path, &validationIssues)
							}
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
						for _, location := range message.Locations {
							// for now we will just fix the values in the opf file and we will handle the other cases separately
							// when that is encountered since it requires keeping track of which files have already been modified
							// and which ones have not been modified yet
							if strings.HasSuffix(location.Path, ".opf") {
								nameToUpdatedContents[opfFilename], deletedLine, err = linter.RemoveEmptyOpfElements(elementName, location.Line-1, nameToUpdatedContents[opfFilename])
								if err != nil {
									return nil, err
								}

								if deletedLine {
									decrementLineNumbersAndRemoveLineReferences(location.Line, location.Path, &validationIssues)
									oneDeleted = true
								}
							}
						}

						if oneDeleted {
							i--
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
					for _, location := range message.Locations {
						if strings.HasSuffix(location.Path, ".opf") {
							nameToUpdatedContents[opfFilename] = linter.RemoveLinkId(nameToUpdatedContents[opfFilename], location.Line-1, location.Column-1)
						} else if strings.HasSuffix(location.Path, ".ncx") {
							nameToUpdatedContents[ncxFilename] = linter.RemoveLinkId(nameToUpdatedContents[ncxFilename], location.Line-1, location.Column-1)
						} else {
							if fileContents, ok := nameToUpdatedContents[location.Path]; ok {
								nameToUpdatedContents[location.Path] = linter.RemoveLinkId(fileContents, location.Line-1, location.Column-1)
							} else {
								zipFile, ok := zipFiles[location.Path]
								if !ok {
									return nil, fmt.Errorf("failed to find %q in the epub", location.Path)
								}

								fileContents, err := filehandler.ReadInZipFileContents(zipFile)
								if err != nil {
									return nil, err
								}

								nameToUpdatedContents[location.Path] = linter.RemoveLinkId(fileContents, location.Line-1, location.Column-1)
							}
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
	EpubCmd.AddCommand(autoFixValidationCmd)

	autoFixValidationCmd.Flags().BoolVarP(&removeJNovelInfo, "cleanup-jnovels", "", false, "whether or not to remove JNovels info if it is present")
	autoFixValidationCmd.Flags().StringVarP(&validationIssuesFilePath, "issue-file", "", "", "the path to the file with the validation issues")
	err := autoFixValidationCmd.MarkFlagRequired("issue-file")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"issue-file\" as required on validation fix command: %v\n", err)
	}

	err = autoFixValidationCmd.MarkFlagFilename("issue-file", "json")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"issue-file\" as looking for specific file types on validation fix command: %v\n", err)
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

	if !strings.HasSuffix(validationIssuesPath, ".json") {
		return ErrIssueFileArgNonJson
	}

	return nil
}

func decrementLineNumbersAndRemoveLineReferences(lineNum int, path string, validationIssues *EpubCheckInfo) {
	for i := 0; i < len(validationIssues.Messages); i++ {
		if len(validationIssues.Messages[i].Locations) != 0 {
			for j := 0; j < len(validationIssues.Messages[i].Locations); j++ {
				if validationIssues.Messages[i].Locations[j].Path == path && validationIssues.Messages[i].Locations[j].Line == lineNum {
					validationIssues.Messages[i].Locations = append(validationIssues.Messages[i].Locations[:j], validationIssues.Messages[i].Locations[j+1:]...)

					j--
				} else if validationIssues.Messages[i].Locations[j].Path == path && validationIssues.Messages[i].Locations[j].Line > lineNum {
					validationIssues.Messages[i].Locations[j].Line--
				}
			}

			if len(validationIssues.Messages[i].Locations) == 0 {
				validationIssues.Messages = slices.Delete(validationIssues.Messages, i, i+1)

				i--
			}
		}
	}
}

func incrementLineNumbers(lineNum int, path string, validationIssues *EpubCheckInfo) {
	for i := range validationIssues.Messages {
		if len(validationIssues.Messages[i].Locations) != 0 {
			for j := range validationIssues.Messages[i].Locations {
				if validationIssues.Messages[i].Locations[j].Path == path && validationIssues.Messages[i].Locations[j].Line > lineNum {
					validationIssues.Messages[i].Locations[j].Line++
				}
			}
		}
	}
}

type EpubCheckInfo struct {
	Messages []struct {
		ID                  string `json:"ID"`
		Severity            string `json:"severity"`
		Message             string `json:"message"`
		AdditionalLocations int    `json:"additionalLocations"`
		Locations           []struct {
			URL struct {
				Opaque       bool `json:"opaque"`
				Hierarchical bool `json:"hierarchical"`
			} `json:"url"`
			Path    string      `json:"path"`
			Line    int         `json:"line"`
			Column  int         `json:"column"`
			Context interface{} `json:"context"`
		} `json:"locations"`
		Suggestion interface{} `json:"suggestion"`
	} `json:"messages"`
	CustomMessageFileName interface{} `json:"customMessageFileName"`
	Checker               struct {
		Path           string `json:"path"`
		Filename       string `json:"filename"`
		CheckerVersion string `json:"checkerVersion"`
		CheckDate      string `json:"checkDate"`
		ElapsedTime    int    `json:"elapsedTime"`
		NFatal         int    `json:"nFatal"`
		NError         int    `json:"nError"`
		NWarning       int    `json:"nWarning"`
		NUsage         int    `json:"nUsage"`
	} `json:"checker"`
	Publication struct {
		Publisher            string        `json:"publisher"`
		Title                string        `json:"title"`
		Creator              []string      `json:"creator"`
		Date                 time.Time     `json:"date"`
		Subject              []interface{} `json:"subject"`
		Description          interface{}   `json:"description"`
		Rights               string        `json:"rights"`
		Identifier           string        `json:"identifier"`
		Language             string        `json:"language"`
		NSpines              int           `json:"nSpines"`
		CheckSum             int           `json:"checkSum"`
		RenditionLayout      string        `json:"renditionLayout"`
		RenditionOrientation string        `json:"renditionOrientation"`
		RenditionSpread      string        `json:"renditionSpread"`
		EPubVersion          string        `json:"ePubVersion"`
		IsScripted           bool          `json:"isScripted"`
		HasFixedFormat       bool          `json:"hasFixedFormat"`
		IsBackwardCompatible bool          `json:"isBackwardCompatible"`
		HasAudio             bool          `json:"hasAudio"`
		HasVideo             bool          `json:"hasVideo"`
		CharsCount           int           `json:"charsCount"`
		EmbeddedFonts        []interface{} `json:"embeddedFonts"`
		RefFonts             []interface{} `json:"refFonts"`
		HasEncryption        bool          `json:"hasEncryption"`
		HasSignatures        bool          `json:"hasSignatures"`
		Contributors         []interface{} `json:"contributors"`
	} `json:"publication"`
	Items []struct {
		ID                   string        `json:"id"`
		FileName             string        `json:"fileName"`
		MediaType            string        `json:"media_type"`
		CompressedSize       int           `json:"compressedSize"`
		UncompressedSize     int           `json:"uncompressedSize"`
		CompressionMethod    string        `json:"compressionMethod"`
		CheckSum             string        `json:"checkSum"`
		IsSpineItem          bool          `json:"isSpineItem"`
		SpineIndex           interface{}   `json:"spineIndex"`
		IsLinear             bool          `json:"isLinear"`
		IsFixedFormat        interface{}   `json:"isFixedFormat"`
		IsScripted           bool          `json:"isScripted"`
		RenditionLayout      interface{}   `json:"renditionLayout"`
		RenditionOrientation interface{}   `json:"renditionOrientation"`
		RenditionSpread      interface{}   `json:"renditionSpread"`
		ReferencedItems      []interface{} `json:"referencedItems"`
	} `json:"items"`
}
