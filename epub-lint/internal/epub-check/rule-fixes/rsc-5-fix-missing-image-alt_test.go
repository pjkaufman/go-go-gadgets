//go:build unit

package rulefixes_test

import (
	"testing"

	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
	"github.com/stretchr/testify/assert"
)

type fixMissingImageAltTestCase struct {
	input          string
	line, column   int
	expectedChange rulefixes.TextEdit
}

var fixMissingImageAltTestCases = map[string]fixMissingImageAltTestCase{
	"An image with a missing alt should get it added correctly when other elements are on the same line": {
		input:  `<html><body><img src="test.png"/></body></html>`,
		line:   1,
		column: 34,
		expectedChange: rulefixes.TextEdit{
			Range: rulefixes.Range{
				Start: rulefixes.Position{
					Line:   1,
					Column: 32,
				},
				End: rulefixes.Position{
					Line:   1,
					Column: 32,
				},
			},
			NewText: ` alt=""`,
		},
	},
	"An image with a missing alt should get it added correctly when no other elements are on the same line": {
		input: `<?xml version="1.0" encoding="utf-8" standalone="no"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN"
  "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">

<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" xmlns:xml="http://www.w3.org/XML/1998/namespace">
<head>
  <meta content="true" name="calibre:cover" />

  <title>Cover</title>
  <style type="text/css">
/*<![CDATA[*/

  p.sgc-1 {text-align: center;}
  /*]]>*/
  </style>
</head>

<body>
  <p class="sgc-1"><img height="100%" src="../Images/cover.jpg" /></p>
</body>
</html>
`,
		line:   19,
		column: 67,
		expectedChange: rulefixes.TextEdit{
			Range: rulefixes.Range{
				Start: rulefixes.Position{
					Line:   19,
					Column: 65,
				},
				End: rulefixes.Position{
					Line:   19,
					Column: 65,
				},
			},
			NewText: `alt=""`,
		},
	},
}

func TestFixMissingImageAlt(t *testing.T) {
	for name, args := range fixMissingImageAltTestCases {
		t.Run(name, func(t *testing.T) {
			actual := rulefixes.FixMissingImageAlt(args.line, args.column, args.input)

			assert.Equal(t, args.expectedChange, actual)
		})
	}
}
