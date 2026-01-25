//go:build unit

package rulefixes_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
	"github.com/stretchr/testify/assert"
)

type fixSectionElementUnexpectedTestCase struct {
	contents        string
	line, column    int
	expectedChanges []positions.TextEdit
}

var fixSectionElementUnexpectedTestCases = map[string]fixSectionElementUnexpectedTestCase{
	"When there is an unexpected section inside an span and paragraph it should get moved outside of it": {
		contents: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xml:lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
<head>
<meta charset="utf-8"/>
<link href="../Styles/styles.css" rel="stylesheet" type="text/css"/>
<title>Chapter 14: Our Whole Family! image</title>
</head>
<body>
<p class="P_TEXTBODY_CENTERALIGN"><span><section epub:type="frontmatter titlepage"><img alt="Front Image1" class="insert" src="../Images/INTERIORIMAGES_10.jpg"/></section></span></p>
</body>
</html>`,
		line:   10,
		column: 84,
		expectedChanges: []positions.TextEdit{
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   10,
						Column: 41,
					},
					End: positions.Position{
						Line:   10,
						Column: 84,
					},
				},
			},
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   10,
						Column: 162,
					},
					End: positions.Position{
						Line:   10,
						Column: 172,
					},
				},
			},
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   10,
						Column: 1,
					},
					End: positions.Position{
						Line:   10,
						Column: 1,
					},
				},
				NewText: `<section epub:type="frontmatter titlepage">`,
			},
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   11,
						Column: 1,
					},
					End: positions.Position{
						Line:   11,
						Column: 1,
					},
				},
				NewText: "</section>",
			},
		},
	},
	//	"When there is an unexpected section inside an span, paragraph, and div it should get moved outside of the span and paragraph, but not the div": {
	//		contents: `<?xml version="1.0" encoding="utf-8"?>
	//
	// <!DOCTYPE html>
	// <html xml:lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
	// <head>
	// <meta charset="utf-8"/>
	// <link href="../Styles/styles.css" rel="stylesheet" type="text/css"/>
	// <title>Chapter 14: Our Whole Family! image</title>
	// </head>
	// <body>
	// <div><p class="P_TEXTBODY_CENTERALIGN"><span><section epub:type="frontmatter titlepage"><img alt="Front Image1" class="insert" src="../Images/INTERIORIMAGES_10.jpg"/></section></span></p></div>
	// </body>
	// </html>`,
	//
	//		line:   10,
	//		column: 89,
	//		expectedChanges: []positions.TextEdit{
	//			{
	//				Range: positions.Range{
	//					Start: positions.Position{
	//						Line:   10,
	//						Column: 46,
	//					},
	//					End: positions.Position{
	//						Line:   10,
	//						Column: 89,
	//					},
	//				},
	//			},
	//			{
	//				Range: positions.Range{
	//					Start: positions.Position{
	//						Line:   10,
	//						Column: 167,
	//					},
	//					End: positions.Position{
	//						Line:   10,
	//						Column: 177,
	//					},
	//				},
	//			},
	//			{
	//				Range: positions.Range{
	//					Start: positions.Position{
	//						Line:   10,
	//						Column: 6,
	//					},
	//					End: positions.Position{
	//						Line:   10,
	//						Column: 6,
	//					},
	//				},
	//				NewText: `<section epub:type="frontmatter titlepage">`,
	//			},
	//			{
	//				Range: positions.Range{
	//					Start: positions.Position{
	//						Line:   10,
	//						Column: 188,
	//					},
	//					End: positions.Position{
	//						Line:   10,
	//						Column: 188,
	//					},
	//				},
	//				NewText: "</section>",
	//			},
	//		},
	//	},
}

func TestFixSectionElementUnexpected(t *testing.T) {
	for name, args := range fixSectionElementUnexpectedTestCases {
		t.Run(name, func(t *testing.T) {
			actual := rulefixes.FixSectionElementUnexpected(args.line, args.column, args.contents)

			assert.Equal(t, args.expectedChanges, actual)
		})
	}
}
