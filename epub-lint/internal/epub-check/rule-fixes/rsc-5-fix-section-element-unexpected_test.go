//go:build unit

package rulefixes_test

import (
	"testing"

	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
)

type fixSectionElementUnexpectedTestCase struct {
	contents       string
	line, column   int
	expectedOutput string
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
		expectedOutput: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xml:lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
<head>
<meta charset="utf-8"/>
<link href="../Styles/styles.css" rel="stylesheet" type="text/css"/>
<title>Chapter 14: Our Whole Family! image</title>
</head>
<body>
<section epub:type="frontmatter titlepage"><p class="P_TEXTBODY_CENTERALIGN"><span><img alt="Front Image1" class="insert" src="../Images/INTERIORIMAGES_10.jpg"/></span></p></section>
</body>
</html>`,
	},
	"When there is an unexpected section inside an span, paragraph, and div it should get moved outside of the span and paragraph, but not the div": {
		contents: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xml:lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
<head>
<meta charset="utf-8"/>
<link href="../Styles/styles.css" rel="stylesheet" type="text/css"/>
<title>Chapter 14: Our Whole Family! image</title>
</head>
<body>
<div><p class="P_TEXTBODY_CENTERALIGN"><span><section epub:type="frontmatter titlepage"><img alt="Front Image1" class="insert" src="../Images/INTERIORIMAGES_10.jpg"/></section></span></p></div>
</body>
</html>`,
		line:   10,
		column: 89,
		expectedOutput: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xml:lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
<head>
<meta charset="utf-8"/>
<link href="../Styles/styles.css" rel="stylesheet" type="text/css"/>
<title>Chapter 14: Our Whole Family! image</title>
</head>
<body>
<div><section epub:type="frontmatter titlepage"><p class="P_TEXTBODY_CENTERALIGN"><span><img alt="Front Image1" class="insert" src="../Images/INTERIORIMAGES_10.jpg"/></span></p></section></div>
</body>
</html>`,
	},
}

func TestFixSectionElementUnexpected(t *testing.T) {
	for name, args := range fixSectionElementUnexpectedTestCases {
		t.Run(name, func(t *testing.T) {
			edits := rulefixes.FixSectionElementUnexpected(args.line, args.column, args.contents)

			checkFinalOutputMatches(t, args.contents, args.expectedOutput, edits...)
		})
	}
}
