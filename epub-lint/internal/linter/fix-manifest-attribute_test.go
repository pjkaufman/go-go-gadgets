//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

type fixManifestAttributeTestCase struct {
	opfContents           string
	attribute             string
	line                  int
	attributeNameToNumber map[string]int
	expectedOutput        string
}

var fixManifestAttributeTestCases = map[string]fixManifestAttributeTestCase{
	"Creator element with role and no id should get the proper id and have the proper meta element added": {
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:creator opf:role="aut">Author Name</dc:creator>
</metadata>`,
		attribute:             "opf:role",
		line:                  1,
		attributeNameToNumber: map[string]int{},
		expectedOutput: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:creator id="creator1">Author Name</dc:creator>
    <meta refines="#creator1" property="role">aut</meta>
</metadata>`,
	},
	"Creator element with role and an id should have the proper meta element added referencing the existing id": {
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
			<dc:creator id="creator-existing" opf:role="aut">Author Name</dc:creator>
</metadata>`,

		attribute:             "opf:role",
		line:                  1,
		attributeNameToNumber: map[string]int{},
		expectedOutput: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
			<dc:creator id="creator-existing">Author Name</dc:creator>
			<meta refines="#creator-existing" property="role">aut</meta>
</metadata>`,
	},
	"Contributor element with file-as and no id should get the proper id and have the proper meta element added when a contributor has already been handled so far": {
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
			<dc:contributor opf:file-as="Contributor Name">Contributor Name</dc:contributor>
</metadata>`,

		attribute: "opf:file-as",
		line:      1,
		attributeNameToNumber: map[string]int{
			"contributor": 2,
		},
		expectedOutput: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
			<dc:contributor id="contributor2">Contributor Name</dc:contributor>
			<meta refines="#contributor2" property="file-as">Contributor Name</meta>
</metadata>`,
	},
	"Contributor element with file-as and an id should have the proper meta element added referencing the proper id when a contributor has already been handled so far": {
		opfContents: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
			<dc:contributor id="contributor-existing" opf:file-as="Contributor Name">Contributor Name</dc:contributor>
</metadata>`,
		attribute: "opf:file-as",
		line:      1,
		attributeNameToNumber: map[string]int{
			"contributor": 2,
		},
		expectedOutput: `<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
			<dc:contributor id="contributor-existing">Contributor Name</dc:contributor>
			<meta refines="#contributor-existing" property="file-as">Contributor Name</meta>
</metadata>`,
	},
}

func TestFixManifestAttribute(t *testing.T) {
	for name, args := range fixManifestAttributeTestCases {
		t.Run(name, func(t *testing.T) {
			actual, err := linter.FixManifestAttribute(args.opfContents, args.attribute, args.line, args.attributeNameToNumber)

			assert.Nil(t, err)
			assert.Equal(t, args.expectedOutput, actual)
		})
	}
}
