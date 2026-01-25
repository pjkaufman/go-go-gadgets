//go:build unit

package rulefixes_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
	"github.com/stretchr/testify/assert"
)

type addScriptedToManifest struct {
	inputText      string
	inputPath      string
	expectedChange positions.TextEdit
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
		expectedChange: positions.TextEdit{
			Range: positions.Range{
				Start: positions.Position{
					Line:   4,
					Column: 59,
				},
				End: positions.Position{
					Line:   4,
					Column: 59,
				},
			},
			NewText: ` properties="scripted"`,
		},
	},
	"Add scripted to properties attribute if the attribute is already present for item matching path file name": {
		inputText: `
<package version="3.0">
<manifest>
<item href="OEBPS/nav.xhtml" media-type="application/xhtml+xml" properties="nav"/>
</manifest>
</package>`,
		inputPath: "OEBPS/nav.xhtml",
		expectedChange: positions.TextEdit{
			Range: positions.Range{
				Start: positions.Position{
					Line:   4,
					Column: 67,
				},
				End: positions.Position{
					Line:   4,
					Column: 67,
				},
			},
			NewText: `scripted `,
		},
	},
	"Add scripted to properties attribute if it is empty for the path file name": {
		inputText: `
<package version="3.0">
<manifest>
<item href="OEBPS/chapter2.xhtml" media-type="application/xhtml+xml" properties=""/>
</manifest>
</package>`,

		inputPath: "OEBPS/chapter2.xhtml",
		expectedChange: positions.TextEdit{
			Range: positions.Range{
				Start: positions.Position{
					Line:   4,
					Column: 72,
				},
				End: positions.Position{
					Line:   4,
					Column: 72,
				},
			},
			NewText: `scripted`,
		},
	},
}

func TestAddScriptedToManifest(t *testing.T) {
	for name, args := range addScriptedToManifestTestCases {
		t.Run(name, func(t *testing.T) {
			actual, err := rulefixes.AddPropertyToManifest(args.inputText, args.inputPath, "scripted")

			assert.Nil(t, err)
			assert.Equal(t, args.expectedChange, actual)
		})
	}
}
