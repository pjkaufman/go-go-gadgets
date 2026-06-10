//go:build unit

package epubhandler_test

import (
	"testing"

	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
	"github.com/stretchr/testify/assert"
)

type addFileToOpfTestCase struct {
	inputText string
	filename  string
	id        string
	mediaType string
	expected  string
}

var addFileToOpfTestCases = map[string]addFileToOpfTestCase{
	"When manifest and spine are present and empty, the file is properly added to both": {
		inputText: `<package>
  <manifest>
  </manifest>
  <spine>
  </spine>
</package>`,
		filename:  "test.xhtml",
		id:        "test-id",
		mediaType: "application/xhtml+xml",
		expected: `<package>
  <manifest>
    <item id="test-id" href="test.xhtml" media-type="application/xhtml+xml"/>
</manifest>
  <spine>
    <itemref idref="test-id"/>
</spine>
</package>`,
	},
	"When the manifest is on a single line, the file should be added and the ending manifest tag should now be on its own line": {
		inputText: `<package><manifest></manifest><spine></spine></package>`,
		filename:  "test.xhtml",
		id:        "test-id",
		mediaType: "application/xhtml+xml",
		expected: `<package><manifest>  <item id="test-id" href="test.xhtml" media-type="application/xhtml+xml"/>
</manifest><spine>  <itemref idref="test-id"/>
</spine></package>`,
	},
	"When the spine is on a single line, the file should be added and the ending spine tag should now be on its own line": {
		inputText: `<package>
  <manifest>
    <item id="item1" href="chapter1.xhtml" media-type="application/xhtml+xml"/>
  </manifest>
  <spine></spine>
</package>`,
		filename:  "test.xhtml",
		id:        "test-id",
		mediaType: "application/xhtml+xml",
		expected: `<package>
  <manifest>
    <item id="item1" href="chapter1.xhtml" media-type="application/xhtml+xml"/>
    <item id="test-id" href="test.xhtml" media-type="application/xhtml+xml"/>
</manifest>
  <spine>  <itemref idref="test-id"/>
</spine>
</package>`,
	},
	"When the manifest and spine each have entries on their own line, the file should be added correctly": {
		inputText: `<package>
  <manifest>
    <item id="item1" href="chapter1.xhtml" media-type="application/xhtml+xml"/>
    <item id="item2" href="chapter2.xhtml" media-type="application/xhtml+xml"/>
  </manifest>
  <spine>
    <itemref idref="item1"/>
    <itemref idref="item2"/>
  </spine>
</package>`,
		filename:  "test.xhtml",
		id:        "test-id",
		mediaType: "application/xhtml+xml",
		expected: `<package>
  <manifest>
    <item id="item1" href="chapter1.xhtml" media-type="application/xhtml+xml"/>
    <item id="item2" href="chapter2.xhtml" media-type="application/xhtml+xml"/>
    <item id="test-id" href="test.xhtml" media-type="application/xhtml+xml"/>
</manifest>
  <spine>
    <itemref idref="item1"/>
    <itemref idref="item2"/>
    <itemref idref="test-id"/>
</spine>
</package>`,
	},
}

func TestAddFileToOpf(t *testing.T) {
	t.Parallel()

	for name, tc := range addFileToOpfTestCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			result := epubhandler.AddFileToOpf(tc.inputText, tc.filename, tc.id, tc.mediaType)
			assert.Equal(t, tc.expected, result)
		})
	}
}
