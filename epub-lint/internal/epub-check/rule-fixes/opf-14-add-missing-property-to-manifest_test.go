//go:build unit

package rulefixes_test

import (
	"testing"

	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
	"github.com/stretchr/testify/require"
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
		expectedOutput: `
<package version="3.0">
<manifest>
<item href="OEBPS/chapter1.xhtml" media-type="application/xhtml+xml" properties="scripted"/>
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
		expectedOutput: `
<package version="3.0">
<manifest>
<item href="OEBPS/nav.xhtml" media-type="application/xhtml+xml" properties="scripted nav"/>
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
		expectedOutput: `
<package version="3.0">
<manifest>
<item href="OEBPS/chapter2.xhtml" media-type="application/xhtml+xml" properties="scripted"/>
</manifest>
</package>`,
	},
	"Add scripted to properties attribute if the path is relative using `../` rather than an absolute style relative path": {
		inputText: `
<package version="3.0">
<manifest>
<item id="id-7" href="../titlepage.xhtml" media-type="application/xhtml+xml" properties="calibre:title-page"/>
</manifest>
</package>`,
		inputPath: "../titlepage.xhtml",
		expectedOutput: `
<package version="3.0">
<manifest>
<item id="id-7" href="../titlepage.xhtml" media-type="application/xhtml+xml" properties="scripted calibre:title-page"/>
</manifest>
</package>`,
	},
	"Add scripted to item with URL encoded filepath": {
		inputText: `
<package version="3.0">
<manifest>
<item href="OEBPS/my%20chapter.xhtml" media-type="application/xhtml+xml"/>
</manifest>
</package>`,
		inputPath: "OEBPS/my chapter.xhtml",
		expectedOutput: `
<package version="3.0">
<manifest>
<item href="OEBPS/my%20chapter.xhtml" media-type="application/xhtml+xml" properties="scripted"/>
</manifest>
</package>`,
	},
}

func TestAddPropertyToManifest(t *testing.T) {
	t.Parallel()

	for name, args := range addScriptedToManifestTestCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			edit, err := rulefixes.AddPropertyToManifest(args.inputText, args.inputPath, "scripted")

			require.NoError(t, err)
			checkFinalOutputMatches(t, args.inputText, args.expectedOutput, edit)
		})
	}
}
