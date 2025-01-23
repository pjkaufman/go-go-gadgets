//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

// Struct for test cases
type identifierTestCase struct {
	name           string
	opfContents    string
	ncxContents    string
	expectedOutput string
}

func TestFixIdentifiers(t *testing.T) {
	testCases := []identifierTestCase{
		{
			name: "When no unique identifier is in the OPF, but it is present in the NCX, the unique identifier should be added",
			opfContents: `
<package xmlns="http://www.idpf.org/2007/opf" unique-identifier="pub-id">
  <metadata>
    <dc:title>Example Book</dc:title>
  </metadata>
  <manifest></manifest>
  <spine></spine>
</package>`,
			ncxContents: `
<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/">
  <head>
    <meta name="dtb:uid" content="12345" />
  </head>
</ncx>`,
			expectedOutput: `
<package xmlns="http://www.idpf.org/2007/opf" unique-identifier="pub-id">
  <metadata>
    <dc:title>Example Book</dc:title>
  <dc:identifier id="pub-id">12345</dc:identifier>
</metadata>
  <manifest></manifest>
  <spine></spine>
</package>`,
		},
		{
			name: "When no unique identifier is in the OPF, but it is present in the NCX, the unique identifier should be added",
			opfContents: `
<package xmlns="http://www.idpf.org/2007/opf" unique-identifier="pub-id">
  <metadata>
    <dc:title>Example Book</dc:title>
  </metadata>
  <manifest></manifest>
  <spine></spine>
</package>`,
			ncxContents: `
<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/">
  <head>
    <meta name="dtb:uid" content="9aedca49-923e-4a61-abca-8c1c88d6f868" />
  </head>
</ncx>`,
			expectedOutput: `
<package xmlns="http://www.idpf.org/2007/opf" unique-identifier="pub-id">
  <metadata>
    <dc:title>Example Book</dc:title>
  <dc:identifier id="pub-id">9aedca49-923e-4a61-abca-8c1c88d6f868</dc:identifier>
</metadata>
  <manifest></manifest>
  <spine></spine>
</package>`,
		},
		{
			name: "When no unique identifier is in the OPF, but it is present in the NCX and it is an ISBN, the unique identifier should be added as an ISBN",
			opfContents: `
<package xmlns="http://www.idpf.org/2007/opf" unique-identifier="pub-id">
  <metadata>
    <dc:title>Example Book</dc:title>
  </metadata>
  <manifest></manifest>
  <spine></spine>
</package>`,
			ncxContents: `
<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/">
  <head>
    <meta name="dtb:uid" content="9781975392543" />
  </head>
</ncx>`,
			expectedOutput: `
<package xmlns="http://www.idpf.org/2007/opf" unique-identifier="pub-id">
  <metadata>
    <dc:title>Example Book</dc:title>
  <dc:identifier id="pub-id">9781975392543</dc:identifier>
</metadata>
  <manifest></manifest>
  <spine></spine>
</package>`,
		},
		{
			name: "When the OPF and NCX have two different unique identifier values and the opf should have a unique identifier added and the id should have been moved from the original identifier to the new one",
			opfContents: `
<package xmlns="http://www.idpf.org/2007/opf" unique-identifier="pub-id">
  <metadata>
    <dc:title>Example Book</dc:title>
    <dc:identifier id="pub-id">67890</dc:identifier>
  </metadata>
  <manifest></manifest>
  <spine></spine>
</package>`,
			ncxContents: `
<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/">
  <head>
    <meta name="dtb:uid" content="9aedca49-923e-4a61-abca-8c1c88d6f868" />
  </head>
</ncx>`,
			expectedOutput: `
<package xmlns="http://www.idpf.org/2007/opf" unique-identifier="pub-id">
  <metadata>
    <dc:title>Example Book</dc:title>
    <dc:identifier>67890</dc:identifier>
    <dc:identifier id="pub-id">9aedca49-923e-4a61-abca-8c1c88d6f868</dc:identifier>
  </metadata>
  <manifest></manifest>
  <spine></spine>
</package>`,
		},
		{
			name: "Different unique identifier in OPF and NCX where the OPF has the identifier from the NCX, but it is not the identifier specified in the OPF and there is no id already gets the unique identifier moved to the one that is in the NCX",
			opfContents: `
<package xmlns="http://www.idpf.org/2007/opf" unique-identifier="MainId">
  <metadata>
    <dc:title>Example Book</dc:title>
    <dc:identifier id="MainId">ef932546-7cf7-4ded-a0ea-5a069fbb8abc</dc:identifier>
    <dc:identifier>12345</dc:identifier>
  </metadata>
  <manifest></manifest>
  <spine></spine>
</package>`,
			ncxContents: `
<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/">
  <head>
    <meta name="dtb:uid" content="12345" />
  </head>
</ncx>`,
			expectedOutput: `
<package xmlns="http://www.idpf.org/2007/opf" unique-identifier="MainId">
  <metadata>
    <dc:title>Example Book</dc:title>
    <dc:identifier>ef932546-7cf7-4ded-a0ea-5a069fbb8abc</dc:identifier>
    <dc:identifier id="MainId">12345</dc:identifier>
  </metadata>
  <manifest></manifest>
  <spine></spine>
</package>`,
		},
		{
			name: "Different unique identifier in OPF and NCX where the OPF has the identifier from the NCX, but it is not the identifier specified in the OPF and there is an id already gets the unique identifier as the replacement for the one that is in the NCX",
			opfContents: `
<package xmlns="http://www.idpf.org/2007/opf" unique-identifier="MainId">
  <metadata>
    <dc:title>Example Book</dc:title>
    <dc:identifier id="MainId">ef932546-7cf7-4ded-a0ea-5a069fbb8abc</dc:identifier>
    <dc:identifier id="secondaryId">12345</dc:identifier>
  </metadata>
  <manifest></manifest>
  <spine></spine>
</package>`,
			ncxContents: `
<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/">
  <head>
    <meta name="dtb:uid" content="12345" />
  </head>
</ncx>`,
			expectedOutput: `
<package xmlns="http://www.idpf.org/2007/opf" unique-identifier="MainId">
  <metadata>
    <dc:title>Example Book</dc:title>
    <dc:identifier>ef932546-7cf7-4ded-a0ea-5a069fbb8abc</dc:identifier>
    <dc:identifier id="MainId">12345</dc:identifier>
  </metadata>
  <manifest></manifest>
  <spine></spine>
</package>`,
		},
	}

	for _, args := range testCases {
		t.Run(args.name, func(t *testing.T) {
			actual, err := linter.FixIdentifierDiscrepancy(args.opfContents, args.ncxContents)

			assert.Nil(t, err)
			assert.Equal(t, args.expectedOutput, actual)
		})
	}
}
