//go:build unit

package rulefixes_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
	"github.com/stretchr/testify/assert"
)

type fixMissingUniqueIdentifierIdTestCase struct {
	opfContents    string
	id             string
	expectedChange positions.TextEdit
}

var fixMissingUniqueIdentifierIdTestCases = map[string]fixMissingUniqueIdentifierIdTestCase{
	"When the first identifier does not have an id, it should have one added": {
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:identifier>12345</dc:identifier>
    <dc:identifier id="existing-id">67890</dc:identifier>
</metadata>`,
		id: "unique-id",
		expectedChange: positions.TextEdit{
			Range: positions.Range{
				Start: positions.Position{
					Line:   2,
					Column: 10,
				},
				End: positions.Position{
					Line:   2,
					Column: 10,
				},
			},
			NewText: ` id="unique-id"`,
		},
	},
	"When the first identifier does have an id and there is a second identifier without it, then the second one should have the id added": {
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:identifier id="existing-id-1">12345</dc:identifier>
    <dc:identifier>67890</dc:identifier>
</metadata>`,
		id: "unique-id",
		expectedChange: positions.TextEdit{
			Range: positions.Range{
				Start: positions.Position{
					Line:   3,
					Column: 10,
				},
				End: positions.Position{
					Line:   3,
					Column: 10,
				},
			},
			NewText: ` id="unique-id"`,
		},
	},
	"When there is no identifier element, then no id will be added": {
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:title>Example Book</dc:title>
    <dc:creator>Author Name</dc:creator>
</metadata>`,
		id:             "unique-id",
		expectedChange: positions.TextEdit{},
	},
	"When all identifier elements have ids, no id will be added": {
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:identifier id="existing-id">12345</dc:identifier>
</metadata>`,
		id:             "unique-id",
		expectedChange: positions.TextEdit{},
	},
}

func TestFixMissingUniqueIdentifierId(t *testing.T) {
	for name, args := range fixMissingUniqueIdentifierIdTestCases {
		t.Run(name, func(t *testing.T) {
			actual, err := rulefixes.FixMissingUniqueIdentifierId(args.opfContents, args.id)

			assert.Nil(t, err)
			assert.Equal(t, args.expectedChange, actual)
		})
	}
}
