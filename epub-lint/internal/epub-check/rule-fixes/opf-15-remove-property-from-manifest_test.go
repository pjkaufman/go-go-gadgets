//go:build unit

package rulefixes_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
	"github.com/stretchr/testify/assert"
)

type removeScriptedFromManifest struct {
	inputText      string
	inputPath      string
	property       string
	expectedChange positions.TextEdit
}

var removeScriptedFromManifestTestCases = map[string]removeScriptedFromManifest{
	"Remove properties attribute if attribute is already present for item matching path file name and only has property that is to be removed present": {
		inputText: `
<package version="3.0">
<manifest>
<item href="OEBPS/chapter1.xhtml" media-type="application/xhtml+xml" properties="scripted"/>
</manifest>
</package>`,
		inputPath: "OEBPS/chapter1.xhtml",
		property:  "scripted",
		expectedChange: positions.TextEdit{
			Range: positions.Range{
				Start: positions.Position{
					Line:   4,
					Column: 59,
				},
				End: positions.Position{
					Line:   4,
					Column: 80,
				},
			},
		},
	},
	"Remove the specified property from properties attribute if the attribute is already present for item matching path file name and is not the only value": {
		inputText: `
<package version="3.0">
<manifest>
<item href="OEBPS/nav.xhtml" media-type="application/xhtml+xml" properties="nav scripted"/>
</manifest>
</package>`,
		inputPath: "OEBPS/nav.xhtml",
		property:  "scripted",
		expectedChange: positions.TextEdit{
			Range: positions.Range{
				Start: positions.Position{
					Line:   4,
					Column: 70,
				},
				End: positions.Position{
					Line:   4,
					Column: 80,
				},
			},
		},
	},
}

func TestRemoveScriptedFromManifest(t *testing.T) {
	for name, args := range removeScriptedFromManifestTestCases {
		t.Run(name, func(t *testing.T) {
			actual, err := rulefixes.RemovePropertyFromManifest(args.inputText, args.inputPath, args.property)

			assert.Nil(t, err)
			assert.Equal(t, args.expectedChange, actual)
		})
	}
}
