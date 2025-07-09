//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

// fixXmlIdValueTestCase structure to hold the OPF/NCX content, line number, and attribute to update
type fixXmlIdValueTestCase struct {
	inputText      string
	lineNumber     int
	attribute      string
	expectedOutput string
}

var fixXmlIdValueTestCases = map[string]fixXmlIdValueTestCase{
	"EPUB 2 OPF with an invalid starting character in the id should be replaced with an underscore": {
		inputText: `<metadata>
								<dc:identifier id="!invalidStartChar">urn:isbn:9781234567890</dc:identifier>
						</metadata>`,
		lineNumber: 2,
		attribute:  "id",
		expectedOutput: `<metadata>
								<dc:identifier id="_invalidStartChar">urn:isbn:9781234567890</dc:identifier>
						</metadata>`,
	},
	"EPUB 2 OPF with a number starting the id should have an underscore added at the start": {
		inputText: `<metadata>
								<dc:identifier id="123numberStart">urn:isbn:9781234567890</dc:identifier>
						</metadata>`,
		lineNumber: 2,
		attribute:  "id",
		expectedOutput: `<metadata>
								<dc:identifier id="_123numberStart">urn:isbn:9781234567890</dc:identifier>
						</metadata>`,
	},
	"EPUB 2 OPF with invalid characters in the id should have them replaced with an underscore": {
		inputText: `<metadata>
								<dc:identifier id="invalid!char#id">urn:isbn:9781234567890</dc:identifier>
						</metadata>`,
		lineNumber: 2,
		attribute:  "id",
		expectedOutput: `<metadata>
								<dc:identifier id="invalid_char_id">urn:isbn:9781234567890</dc:identifier>
						</metadata>`,
	},
	"EPUB 2 OPF with colons in the id value should have them replaced with underscores": {
		inputText: `<metadata>
								<dc:identifier id="id:with:colon">urn:isbn:9781234567890</dc:identifier>
						</metadata>`,
		lineNumber: 2,
		attribute:  "id",
		expectedOutput: `<metadata>
								<dc:identifier id="id_with_colon">urn:isbn:9781234567890</dc:identifier>
						</metadata>`,
	},
	"EPUB 3 OPF with an invalid start character in the id should have it replaced with an underscore": {
		inputText: `<manifest>
								<item id="!invalidStartChar" href="chapter1.xhtml" media-type="application/xhtml+xml"/>
						</manifest>`,
		lineNumber: 2,
		attribute:  "id",
		expectedOutput: `<manifest>
								<item id="_invalidStartChar" href="chapter1.xhtml" media-type="application/xhtml+xml"/>
						</manifest>`,
	},
	"EPUB 3 OPF with an id starting with a number should have an underscore added at the start": {
		inputText: `<manifest>
								<item id="123numberStart" href="chapter1.xhtml" media-type="application/xhtml+xml"/>
						</manifest>`,
		lineNumber: 2,
		attribute:  "id",
		expectedOutput: `<manifest>
								<item id="_123numberStart" href="chapter1.xhtml" media-type="application/xhtml+xml"/>
						</manifest>`,
	},
	"EPUB 3 OPF with an id with invalid characters should be replaced with underscores": {
		inputText: `<manifest>
								<item id="invalid!char#id" href="chapter1.xhtml" media-type="application/xhtml+xml"/>
						</manifest>`,
		lineNumber: 2,
		attribute:  "id",
		expectedOutput: `<manifest>
								<item id="invalid_char_id" href="chapter1.xhtml" media-type="application/xhtml+xml"/>
						</manifest>`,
	},
	"EPUB 3 OPF with an id with colons in the value should replace them with underscores": {
		inputText: `<manifest>
								<item id="id:with:colon" href="chapter1.xhtml" media-type="application/xhtml+xml"/>
						</manifest>`,
		lineNumber: 2,
		attribute:  "id",
		expectedOutput: `<manifest>
								<item id="id_with_colon" href="chapter1.xhtml" media-type="application/xhtml+xml"/>
						</manifest>`,
	},
	"EPUB 3 OPF with an idref with an invalid starting character should replace it with an underscore": {
		inputText: `<spine>
								<itemref idref="!invalidStartChar"/>
						</spine>`,
		lineNumber: 2,
		attribute:  "idref",
		expectedOutput: `<spine>
								<itemref idref="_invalidStartChar"/>
						</spine>`,
	},
	"EPUB 3 OPF with a number starting the idref should have an underscore added before it": {
		inputText: `<spine>
								<itemref idref="123numberStart"/>
						</spine>`,
		lineNumber: 2,
		attribute:  "idref",
		expectedOutput: `<spine>
								<itemref idref="_123numberStart"/>
						</spine>`,
	},
	"EPUB 3 OPF with invalid characters in the idref should be replaced with an underscore": {
		inputText: `<spine>
								<itemref idref="invalid!char#id"/>
						</spine>`,
		lineNumber: 2,
		attribute:  "idref",
		expectedOutput: `<spine>
								<itemref idref="invalid_char_id"/>
						</spine>`,
	},
	"EPUB 3 OPF with colons in the value of an idref should be replaced with an underscore": {
		inputText: `<spine>
								<itemref idref="id:with:colon"/>
						</spine>`,
		lineNumber: 2,
		attribute:  "idref",
		expectedOutput: `<spine>
								<itemref idref="id_with_colon"/>
						</spine>`,
	},
	"NCX an invalid start character should be replaced with an underscore": {
		inputText: `<navMap>
								<navPoint id="!invalidStartChar" class="chapter" playOrder="1">
										<navLabel>
												<text>Introduction</text>
										</navLabel>
										<content src="chapter1.xhtml"/>
								</navPoint>
						</navMap>`,
		lineNumber: 2,
		attribute:  "id",
		expectedOutput: `<navMap>
								<navPoint id="_invalidStartChar" class="chapter" playOrder="1">
										<navLabel>
												<text>Introduction</text>
										</navLabel>
										<content src="chapter1.xhtml"/>
								</navPoint>
						</navMap>`,
	},
	"NCX a number starting an id should get an underscore added before it": {
		inputText: `<navMap>
								<navPoint id="123numberStart" class="chapter" playOrder="1">
										<navLabel>
												<text>Introduction</text>
										</navLabel>
										<content src="chapter1.xhtml"/>
								</navPoint>
						</navMap>`,
		lineNumber: 2,
		attribute:  "id",
		expectedOutput: `<navMap>
								<navPoint id="_123numberStart" class="chapter" playOrder="1">
										<navLabel>
												<text>Introduction</text>
										</navLabel>
										<content src="chapter1.xhtml"/>
								</navPoint>
						</navMap>`,
	},
	"NCX invalid characters should be replaced with an underscore": {
		inputText: `<navMap>
								<navPoint id="invalid!char#id" class="chapter" playOrder="1">
										<navLabel>
												<text>Introduction</text>
										</navLabel>
										<content src="chapter1.xhtml"/>
								</navPoint>
						</navMap>`,
		lineNumber: 2,
		attribute:  "id",
		expectedOutput: `<navMap>
								<navPoint id="invalid_char_id" class="chapter" playOrder="1">
										<navLabel>
												<text>Introduction</text>
										</navLabel>
										<content src="chapter1.xhtml"/>
								</navPoint>
						</navMap>`,
	},
	"NCX colon in value should get replaced with underscore": {
		inputText: `<navMap>
								<navPoint id="id:with:colon" class="chapter" playOrder="1">
										<navLabel>
												<text>Introduction</text>
										</navLabel>
										<content src="chapter1.xhtml"/>
								</navPoint>
						</navMap>`,
		lineNumber: 2,
		attribute:  "id",
		expectedOutput: `<navMap>
								<navPoint id="id_with_colon" class="chapter" playOrder="1">
										<navLabel>
												<text>Introduction</text>
										</navLabel>
										<content src="chapter1.xhtml"/>
								</navPoint>
						</navMap>`,
	},
}

func TestFixXmlIdValue(t *testing.T) {
	for name, args := range fixXmlIdValueTestCases {
		t.Run(name, func(t *testing.T) {
			actual := linter.FixXmlIdValue(args.inputText, args.lineNumber, args.attribute)

			assert.Equal(t, args.expectedOutput, actual)
		})
	}
}
