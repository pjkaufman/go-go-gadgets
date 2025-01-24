package epub

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"
	epubhandler "github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/epub-handler"
	"github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/linter"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

const invalidIdPrefix = "Error while parsing file: value of attribute \""

var (
	validationIssuesFilePath string
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

			var handledFiles []string
			for _, message := range validationIssues.Messages {
				switch message.ID {
				case "OPF-014":
					opfFileContents, err = linter.AddScriptedToManifest(opfFileContents, strings.TrimLeft(message.Locations[0].Path, opfFolder+"/"))

					if err != nil {
						return nil, err
					}
				case "OPF-015":
					opfFileContents, err = linter.RemoveScriptedFromManifest(opfFileContents, strings.TrimLeft(message.Locations[0].Path, opfFolder+"/"))

					if err != nil {
						return nil, err
					}
				case "NCX-001":
					opfFileContents, err = linter.FixIdentifierDiscrepancy(opfFileContents, ncxFileContents)

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
						if strings.HasSuffix(message.Locations[0].Path, ".opf") {
							opfFileContents = linter.FixXmlIdValue(opfFileContents, message.Locations[0].Line, attribute)
						} else if strings.HasSuffix(message.Locations[0].Path, ".ncx") {
							ncxFileContents = linter.FixXmlIdValue(ncxFileContents, message.Locations[0].Line, attribute)
						}
					}
				}
			}

			err = filehandler.WriteZipCompressedString(w, opfFilename, opfFileContents)
			if err != nil {
				return nil, err
			}

			handledFiles = append(handledFiles, opfFilename)

			err = filehandler.WriteZipCompressedString(w, ncxFilename, ncxFileContents)
			if err != nil {
				return nil, err
			}

			handledFiles = append(handledFiles, ncxFilename)

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
