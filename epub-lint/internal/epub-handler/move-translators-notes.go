package epubhandler

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/linter"
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

func MoveTranslatorsNotes(spineOrder []string, opfFolder, ncxFilename, opfFilename string, nameToUpdatedContents map[string]string, getContentByFileName func(string) (string, error)) (int, error) {
	var (
		translatorNoteListItems []string
		fileTranslatorNotes     []string
		startingNumber          int
		fullFilePath            string
	)
	for _, file := range spineOrder {
		fullFilePath = filepath.Join(opfFolder, file)

		contents, err := getContentByFileName(fullFilePath)
		if err != nil {
			return 0, err
		}

		var nameParts = strings.Split(file, "/")
		contents, fileTranslatorNotes, startingNumber = linter.GetTranslatorsNotes(contents, nameParts[len(nameParts)-1], tlNoteFileName, startingNumber)
		translatorNoteListItems = append(translatorNoteListItems, fileTranslatorNotes...)

		nameToUpdatedContents[fullFilePath] = contents
	}

	/*TODO:
	If there are translator's notes:
	- Add tl_notes.xhtml to the nav
	*/
	if len(translatorNoteListItems) > 0 {
		if opfFolder == "." {
			opfFolder = ""
		}

		var (
			pathParts      = strings.Split(fullFilePath, "/")
			htmlFolderPath = opfFolder
			relativePath   = tlNoteFileName
		)
		if len(pathParts) > 1 {
			htmlFolderPath = strings.Join(pathParts[0:len(pathParts)-1], "/")
		}

		if len(pathParts) > 2 {
			relativePath = filepath.Join(strings.Join(pathParts[1:len(pathParts)-1], "/"), tlNoteFileName)
		}

		var tlNotesFilePath = tlNoteFileName
		if htmlFolderPath != "" {
			tlNotesFilePath = filepath.Join(htmlFolderPath, tlNoteFileName)
		}

		nameToUpdatedContents[tlNotesFilePath] = fmt.Sprintf(defaultTLNoteContents, strings.Join(translatorNoteListItems, "				"))

		opfFileContents, err := getContentByFileName(opfFilename)
		if err != nil {
			return 0, err
		}

		if opfFolder == "" {
			relativePath = tlNotesFilePath
		}
		opfFileContents = AddFileToOpf(opfFileContents, relativePath, "tl_notes", "application/xhtml+xml")
		nameToUpdatedContents[opfFilename] = opfFileContents

		if ncxFilename != "" {
			ncxFileContents, err := getContentByFileName(ncxFilename)
			if err != nil {
				return 0, err
			}

			ncxFileContents = AddFileToNcx(ncxFileContents, relativePath, "Translator's Notes", "tl_notes")
			nameToUpdatedContents[ncxFilename] = ncxFileContents
		}
	}

	return len(translatorNoteListItems), nil

	// return epubhandler.UpdateEpub(epubFile, func(zipFiles map[string]*zip.File, w *zip.Writer, epubInfo epubhandler.EpubInfo, opfFolder string) ([]string, error) {
	// 	var (
	// 		handledFiles            []string
	// 		translatorNoteListItems []string
	// 		fileTranslatorNotes     []string
	// 		startingNumber          int
	// 		fullFilePath            string
	// 	)
	// 	for _, file := range epubInfo.FilePathsInSpineOrder {
	// 		if _, ok := epubInfo.HtmlFiles[file]; !ok {
	// 			continue
	// 		}

	// 		var (
	// 			filePath = getFilePath(opfFolder, file)
	// 			zipFile  = zipFiles[filePath]
	// 		)

	// 		fullFilePath = filePath

	// 		fileText, err := filehandler.ReadInZipFileContents(zipFile)
	// 		if err != nil {
	// 			return nil, err
	// 		}

	// var nameParts = strings.Split(file, "/")
	// fileText, fileTranslatorNotes, startingNumber = linter.GetTranslatorsNotes(fileText, nameParts[len(nameParts)-1], tlNoteFileName, startingNumber)

	// 		translatorNoteListItems = append(translatorNoteListItems, fileTranslatorNotes...)
	// 	}

	// /*TODO:
	// If there are translator's notes:
	// - Add tl_notes.xhtml to the nav
	// */
	// if len(translatorNoteListItems) > 0 {
	// 	if opfFolder == "." {
	// 		opfFolder = ""
	// 	}

	// 	var (
	// 		pathParts      = strings.Split(fullFilePath, "/")
	// 		htmlFolderPath = opfFolder
	// 		relativePath   = tlNoteFileName
	// 	)
	// 	if len(pathParts) > 1 {
	// 		htmlFolderPath = strings.Join(pathParts[0:len(pathParts)-1], "/")
	// 	}

	// 	if len(pathParts) > 2 {
	// 		relativePath = getFilePath(strings.Join(pathParts[1:len(pathParts)-1], "/"), tlNoteFileName)
	// 	}

	// 	var tlNotesFilePath = tlNoteFileName
	// 	if htmlFolderPath != "" {
	// 		tlNotesFilePath = getFilePath(htmlFolderPath, tlNoteFileName)
	// 	}

	// 	err = filehandler.WriteZipCompressedString(w, tlNotesFilePath, fmt.Sprintf(defaultTLNoteContents, strings.Join(translatorNoteListItems, "				")))
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	handledFiles = append(handledFiles, tlNotesFilePath)

	// 	opfFile := zipFiles[epubInfo.OpfFile]
	// 	opfFileContents, err := filehandler.ReadInZipFileContents(opfFile)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	if opfFolder == "" && relativePath == tlNoteFileName {
	// 		relativePath = tlNotesFilePath
	// 	}
	// 	opfFileContents = epubhandler.AddFileToOpf(opfFileContents, relativePath, "tl_notes", "application/xhtml+xml")
	// 	err = filehandler.WriteZipCompressedString(w, epubInfo.OpfFile, opfFileContents)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	handledFiles = append(handledFiles, epubInfo.OpfFile)

	// 	if epubInfo.NcxFile != "" {
	// 		var ncxFilePath = getFilePath(opfFolder, epubInfo.NcxFile)
	// 		ncxFileContents, err := filehandler.ReadInZipFileContents(zipFiles[ncxFilePath])
	// 		if err != nil {
	// 			return nil, err
	// 		}

	// 		ncxFileContents = epubhandler.AddFileToNcx(ncxFileContents, relativePath, "Translator's Notes", "tl_notes")

	// 		err = filehandler.WriteZipCompressedString(w, ncxFilePath, ncxFileContents)
	// 		if err != nil {
	// 			return nil, err
	// 		}

	// 			handledFiles = append(handledFiles, ncxFilePath)
	// 		}

	// 		var notesPluralization = "s"
	// 		if len(translatorNoteListItems) == 1 {
	// 			notesPluralization = ""
	// 		}
	// 		logger.WriteInfof("Found %d translator's note%s.\nAdding translator's notes file.\n", len(translatorNoteListItems), notesPluralization)

	// 	} else {
	// 		logger.WriteInfo("No translator's notes found.")
	// 	}

	// 	return handledFiles, nil
	// })
	// return 0, nil
}
