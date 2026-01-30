//go:build unit

package rulefixes_test

import (
	"testing"

	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
)

type fixIrregularDoctypeTestCase struct {
	input           string
	expectedDoctype string
	expected        string
}

var fixIrregularDoctypeTestCases = map[string]fixIrregularDoctypeTestCase{
	"A blockquote with another blockquote inside of it with text in that blockquote does not get a div tag inserted": {
		input: `<?xml version="1.0" encoding="utf-8" standalone="no"?>
<!DOCTYPE html>
  
<html xmlns="http://www.w3.org/1999/xhtml">
    <head>
        <title>Title</title>
        <link href="../Styles/default.css" rel="stylesheet" type="text/css" />
        <meta content="width=1150, height=1725" name="viewport" />
    </head>
    
    <body></body>
</html>
`,
		expectedDoctype: `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">`,
		expected: `<?xml version="1.0" encoding="utf-8" standalone="no"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">
  
<html xmlns="http://www.w3.org/1999/xhtml">
    <head>
        <title>Title</title>
        <link href="../Styles/default.css" rel="stylesheet" type="text/css" />
        <meta content="width=1150, height=1725" name="viewport" />
    </head>
    
    <body></body>
</html>
`,
	},
}

func TestFixIrregularDoctype(t *testing.T) {
	for name, args := range fixIrregularDoctypeTestCases {
		t.Run(name, func(t *testing.T) {
			edit := rulefixes.FixIrregularDoctype(args.input, args.expectedDoctype)

			checkFinalOutputMatches(t, args.input, args.expected, edit)
		})
	}
}
