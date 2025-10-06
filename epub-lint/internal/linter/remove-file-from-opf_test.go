//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

type removeFileFromOpfTestCase struct {
	inputText      string
	inputPath      string
	expectedOutput string
	expectedError  error
}

var removeFileFromOpfTestCases = map[string]removeFileFromOpfTestCase{
	"Remove a file that exists in both the manifest and the spine": {
		inputText: `
<package version="3.0">
<manifest>
<item id="item1" href="OEBPS/chapter1.xhtml" media-type="application/xhtml+xml"/>
</manifest>
<spine>
<itemref idref="item1"/>
</spine>
</package>`,
		inputPath: "OEBPS/chapter1.xhtml",
		expectedOutput: `
<package version="3.0">
<manifest>
</manifest>
<spine>
</spine>
</package>`,
		expectedError: nil,
	},
	"Remove a file that exists just in the manifest": {
		inputText: `
<package version="3.0">
<manifest>
<item id="item1" href="OEBPS/chapter1.xhtml" media-type="application/xhtml+xml"/>
</manifest>
<spine>
</spine>
</package>`,
		inputPath: "OEBPS/chapter1.xhtml",
		expectedOutput: `
<package version="3.0">
<manifest>
</manifest>
<spine>
</spine>
</package>`,
		expectedError: nil,
	},
	"Remove a file that is in the manifest, but has no id": {
		inputText: `
<package version="3.0">
<manifest>
<item href="OEBPS/chapter1.xhtml" media-type="application/xhtml+xml"/>
</manifest>
<spine>
</spine>
</package>`,
		inputPath: "OEBPS/chapter1.xhtml",
		expectedOutput: `
<package version="3.0">
<manifest>
</manifest>
<spine>
</spine>
</package>`,
		expectedError: nil,
	},
	"Remove a file that does not exist in the manifest": {
		inputText: `
<package version="3.0">
<manifest>
<item id="item1" href="OEBPS/chapter1.xhtml" media-type="application/xhtml+xml"/>
</manifest>
<spine>
<itemref idref="item1"/>
</spine>
</package>`,
		inputPath: "OEBPS/chapter2.xhtml",
		expectedOutput: `
<package version="3.0">
<manifest>
<item id="item1" href="OEBPS/chapter1.xhtml" media-type="application/xhtml+xml"/>
</manifest>
<spine>
<itemref idref="item1"/>
</spine>
</package>`,
		expectedError: nil,
	},
	"Remove a file that does not exist in the manifest but has a name that is the suffix of the file and make sure it does not affect the file with the suffix": {
		inputText: `
<package version="3.0">
<manifest>
<item id="item1" href="1-1.png" media-type="image/png"/>
</manifest>
<spine>
<itemref idref="item1"/>
</spine>
</package>`,
		inputPath: "1.png",
		expectedOutput: `
<package version="3.0">
<manifest>
<item id="item1" href="1-1.png" media-type="image/png"/>
</manifest>
<spine>
<itemref idref="item1"/>
</spine>
</package>`,
		expectedError: nil,
	},
	"Remove a file that does exist in the manifest, and there is no spine present": {
		inputText: `
<package version="3.0">
<manifest>
<item id="item1" href="OEBPS/chapter1.xhtml" media-type="application/xhtml+xml"/>
</manifest>
</package>`,
		inputPath:      "OEBPS/chapter1.xhtml",
		expectedOutput: "",
		expectedError:  linter.ErrNoSpine,
	},
}

func TestRemoveFileFromOpf(t *testing.T) {
	for name, args := range removeFileFromOpfTestCases {
		t.Run(name, func(t *testing.T) {
			actual, err := linter.RemoveFileFromOpf(args.inputText, args.inputPath)

			assert.Equal(t, args.expectedError, err)
			assert.Equal(t, args.expectedOutput, actual)
		})
	}
}
