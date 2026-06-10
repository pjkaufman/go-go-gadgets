//go:build unit

package epubhandler_test

import (
	"testing"

	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
	"github.com/stretchr/testify/assert"
)

type removeFileFromNxcTestCase struct {
	input            string
	relativeFilePath string
	expected         string
}

var removeFileFromNcxTestCases = map[string]removeFileFromNxcTestCase{
	"When the NCX file does not have the relative path file referenced, no changes are made": {
		input: `<?xml version="1.0" encoding="UTF-8"?>
<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/" version="2005-1">
  <head>
    <meta name="dtb:uid" content="id"/>
  </head>
  <docTitle><text>Sample</text></docTitle>
  <navMap>
    <navPoint id="navPoint-1" playOrder="1">
      <navLabel><text>Chapter 1</text></navLabel>
      <content src="chapter1.xhtml"/>
    </navPoint>
  </navMap>
</ncx>`,
		relativeFilePath: "chapter2.xhtml",
		expected: `<?xml version="1.0" encoding="UTF-8"?>
<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/" version="2005-1">
  <head>
    <meta name="dtb:uid" content="id"/>
  </head>
  <docTitle><text>Sample</text></docTitle>
  <navMap>
    <navPoint id="navPoint-1" playOrder="1">
      <navLabel><text>Chapter 1</text></navLabel>
      <content src="chapter1.xhtml"/>
    </navPoint>
  </navMap>
</ncx>`,
	},
	"When the NCX file does have the relative path file referenced, the matching NavPoint is removed": {
		input: `<?xml version="1.0" encoding="UTF-8"?>
<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/" version="2005-1">
  <head>
    <meta name="dtb:uid" content="id"/>
  </head>
  <docTitle><text>Sample</text></docTitle>
  <navMap>
    <navPoint id="navPoint-1" playOrder="1">
      <navLabel><text>Chapter 1</text></navLabel>
      <content src="chapter1.xhtml"/>
    </navPoint>
    <navPoint id="navPoint-2" playOrder="2">
      <navLabel><text>Chapter 2</text></navLabel>
      <content src="chapter2.xhtml"/>
    </navPoint>
  </navMap>
</ncx>`,
		relativeFilePath: "chapter2.xhtml",
		expected: `<?xml version="1.0" encoding="UTF-8"?>
<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/" version="2005-1">
  <head>
    <meta name="dtb:uid" content="id"/>
  </head>
  <docTitle><text>Sample</text></docTitle>
  <navMap>
    <navPoint id="navPoint-1" playOrder="1">
      <navLabel><text>Chapter 1</text></navLabel>
      <content src="chapter1.xhtml"/>
    </navPoint>
  </navMap>
</ncx>`,
	},
	"When the NCX file does have the relative path file referenced and the referenced NavPoint has its starting and elements not beginning and ending the line, the matching NavPoint is removed, but up until the start and end of the line is left alone": {
		input:            `<?xml version="1.0" encoding="UTF-8"?><ncx xmlns="http://www.daisy.org/z3986/2005/ncx/" version="2005-1"><head><meta name="dtb:uid" content="id"/></head><docTitle><text>Sample</text></docTitle><navMap><navPoint id="np1" playOrder="1"><navLabel><text>Ch1</text></navLabel><content src="ch1.xhtml"/></navPoint><navPoint id="np2" playOrder="2"><navLabel><text>Ch2</text></navLabel><content src="ch2.xhtml"/></navPoint></navMap></ncx>`,
		relativeFilePath: "ch2.xhtml",
		expected:         `<?xml version="1.0" encoding="UTF-8"?><ncx xmlns="http://www.daisy.org/z3986/2005/ncx/" version="2005-1"><head><meta name="dtb:uid" content="id"/></head><docTitle><text>Sample</text></docTitle><navMap><navPoint id="np1" playOrder="1"><navLabel><text>Ch1</text></navLabel><content src="ch1.xhtml"/></navPoint></navMap></ncx>`,
	},
}

func TestRemoveFileFromNcx(t *testing.T) {
	t.Parallel()

	for name, tc := range removeFileFromNcxTestCases {
		t.Parallel()

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actual := epubhandler.RemoveFileFromNcx(tc.input, tc.relativeFilePath)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
