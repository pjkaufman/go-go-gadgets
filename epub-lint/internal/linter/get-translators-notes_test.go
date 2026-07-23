//go:build unit

package linter_test

import (
	"fmt"
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/linter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type getTranslatorsNotesTestCase struct {
	inputText      string
	fileName       string
	noteFileName   string
	startingNumber int
	expectedText   string
	expectedNotes  []string
	expectedNext   int
	expectedError  error
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
	`starting number is 0 and there is a single translator's note that makes up the entirety of the paragraph`: {
		inputText:      `<p class="block_16"><span class="text_4">TL Note: This is a pun that I unfortunately couldn't properly translate to English. The word that was used for break was "</span><span class="text_5">水入り</span><span class="text_4">". The pun is, she said 'literally'. So it translates as "let's get some water in there".</span></p>`,
		fileName:       "main.xhtml",
		noteFileName:   "notes.xhtml",
		startingNumber: 0,
		expectedText:   `<p class="block_16"><a id="note_ref_1" href="notes.xhtml#tl_note_1"><sup>1</sup></a></p>`,
		expectedNotes: []string{
			`<li id="tl_note_1"><span class="text_4">This is a pun that I unfortunately couldn't properly translate to English. The word that was used for break was "</span><span class="text_5">水入り</span><span class="text_4">". The pun is, she said 'literally'. So it translates as "let's get some water in there".</span><br/><a href="main.xhtml#note_ref_1">Back to Reference</a></li>
`,
		},
		expectedNext: 1,
	},
	`a translator's note in square brackets should get its content properly extracted`: {
		inputText:      `<p>　Audrey is the legitimate wife of the Prince, in other words, the Princess.[TN: it says queen (王妃 'ohi'), but I didn't feel like using it, so changed it to Princess, suggestions are welcomed though.] </p>`,
		fileName:       "main.xhtml",
		noteFileName:   "notes.xhtml",
		startingNumber: 0,
		expectedText:   `<p>　Audrey is the legitimate wife of the Prince, in other words, the Princess.<a id="note_ref_1" href="notes.xhtml#tl_note_1"><sup>1</sup></a> </p>`,
		expectedNotes: []string{
			`<li id="tl_note_1">it says queen (王妃 'ohi'), but I didn't feel like using it, so changed it to Princess, suggestions are welcomed though.<br/><a href="main.xhtml#note_ref_1">Back to Reference</a></li>
`,
		},
	},
	`a translator's note with an html entity in it causes an error`: {
		inputText:      `<p class="block_16"><span class="text_4">TL Note: This is a pun that I unfortunately couldn&#8216;t properly translate to English. The word that was used for break was "</span><span class="text_5">水入り</span><span class="text_4">". The pun is, she said 'literally'. So it translates as "let's get some water in there".</span></p>`,
		fileName:       "main.xhtml",
		noteFileName:   "notes.xhtml",
		startingNumber: 0,
		expectedText:   "",
		expectedNotes:  []string{},
		expectedNext:   0,
		expectedError:  fmt.Errorf(`file %q had issues determining translator's notes: attempting to find translator's note text %q failed. This likely means that the source text has html entities. Please convert them to the corresponding character and then try again.`, "main.xhtml", `<span class="text_4">TL Note: This is a pun that I unfortunately couldn‘t properly translate to English. The word that was used for break was "</span><span class="text_5">水入り</span><span class="text_4">". The pun is, she said 'literally'. So it translates as "let's get some water in there".</span>`),
	},
}

func TestGetTranslatorsNotes(t *testing.T) {
	t.Parallel()

	for name, args := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			updatedText, notes, next, err := linter.GetTranslatorsNotes(args.inputText, args.fileName, args.noteFileName, args.startingNumber)

			if args.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, args.expectedError.Error(), err.Error(), "Expected and actual error text differ")
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, args.expectedText, updatedText, "expected text output mismatch")
			assert.Equal(t, args.expectedNotes, notes, "expected notes mismatch")
			assert.Equal(t, args.expectedNext, next, "expected next note number mismatch")
		})
	}
}
