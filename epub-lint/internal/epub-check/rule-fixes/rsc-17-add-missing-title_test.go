//go:build unit

package rulefixes_test

import (
	"testing"

	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
)

type addMissingTitleTestCase struct {
	input        string
	line, column int
	expected     string
}

var addMissingTitleTestCases = map[string]addMissingTitleTestCase{
	"A head element without a title element should use the first heading present as the title": {
		input: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en" xml:lang="en">
  <head>
    <meta charset="utf-8" />
  </head>
  <body>
  <h1>Table of Contents</h1>
  <h2>Subtitle</h2>
  </body>
</html>`,
		line:   4,
		column: 9,
		expected: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en" xml:lang="en">
  <head>
    <title>Table of Contents</title>
    <meta charset="utf-8" />
  </head>
  <body>
  <h1>Table of Contents</h1>
  <h2>Subtitle</h2>
  </body>
</html>`,
	},
	"A head element without a title element should use the first paragraph present as the title when no headings exist": {
		input: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en" xml:lang="en">
  <head>
    <meta charset="utf-8" />
  </head>
  <body>
  <p>Paragraph 1</p>
  <p>Paragraph 2</p>
  </body>
</html>`,
		line:   4,
		column: 9,
		expected: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en" xml:lang="en">
  <head>
    <title>Paragraph 1</title>
    <meta charset="utf-8" />
  </head>
  <body>
  <p>Paragraph 1</p>
  <p>Paragraph 2</p>
  </body>
</html>`,
	},
	"A self-closing head element should use the first heading present as the title": {
		input: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en" xml:lang="en">
  <head/>
  <body>
	<h1>Table of Contents</h1>
  <h2>Subtitle</h2>
  <p>Paragraph 1</p>
  <p>Paragraph 2</p>
  </body>
</html>`,
		line:   4,
		column: 10,
		expected: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en" xml:lang="en">
  <head><title>Table of Contents</title></head>
  <body>
	<h1>Table of Contents</h1>
  <h2>Subtitle</h2>
  <p>Paragraph 1</p>
  <p>Paragraph 2</p>
  </body>
</html>`,
	},
	"A head element without a title element, no headings, and no paragraphs will not get modified": {
		input: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en" xml:lang="en">
  <head>
    <meta charset="utf-8" />
  </head>
  <body>
  </body>
</html>`,
		line:   4,
		column: 9,
		expected: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en" xml:lang="en">
  <head>
    <meta charset="utf-8" />
  </head>
  <body>
  </body>
</html>`,
	},
	"A head element that has both the start and end tag element present on the same line, should have the title inserted on the same line": {
		input:    `<?xml version="1.0" encoding="utf-8"?><!DOCTYPE html><html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en" xml:lang="en"><head><meta charset="utf-8" /></head><body><h1>Table of Contents</h1><h2>Subtitle</h2></body></html>`,
		line:     1,
		column:   169,
		expected: `<?xml version="1.0" encoding="utf-8"?><!DOCTYPE html><html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en" xml:lang="en"><head><title>Table of Contents</title><meta charset="utf-8" /></head><body><h1>Table of Contents</h1><h2>Subtitle</h2></body></html>`,
	},
}

func TestAddMissingTitle(t *testing.T) {
	for name, args := range addMissingTitleTestCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			edit := rulefixes.AddMissingTitle(args.line, args.column, args.input)

			checkFinalOutputMatches(t, args.input, args.expected, edit)
		})
	}
}
