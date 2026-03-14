//go:build unit

package epubhandler_test

import (
	"testing"

	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
	"github.com/stretchr/testify/assert"
)

type removeFileFromNavTestCase struct {
	input    string
	file     string
	expected string
}

const (
	simpleNav = `<ol>
	<li>
		<a href="Text/cover.xhtml">Cover</a>
	</li>
	<li>
		<a href="Text/remove.xhtml">Remove Me</a>
	</li>
	<li>
		<a href="Text/prologue.xhtml">Prologue</a>
	</li>
</ol>`
	updatedSimpleNav = `<ol>
	<li>
		<a href="Text/cover.xhtml">Cover</a>
	</li>
	<li>
		<a href="Text/prologue.xhtml">Prologue</a>
	</li>
</ol>`
	sameLineSimpleNav = `<ol>
	<li>
		<a href="Text/cover.xhtml">Cover</a>
	</li><li><a href="Text/remove.xhtml">Remove Me</a></li>
	<li>
		<a href="Text/prologue.xhtml">Prologue</a>
	</li>
</ol>`
	unUpdatedNav = `<p>
	<a href="Text/remove.xhtml">Remove Me</a>
</p>`
)

var removeFileFromNavTestCases = map[string]removeFileFromNavTestCase{
	"When the nav file doesn't include the specified file, nothing will happen to the nav file": {
		input:    simpleNav,
		expected: simpleNav,
		file:     "notpresent.xhtml",
	},
	"When the nav file includes the specified file, the list item and the anchor tag referencing the file should be removed": {
		input:    simpleNav,
		expected: updatedSimpleNav,
		file:     "remove.xhtml",
	},
	"When the nav file includes the specified file and it is on the same line as another list item, the list item and the anchor tag referencing the file should be removed with ending whitespace kept": {
		input:    sameLineSimpleNav,
		expected: updatedSimpleNav,
		file:     "remove.xhtml",
	},
	"When the nav file includes the specified file, but it is not an anchor tag inside a list item, it does not remove the file": {
		input:    unUpdatedNav,
		expected: unUpdatedNav,
		file:     "remove.xhtml",
	},
}

func TestRemoveFileFromNav(t *testing.T) {
	for name, tc := range removeFileFromNavTestCases {
		t.Run(name, func(t *testing.T) {
			actual := epubhandler.RemoveFileFromNav(tc.input, tc.file)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
