//go:build unit

package epubhandler_test

import (
	"fmt"
	"testing"

	_ "embed"

	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type moveTranslatorsNotesTestCase struct {
	opfFolder, ncxFilename, opfFilename, navFilename string
	expectedTranslatorNoteCount                      int
	expectedFileState, validFilesToInitialContent    map[string]string // filename to content
	spineOrder                                       []string
}

var (
	//go:embed testdata/move-translators-notes/simple-html-file-no-translators-notes.xhtml
	noTranslatorNotesXhtml string
	//go:embed testdata/move-translators-notes/html-file-with-single-translators-note.html
	htmlSingleTranslatorNoteOriginal string
	//go:embed testdata/move-translators-notes/html-file-with-single-translators-note_updated.html
	htmlSingleTranslatorNoteExpected string
	//go:embed testdata/move-translators-notes/first-multiple-translators-notes.html
	htmlFirstMultipleTranslatorsNotesOriginal string
	//go:embed testdata/move-translators-notes/first-multiple-translators-notes_updated.html
	htmlFirstMultipleTranslatorsNotesExpected string
	//go:embed testdata/move-translators-notes/second-multiple-translators-notes.html
	htmlSecondMultipleTranslatorsNotesOriginal string
	//go:embed testdata/move-translators-notes/second-multiple-translators-notes_updated.html
	htmlSecondMultipleTranslatorsNotesExpected string
	//go:embed testdata/move-translators-notes/simple.opf
	opfSimpleOriginal string
	//go:embed testdata/move-translators-notes/simple_updated.opf
	opfSimpleExpected string
	//go:embed testdata/move-translators-notes/simple.ncx
	ncxSimpleOriginal string
	//go:embed testdata/move-translators-notes/simple_updated.ncx
	ncxSimpleExpected string
	//go:embed testdata/move-translators-notes/no-opf-folder.opf
	opfNoOpfFolderOriginal string
	//go:embed testdata/move-translators-notes/no-opf-folder_updated.opf
	opfNoOpfFolderExpected string
	//go:embed testdata/move-translators-notes/no-opf-folder.ncx
	ncxNoOpfFolderOriginal string
	//go:embed testdata/move-translators-notes/no-opf-folder_updated.ncx
	ncxNoOpfFolderExpected string
	//go:embed testdata/move-translators-notes/translators-notes-simple.xhtml
	xhtmlSimpleTranslatorsNotes string
	//go:embed testdata/move-translators-notes/translators-notes-complex.xhtml
	xhtmlComplexTranslatorsNotes string
	//go:embed testdata/move-translators-notes/no-opf-folder-nav.xhtml
	xhtmlNoOpfFolderNavOriginal string
	//go:embed testdata/move-translators-notes/no-opf-folder-nav_updated.xhtml
	xhtmlNoOpfFolderNavExpected string
)

var moveTranslatorsNotesTestCases = map[string]moveTranslatorsNotesTestCase{
	"When no translator's notes are present in any files, no changes are made": {
		opfFolder:   "OPS",
		ncxFilename: "OPS/toc.ncx",
		opfFilename: "OPS/content.opf",
		spineOrder: []string{
			"Text/section-0001.html",
			"Text/section-0002.html",
		},
		expectedTranslatorNoteCount: 0,
		expectedFileState: map[string]string{
			"OPS/Text/section-0001.html": noTranslatorNotesXhtml,
			"OPS/Text/section-0002.html": noTranslatorNotesXhtml,
		},
		validFilesToInitialContent: map[string]string{
			"OPS/Text/section-0001.html": noTranslatorNotesXhtml,
			"OPS/Text/section-0002.html": noTranslatorNotesXhtml,
		},
	},
	"When a single translator's note is found, then it should be replaced with a reference and the corresponding changes should be made to the OPF and NCX files": {
		opfFolder:   "OPS",
		ncxFilename: "OPS/toc.ncx",
		opfFilename: "OPS/content.opf",
		spineOrder: []string{
			"Text/section-0001.html",
			"Text/section-0002.html",
		},
		expectedTranslatorNoteCount: 1,
		expectedFileState: map[string]string{
			"OPS/Text/section-0001.html": noTranslatorNotesXhtml,
			"OPS/Text/section-0002.html": htmlSingleTranslatorNoteExpected,
			"OPS/toc.ncx":                ncxSimpleExpected,
			"OPS/content.opf":            opfSimpleExpected,
			"OPS/Text/tl_notes.xhtml":    xhtmlSimpleTranslatorsNotes,
		},
		validFilesToInitialContent: map[string]string{
			"OPS/Text/section-0001.html": noTranslatorNotesXhtml,
			"OPS/Text/section-0002.html": htmlSingleTranslatorNoteOriginal,
			"OPS/toc.ncx":                ncxSimpleOriginal,
			"OPS/content.opf":            opfSimpleOriginal,
		},
	},
	"When multiple translator's notes are found, then they should be replaced with a reference and the corresponding changes should be made to the OPF and NCX files": {
		opfFolder:   "OPS",
		ncxFilename: "OPS/toc.ncx",
		opfFilename: "OPS/content.opf",
		spineOrder: []string{
			"Text/section-0001.html",
			"Text/section-0002.html",
		},
		expectedTranslatorNoteCount: 4,
		expectedFileState: map[string]string{
			"OPS/Text/section-0001.html": htmlFirstMultipleTranslatorsNotesExpected,
			"OPS/Text/section-0002.html": htmlSecondMultipleTranslatorsNotesExpected,
			"OPS/toc.ncx":                ncxSimpleExpected,
			"OPS/content.opf":            opfSimpleExpected,
			"OPS/Text/tl_notes.xhtml":    xhtmlComplexTranslatorsNotes,
		},
		validFilesToInitialContent: map[string]string{
			"OPS/Text/section-0001.html": htmlFirstMultipleTranslatorsNotesOriginal,
			"OPS/Text/section-0002.html": htmlSecondMultipleTranslatorsNotesOriginal,
			"OPS/toc.ncx":                ncxSimpleOriginal,
			"OPS/content.opf":            opfSimpleOriginal,
		},
	},
	"When the OPF folder is the base of the zip file and translator's notes are found, the generated file paths for tl_notes, the NCX, Nav, and OPF files is correct": {
		opfFolder:   ".", // for some reason it is represented by a period in the examples I have encountered
		ncxFilename: "toc.ncx",
		opfFilename: "content.opf",
		navFilename: "OPS/Text/nav.xhtml",
		spineOrder: []string{
			"OPS/Text/section-0001.html",
			"OPS/Text/section-0002.html",
		},
		expectedTranslatorNoteCount: 1,
		expectedFileState: map[string]string{
			"OPS/Text/section-0001.html": noTranslatorNotesXhtml,
			"OPS/Text/section-0002.html": htmlSingleTranslatorNoteExpected,
			"toc.ncx":                    ncxNoOpfFolderExpected,
			"content.opf":                opfNoOpfFolderExpected,
			"OPS/Text/tl_notes.xhtml":    xhtmlSimpleTranslatorsNotes,
			"OPS/Text/nav.xhtml":         xhtmlNoOpfFolderNavExpected,
		},
		validFilesToInitialContent: map[string]string{
			"OPS/Text/section-0001.html": noTranslatorNotesXhtml,
			"OPS/Text/section-0002.html": htmlSingleTranslatorNoteOriginal,
			"toc.ncx":                    ncxNoOpfFolderOriginal,
			"content.opf":                opfNoOpfFolderOriginal,
			"OPS/Text/nav.xhtml":         xhtmlNoOpfFolderNavOriginal,
		},
	},
}

func createTestCaseFileHandlerFunction(validFilesToContent map[string]string, currentContents map[string]string) func(string) (string, error) {
	return func(s string) (string, error) {
		if content, ok := currentContents[s]; ok {
			return content, nil
		} else if content, ok := validFilesToContent[s]; ok {
			return content, nil
		}

		return "", fmt.Errorf("unexpected attempt to get file contents for file %q", s)
	}
}

func TestMoveTranslatorsNotes(t *testing.T) {
	for name, tc := range moveTranslatorsNotesTestCases {
		t.Run(name, func(t *testing.T) {
			var nameToUpdatedFileContents = map[string]string{}

			actualTranslatorNoteCount, err := epubhandler.MoveTranslatorsNotes(tc.spineOrder, tc.opfFolder, tc.ncxFilename, tc.opfFilename, tc.navFilename, nameToUpdatedFileContents, createTestCaseFileHandlerFunction(tc.validFilesToInitialContent, nameToUpdatedFileContents))
			require.NoError(t, err)
			assert.Equal(t, tc.expectedTranslatorNoteCount, actualTranslatorNoteCount)

			for name, expectedContents := range tc.expectedFileState {
				actualContents, ok := nameToUpdatedFileContents[name]
				assert.True(t, ok, "expected file %q to be updated, but it was not", name)
				assert.Equal(t, expectedContents, actualContents, "expected file contents for %q did not match actual contents", name)
			}
		})
	}
}
