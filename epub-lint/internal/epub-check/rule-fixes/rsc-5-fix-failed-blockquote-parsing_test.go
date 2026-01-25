//go:build unit

package rulefixes_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
	"github.com/stretchr/testify/assert"
)

type fixFailedBlockquoteParsingTestCase struct {
	input           string
	line, column    int
	expectedChanges []positions.TextEdit
}

var fixFailedBlockquoteParsingTestCases = map[string]fixFailedBlockquoteParsingTestCase{
	"A blockquote with another blockquote inside of it with text in that blockquote does not make any changes": {
		input:           `<html><body><blockquote><blockquote>text</blockquote></blockquote></body></html>`,
		line:            1,
		column:          67,
		expectedChanges: nil,
	},
	"A blockquote with an `img` tag in it does get the paragraph tags inserted": {
		input:  `<blockquote><img src="foo" /></blockquote>`,
		line:   1,
		column: 43,
		expectedChanges: []positions.TextEdit{
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   1,
						Column: 13,
					},
					End: positions.Position{
						Line:   1,
						Column: 13,
					},
				},
				NewText: "<p>",
			},
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   1,
						Column: 30,
					},
					End: positions.Position{
						Line:   1,
						Column: 30,
					},
				},
				NewText: "</p>",
			},
		},
	},
	"A blockquote that ends in a `</span> ` gets the paragraph tags inserted": {
		input:  `<blockquote><span>foo</span> </blockquote>`,
		line:   1,
		column: 50,
		expectedChanges: []positions.TextEdit{
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   1,
						Column: 13,
					},
					End: positions.Position{
						Line:   1,
						Column: 13,
					},
				},
				NewText: "<p>",
			},
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   1,
						Column: 30,
					},
					End: positions.Position{
						Line:   1,
						Column: 30,
					},
				},
				NewText: "</p>",
			},
		},
	},
	"A blockquote that has no html tags in it and just text and whitespace present does have the paragraph tags inserted": {
		input:  `<blockquote>   some text    </blockquote>`,
		line:   1,
		column: 42,
		expectedChanges: []positions.TextEdit{
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   1,
						Column: 13,
					},
					End: positions.Position{
						Line:   1,
						Column: 13,
					},
				},
				NewText: "<p>",
			},
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   1,
						Column: 29,
					},
					End: positions.Position{
						Line:   1,
						Column: 29,
					},
				},
				NewText: "</p>",
			},
		},
	},
	"A blockquote with a paragraph tag that starts and ends the content of the blockquote does not get paragraph tags inserted": {
		input:           `<blockquote><p>content</p></blockquote>`,
		line:            1,
		column:          40,
		expectedChanges: nil,
	},
	"A blockquote with an `img` tag in it on line 2 does get the paragraph tags inserted": {
		input: `<html><body>
<blockquote><img src="foo" /></blockquote>
</body></html>`,
		line:   2,
		column: 43,
		expectedChanges: []positions.TextEdit{
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   2,
						Column: 13,
					},
					End: positions.Position{
						Line:   2,
						Column: 13,
					},
				},
				NewText: "<p>",
			},
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   2,
						Column: 30,
					},
					End: positions.Position{
						Line:   2,
						Column: 30,
					},
				},
				NewText: "</p>",
			},
		},
	},
	"A blockquote that has no html tags in it and just text and whitespace present on line 3 does have the paragraph tags inserted": {
		input: `<html>
<body>
<blockquote>   some text    </blockquote>
</body></html>`,
		line:   3,
		column: 42,
		expectedChanges: []positions.TextEdit{
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   3,
						Column: 13,
					},
					End: positions.Position{
						Line:   3,
						Column: 13,
					},
				},
				NewText: "<p>",
			},
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   3,
						Column: 29,
					},
					End: positions.Position{
						Line:   3,
						Column: 29,
					},
				},
				NewText: "</p>",
			},
		},
	},
	"A blockquote with a paragraph tag that starts and ends the content of the blockquote on line 2 does not get paragraph tags inserted": {
		input: `<html>
<blockquote><p>content</p></blockquote>
</html>`,
		line:            2,
		column:          40,
		expectedChanges: nil,
	},
}

func TestFixFailedBlockquoteParsing(t *testing.T) {
	for name, args := range fixFailedBlockquoteParsingTestCases {
		t.Run(name, func(t *testing.T) {
			actual := rulefixes.FixFailedBlockquoteParsing(args.line, args.column, args.input)

			assert.Equal(t, args.expectedChanges, actual)
		})
	}
}
