package cmd

import (
	"archive/zip"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
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

// organizeNotesCmd represents the move translator's notes command
var organizeNotesCmd = &cobra.Command{
	Use:   "organize-notes",
	Short: "Moves translator's notes to their own file at the end of the epub.",
	Example: heredoc.Doc(`Finds all translator's notes and moves them to their own file if present
	epub-lint organize-notes -f test.epub
	`),
	Long: heredoc.Doc(`Goes through all of the content files and looks for "TL Note:", "Translator's Note:", "T/N:", or "Note:"
	and moves any matches to their own file with bidirectional linking between the footnote and its reference location.
	It also adds an entry to the TOC and spine of the epub so the "tl_notes.xhtml" file is at the end of the file's contents.
`),
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
		if err != nil {
			logger.WriteError(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(organizeNotesCmd)
	organizeNotesCmd.Flags().StringVarP(&epubFile, "file", "f", "", "the epub file to move translator's notes to their own file in")
	err := organizeNotesCmd.MarkFlagRequired("file")
	if err != nil {
		logger.WriteErrorf(`failed to mark flag "file" as required on tl_notes command: %v\n`, err)
	}

	err = organizeNotesCmd.MarkFlagFilename("file", "epub")
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

			var notesPluralization = "s"
			if len(translatorNoteListItems) == 1 {
				notesPluralization = ""
			}
			logger.WriteInfof("Found %d translator's note%s.\nAdding translator's notes file.\n", len(translatorNoteListItems), notesPluralization)

		} else {
			logger.WriteInfo("No translator's notes found.")
		}

		return handledFiles, nil
	})
}
