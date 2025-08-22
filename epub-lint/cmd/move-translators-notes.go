package cmd

import (
	"archive/zip"
	"fmt"
	"strings"

	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/linter"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

const tlNoteFileName = "tl_notes.xhtml"

var defaultTLNoteContents = `<?xml version='1.0' encoding='utf-8'?>
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
    <title>Translator's Notes</title>
</head>
<body>
    <h3>Translator's Notes</h3>
    <ol>
				%s</ol>
</body>
</html>
`

// translatorsNotesCmd represents the move translator's notes command
var translatorsNotesCmd = &cobra.Command{
	Use:   "tl_notes",
	Short: "Moves translator's notes to their own file at the end of the epub.",
	// Example: heredoc.Doc(`To run all of the possible potential fixes:
	// epub-lint fixable -f test.epub -a
	// Note: this will require a css file to already exist in the epub

	// To just fix broken paragraph endings:
	// epub-lint fixable -f test.epub --broken-lines

	// To just update section breaks:
	// epub-lint fixable -f test.epub --section-breaks
	// Note: this will require a css file to already exist in the epub

	// To just update page breaks:
	// epub-lint fixable -f test.epub --page-breaks
	// Note: this will require a css file to already exist in the epub

	// To just fix missing oxford commas:
	// epub-lint fixable -f test.epub --oxford-commas

	// To just fix potentially lacking subordinate clause instances:
	// epub-lint fixable -f test.epub --lacking-subordinate-clause

	// To just fix instances of thoughts in parentheses:
	// epub-lint fixable -f test.epub --thoughts

	// To run a combination of options:
	// epub-lint fixable -f test.epub -oxford-commas --thoughts --necessary-words
	// `),
	// Long: heredoc.Doc(`Goes through all of the content files and runs the specified fixable actions on them asking
	// for user input on each value found that matches the potential fix criteria.
	// Potential things that can be fixed:
	// - Broken paragraph endings
	// - Section breaks being hardcoded instead of an hr
	// - Page breaks being hardcoded instead of an hr
	// - Oxford commas being missing before or's or and's
	// - Possible instances of sentences with two subordinate clauses (i.e. have although..., but)
	// - Possible instances of thoughts that are in parentheses
	// - Possible instances of conversation encapsulated in square brackets
	// - Possible instances of words in square brackets that may be necessary for the sentence (i.e. need to have the brackets removed)
	// - Possible instances of single quotes that should actually be double quotes (i.e. when a word is in single quotes, but is not inside of double quotes)
	// `),
	Run: func(cmd *cobra.Command, args []string) {
		err := validateCommonEpubFlags(epubFile)
		if err != nil {
			logger.WriteError(err.Error())
		}

		err = filehandler.FileArgExists(epubFile, "file")
		if err != nil {
			logger.WriteError(err.Error())
		}

		err = moveTranslatorsNotes(epubFile)
	},
}

func init() {
	rootCmd.AddCommand(translatorsNotesCmd)
	translatorsNotesCmd.Flags().StringVarP(&epubFile, "file", "f", "", "the epub file to move translator's notes to their own file in")
	err := translatorsNotesCmd.MarkFlagRequired("file")
	if err != nil {
		logger.WriteErrorf(`failed to mark flag "file" as required on tl_notes command: %v\n`, err)
	}

	err = translatorsNotesCmd.MarkFlagFilename("file", "epub")
	if err != nil {
		logger.WriteErrorf(`failed to mark flag "file" as looking for specific file types on tl_notes command: %v\n`, err)
	}
}

func moveTranslatorsNotes(epubFile string) error {
	return epubhandler.UpdateEpub(epubFile, func(zipFiles map[string]*zip.File, w *zip.Writer, epubInfo epubhandler.EpubInfo, opfFolder string) ([]string, error) {
		err := validateFilesExist(opfFolder, epubInfo.HtmlFiles, zipFiles)
		if err != nil {
			return nil, err
		}

		var (
			handledFiles            []string
			translatorNoteListItems []string
			fileTranslatorNotes     []string
			startingNumber          int
			fullFilePath            string
		)
		for _, file := range epubInfo.FilePathsInSpineOrder {
			if _, ok := epubInfo.HtmlFiles[file]; !ok {
				continue
			}

			var (
				filePath = getFilePath(opfFolder, file)
				zipFile  = zipFiles[filePath]
			)

			fullFilePath = filePath

			fileText, err := filehandler.ReadInZipFileContents(zipFile)
			if err != nil {
				return nil, err
			}

			var nameParts = strings.Split(file, "/")
			fileText, fileTranslatorNotes, startingNumber = linter.GetTranslatorsNotes(fileText, nameParts[len(nameParts)-1], tlNoteFileName, startingNumber)

			translatorNoteListItems = append(translatorNoteListItems, fileTranslatorNotes...)

			err = filehandler.WriteZipCompressedString(w, filePath, fileText)
			if err != nil {
				return nil, err
			}

			handledFiles = append(handledFiles, filePath)
		}

		/*TODO:
		If there are translator's notes:
		- Add tl_notes.xhtml to the nav
		*/
		if len(translatorNoteListItems) > 0 {
			var (
				pathParts      = strings.Split(fullFilePath, "/")
				htmlFolderPath = opfFolder
				relativePath   = tlNoteFileName
			)
			if len(pathParts) > 1 {
				htmlFolderPath = strings.Join(pathParts[0:len(pathParts)-1], "/")
			}

			if len(pathParts) > 2 {
				relativePath = strings.Join(pathParts[1:len(pathParts)-1], "/")
			}

			var tlNotesFilePath = getFilePath(htmlFolderPath, tlNoteFileName)
			err = filehandler.WriteZipCompressedString(w, tlNotesFilePath, fmt.Sprintf(defaultTLNoteContents, strings.Join(translatorNoteListItems, "				")))
			if err != nil {
				return nil, err
			}

			handledFiles = append(handledFiles, tlNotesFilePath)

			opfFile := zipFiles[epubInfo.OpfFile]
			opfFileContents, err := filehandler.ReadInZipFileContents(opfFile)
			if err != nil {
				return nil, err
			}

			relativePath = getFilePath(relativePath, tlNoteFileName)

			opfFileContents = epubhandler.AddFileToOpf(opfFileContents, relativePath, "tl_notes", "application/xhtml+xml")
			err = filehandler.WriteZipCompressedString(w, epubInfo.OpfFile, opfFileContents)
			if err != nil {
				return nil, err
			}

			handledFiles = append(handledFiles, epubInfo.OpfFile)

			if epubInfo.NcxFile != "" {
				var ncxFilePath = getFilePath(opfFolder, epubInfo.NcxFile)
				ncxFileContents, err := filehandler.ReadInZipFileContents(zipFiles[ncxFilePath])
				if err != nil {
					return nil, err
				}

				ncxFileContents = epubhandler.AddFileToNcx(ncxFileContents, relativePath, "Translator's Notes", "tl_notes")

				err = filehandler.WriteZipCompressedString(w, ncxFilePath, ncxFileContents)
				if err != nil {
					return nil, err
				}

				handledFiles = append(handledFiles, ncxFilePath)
			}
		}

		return handledFiles, nil
	})
}
