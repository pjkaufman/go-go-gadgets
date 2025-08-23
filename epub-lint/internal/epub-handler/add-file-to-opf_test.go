//go:build unit

package epubhandler_test

import (
	"testing"

	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
	"github.com/stretchr/testify/assert"
)

type addFileToOpfTestCase struct {
	name      string
	inputText string
	filename  string
	id        string
	mediaType string
	expected  string
}

func TestAddFileToOpf(t *testing.T) {
	tests := []addFileToOpfTestCase{
		{
			name: "manifest and spine present and empty",
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
		{
			name:      "manifest is on a single line",
			inputText: `<package><manifest></manifest><spine></spine></package>`,
			filename:  "test.xhtml",
			id:        "test-id",
			mediaType: "application/xhtml+xml",
			expected: `<package><manifest>  <item id="test-id" href="test.xhtml" media-type="application/xhtml+xml"/>
</manifest><spine>  <itemref idref="test-id"/>
</spine></package>`,
		},
		{
			name: "spine is on a single line",
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
		{
			name: "manifest and spine each have entries on their own line",
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

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := epubhandler.AddFileToOpf(tc.inputText, tc.filename, tc.id, tc.mediaType)
			assert.Equal(t, tc.expected, result)
		})
	}
}
