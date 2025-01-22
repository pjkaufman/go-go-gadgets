//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

type addScriptedToManifest struct {
	inputText      string
	inputPath      string
	expectedOutput string
}

var addScriptedToManifestTestCases = map[string]addScriptedToManifest{
	"Add properties attribute if no attribute is already present for item matching path file name": {
		inputText: `
<package version="3.0">
  <manifest>
    <item href="OEBPS/chapter1.xhtml" media-type="application/xhtml+xml"/>
  </manifest>
</package>`,
		inputPath: "OEBPS/chapter1.xhtml",
		expectedOutput: `<package version="3.0">
  <manifest>
    <item href="OEBPS/chapter1.xhtml" media-type="application/xhtml+xml" properties="scripted"></item>
  </manifest>
</package>`,
	},
	"Add scripted to properties attribute if the attribute is already present for item matching path file name": {
		inputText: `
<package version="3.0">
<manifest>
<item href="OEBPS/nav.xhtml" media-type="application/xhtml+xml" properties="nav"/>
</manifest>
</package>`,
		inputPath: "OEBPS/nav.xhtml",
		expectedOutput: `<package version="3.0">
  <manifest>
    <item href="OEBPS/nav.xhtml" media-type="application/xhtml+xml" properties="nav scripted"></item>
  </manifest>
</package>`,
	},
	"Add scripted to properties attribute if it is empty for the path file name": {
		inputText: `
<package version="3.0">
<manifest>
<item href="OEBPS/chapter2.xhtml" media-type="application/xhtml+xml" properties=""/>
</manifest>
</package>`,
		inputPath: "OEBPS/chapter2.xhtml",
		expectedOutput: `<package version="3.0">
  <manifest>
    <item href="OEBPS/chapter2.xhtml" media-type="application/xhtml+xml" properties="scripted"></item>
  </manifest>
</package>`,
	},
	"Add properties attribute to manifest item if no properties tag exists even if the paths for the href and input path are different": {
		inputText: `
<package version="3.0">
<manifest>
<item href="chapter1.xhtml" media-type="application/xhtml+xml"/>
</manifest>
</package>`,
		inputPath: "OEBPS/chapter1.xhtml",
		expectedOutput: `<package version="3.0">
  <manifest>
    <item href="chapter1.xhtml" media-type="application/xhtml+xml" properties="scripted"></item>
  </manifest>
</package>`,
	},
}

func TestAddScriptedToManifest(t *testing.T) {
	for name, args := range addScriptedToManifestTestCases {
		t.Run(name, func(t *testing.T) {
			actual, err := linter.AddScriptedToManifest(args.inputText, args.inputPath)

			assert.Nil(t, err)
			assert.Equal(t, args.expectedOutput, actual)
		})
	}
}
