package cmd

import (
	"archive/zip"
	"fmt"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/cli/flags"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var organizeNotesFlags = flags.Flags{
	Flags: []flags.Flag{
		flags.NewFileFlag(true, false, &epubFile, "file", "f", "", "the epub file to move translator's notes to their own file in", []string{"epub"}, true),
	},
}

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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return organizeNotesFlags.Validate()
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := moveTranslatorsNotes(epubFile)

		if err != nil {
			logger.WriteError(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(organizeNotesCmd)

	err := organizeNotesFlags.AddToCmd(organizeNotesCmd)
	if err != nil {
		logger.WriteError(err.Error())
	}
}

func moveTranslatorsNotes(epubFile string) error {
	return epubhandler.UpdateEpub(epubFile, func(zipFiles map[string]*zip.File, w *zip.Writer, epubInfo epubhandler.EpubInfo, opfFolder string) ([]string, error) {
		err := validateFilesExist(opfFolder, epubInfo.HtmlFiles, zipFiles)
		if err != nil {
			return nil, err
		}

		var (
			ncxFilename           = filepath.Join(opfFolder, epubInfo.NcxFile)
			nameToUpdatedContents = map[string]string{}
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

		var navFilename = epubInfo.NavFile
		if opfFolder != "." && opfFolder != "" && navFilename != "" {
			navFilename = filepath.Join(opfFolder, navFilename)
		}

		numberOfTranslatorsNotes, err := epubhandler.MoveTranslatorsNotes(epubInfo.FilePathsInSpineOrder, opfFolder, ncxFilename, epubInfo.OpfFile, navFilename, nameToUpdatedContents, getFileContentsByName)
		if err != nil {
			return nil, err
		}

		for filename, updatedContents := range nameToUpdatedContents {
			handledFiles = append(handledFiles, filename)

			err = filehandler.WriteZipCompressedString(w, filename, updatedContents)
			if err != nil {
				return nil, err
			}
		}

		if numberOfTranslatorsNotes > 0 {
			var notesPluralization = "s"
			if numberOfTranslatorsNotes == 1 {
				notesPluralization = ""
			}

			logger.WriteInfof("Found %d translator's note%s.\nAdding translator's notes file.\n", numberOfTranslatorsNotes, notesPluralization)
		} else {
			logger.WriteInfo("No translator's notes found.")
		}

		return handledFiles, nil
	})
}
