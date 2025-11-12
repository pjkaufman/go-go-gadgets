//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

type getTranslatorsNotesTestCase struct {
	inputText      string
	fileName       string
	noteFileName   string
	startingNumber int
	expectedText   string
	expectedNotes  []string
	expectedNext   int
}

var testCases = map[string]getTranslatorsNotesTestCase{
	"3 author notes with two lines in between each other and starting number is 3": {
		inputText: `<p>Some content before.</p>
<p>TL Note: First author note.</p>
<p>Line between first and second.</p>
<p>Translator's Note: Second author note.</p>
<p>Another line.</p>
<p>Note: Third author note.</p>
<p>Some content after.</p>`,
		fileName:       "currentfile.xhtml",
		noteFileName:   "notes.xhtml",
		startingNumber: 3,
		expectedText: `<p>Some content before.</p>
<p><a id="note_ref_4" href="notes.xhtml#tl_note_4"><sup>4</sup></a></p>
<p>Line between first and second.</p>
<p><a id="note_ref_5" href="notes.xhtml#tl_note_5"><sup>5</sup></a></p>
<p>Another line.</p>
<p><a id="note_ref_6" href="notes.xhtml#tl_note_6"><sup>6</sup></a></p>
<p>Some content after.</p>`,
		expectedNotes: []string{
			`<li id="tl_note_4">First author note.<br/><a href="currentfile.xhtml#note_ref_4">Back to Reference</a></li>
`,
			`<li id="tl_note_5">Second author note.<br/><a href="currentfile.xhtml#note_ref_5">Back to Reference</a></li>
`,
			`<li id="tl_note_6">Third author note.<br/><a href="currentfile.xhtml#note_ref_6">Back to Reference</a></li>
`,
		},
		expectedNext: 6,
	},
	"starting number is 0 and there are no translator's notes meaning no changes to the text": {
		inputText: `<p>This is a normal paragraph.</p>
		<p>Here is some other text.</p>
		<p>Nothing special here.</p>`,
		fileName:       "main.xhtml",
		noteFileName:   "notes.xhtml",
		startingNumber: 0,
		expectedText: `<p>This is a normal paragraph.</p>
		<p>Here is some other text.</p>
		<p>Nothing special here.</p>`,
		expectedNotes: []string{},
		expectedNext:  0,
	},
	"starting number is 0 and there is a single translator's note that is encapsulated in parentheses": {
		inputText: `<p>Something here.</p>
		<p>(TL Note: This is a parenthetical translator's note.)</p>
		<p>The end.</p>`,
		fileName:       "main.xhtml",
		noteFileName:   "notes.xhtml",
		startingNumber: 0,
		expectedText: `<p>Something here.</p>
		<p><a id="note_ref_1" href="notes.xhtml#tl_note_1"><sup>1</sup></a></p>
		<p>The end.</p>`,
		expectedNotes: []string{
			`<li id="tl_note_1">This is a parenthetical translator's note.<br/><a href="main.xhtml#note_ref_1">Back to Reference</a></li>
`,
		},
		expectedNext: 1,
	},
}

func TestGetTranslatorsNotes(t *testing.T) {
	for name, args := range testCases {
		t.Run(name, func(t *testing.T) {
			updatedText, notes, next := linter.GetTranslatorsNotes(args.inputText, args.fileName, args.noteFileName, args.startingNumber)
			assert.Equal(t, args.expectedText, updatedText, "expected text output mismatch")
			assert.Equal(t, args.expectedNotes, notes, "expected notes mismatch")
			assert.Equal(t, args.expectedNext, next, "expected next note number mismatch")
		})
	}
}
