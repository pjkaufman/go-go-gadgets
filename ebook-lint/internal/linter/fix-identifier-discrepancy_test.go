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
  <dc:identifier id="pub-id" opf:scheme="UUID">12345</dc:identifier>
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
  <dc:identifier id="pub-id" opf:scheme="ISBN">9781975392543</dc:identifier>
</metadata>
  <manifest></manifest>
  <spine></spine>
</package>`,
		},
		// 		{
		// 			name: "When the OPF and NCX have two different unique identifier values and the ",
		// 			opfContents: `
		// <package xmlns="http://www.idpf.org/2007/opf" unique-identifier="pub-id">
		// 	<metadata>
		// 		<dc:title>Example Book</dc:title>
		// 		<dc:identifier id="pub-id" opf:scheme="UUID">67890</dc:identifier>
		// 	</metadata>
		// 	<manifest></manifest>
		// 	<spine></spine>
		// </package>`,
		// 			ncxContents: `
		// <ncx xmlns="http://www.daisy.org/z3986/2005/ncx/">
		// 	<head>
		// 		<meta name="dtb:uid" content="12345" />
		// 	</head>
		// </ncx>`,
		// 			expectedOutput: `
		// <package xmlns="http://www.idpf.org/2007/opf" unique-identifier="pub-id">
		// 	<metadata>
		// 		<dc:title>Example Book</dc:title>
		// 		<dc:identifier id="pub-id" opf:scheme="UUID">12345</dc:identifier>
		// 	</metadata>
		// 	<manifest></manifest>
		// 	<spine></spine>
		// </package>`,
		// 		},
		// 		{
		// 			name: "Different unique identifier in OPF and NCX where the OPF has the identifier from the NCX, but it is not the identifier specified in the OPF",
		// 			opfContents: `
		// <package xmlns="http://www.idpf.org/2007/opf" unique-identifier="MainId">
		//   <metadata>
		//     <dc:title>Example Book</dc:title>
		//     <dc:identifier id="MainId" opf:scheme="UUID">67890</dc:identifier>
		//     <dc:identifier id="pub-id" opf:scheme="UUID">12345</dc:identifier>
		//   </metadata>
		//   <manifest></manifest>
		//   <spine></spine>
		// </package>`,
		// 			ncxContents: `
		// <ncx xmlns="http://www.daisy.org/z3986/2005/ncx/">
		//   <head>
		//     <meta name="dtb:uid" content="12345" />
		//   </head>
		// </ncx>`,
		// 			expectedOutput: `
		// <package xmlns="http://www.idpf.org/2007/opf" unique-identifier="pub-id">
		//   <metadata>
		//     <dc:title>Example Book</dc:title>
		//     <dc:identifier id="MainId" opf:scheme="UUID">67890</dc:identifier>
		//     <dc:identifier id="pub-id" opf:scheme="UUID">12345</dc:identifier>
		//   </metadata>
		//   <manifest></manifest>
		//   <spine></spine>
		// </package>`,
		// 		},
	}

	for _, args := range testCases {
		t.Run(args.name, func(t *testing.T) {
			actual, err := linter.FixIdentifierDiscrepancy(args.opfContents, args.ncxContents)

			assert.Nil(t, err)
			assert.Equal(t, args.expectedOutput, actual)
		})
	}
}
