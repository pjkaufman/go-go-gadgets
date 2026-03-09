//go:build unit

package compare_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/compare"
	"github.com/stretchr/testify/assert"
)

type compareLinesTestCase struct {
	pdfLines, htmlLines []string
	differences         []compare.Difference
}

var sampleLines = []string{
	"The morning sun cast a warm glow across the quiet street.",
	"A gentle breeze rustled the leaves outside the window.",
	"She paused for a moment to appreciate the peaceful silence.",
	"The library was filled with the soft hum of turning pages.",
	"He brewed a fresh cup of coffee before starting his day.",
	"The dog wagged its tail excitedly at the sound of footsteps.",
	"A colorful bird perched on the fence and began to sing.",
	"They planned a weekend trip to explore the nearby mountains.",
	"The classroom buzzed with anticipation before the lesson began.",
	"A light rain tapped rhythmically against the rooftop.",
	"She found an old photograph tucked inside a dusty book.",
	"The aroma of freshly baked bread drifted through the kitchen.",
	"He took a deep breath and stepped onto the stage.",
	"The park was alive with families enjoying the sunny afternoon.",
	"A small boat drifted lazily across the calm lake.",
	"They shared stories around the campfire late into the night.",
	"The cat curled up in a warm patch of sunlight.",
	"He admired the intricate patterns on the handmade quilt.",
	"A distant train whistle echoed through the valley.",
	"She planted a row of flowers along the garden path.",
	"The museum displayed artifacts from ancient civilizations.",
	"A soft melody played from the radio in the background.",
	"They watched the stars appear one by one in the night sky.",
	"The bakery opened early to serve the morning crowd.",
	"He sketched the landscape with quick, confident strokes.",
	"A friendly stranger offered directions to the lost traveler.",
	"The festival lights shimmered brightly against the evening sky.",
	"She organized her desk before beginning the new project.",
	"The waves rolled gently onto the sandy shore.",
	"He discovered a new hiking trail hidden behind the old bridge.",
	"A group of children laughed as they chased bubbles in the yard.",
	"The clock tower chimed softly at the top of the hour.",
	"She wrapped the gift carefully with bright paper and ribbon.",
	"The scent of pine trees filled the crisp mountain air.",
	"He paused to watch the clouds drift slowly overhead.",
	"A handwritten note was left on the kitchen table.",
	"They enjoyed a quiet picnic under the shade of an oak tree.",
	"The bookstore owner recommended a novel she thought he'd enjoy.",
	"A lantern glowed warmly on the porch as evening approached.",
	"She practiced the piano piece until it sounded just right.",
	"The town square was decorated for the upcoming celebration.",
	"He noticed the first signs of spring blooming in the garden.",
	"A gentle laugh echoed from the other side of the room.",
	"They took a scenic route to enjoy the countryside views.",
	"The old clock ticked steadily in the corner of the room.",
	"She watched as snowflakes drifted softly to the ground.",
	"He packed his backpack carefully before the long journey.",
	"A pair of ducks glided gracefully across the pond.",
	"They ended the day with a warm meal and good conversation.",
}

