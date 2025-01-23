package epub

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"path/filepath"
	"strings"
	"time"

	epubhandler "github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/epub-handler"
	"github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/linter"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	validationIssuesFilePath string
	ErrIssueFileArgNonJson   = errors.New("issue-file must be a JSON file")
	ErrIssueFileArgEmpty     = errors.New("issue-file must have a non-whitespace value")
)

// autoFixValidationCmd represents the auto fix validation command
var autoFixValidationCmd = &cobra.Command{
	Use:   "fix-validation",
	Short: "Reads in the output of EpubCheck and fixes as many issues as are able to be fixed without the user making any changes.",
	// Long: heredoc.Doc(`Uses the provided epub and extra replace Markdown file to replace a common set of strings and any extra instances specified in the extra file replace. After all replacements are made, the original epub will be moved to a .original file and the new file will take the place of the old file. It will also print out the successful extra replacements with the number of replacements made followed by warnings for any extra strings that it tried to find and replace values for, but did not find any instances to replace.
	// 	Note: it only replaces strings in content/xhtml files listed in the opf file.`),
	// Example: heredoc.Doc(`
	// 	ebook-lint epub replace-strings -f test.epub -e replacements.md
	// 	will replace the common strings and extra strings parsed out of replacements.md in content/xhtml files located in test.epub.
	// 	The original test.epub will be moved to test.epub.original and test.epub will have the updated files.

	// 	replacements.md is expected to be in the following format:
	// 	| Text to replace | Text to replace with |
	// 	| --------------- | -------------------- |
	// 	| I am typo | I the correct value |
	// 	...
	// 	| I am another issue to correct | the correction |
	// `),
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

		logger.WriteInfo("Starting epub validation fixes..\n")

		validationBytes, err := filehandler.ReadInBinaryFileContents(validationIssuesFilePath)
		if err != nil {
			logger.WriteError(err.Error())
		}

		var validationIssues EpubCheckInfo
		err = json.Unmarshal(validationBytes, &validationIssues)
		if err != nil {
			logger.WriteErrorf("failed to unmarshal validation issues: %s", err)
		}

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
					ncxFileContents, err := filehandler.ReadInZipFileContents(zipFiles[filepath.Join(opfFolder, epubInfo.NcxFile)])
					if err != nil {
						return nil, err
					}

					opfFileContents, err = linter.FixIdentifierDiscrepancy(opfFileContents, ncxFileContents)

					if err != nil {
						return nil, err
					}
				}
			}

			err = filehandler.WriteZipCompressedString(w, opfFilename, opfFileContents)
			if err != nil {
				return nil, err
			}

			handledFiles = append(handledFiles, opfFilename)

			return handledFiles, nil
		})
		if err != nil {
			logger.WriteErrorf("failed to fix validation issues in %q: %s", epubFile, err)
		}

		logger.WriteInfo("\nFinished fixing epub validation issues...")
	},
}

func init() {
	EpubCmd.AddCommand(autoFixValidationCmd)

	autoFixValidationCmd.Flags().StringVarP(&validationIssuesFilePath, "issue-file", "e", "", "the path to the file with the validation issues")
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
