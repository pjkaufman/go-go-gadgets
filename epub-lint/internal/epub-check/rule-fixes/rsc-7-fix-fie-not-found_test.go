//go:build unit

package rulefixes_test

import (
	"testing"

	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
	"github.com/stretchr/testify/assert"
)

type fixFileNotFoundTestCase struct {
	contents            string
	currentFile         string
	requestedResource   string
	basenameToFilePaths map[string][]string
	line, column        int
	expectedOutput      string
}

var fixFileNotFoundTestCases = map[string]fixFileNotFoundTestCase{
	"When there is a stylesheet link that is not found and there is no file matching the name in the epub, remove it": {
		contents: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xml:lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
<head>
<meta charset="utf-8"/>
<link href="../Styles/styles.css" rel="stylesheet" type="text/css"/>
<title>Some Title</title>
</head>
<body>
<p>Some text...</p>
</body>
</html>`,
		line:              6,
		column:            69,
		requestedResource: "OEBPS/Styles/styles.css",
		currentFile:       "OEBPS/Text/file.html",
		expectedOutput: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xml:lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
<head>
<meta charset="utf-8"/>
<title>Some Title</title>
</head>
<body>
<p>Some text...</p>
</body>
</html>`,
	},
	"When there is a stylesheet link that is not found and there is a file matching the name in the epub, replace the reference with the correct relative path": {
		contents: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xml:lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
<head>
<meta charset="utf-8"/>
<link href="../Styles/styles.css" rel="stylesheet" type="text/css"/>
<title>Some Title</title>
</head>
<body>
<p>Some text...</p>
</body>
</html>`,
		line:              6,
		column:            69,
		requestedResource: "OEBPS/Styles/styles.css",
		currentFile:       "OEBPS/Text/file.html",
		basenameToFilePaths: map[string][]string{
			"styles.css": {"OEBPS/styles.css"},
		},
		expectedOutput: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xml:lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
<head>
<meta charset="utf-8"/>
<link href="../styles.css" rel="stylesheet" type="text/css"/>
<title>Some Title</title>
</head>
<body>
<p>Some text...</p>
</body>
</html>`,
	},
	"When there is an image that is not found and there is a file matching the name in the epub, remove it": {
		contents: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xml:lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
<head>
<meta charset="utf-8"/>
<link href="../Styles/styles.css" rel="stylesheet" type="text/css"/>
<title>Some Title</title>
</head>
<body>
<p><img src="../Images/graphic.png"/></p>
<p>Some text...</p>
</body>
</html>`,
		line:              10,
		column:            38,
		requestedResource: "OEBPS/Images/graphic.png",
		currentFile:       "OEBPS/Text/file.html",
		expectedOutput: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xml:lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
<head>
<meta charset="utf-8"/>
<link href="../Styles/styles.css" rel="stylesheet" type="text/css"/>
<title>Some Title</title>
</head>
<body>
<p></p>
<p>Some text...</p>
</body>
</html>`,
	},
	"When there is an image that is not found and there is no file matching the name in the epub, replace the reference with the correct relative path": {
		contents: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xml:lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
<head>
<meta charset="utf-8"/>
<link href="../Styles/styles.css" rel="stylesheet" type="text/css"/>
<title>Some Title</title>
</head>
<body>
<p><img src="../Images/graphic.png"/></p>
<p>Some text...</p>
</body>
</html>`,
		line:              10,
		column:            38,
		requestedResource: "OEBPS/Images/graphic.png",
		currentFile:       "OEBPS/Text/file.html",
		basenameToFilePaths: map[string][]string{
			"graphic.png": {"OEBPS/graphic.png"},
		},
		expectedOutput: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xml:lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
<head>
<meta charset="utf-8"/>
<link href="../Styles/styles.css" rel="stylesheet" type="text/css"/>
<title>Some Title</title>
</head>
<body>
<p><img src="../graphic.png"/></p>
<p>Some text...</p>
</body>
</html>`,
	},
	"When there is a file that is not found and there is a file matching the name in the epub, remove it": {
		contents: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xml:lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
<head>
<meta charset="utf-8"/>
<link href="../Styles/styles.css" rel="stylesheet" type="text/css"/>
<title>Some Title</title>
</head>
<body>
<p><a href="../files.html">Some text</a></p>
<p>Some text...</p>
</body>
</html>`,
		line:              10,
		column:            29,
		requestedResource: "OEBPS/files.html",
		currentFile:       "OEBPS/Text/file.html",
		expectedOutput: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xml:lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
<head>
<meta charset="utf-8"/>
<link href="../Styles/styles.css" rel="stylesheet" type="text/css"/>
<title>Some Title</title>
</head>
<body>
<p></p>
<p>Some text...</p>
</body>
</html>`,
	},
	"When there is a file that is not found and there is no file matching the name in the epub, replace the reference with the correct relative path": {
		contents: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xml:lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
<head>
<meta charset="utf-8"/>
<link href="../Styles/styles.css" rel="stylesheet" type="text/css"/>
<title>Some Title</title>
</head>
<body>
<p><a href="../files.html">Some text</a></p>
<p>Some text...</p>
</body>
</html>`,
		line:              10,
		column:            29,
		requestedResource: "OEBPS/Text/files.html",
		currentFile:       "OEBPS/Text/file.html",
		basenameToFilePaths: map[string][]string{
			"files.html": {"OEBPS/Text/files.html"},
		},
		expectedOutput: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xml:lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
<head>
<meta charset="utf-8"/>
<link href="../Styles/styles.css" rel="stylesheet" type="text/css"/>
<title>Some Title</title>
</head>
<body>
<p><a href="files.html">Some text</a></p>
<p>Some text...</p>
</body>
</html>`,
	},
}

func TestFixFileNotFound(t *testing.T) {
	for name, args := range fixFileNotFoundTestCases {
		t.Run(name, func(t *testing.T) {
			if args.basenameToFilePaths == nil {
				args.basenameToFilePaths = map[string][]string{}
			}

			edit, err := rulefixes.FixFileNotFound(args.contents, args.requestedResource, args.currentFile, args.line, args.column, args.basenameToFilePaths)

			assert.Nil(t, err)
			checkFinalOutputMatches(t, args.contents, args.expectedOutput, edit)
		})
	}
}