var compareLinesTestCases = map[string]compareLinesTestCase{
	// "When there is a difference in lines between the PDF and HTML content, there should be a difference mentioning that differing line count": {
	// 	pdfLines:  []string{"Line 1", "Line 2"},
	// 	htmlLines: []string{"Line 1"},
	// 	differences: []compare.Difference{
	// 		{
	// 			Message:  "Line count mismatch for HTML and PDF file: expected 1 but was 2",
	// 			DiffType: compare.LikelyMismatch,
	// 		},
	// 		{
	// 			Message:  "Ran out of lines in the HTML to compare to the PDF: had 1 line to go",
	// 			DiffType: compare.DefiniteMismatch,
	// 		},
	// 	},
	// },
	// "When there is a difference in lines between the PDF and HTML content and the PDF has more HTML has more lines than the PDF, there should be a difference mentioning that differing line count": {
	// 	pdfLines:  []string{"Line 1", "Line 2"},
	// 	htmlLines: []string{"Line 1", "Line 2", "Line 3", "Line 4", "Line 5", "Line 6", "Line 7"},
	// 	differences: []compare.Difference{
	// 		{
	// 			Message:  "Line count mismatch for HTML and PDF file: expected 7 but was 2",
	// 			DiffType: compare.LikelyMismatch,
	// 		},
	// 		{
	// 			Message:  "Ran out of lines in the PDF to compare to the HTML: had 5 lines to go",
	// 			DiffType: compare.DefiniteMismatch,
	// 		},
	// 	},
	// },
	// "When a there is a line in the HTML that gets broken into multiple in the PDF, it should get reported as a line wrapped line": {
	// 	pdfLines:  []string{"Line 1", "Line 2", "Here is a line", "that is being broken in two"},
	// 	htmlLines: []string{"Line 1", "Line 2", "Here is a line that is being broken in two"},
	// 	differences: []compare.Difference{
	// 		{
	// 			Message:  "Line count mismatch for HTML and PDF file: expected 3 but was 4",
	// 			DiffType: compare.LikelyMismatch,
	// 		},
	// 		{
	// 			Message:  `HTML line 3 matches across 2 PDF lines: "Here is a line that is being broken in two"`,
	// 			DiffType: compare.WrappedLine,
	// 		},
	// 	},
	// },
	"When a there is a line in the HTML that gets broken into multiple in the PDF, but it only partially matches the lines in the PDF, it should get reported as a partially wrapped line": {
		pdfLines:  []string{"Line 1", "Line 2", "Here is a line", "that is being broken"},
		htmlLines: []string{"Line 1", "Line 2", "Here is a line that is being broken in two"},
		differences: []compare.Difference{
			{
				Message:  "Line count mismatch for HTML and PDF file: expected 3 but was 4",
				DiffType: compare.LikelyMismatch,
			},
			{
				Message:  `HTML line 3 partially across 2 PDF lines: "Here is a line that is being broken in two"`,
				DiffType: compare.PartiallyWrappedLine,
			},
		},
	},
	// "When the only difference in a line is a whitespace character, the difference reported should be a whitespace difference": {
	// 	pdfLines:  []string{"Line  1", "Line 2 ", " Here is a line"},
	// 	htmlLines: []string{"Line 1", "Line 2", "Here is a line"},
	// 	differences: []compare.Difference{
	// 		{
	// 			Message:  `Line 1 vs. 1 differs only by whitespace (HTML: "Line 1" | PDF: "Line  1")`,
	// 			DiffType: compare.Whitespace,
	// 		},
	// 		{
	// 			Message:  `Line 2 vs. 2 differs only by whitespace (HTML: "Line 2" | PDF: "Line 2 ")`,
	// 			DiffType: compare.Whitespace,
	// 		},
	// 		{
	// 			Message:  `Line 3 vs. 3 differs only by whitespace (HTML: "Here is a line" | PDF: " Here is a line")`,
	// 			DiffType: compare.Whitespace,
	// 		},
	// 	},
	// },
	// "When a line does not match at all, it should get reported as a line difference": {
	// 	pdfLines:  []string{"Line 1", "Line 2", "Here is a line", "Line 4", "Yet another line"},
	// 	htmlLines: []string{"Line 1", "Line 2", "A different line", "Line 4", "Another different line"},
	// 	differences: []compare.Difference{
	// 		{
	// 			Message: `Line 3 does not match:
	// HTML: "A different line"
	// PDF:  "Here is a line"`,
	// 			DiffType: compare.Line,
	// 		},
	// 		{
	// 			Message: `Line 5 does not match:
	// HTML: "Another different line"
	// PDF:  "Yet another line"`,
	// 			DiffType: compare.Line,
	// 		},
	// 	},
	// },
	// "When the contents of the HTML and PDF lines match exactly, no differences should be reported": {
	// 	pdfLines:  sampleLines,
	// 	htmlLines: sampleLines,
	// },
	// "When there are various differences between the HTML and PDF lines, they should all get reported": {
	// 	pdfLines:  []string{"Line 1", "This line", "is getting wrapped", "not same", "broken", "into", "words", "line", "Whitespace  diff", "partial", "wraps", "are weird 52617"},
	// 	htmlLines: []string{"Line 1", "This line is getting wrapped", "these lines differ", "broken into words line", "Whitespace diff", "partial wraps are weird"},
	// 	differences: []compare.Difference{
	// 		{
	// 			Message:  "Line count mismatch for HTML and PDF file: expected 6 but was 12",
	// 			DiffType: compare.LikelyMismatch,
	// 		},
	// 		{
	// 			Message:  `HTML line 2 matches across 2 PDF lines: "This line is getting wrapped"`,
	// 			DiffType: compare.WrappedLine,
	// 		},
	// 		{
	// 			Message: `Line 3 does not match:
	// HTML: "these lines differ"
	// PDF:  "not same"`,
	// 			DiffType: compare.Line,
	// 		},
	// 		{
	// 			Message:  `HTML line 4 matches across 4 PDF lines: "broken into words line"`,
	// 			DiffType: compare.WrappedLine,
	// 		},
	// 		{
	// 			Message:  `Line 5 vs. 9 differs only by whitespace (HTML: "Whitespace diff" | PDF: "Whitespace  diff")`,
	// 			DiffType: compare.Whitespace,
	// 		},
	// 		{
	// 			Message:  `HTML line 6 partially across 3 PDF lines: "partial wraps are weird"`,
	// 			DiffType: compare.PartiallyWrappedLine,
	// 		},
	// 		{
	// 			Message:  "Ran out of lines in the HTML to compare to the PDF: had 1 line to go",
	// 			DiffType: compare.DefiniteMismatch,
	// 		},
	// },
	// },
}

func TestCompareLines(t *testing.T) {
	for name, args := range compareLinesTestCases {
		t.Run(name, func(t *testing.T) {
			actual := compare.CompareLines(args.pdfLines, args.htmlLines)

			assert.Equal(t, args.differences, actual)
		})
	}
}
