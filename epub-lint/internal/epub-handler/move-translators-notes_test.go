//go:build unit

package epubhandler_test

import (
	"fmt"
	"testing"

	_ "embed"

	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
	"github.com/stretchr/testify/assert"
)

type moveTranslatorsNotesTestCase struct {
	opfFolder, ncxFilename, opfFilename           string
	expectedTranslatorNoteCount                   int
	expectedFileState, validFilesToInitialContent map[string]string // filename to content
	spineOrder                                    []string
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
	//go:embed testdata/move-translators-notes/translators-notes-simple.xhtml
	xhtmlSimpleTranslatorsNotes string
	//go:embed testdata/move-translators-notes/translators-notes-complex.xhtml
	xhtmlComplexTranslatorsNotes string
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
	/* Test cases to add:
	3. Multiple translator's notes present get replaced with the proper references and corresponding tl_notes file is created
	4. No OPF folder (i.e. the OPF file is at the base of the zip) when translator's notes are found generates the correct paths and content for the tl_notes, NCX, and OPF
	5. Current bug scenario gets found and fixed (see #88)
	6. See about the scenario where all of the body content is inside a div and there is a translator's note
	7. Make sure that it handles partial element translator's notes correctly
	8. Make sure that it handles full element translator's notes correctly
	*/
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

			actualTranslatorNoteCount, err := epubhandler.MoveTranslatorsNotes(tc.spineOrder, tc.opfFolder, tc.ncxFilename, tc.opfFilename, nameToUpdatedFileContents, createTestCaseFileHandlerFunction(tc.validFilesToInitialContent, nameToUpdatedFileContents))
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedTranslatorNoteCount, actualTranslatorNoteCount)

			for name, expectedContents := range tc.expectedFileState {
				actualContents, ok := nameToUpdatedFileContents[name]
				assert.True(t, ok, fmt.Sprintf("expected file %q to be updated, but it was not", name))
				assert.Equal(t, expectedContents, actualContents, fmt.Sprintf("expected file contents for %q did not match actual contents", name))
			}
		})
	}
}
