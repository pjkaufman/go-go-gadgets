//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

type fixMissingUniqueIdentifierIdTestCase struct {
	opfContents    string
	id             string
	expectedOutput string
}

var fixMissingUniqueIdentifierIdTestCases = map[string]fixMissingUniqueIdentifierIdTestCase{
	"When the first identifier does not have an id, it should have one added": {
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:identifier>12345</dc:identifier>
    <dc:identifier id="existing-id">67890</dc:identifier>
</metadata>`,
		id: "unique-id",
		expectedOutput: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:identifier id="unique-id">12345</dc:identifier>
    <dc:identifier id="existing-id">67890</dc:identifier>
</metadata>`,
	},
	"When the first identifier does have an id and there is a second identifier without it, then the second one should have the id added": {
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:identifier id="existing-id-1">12345</dc:identifier>
    <dc:identifier>67890</dc:identifier>
</metadata>`,
		id: "unique-id",
		expectedOutput: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:identifier id="existing-id-1">12345</dc:identifier>
    <dc:identifier id="unique-id">67890</dc:identifier>
</metadata>`,
	},
	"When there is no identifier element, then no id will be added": {
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:title>Example Book</dc:title>
    <dc:creator>Author Name</dc:creator>
</metadata>`,
		id: "unique-id",
		expectedOutput: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:title>Example Book</dc:title>
    <dc:creator>Author Name</dc:creator>
</metadata>`,
	},
	"When all identifier elements have ids, no id will be added": {
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:identifier id="existing-id">12345</dc:identifier>
</metadata>`,
		id: "unique-id",
		expectedOutput: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:identifier id="existing-id">12345</dc:identifier>
</metadata>`,
	},
}

func TestFixMissingUniqueIdentifierId(t *testing.T) {
	for name, args := range fixMissingUniqueIdentifierIdTestCases {
		t.Run(name, func(t *testing.T) {
			actual, err := linter.FixMissingUniqueIdentifierId(args.opfContents, args.id)

			assert.Nil(t, err)
			assert.Equal(t, args.expectedOutput, actual)
		})
	}
}
