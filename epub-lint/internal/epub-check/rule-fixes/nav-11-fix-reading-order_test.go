//go:build unit

package rulefixes_test

import (
	"testing"

	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
)

type fixReadingOrderTestCase struct {
	spineOrder         []string
	input              string
	navPath, opfFolder string
	expected           string
}

var fixReadingOrderTestCases = map[string]fixReadingOrderTestCase{
	"When there is an issue with the reading order from the spine and the nav toc order, then the nav links should be reordered to match the spine's reading order": {
		input: `<?xml version="1.0" encoding="utf-8"?>
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
          <a href="../Text/section-0005.html">Table of Contents Page</a>
        </li>
               <li>
          <a href="../Text/section-0002.html">Color Gallery</a>
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
		expected: `<?xml version="1.0" encoding="utf-8"?>
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
		navPath:   "OEBPS/Text/",
		opfFolder: "OEBPS",
		spineOrder: []string{
			"Text/CoverPage.html",
			"Text/section-0001.html",
			"Text/section-0002.html",
			"Text/section-0003.html",
			"Text/section-0004.html",
			"Text/section-0005.html",
			"Text/section-0006.html",
			"Text/section-0007.html",
			"Text/section-0008.html",
			"Text/section-0009.html",
			"Text/section-0010.html",
			"Text/section-0011.html",
			"Text/section-0012.html",
			"Text/section-0013.html",
			"Text/section-0014.html",
			"Text/section-0015.html",
			"Text/section-0016.html",
			"Text/section-0017.html",
			"Text/section-0018.html",
			"Text/section-0019.html",
			"Text/section-0020.html",
			"Text/section-0021.html",
			"Text/section-0022.html",
			"Text/section-0023.html",
			"Text/section-0024.html",
			"Text/section-0025.html",
			"Text/section-0026.html",
			"Text/section-0027.html",
			"Text/section-0028.html",
			"Text/section-0029.html",
			"Text/nav.xhtml",
		},
	},
}

func TestFixReadingOrder(t *testing.T) {
	t.Parallel()

	for name, args := range fixReadingOrderTestCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			edits := rulefixes.FixReadingOrder(args.spineOrder, args.input, args.navPath, args.opfFolder)

			checkFinalOutputMatches(t, args.input, args.expected, edits...)
		})
	}
}
