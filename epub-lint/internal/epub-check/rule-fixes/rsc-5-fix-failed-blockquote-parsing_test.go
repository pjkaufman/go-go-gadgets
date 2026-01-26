//go:build unit

package rulefixes_test

import (
	"testing"

	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
)

type fixFailedBlockquoteParsingTestCase struct {
	input        string
	line, column int
	expected     string
}

var fixFailedBlockquoteParsingTestCases = map[string]fixFailedBlockquoteParsingTestCase{
	"A blockquote with another blockquote inside of it with text in that blockquote does not make any changes": {
		input:    `<html><body><blockquote><blockquote>text</blockquote></blockquote></body></html>`,
		line:     1,
		column:   67,
		expected: `<html><body><blockquote><blockquote>text</blockquote></blockquote></body></html>`,
	},
	"A blockquote with an `img` tag in it does get the paragraph tags inserted": {
		input:    `<blockquote><img src="foo" /></blockquote>`,
		line:     1,
		column:   43,
		expected: `<blockquote><p><img src="foo" /></p></blockquote>`,
	},
	"A blockquote that ends in a `</span> ` gets the paragraph tags inserted": {
		input:    `<blockquote><span>foo</span> </blockquote>`,
		line:     1,
		column:   50,
		expected: `<blockquote><p><span>foo</span> </p></blockquote>`,
	},
	"A blockquote that has no html tags in it and just text and whitespace present does have the paragraph tags inserted": {
		input:    `<blockquote>   some text    </blockquote>`,
		line:     1,
		column:   42,
		expected: `<blockquote><p>   some text    </p></blockquote>`,
	},
	"A blockquote with a paragraph tag that starts and ends the content of the blockquote does not get paragraph tags inserted": {
		input:    `<blockquote><p>content</p></blockquote>`,
		line:     1,
		column:   40,
		expected: `<blockquote><p>content</p></blockquote>`,
	},
	"A blockquote with an `img` tag in it on line 2 does get the paragraph tags inserted": {
		input: `<html><body>
<blockquote><img src="foo" /></blockquote>
</body></html>`,
		line:   2,
		column: 43,
		expected: `<html><body>
<blockquote><p><img src="foo" /></p></blockquote>
</body></html>`,
	},
	"A blockquote that has no html tags in it and just text and whitespace present on line 3 does have the paragraph tags inserted": {
		input: `<html>
<body>
<blockquote>   some text    </blockquote>
</body></html>`,
		line:   3,
		column: 42,
		expected: `<html>
<body>
<blockquote><p>   some text    </p></blockquote>
</body></html>`,
	},
	"A blockquote with a paragraph tag that starts and ends the content of the blockquote on line 2 does not get paragraph tags inserted": {
		input: `<html>
<blockquote><p>content</p></blockquote>
</html>`,
		line:   2,
		column: 40,
		expected: `<html>
<blockquote><p>content</p></blockquote>
</html>`,
	},
}

func TestFixFailedBlockquoteParsing(t *testing.T) {
	for name, args := range fixFailedBlockquoteParsingTestCases {
		t.Run(name, func(t *testing.T) {
			edits := rulefixes.FixFailedBlockquoteParsing(args.line, args.column, args.input)

			checkFinalOutputMatches(t, args.input, args.expected, edits...)
		})
	}
}
