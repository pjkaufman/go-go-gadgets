//go:build unit

package epubhandler_test

import (
	"testing"

	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
	"github.com/stretchr/testify/assert"
)

type addFileToNavTestCase struct {
	inputText      string
	inputPath      string
	inputTitle     string
	expectedOutput string
}

const (
	simpleNoTOCNav = `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>

<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en" xml:lang="en">
<head>
  <meta charset="utf-8"/>
  <title>ePub Nav</title>
  <style type="text/css">
  ol { list-style-type: none; }
  </style>
  </head>
  <body epub:type="frontmatter">
  <nav epub:type="landmarks" id="landmarks" hidden="">
  <h2>Guide</h2>
  <ol>
  <li>
  <a epub:type="cover" href="../Text/CoverPage.html">Cover Page</a>
  </li>
  <li>
  <a epub:type="toc" href="../Text/section-0001.html">Table of Contents</a>
  </li>
  </ol>
  </nav>
</body>
</html>`
	simpleNoOlInTOCNav = `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>

<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en" xml:lang="en">
<head>
  <meta charset="utf-8"/>
  <title>ePub Nav</title>
  <style type="text/css">
  ol { list-style-type: none; }
  </style>
  </head>
  <body epub:type="frontmatter">
	<nav epub:type="toc" id="toc">
  <h1>Table of Contents</h1>
	</nav>
  <nav epub:type="landmarks" id="landmarks" hidden="">
  <h2>Guide</h2>
  <ol>
  <li>
  <a epub:type="cover" href="../Text/CoverPage.html">Cover Page</a>
  </li>
  <li>
  <a epub:type="toc" href="../Text/section-0001.html">Table of Contents</a>
  </li>
  </ol>
  </nav>
</body>
</html>`
)

