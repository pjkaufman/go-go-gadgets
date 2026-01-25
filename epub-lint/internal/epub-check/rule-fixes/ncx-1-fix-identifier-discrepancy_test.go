//go:build unit

package rulefixes_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
	"github.com/stretchr/testify/assert"
)

type identifierTestCase struct {
	opfContents     string
	ncxContents     string
	expectedChanges []positions.TextEdit
}

var fixIdentifierTestCases = map[string]identifierTestCase{
	"When no unique identifier is in the OPF, but it is present in the NCX, the unique identifier should be added as a number": {
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
		expectedChanges: []positions.TextEdit{
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   5,
						Column: 2,
					},
					End: positions.Position{
						Line:   5,
						Column: 2,
					},
				},
				NewText: `	<dc:identifier id="pub-id">12345</dc:identifier>
	`,
			},
		},
	},
	"When no unique identifier is in the OPF, but it is present in the NCX, the unique identifier should be added as a UUID": {
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
		expectedChanges: []positions.TextEdit{
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   5,
						Column: 2,
					},
					End: positions.Position{
						Line:   5,
						Column: 2,
					},
				},
				NewText: `	<dc:identifier id="pub-id">9aedca49-923e-4a61-abca-8c1c88d6f868</dc:identifier>
	`,
			},
		},
	},
	"When no unique identifier is in the OPF, but it is present in the NCX and it is an ISBN, the unique identifier should be added as an ISBN": {
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
		expectedChanges: []positions.TextEdit{
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   5,
						Column: 1,
					},
					End: positions.Position{
						Line:   5,
						Column: 1,
					},
				},
				NewText: `	<dc:identifier id="pub-id">9781975392543</dc:identifier>
`,
			},
		},
	},
	"When the OPF and NCX have two different unique identifier values, the opf should have a unique identifier added and the id should have been moved from the original identifier to the new one": {
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
		expectedChanges: []positions.TextEdit{
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   5,
						Column: 16,
					},
					End: positions.Position{
						Line:   5,
						Column: 28,
					},
				},
			},
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   5,
						Column: 50,
					},
					End: positions.Position{
						Line:   5,
						Column: 50,
					},
				},
				NewText: `
	<dc:identifier id="pub-id">9aedca49-923e-4a61-abca-8c1c88d6f868</dc:identifier>`,
			},
		},
	},
	"Different unique identifier in OPF and NCX where the OPF has the identifier from the NCX, but it is not the identifier specified in the OPF and there is no id already gets the unique identifier moved to the one that is in the NCX": {
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
		expectedChanges: []positions.TextEdit{
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   5,
						Column: 15,
					},
					End: positions.Position{
						Line:   5,
						Column: 27,
					},
				},
			},
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   6,
						Column: 15,
					},
					End: positions.Position{
						Line:   6,
						Column: 15,
					},
				},
				NewText: ` id="MainId"`,
			},
		},
	},
	"Different unique identifier in OPF and NCX where the OPF has the identifier from the NCX, but it is not the identifier specified in the OPF and there is an id already gets the unique identifier as the replacement for the one that is in the NCX": {
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
		expectedChanges: []positions.TextEdit{
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   5,
						Column: 16,
					},
					End: positions.Position{
						Line:   5,
						Column: 28,
					},
				},
			},
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   6,
						Column: 21,
					},
					End: positions.Position{
						Line:   6,
						Column: 32,
					},
				},
				NewText: `MainId`,
			},
		},
	},
	"When the OPF and NCX have two different unique identifier values, the OPF should have a unique identifier added and the id should have been moved from the original identifier to the new one make sure that the ending metadata tag is not put before the added identifier": {
		opfContents: `
<package xmlns="http://www.idpf.org/2007/opf" unique-identifier="MainId">
<metadata>
  <dc:title>Example Book</dc:title>
  <dc:identifier id="MainId">ef932546-7cf7-4ded-a0ea-5a069fbb8abc</dc:identifier></metadata>
<manifest></manifest>
<spine></spine>
</package>`,
		ncxContents: `
<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/">
<head>
  <meta name="dtb:uid" content="12345" />
</head>
</ncx>`,
		expectedChanges: []positions.TextEdit{
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   5,
						Column: 17,
					},
					End: positions.Position{
						Line:   5,
						Column: 29,
					},
				},
			},
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   5,
						Column: 93,
					},
					End: positions.Position{
						Line:   5,
						Column: 93,
					},
				},
				NewText: `
  <dc:identifier id="MainId">12345</dc:identifier>`,
			},
		},
	},
	"When the OPF file has no unique identifier set, but the NCX id is present as an identifier, set the id for the identifier instead of adding a new identifier": {
		opfContents: `
<package xmlns="http://www.idpf.org/2007/opf" unique-identifier="MainId">
<metadata>
  <dc:title>Example Book</dc:title>
  <dc:identifier>ef932546-7cf7-4ded-a0ea-5a069fbb8abc</dc:identifier>
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

		expectedChanges: []positions.TextEdit{
			{
				Range: positions.Range{
					Start: positions.Position{Line: 6, Column: 17},
					End:   positions.Position{Line: 6, Column: 17},
				},
				NewText: ` id="MainId"`,
			},
		},
	},
}

func TestFixIdentifiers(t *testing.T) {
	for name, args := range fixIdentifierTestCases {
		t.Run(name, func(t *testing.T) {
			actual, err := rulefixes.FixIdentifierDiscrepancy(args.opfContents, args.ncxContents)

			assert.Nil(t, err)
			assert.Equal(t, args.expectedChanges, actual)
		})
	}
}
