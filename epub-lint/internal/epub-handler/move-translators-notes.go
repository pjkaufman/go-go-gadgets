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

func MoveTranslatorsNotes(spineOrder []string, opfFolder, ncxFilename, opfFilename, navFilename string, nameToUpdatedContents map[string]string, getContentByFileName func(string) (string, error)) (int, error) {
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

		if navFilename != "" {
			var (
				navFolderPath = filepath.Dir(navFilename) // used instead of the file path as that results in an additional "../" being added
			)
			navFileContents, err := getContentByFileName(navFilename)
			if err != nil {
				return 0, err
			}

			var relativeTlNotesPath string
			relativeTlNotesPath, err = filepath.Rel(navFolderPath, tlNotesFilePath)
			if err != nil {
				return 0, fmt.Errorf("Failed to determine relative path between nav file %q and file %q: %w", navFilename, tlNotesFilePath, err)
			}

			nameToUpdatedContents[navFilename] = AddFileToNav(navFileContents, relativeTlNotesPath, "Translator's Notes")
		}
	}

	return len(translatorNoteListItems), nil
}
