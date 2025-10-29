//go:build unit

package rulefixes_test

import (
	"testing"

	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
	"github.com/stretchr/testify/assert"
)

type handleDuplicateIDTestCase struct {
	name     string
	contents string
	id       string
	expected string
}

var handleDuplicateIDTestCases = []handleDuplicateIDTestCase{
	{
		name: "id not present returns original and zero",
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
	{
		name: "two duplicate ids get _2 suffix on second occurrence and diff of 2",
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
	{
		name: "three duplicate ids get _2 and _3, total diff of 4 and no double _2",
		contents: `<div id="chapter1"></div>
<div id="chapter1"></div>
<div id="chapter1"></div>`,
		id: "chapter1",
		expected: `<div id="chapter1"></div>
<div id="chapter1_2"></div>
<div id="chapter1_3"></div>`,
	},
	{
		name: "only exact matches updated even if the id to update is a subset of the id to remove duplicates for",
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
	for _, tc := range handleDuplicateIDTestCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := rulefixes.UpdateDuplicateIds(tc.contents, tc.id)

			assert.Equal(t, tc.expected, actual)
		})
	}
}
