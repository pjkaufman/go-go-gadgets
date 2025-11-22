package cmd

import (
	"archive/zip"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/MakeNowJust/heredoc"
	epubcheck "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check"
	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

const (
	jnovelsFile  = "jnovels.xhtml"
	jnovelsImage = "1.png"
)

var (
	validationIssuesFilePath string
	removeJNovelInfo         bool
	ErrIssueFileArgEmpty     = errors.New("issues must have a non-whitespace value")
)

// autoFixValidationCmd represents the auto fix validation command
var autoFixValidationCmd = &cobra.Command{
	Use:   "validation",
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
		- Add an empty alt attribute to img elements that are missing them
	- RSC-012: try to fix broken links by removing the id link in the href attribute
	`),
	Example: heredoc.Doc(`
		epub-lint fix validation -f test.epub --issues epubCheckOutput.txt
		will read in the contents of the file and try to fix any of the fixable
		validation issues

		epub-lint fix validation -f test.epub --issues epubCheckOutput.txt --cleanup-jnovels
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

		err = filehandler.FileArgExists(validationIssuesFilePath, "issues")
		if err != nil {
			logger.WriteError(err.Error())
		}

		logger.WriteInfo("Starting epub validation fixes...")

		validationOutput, err := filehandler.ReadInFileContents(validationIssuesFilePath)
		if err != nil {
			logger.WriteError(err.Error())
		}

		validationErrors, err := epubcheck.ParseEPUBCheckOutput(validationOutput)
		if err != nil {
			logger.WriteError(err.Error())
		}

		validationErrors.Sort()

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
				handledFiles          []string
				getFileContentsByName = func(filename string) (string, error) {
					fileContents, ok := nameToUpdatedContents[filename]
					if !ok {
						zipFile, ok := zipFiles[filename]
						if !ok {
							return "", fmt.Errorf("failed to find %q in the epub", filename)
						}

						fileContents, err = filehandler.ReadInZipFileContents(zipFile)
						if err != nil {
							return "", err
						}
					}

					return fileContents, nil
				}
			)
			err = epubcheck.HandleValidationErrors(opfFolder, ncxFilename, opfFilename, nameToUpdatedContents, &validationErrors, getFileContentsByName)
			if err != nil {
				return nil, err
			}

			if removeJNovelInfo {
				for filename := range zipFiles {
					var name = filepath.Base(filename)
					if name == jnovelsFile || name == jnovelsImage {
						handledFiles = append(handledFiles, filename)
					} else {
						continue
					}

					updatedOpfContents, err := epubhandler.RemoveFileFromOpf(nameToUpdatedContents[opfFilename], filename)
					if err != nil {
						logger.WriteErrorf("Failed to remove file %q from the opf contents: %s", filename, err)
					}

					nameToUpdatedContents[opfFilename] = updatedOpfContents
				}
			}

			for filename, updatedContents := range nameToUpdatedContents {
				var name = filepath.Base(filename)
				if removeJNovelInfo && (name == jnovelsFile || name == jnovelsImage) {
					continue
				}

				handledFiles = append(handledFiles, filename)

				err = filehandler.WriteZipCompressedString(w, filename, updatedContents)
				if err != nil {
					return nil, err
				}
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
	fixCmd.AddCommand(autoFixValidationCmd)

	autoFixValidationCmd.Flags().BoolVarP(&removeJNovelInfo, "cleanup-jnovels", "", false, "whether or not to remove JNovels info if it is present")
	autoFixValidationCmd.Flags().StringVarP(&validationIssuesFilePath, "issues", "", "", "the path to the file with the validation issues")
	err := autoFixValidationCmd.MarkFlagRequired("issues")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"issues\" as required on validation fix command: %v\n", err)
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
