//go:build unit

package rulefixes_test

import (
	"testing"

	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
)

type handleDuplicateIDTestCase struct {
	name     string
	contents string
	id       string
	expected string
}

var handleDuplicateIDTestCases = map[string]handleDuplicateIDTestCase{
	"Id not present returns original and zero": {
		contents: `<html>
  <body>
    <div id="something"></div>
  </body>
</html>`,
		id: "chapter1",
		expected: `<html>
  <body>
    <div id="something"></div>
  </body>
</html>`,
	},
	"Two duplicate ids get _2 suffix on second occurrence and diff of 2": {
		contents: `<html>
  <body>
    <div id="chapter1"></div>
    <span id="chapter1"></span>
  </body>
</html>`,
		id: "chapter1",
		expected: `<html>
  <body>
    <div id="chapter1"></div>
    <span id="chapter1_2"></span>
  </body>
</html>`,
	},
	"Three duplicate ids get _2 and _3, total diff of 4 and no double _2": {
		contents: `<div id="chapter1"></div>
<div id="chapter1"></div>
<div id="chapter1"></div>`,
		id: "chapter1",
		expected: `<div id="chapter1"></div>
<div id="chapter1_2"></div>
<div id="chapter1_3"></div>`,
	},
	"Only exact matches updated even if the id to update is a subset of the id to remove duplicates for": {
		contents: `<div id="chapter1"></div>
<div id="chapter1"></div>
<div id="chapter1-long"></div>`,
		id: "chapter1",
		expected: `<div id="chapter1"></div>
<div id="chapter1_2"></div>
<div id="chapter1-long"></div>`,
	},
}

func TestHandleDuplicateID(t *testing.T) {
	for _, args := range handleDuplicateIDTestCases {
		t.Run(args.name, func(t *testing.T) {
			edits := rulefixes.UpdateDuplicateIds(args.contents, args.id)

			checkFinalOutputMatches(t, args.contents, args.expected, edits...)
		})
	}
}