var addFileToNavTestCases = map[string]addFileToNavTestCase{
	"When there is no epub type toc, there is no change made to the corresponding nav": {
		inputText:      simpleNoTOCNav,
		expectedOutput: simpleNoTOCNav,
	},
	"When there is no order list in the toc, there is no change made to the corresponding nav": {
		inputText:      simpleNoOlInTOCNav,
		expectedOutput: simpleNoOlInTOCNav,
	},
	"When there is an empty order list in the toc, the new list item is properly inserted at the end": {
		inputText: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>

<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en" xml:lang="en">
<head>
  <meta charset="utf-8"/>
  <title>ePub Nav</title>
  <style type="text/css">
  ol { list-style-type: none; }
  </style>
  </head>
  <body epub:type="frontmatter">
	<nav epub:type="toc" id="toc">
  <h1>Table of Contents</h1>
	<ol>
	</ol>
	</nav>
  <nav epub:type="landmarks" id="landmarks" hidden="">
  <h2>Guide</h2>
  <ol>
  <li>
  <a epub:type="cover" href="../Text/CoverPage.html">Cover Page</a>
  </li>
  <li>
  <a epub:type="toc" href="../Text/section-0001.html">Table of Contents</a>
  </li>
  </ol>
  </nav>
</body>
</html>`,
		inputPath:  "Text/tl_notes.xhtml",
		inputTitle: "Translator's Notes",
		expectedOutput: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>

<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en" xml:lang="en">
<head>
  <meta charset="utf-8"/>
  <title>ePub Nav</title>
  <style type="text/css">
  ol { list-style-type: none; }
  </style>
  </head>
  <body epub:type="frontmatter">
	<nav epub:type="toc" id="toc">
  <h1>Table of Contents</h1>
	<ol>
	<li><a href="Text/tl_notes.xhtml">Translator's Notes</a></li>
</ol>
	</nav>
  <nav epub:type="landmarks" id="landmarks" hidden="">
  <h2>Guide</h2>
  <ol>
  <li>
  <a epub:type="cover" href="../Text/CoverPage.html">Cover Page</a>
  </li>
  <li>
  <a epub:type="toc" href="../Text/section-0001.html">Table of Contents</a>
  </li>
  </ol>
  </nav>
</body>
</html>`,
	},
	"When there is a filled out order list in the toc, the new list item is properly inserted at the end": {
		inputText: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>

<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en" xml:lang="en">
<head>
  <meta charset="utf-8"/>
  <title>ePub Nav</title>
  <style type="text/css">
  ol { list-style-type: none; }
  </style>
  </head>
  <body epub:type="frontmatter">
	<nav epub:type="toc" id="toc">
  <h1>Table of Contents</h1>
	<ol>
  <li>
  <a href="../Text/section-0001.html#tableofcontents">Table of Contents</a>
  </li>
  <li>
  <a href="../Text/section-0002.html">Color Gallery</a>
  </li>
  <li>
  <a href="../Text/section-0005.html">Table of Contents Page</a>
  </li>
  <li>
  <a href="../Text/section-0006.html">Title Page</a>
  </li>
  <li>
  <a href="../Text/section-0007.html">Copyrights and Credits</a>
  </li>
  <li>
  <a href="../Text/section-0008.html#auto_bookmark_toc_8">Prologue</a>
  </li>
  <li>
  <a href="../Text/section-0009.html#auto_bookmark_toc_9">Chapter 1</a>
  </li>
  <li>
  <a href="../Text/section-0012.html#auto_bookmark_toc_12">Chapter 2</a>
  </li>
  <li>
  <a href="../Text/section-0015.html#auto_bookmark_toc_15">Chapter 3</a>
  </li>
  <li>
  <a href="../Text/section-0018.html#auto_bookmark_toc_18">Chapter 4</a>
  </li>
  <li>
  <a href="../Text/section-0021.html#auto_bookmark_toc_21">Chapter 5</a>
  </li>
  <li>
  <a href="../Text/section-0026.html#auto_bookmark_toc_26">Epilogue</a>
  </li>
  <li>
  <a href="../Text/section-0028.html#auto_bookmark_toc_28">Report: Cohort 452 Year 1, Term 1</a>
  </li>
  <li>
  <a href="../Text/section-0029.html">Newsletter</a>
  </li>
	</ol>
	</nav>
  <nav epub:type="landmarks" id="landmarks" hidden="">
  <h2>Guide</h2>
  <ol>
  <li>
  <a epub:type="cover" href="../Text/CoverPage.html">Cover Page</a>
  </li>
  <li>
  <a epub:type="toc" href="../Text/section-0001.html">Table of Contents</a>
  </li>
  </ol>
  </nav>
</body>
</html>`,
		inputPath:  "Text/tl_notes.xhtml",
		inputTitle: "Translator's Notes",
		expectedOutput: `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>

<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en" xml:lang="en">
<head>
  <meta charset="utf-8"/>
  <title>ePub Nav</title>
  <style type="text/css">
  ol { list-style-type: none; }
  </style>
  </head>
  <body epub:type="frontmatter">
	<nav epub:type="toc" id="toc">
  <h1>Table of Contents</h1>
	<ol>
  <li>
  <a href="../Text/section-0001.html#tableofcontents">Table of Contents</a>
  </li>
  <li>
  <a href="../Text/section-0002.html">Color Gallery</a>
  </li>
  <li>
  <a href="../Text/section-0005.html">Table of Contents Page</a>
  </li>
  <li>
  <a href="../Text/section-0006.html">Title Page</a>
  </li>
  <li>
  <a href="../Text/section-0007.html">Copyrights and Credits</a>
  </li>
  <li>
  <a href="../Text/section-0008.html#auto_bookmark_toc_8">Prologue</a>
  </li>
  <li>
  <a href="../Text/section-0009.html#auto_bookmark_toc_9">Chapter 1</a>
  </li>
  <li>
  <a href="../Text/section-0012.html#auto_bookmark_toc_12">Chapter 2</a>
  </li>
  <li>
  <a href="../Text/section-0015.html#auto_bookmark_toc_15">Chapter 3</a>
  </li>
  <li>
  <a href="../Text/section-0018.html#auto_bookmark_toc_18">Chapter 4</a>
  </li>
  <li>
  <a href="../Text/section-0021.html#auto_bookmark_toc_21">Chapter 5</a>
  </li>
  <li>
  <a href="../Text/section-0026.html#auto_bookmark_toc_26">Epilogue</a>
  </li>
  <li>
  <a href="../Text/section-0028.html#auto_bookmark_toc_28">Report: Cohort 452 Year 1, Term 1</a>
  </li>
  <li>
  <a href="../Text/section-0029.html">Newsletter</a>
  </li>
	<li><a href="Text/tl_notes.xhtml">Translator's Notes</a></li>
</ol>
	</nav>
  <nav epub:type="landmarks" id="landmarks" hidden="">
  <h2>Guide</h2>
  <ol>
  <li>
  <a epub:type="cover" href="../Text/CoverPage.html">Cover Page</a>
  </li>
  <li>
  <a epub:type="toc" href="../Text/section-0001.html">Table of Contents</a>
  </li>
  </ol>
  </nav>
</body>
</html>`,
	},
}

func TestAddFileToNav(t *testing.T) {
	t.Parallel()

	for name, args := range addFileToNavTestCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actual := epubhandler.AddFileToNav(args.inputText, args.inputPath, args.inputTitle)

			assert.Equal(t, args.expectedOutput, actual)
		})
	}
}
