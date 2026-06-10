//go:build unit

package rulefixes_test

import (
	"testing"

	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
)

type fixEmptyTitleTestCase struct {
	input        string
	line, column int
	expected     string
}

var fixEmptyTitleTestCases = map[string]fixEmptyTitleTestCase{
	"An empty title element should get properly filled in by the first heading in the file": {
		input: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en" xml:lang="en">
  <head>
  <meta charset="utf-8" />
  <title></title>
  </head>
  <body>
  <h1>Table of Contents</h1>
  <h2>Subtitle</h2>
  </body>
</html>`,
		line:   6,
		column: 10,
		expected: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en" xml:lang="en">
  <head>
  <meta charset="utf-8" />
  <title>Table of Contents</title>
  </head>
  <body>
  <h1>Table of Contents</h1>
  <h2>Subtitle</h2>
  </body>
</html>`,
	},
	"An empty self-closing title element should get properly filled in by the first heading in the file": {
		input: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en" xml:lang="en">
  <head>
  <meta charset="utf-8" />
  <title/>
  </head>
  <body>
  <h4>Table of Contents</h4>
  <h5>Subtitle</h5>
  </body>
</html>`,
		line:   6,
		column: 11,
		expected: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en" xml:lang="en">
  <head>
  <meta charset="utf-8" />
  <title>Table of Contents</title>
  </head>
  <body>
  <h4>Table of Contents</h4>
  <h5>Subtitle</h5>
  </body>
</html>`,
	},
	"An empty title element should get filled in with the contents of the first paragraph when no heading element is present": {
		input: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en" xml:lang="en">
  <head>
  <meta charset="utf-8" />
  <title/>
  </head>
  <body>
  <p>Paragraph 1</p>
  <p>Paragraph 2</p>
  </body>
</html>`,
		line:   6,
		column: 11,
		expected: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en" xml:lang="en">
  <head>
  <meta charset="utf-8" />
  <title>Paragraph 1</title>
  </head>
  <body>
  <p>Paragraph 1</p>
  <p>Paragraph 2</p>
  </body>
</html>`,
	},
	"An empty title element should be left as is when no heading or paragraph element is present": {
		input: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en" xml:lang="en">
  <head>
  <meta charset="utf-8" />
  <title/>
  </head>
  <body>
  </body>
</html>`,
		line:   6,
		column: 11,
		expected: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en" xml:lang="en">
  <head>
  <meta charset="utf-8" />
  <title/>
  </head>
  <body>
  </body>
</html>`,
	},
}

func TestFixEmptyTitle(t *testing.T) {
	t.Parallel()

	for name, args := range fixEmptyTitleTestCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			edit := rulefixes.FixEmptyTitle(args.line, args.column, args.input)

			checkFinalOutputMatches(t, args.input, args.expected, edit)
		})
	}
}
