//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

type removeScriptedFromManifest struct {
	inputText      string
	inputPath      string
	expectedOutput string
}

var removeScriptedFromManifestTestCases = map[string]removeScriptedFromManifest{
	"Remove properties attribute if attribute is already present for item matching path file name and only has scripted present": {
		inputText: `
<package version="3.0">
<manifest>
<item href="OEBPS/chapter1.xhtml" media-type="application/xhtml+xml" properties="scripted"/>
</manifest>
</package>`,
		inputPath: "OEBPS/chapter1.xhtml",
		expectedOutput: `
<package version="3.0">
<manifest>
<item href="OEBPS/chapter1.xhtml" media-type="application/xhtml+xml"/>
</manifest>
</package>`,
	},
	"Remove scripted from properties attribute if the attribute is already present for item matching path file name and is not the only value": {
		inputText: `
<package version="3.0">
<manifest>
<item href="OEBPS/nav.xhtml" media-type="application/xhtml+xml" properties="nav scripted"/>
</manifest>
</package>`,
		inputPath: "OEBPS/nav.xhtml",
		expectedOutput: `
<package version="3.0">
<manifest>
<item href="OEBPS/nav.xhtml" media-type="application/xhtml+xml" properties="nav"/>
</manifest>
</package>`,
	},
}

func TestRemoveScriptedFromManifest(t *testing.T) {
	for name, args := range removeScriptedFromManifestTestCases {
		t.Run(name, func(t *testing.T) {
			actual, err := linter.RemoveScriptedFromManifest(args.inputText, args.inputPath)

			assert.Nil(t, err)
			assert.Equal(t, args.expectedOutput, actual)
		})
	}
}
