//go:build unit

package epubhandler_test

import (
	"testing"

	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
	"github.com/stretchr/testify/assert"
)

type addFileToNcxTestCase struct {
	inputText string
	filePath  string
	title     string
	id        string
	expected  string
}

var addFileToNcxTestCases = map[string]addFileToNcxTestCase{
	"When there is no navMap in ncx, no change should be made": {
		inputText: `<ncx>
  <head>
    <meta name="dtb:uid" content="BookId"/>
  </head>
  <docTitle>
    <text>Book Title</text>
  </docTitle>
</ncx>`,
		filePath: "chapter1.xhtml",
		title:    "Chapter 1",
		id:       "ch1",
		expected: `<ncx>
  <head>
    <meta name="dtb:uid" content="BookId"/>
  </head>
  <docTitle>
    <text>Book Title</text>
  </docTitle>
</ncx>`,
	},
	"When there are no navPoints already, the playOrder of the add file should be 1": {
		inputText: `<ncx>
  <navMap>
  </navMap>
</ncx>`,
		filePath: "chapter1.xhtml",
		title:    "Chapter 1",
		id:       "ch1",
		expected: `<ncx>
  <navMap>
    <navPoint id="ch1" playOrder="1">
    <navLabel>
      <text>Chapter 1</text>
    </navLabel>
    <content src="chapter1.xhtml"/>
  </navPoint>
</navMap>
</ncx>`,
	},
	"When there are 3 navPoints, the added file's playOrder should be 4": {
		inputText: `<ncx>
  <navMap>
    <navPoint id="np1" playOrder="1">
      <navLabel><text>One</text></navLabel>
      <content src="ch1.xhtml"/>
    </navPoint>
    <navPoint id="np2" playOrder="2">
      <navLabel><text>Two</text></navLabel>
      <content src="ch2.xhtml"/>
    </navPoint>
    <navPoint id="np3" playOrder="3">
      <navLabel><text>Three</text></navLabel>
      <content src="ch3.xhtml"/>
    </navPoint>
  </navMap>
</ncx>`,
		filePath: "chapter4.xhtml",
		title:    "Chapter 4",
		id:       "ch4",
		expected: `<ncx>
  <navMap>
    <navPoint id="np1" playOrder="1">
      <navLabel><text>One</text></navLabel>
      <content src="ch1.xhtml"/>
    </navPoint>
    <navPoint id="np2" playOrder="2">
      <navLabel><text>Two</text></navLabel>
      <content src="ch2.xhtml"/>
    </navPoint>
    <navPoint id="np3" playOrder="3">
      <navLabel><text>Three</text></navLabel>
      <content src="ch3.xhtml"/>
    </navPoint>
    <navPoint id="ch4" playOrder="4">
    <navLabel>
      <text>Chapter 4</text>
    </navLabel>
    <content src="chapter4.xhtml"/>
  </navPoint>
</navMap>
</ncx>`,
	},
}

func TestAddFileToNcx(t *testing.T) {
	for name, tc := range addFileToNcxTestCases {
		t.Run(name, func(t *testing.T) {
			result := epubhandler.AddFileToNcx(tc.inputText, tc.filePath, tc.title, tc.id)
			assert.Equal(t, tc.expected, result)
		})
	}
}
