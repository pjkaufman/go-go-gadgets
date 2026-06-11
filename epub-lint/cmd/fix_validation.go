package cmd

import (
	"archive/zip"
	"fmt"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	epubcheck "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check"
	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/jnovels"
	"github.com/pjkaufman/go-go-gadgets/pkg/cli/flags"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	validationIssuesFilePath string
	removeJNovelInfo         bool
	autoFixValidationFlags   = flags.Flags{
		Flags: []flags.Flag{
			flags.NewBoolFlag(false, false, &removeJNovelInfo, "cleanup-jnovels", "", false, "whether or not to remove JNovels info if it is present"),
			flags.NewFileFlag(true, false, &validationIssuesFilePath, "issues", "", "", "the path to the file with the validation issues", nil, true),
			flags.NewFileFlag(true, false, &epubFile, "file", "f", "", "the epub file to replace strings in", []string{"epub"}, true),
		},
	}
)

// autoFixValidationCmd represents the auto fix validation command
var autoFixValidationCmd = &cobra.Command{
	Use:   "validation",
	Short: "Reads in the output of EPUBCheck and fixes as many issues as are able to be fixed without the user making any changes.",
	Long: heredoc.Doc(`Uses the provided epub and EPUBCheck output file to fix auto fixable auto fix issues. Here is a list of all of the error codes that are currently handled:
	- OPF-014: add scripted to the list of values in the properties attribute on the manifest item
	- OPF-015: remove scripted to the list of values in the properties attribute on the manifest item
	- OPF-030: add the unique identifier id to the first dc:identifier element that does not have an id already
	- OPF-074: remove duplicate manifest entries
	- OPF-096: make file content reachable by removing linear attribute
	- NAV-011: fix nav file's table of contents not being in the same order as the specified reading order
	- NCX-001: fix discrepancy in identifier between the OPF and NCX files
	- RSC-005: seems to be a catch all error id, but the following are handled around it
		- Update ids/attributes to have valid xml ids that conform to the xml and epub spec by removing colons and any other invalid characters with an underscore
			and starting the value with an underscore instead of a number if it currently is started by a number
		- Move attribute properties to their own meta elements that refine the element they were on to fix incorrect scheme declarations or other prefixes
		- Remove empty elements that should not be empty but are empty which is typically an identifier or description that has 0 content in it
		- Update duplicate ids to no longer be duplicates
		- Add div tags inside of blockquote elements that were not able to be parsed and do not have a blockquote inside of them
		- Add an empty alt attribute to img elements that are missing them
		- Move section elements from inside of span and paragraph tags to outside of them if they have no other siblings or other parent tags before the span and paragraph
		- Update empty title with the text of the first header in the file or the first paragraph if there is no header and there is a paragraph
	- RSC-007: try to fix broken file links and remove
	- RSC-012: try to fix broken links by removing the id link in the href attribute
	- RSC-017: seems to be a catch all error id, but the following are handled around it
		- Add missing title element with the text of the first header or, if no header is present, the first paragraph present in the file 
	- HTM-004: try to fix broken DOCTYPEs by replacing them with the expected DOCTYPE
	`),
	Example: heredoc.Doc(`
		epub-lint fix validation -f test.epub --issues epubCheckOutput.txt
		will read in the contents of the file and try to fix any of the fixable
		validation issues

		epub-lint fix validation -f test.epub --issues epubCheckOutput.txt --cleanup-jnovels
		will read in the contents of the file and try to fix any of the fixable
		validation issues as well as remove any jnovels specific files
	`),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return autoFixValidationFlags.Validate()
	},
	Run: func(cmd *cobra.Command, args []string) {
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

			var basenameToFilePaths = make(map[string][]string)
			for filename := range zipFiles {
				var basename = filepath.Base(filename)
				if files, ok := basenameToFilePaths[basename]; ok {
					basenameToFilePaths[basename] = append(files, filename)
				} else {
					basenameToFilePaths[basename] = []string{filename}
				}
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
			err = epubcheck.HandleValidationErrors(opfFolder, ncxFilename, opfFilename, nameToUpdatedContents, basenameToFilePaths, &validationErrors, getFileContentsByName, epubInfo.FilePathsInSpineOrder)
			if err != nil {
				return nil, err
			}

			if removeJNovelInfo {
				handledFiles, err = jnovels.CleanupJNovelsFiles(jnovels.JNovelsCleanupContext{
					EpubInfo:            epubInfo,
					OpfFolder:           opfFolder,
					OpfFileName:         opfFilename,
					NcxFileName:         ncxFilename,
					FileBasenameMap:     basenameToFilePaths,
					UpdatedFileContents: nameToUpdatedContents,
					GetFileContents:     getFileContentsByName,
				})

				if err != nil {
					return nil, err
				}
			}

			for filename, updatedContents := range nameToUpdatedContents {
				var name = filepath.Base(filename)
				if removeJNovelInfo && (name == jnovels.JnovelsFile || name == jnovels.JnovelsImage) {
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

	err := autoFixValidationFlags.AddToCmd(autoFixValidationCmd)
	if err != nil {
		logger.WriteError(err.Error())
	}
}
