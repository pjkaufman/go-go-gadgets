//go:build unit

package rulefixes_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
	"github.com/stretchr/testify/assert"
)

type fixManifestAttributeTestCase struct {
	opfContents           string
	attribute             string
	line                  int
	attributeNameToNumber map[string]int
	expectedChanges       []positions.TextEdit
}

var fixManifestAttributeTestCases = map[string]fixManifestAttributeTestCase{
	"Creator element with role and no id should get the proper id and have the proper meta element added": {
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:creator opf:role="aut">Author Name</dc:creator>
</metadata>`,
		attribute:             "opf:role",
		line:                  2,
		attributeNameToNumber: map[string]int{},
		expectedChanges: []positions.TextEdit{
			{ // Insert id="creator1"
				Range: positions.Range{
					Start: positions.Position{Line: 2, Column: 31},
					End:   positions.Position{Line: 2, Column: 31},
				},
				NewText: ` id="creator1"`,
			},
			{ // Remove opf:role="aut"
				Range: positions.Range{
					Start: positions.Position{Line: 2, Column: 16},
					End:   positions.Position{Line: 2, Column: 31},
				},
				NewText: "",
			},
			{ // Insert meta tag
				Range: positions.Range{
					Start: positions.Position{Line: 3, Column: 1},
					End:   positions.Position{Line: 3, Column: 1},
				},
				NewText: `<meta refines="#creator1" property="role">aut</meta>
    `,
			},
		},
	},
	"Creator element with role and an id should have the proper meta element added referencing the existing id": {
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
            <dc:creator id="creator-existing" opf:role="aut">Author Name</dc:creator>
</metadata>`,
		attribute:             "opf:role",
		line:                  2,
		attributeNameToNumber: map[string]int{},
		expectedChanges: []positions.TextEdit{
			{ // Remove opf:role="aut"
				Range: positions.Range{
					Start: positions.Position{Line: 2, Column: 46},
					End:   positions.Position{Line: 2, Column: 61},
				},
				NewText: "",
			},
			{ // Insert meta tag
				Range: positions.Range{
					Start: positions.Position{Line: 3, Column: 1},
					End:   positions.Position{Line: 3, Column: 1},
				},
				NewText: `<meta refines="#creator-existing" property="role">aut</meta>
            `,
			},
		},
	},

	"Contributor element with file-as and no id should get the proper id and have the proper meta element added when a contributor has already been handled so far": {
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
            <dc:contributor opf:file-as="Contributor Name">Contributor Name</dc:contributor>
</metadata>`,
		attribute: "opf:file-as",
		line:      2,
		attributeNameToNumber: map[string]int{
			"contributor": 2,
		},
		expectedChanges: []positions.TextEdit{
			{ // Insert id="contributor2"
				Range: positions.Range{
					Start: positions.Position{Line: 2, Column: 59},
					End:   positions.Position{Line: 2, Column: 59},
				},
				NewText: ` id="contributor2"`,
			},
			{ // Remove opf:file-as="Contributor Name"
				Range: positions.Range{
					Start: positions.Position{Line: 2, Column: 28},
					End:   positions.Position{Line: 2, Column: 59},
				},
				NewText: "",
			},
			{ // Insert meta tag
				Range: positions.Range{
					Start: positions.Position{Line: 3, Column: 1},
					End:   positions.Position{Line: 3, Column: 1},
				},
				NewText: `<meta refines="#contributor2" property="file-as">Contributor Name</meta>
            `,
			},
		},
	},

	"Contributor element with file-as and an id should have the proper meta element added referencing the proper id when a contributor has already been handled so far": {
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
            <dc:contributor id="contributor-existing" opf:file-as="Contributor Name">Contributor Name</dc:contributor>
</metadata>`,
		attribute: "opf:file-as",
		line:      2,
		attributeNameToNumber: map[string]int{
			"contributor": 2,
		},
		expectedChanges: []positions.TextEdit{
			{ // Remove opf:file-as="Contributor Name"
				Range: positions.Range{
					Start: positions.Position{Line: 2, Column: 54},
					End:   positions.Position{Line: 2, Column: 85},
				},
				NewText: "",
			},
			{ // Insert meta tag
				Range: positions.Range{
					Start: positions.Position{Line: 3, Column: 1},
					End:   positions.Position{Line: 3, Column: 1},
				},
				NewText: `<meta refines="#contributor-existing" property="file-as">Contributor Name</meta>
            `,
			},
		},
	},
}

func TestFixManifestAttribute(t *testing.T) {
	for name, args := range fixManifestAttributeTestCases {
		t.Run(name, func(t *testing.T) {
			actual, err := rulefixes.FixManifestAttribute(args.opfContents, args.attribute, args.line, args.attributeNameToNumber)

			assert.Nil(t, err)
			assert.Equal(t, args.expectedChanges, actual)
		})
	}
}
