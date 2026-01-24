//go:build unit

package rulefixes_test

import (
	"testing"

	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
	"github.com/stretchr/testify/assert"
)

// fixXmlIdValueTestCase structure to hold the OPF/NCX content, line number, and attribute to update
type fixXmlIdValueTestCase struct {
	inputText      string
	lineNumber     int
	attribute      string
	expectedChange rulefixes.TextEdit
}

var fixXmlIdValueTestCases = map[string]fixXmlIdValueTestCase{
	"EPUB 2 OPF with an invalid starting character in the id should be replaced with an underscore": {
		inputText: `<metadata>
								<dc:identifier id="!invalidStartChar">urn:isbn:9781234567890</dc:identifier>
						</metadata>`,
		lineNumber: 2,
		attribute:  "id",
		expectedChange: rulefixes.TextEdit{
			Range: rulefixes.Range{
				Start: rulefixes.Position{
					Line:   2,
					Column: 28,
				},
				End: rulefixes.Position{
					Line:   2,
					Column: 45,
				},
			},
			NewText: "_invalidStartChar",
		},
	},
	"EPUB 2 OPF with a number starting the id should have an underscore added at the start": {
		inputText: `<metadata>
								<dc:identifier id="123numberStart">urn:isbn:9781234567890</dc:identifier>
						</metadata>`,
		lineNumber: 2,
		attribute:  "id",
		expectedChange: rulefixes.TextEdit{
			Range: rulefixes.Range{
				Start: rulefixes.Position{
					Line:   2,
					Column: 28,
				},
				End: rulefixes.Position{
					Line:   2,
					Column: 42,
				},
			},
			NewText: "_123numberStart",
		},
	},
	"EPUB 2 OPF with invalid characters in the id should have them replaced with an underscore": {
		inputText: `<metadata>
								<dc:identifier id="invalid!char#id">urn:isbn:9781234567890</dc:identifier>
						</metadata>`,
		lineNumber: 2,
		attribute:  "id",
		expectedChange: rulefixes.TextEdit{
			Range: rulefixes.Range{
				Start: rulefixes.Position{
					Line:   2,
					Column: 28,
				},
				End: rulefixes.Position{
					Line:   2,
					Column: 43,
				},
			},
			NewText: "invalid_char_id",
		},
	},
	"EPUB 2 OPF with colons in the id value should have them replaced with underscores": {
		inputText: `<metadata>
								<dc:identifier id="id:with:colon">urn:isbn:9781234567890</dc:identifier>
						</metadata>`,
		lineNumber: 2,
		attribute:  "id",
		expectedChange: rulefixes.TextEdit{
			Range: rulefixes.Range{
				Start: rulefixes.Position{
					Line:   2,
					Column: 28,
				},
				End: rulefixes.Position{
					Line:   2,
					Column: 41,
				},
			},
			NewText: "id_with_colon",
		},
	},
	"EPUB 3 OPF with an invalid start character in the id should have it replaced with an underscore": {
		inputText: `<manifest>
								<item id="!invalidStartChar" href="chapter1.xhtml" media-type="application/xhtml+xml"/>
						</manifest>`,
		lineNumber: 2,
		attribute:  "id",
		expectedChange: rulefixes.TextEdit{
			Range: rulefixes.Range{
				Start: rulefixes.Position{
					Line:   2,
					Column: 19,
				},
				End: rulefixes.Position{
					Line:   2,
					Column: 36,
				},
			},
			NewText: "_invalidStartChar",
		},
	},
	"EPUB 3 OPF with an id starting with a number should have an underscore added at the start": {
		inputText: `<manifest>
								<item id="123numberStart" href="chapter1.xhtml" media-type="application/xhtml+xml"/>
						</manifest>`,
		lineNumber: 2,
		attribute:  "id",
		expectedChange: rulefixes.TextEdit{
			Range: rulefixes.Range{
				Start: rulefixes.Position{
					Line:   2,
					Column: 19,
				},
				End: rulefixes.Position{
					Line:   2,
					Column: 33,
				},
			},
			NewText: "_123numberStart",
		},
	},
	"EPUB 3 OPF with an id with invalid characters should be replaced with underscores": {
		inputText: `<manifest>
								<item id="invalid!char#id" href="chapter1.xhtml" media-type="application/xhtml+xml"/>
						</manifest>`,
		lineNumber: 2,
		attribute:  "id",
		expectedChange: rulefixes.TextEdit{
			Range: rulefixes.Range{
				Start: rulefixes.Position{
					Line:   2,
					Column: 19,
				},
				End: rulefixes.Position{
					Line:   2,
					Column: 34,
				},
			},
			NewText: "invalid_char_id",
		},
	},
	"EPUB 3 OPF with an id with colons in the value should replace them with underscores": {
		inputText: `<manifest>
								<item id="id:with:colon" href="chapter1.xhtml" media-type="application/xhtml+xml"/>
						</manifest>`,
		lineNumber: 2,
		attribute:  "id",
		expectedChange: rulefixes.TextEdit{
			Range: rulefixes.Range{
				Start: rulefixes.Position{
					Line:   2,
					Column: 19,
				},
				End: rulefixes.Position{
					Line:   2,
					Column: 32,
				},
			},
			NewText: "id_with_colon",
		},
	},
	"EPUB 3 OPF with an idref with an invalid starting character should replace it with an underscore": {
		inputText: `<spine>
								<itemref idref="!invalidStartChar"/>
						</spine>`,
		lineNumber: 2,
		attribute:  "idref",
		expectedChange: rulefixes.TextEdit{
			Range: rulefixes.Range{
				Start: rulefixes.Position{
					Line:   2,
					Column: 25,
				},
				End: rulefixes.Position{
					Line:   2,
					Column: 42,
				},
			},
			NewText: "_invalidStartChar",
		},
	},
	"EPUB 3 OPF with a number starting the idref should have an underscore added before it": {
		inputText: `<spine>
								<itemref idref="123numberStart"/>
						</spine>`,
		lineNumber: 2,
		attribute:  "idref",
		expectedChange: rulefixes.TextEdit{
			Range: rulefixes.Range{
				Start: rulefixes.Position{
					Line:   2,
					Column: 25,
				},
				End: rulefixes.Position{
					Line:   2,
					Column: 39,
				},
			},
			NewText: "_123numberStart",
		},
	},
	"EPUB 3 OPF with invalid characters in the idref should be replaced with an underscore": {
		inputText: `<spine>
								<itemref idref="invalid!char#id"/>
						</spine>`,
		lineNumber: 2,
		attribute:  "idref",
		expectedChange: rulefixes.TextEdit{
			Range: rulefixes.Range{
				Start: rulefixes.Position{
					Line:   2,
					Column: 25,
				},
				End: rulefixes.Position{
					Line:   2,
					Column: 40,
				},
			},
			NewText: "invalid_char_id",
		},
	},
	"EPUB 3 OPF with colons in the value of an idref should be replaced with an underscore": {
		inputText: `<spine>
								<itemref idref="id:with:colon"/>
						</spine>`,
		lineNumber: 2,
		attribute:  "idref",
		expectedChange: rulefixes.TextEdit{
			Range: rulefixes.Range{
				Start: rulefixes.Position{
					Line:   2,
					Column: 25,
				},
				End: rulefixes.Position{
					Line:   2,
					Column: 38,
				},
			},
			NewText: "id_with_colon",
		},
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
		expectedChange: rulefixes.TextEdit{
			Range: rulefixes.Range{
				Start: rulefixes.Position{
					Line:   2,
					Column: 23,
				},
				End: rulefixes.Position{
					Line:   2,
					Column: 40,
				},
			},
			NewText: "_invalidStartChar",
		},
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
		expectedChange: rulefixes.TextEdit{
			Range: rulefixes.Range{
				Start: rulefixes.Position{
					Line:   2,
					Column: 23,
				},
				End: rulefixes.Position{
					Line:   2,
					Column: 37,
				},
			},
			NewText: "_123numberStart",
		},
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
		expectedChange: rulefixes.TextEdit{
			Range: rulefixes.Range{
				Start: rulefixes.Position{
					Line:   2,
					Column: 23,
				},
				End: rulefixes.Position{
					Line:   2,
					Column: 38,
				},
			},
			NewText: "invalid_char_id",
		},
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
		expectedChange: rulefixes.TextEdit{
			Range: rulefixes.Range{
				Start: rulefixes.Position{
					Line:   2,
					Column: 23,
				},
				End: rulefixes.Position{
					Line:   2,
					Column: 36,
				},
			},
			NewText: "id_with_colon",
		},
	},
}

func TestFixXmlIdValue(t *testing.T) {
	for name, args := range fixXmlIdValueTestCases {
		t.Run(name, func(t *testing.T) {
			actual := rulefixes.FixXmlIdValue(args.inputText, args.lineNumber, args.attribute)

			assert.Equal(t, args.expectedChange, actual)
		})
	}
}
