//go:build unit

package rulefixes_test

import (
	"testing"

	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
)

type fixUnreachableFileTestCase struct {
	opfContents    string
	line, column   int
	expectedOutput string
}

var fixUnreachableFileTestCases = map[string]fixUnreachableFileTestCase{
	"When there is a manifest entry that is not reachable, it should have `linear=\"no\"` removed": {
		opfContents: `<?xml version='1.0' encoding='utf-8'?>
<package xmlns="http://www.idpf.org/2007/opf" xmlns:dc="http://purl.org/dc/elements/1.1/" version="3.0" unique-identifier="BookId" prefix="rendition: http://www.idpf.org/vocab/rendition/#">
  <metadata xmlns:opf="http://www.idpf.org/2007/opf" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcterms="http://purl.org/dc/terms/">
    <dc:title id="title1">Title</dc:title>
    <meta refines="#title1" property="title-type">main</meta>
    <dc:language>en</dc:language>
    <meta name="cover" content="CoverDesign.jpg"/>
    <meta name="Sigil version" content="0.9.14"/>
    <meta property="dcterms:modified">2026-03-25T02:01:56Z</meta>
  </metadata>
  <manifest>
    <item id="ncx" href="toc.ncx" media-type="application/x-dtbncx+xml"/>
    <item id="CoverPage_html" href="Text/CoverPage.html" media-type="application/xhtml+xml"/>
    <item id="styles_css" href="Styles/styles.css" media-type="text/css"/>
    <item id="toc" href="Text/section-0001.html" media-type="application/xhtml+xml"/>
    <item id="section-0002_html" href="Text/section-0002.html" media-type="application/xhtml+xml"/>
    <item id="section-0003_html" href="Text/section-0003.html" media-type="application/xhtml+xml"/>
    <item id="section-0004_html" href="Text/section-0004.html" media-type="application/xhtml+xml"/>
    <item id="section-0005_html" href="Text/section-0005.html" media-type="application/xhtml+xml"/>
    <item id="section-0006_html" href="Text/section-0006.html" media-type="application/xhtml+xml"/>
    <item id="section-0007_html" href="Text/section-0007.html" media-type="application/xhtml+xml"/>
    <item id="section-0008_html" href="Text/section-0008.html" media-type="application/xhtml+xml"/>
    <item id="section-0009_html" href="Text/section-0009.html" media-type="application/xhtml+xml"/>
    <item id="section-0010_html" href="Text/section-0010.html" media-type="application/xhtml+xml"/>
    <item id="section-0011_html" href="Text/section-0011.html" media-type="application/xhtml+xml"/>
    <item id="section-0012_html" href="Text/section-0012.html" media-type="application/xhtml+xml"/>
    <item id="section-0013_html" href="Text/section-0013.html" media-type="application/xhtml+xml"/>
    <item id="section-0014_html" href="Text/section-0014.html" media-type="application/xhtml+xml"/>
    <item id="section-0015_html" href="Text/section-0015.html" media-type="application/xhtml+xml"/>
    <item id="section-0016_html" href="Text/section-0016.html" media-type="application/xhtml+xml"/>
    <item id="section-0017_html" href="Text/section-0017.html" media-type="application/xhtml+xml"/>
    <item id="section-0018_html" href="Text/section-0018.html" media-type="application/xhtml+xml"/>
    <item id="section-0019_html" href="Text/section-0019.html" media-type="application/xhtml+xml"/>
    <item id="section-0020_html" href="Text/section-0020.html" media-type="application/xhtml+xml"/>
    <item id="section-0021_html" href="Text/section-0021.html" media-type="application/xhtml+xml"/>
    <item id="section-0022_html" href="Text/section-0022.html" media-type="application/xhtml+xml"/>
    <item id="section-0023_html" href="Text/section-0023.html" media-type="application/xhtml+xml"/>
    <item id="section-0024_html" href="Text/section-0024.html" media-type="application/xhtml+xml"/>
    <item id="section-0025_html" href="Text/section-0025.html" media-type="application/xhtml+xml"/>
    <item id="section-0026_html" href="Text/section-0026.html" media-type="application/xhtml+xml"/>
    <item id="section-0027_html" href="Text/section-0027.html" media-type="application/xhtml+xml"/>
    <item id="section-0028_html" href="Text/section-0028.html" media-type="application/xhtml+xml"/>
    <item id="section-0029_html" href="Text/section-0029.html" media-type="application/xhtml+xml"/>
    <item id="navid" href="Text/nav.xhtml" media-type="application/xhtml+xml" properties="nav"/>
    </manifest>
  <spine toc="ncx">
    <itemref idref="CoverPage_html"/>
    <itemref idref="toc"/>
    <itemref idref="section-0002_html"/>
    <itemref idref="section-0003_html"/>
    <itemref idref="section-0004_html"/>
    <itemref idref="section-0005_html"/>
    <itemref idref="section-0006_html"/>
    <itemref idref="section-0007_html"/>
    <itemref idref="section-0008_html"/>
    <itemref idref="section-0009_html"/>
    <itemref idref="section-0010_html"/>
    <itemref idref="section-0011_html"/>
    <itemref idref="section-0012_html"/>
    <itemref idref="section-0013_html"/>
    <itemref idref="section-0014_html"/>
    <itemref idref="section-0015_html"/>
    <itemref idref="section-0016_html"/>
    <itemref idref="section-0017_html"/>
    <itemref idref="section-0018_html"/>
    <itemref idref="section-0019_html"/>
    <itemref idref="section-0020_html"/>
    <itemref idref="section-0021_html"/>
    <itemref idref="section-0022_html"/>
    <itemref idref="section-0023_html"/>
    <itemref idref="section-0024_html"/>
    <itemref idref="section-0025_html"/>
    <itemref idref="section-0026_html"/>
    <itemref idref="section-0027_html"/>
    <itemref idref="section-0028_html"/>
    <itemref idref="section-0029_html"/>
    <itemref idref="navid" linear="no"/>
    </spine>
  <guide>
    <reference type="cover" title="Cover Page" href="Text/CoverPage.html"/>
    <reference type="toc" title="Table of Contents" href="Text/section-0001.html"/>
  </guide>
</package>`,
		line:   44,
		column: 97,
		expectedOutput: `<?xml version='1.0' encoding='utf-8'?>
<package xmlns="http://www.idpf.org/2007/opf" xmlns:dc="http://purl.org/dc/elements/1.1/" version="3.0" unique-identifier="BookId" prefix="rendition: http://www.idpf.org/vocab/rendition/#">
  <metadata xmlns:opf="http://www.idpf.org/2007/opf" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcterms="http://purl.org/dc/terms/">
    <dc:title id="title1">Title</dc:title>
    <meta refines="#title1" property="title-type">main</meta>
    <dc:language>en</dc:language>
    <meta name="cover" content="CoverDesign.jpg"/>
    <meta name="Sigil version" content="0.9.14"/>
    <meta property="dcterms:modified">2026-03-25T02:01:56Z</meta>
  </metadata>
  <manifest>
    <item id="ncx" href="toc.ncx" media-type="application/x-dtbncx+xml"/>
    <item id="CoverPage_html" href="Text/CoverPage.html" media-type="application/xhtml+xml"/>
    <item id="styles_css" href="Styles/styles.css" media-type="text/css"/>
    <item id="toc" href="Text/section-0001.html" media-type="application/xhtml+xml"/>
    <item id="section-0002_html" href="Text/section-0002.html" media-type="application/xhtml+xml"/>
    <item id="section-0003_html" href="Text/section-0003.html" media-type="application/xhtml+xml"/>
    <item id="section-0004_html" href="Text/section-0004.html" media-type="application/xhtml+xml"/>
    <item id="section-0005_html" href="Text/section-0005.html" media-type="application/xhtml+xml"/>
    <item id="section-0006_html" href="Text/section-0006.html" media-type="application/xhtml+xml"/>
    <item id="section-0007_html" href="Text/section-0007.html" media-type="application/xhtml+xml"/>
    <item id="section-0008_html" href="Text/section-0008.html" media-type="application/xhtml+xml"/>
    <item id="section-0009_html" href="Text/section-0009.html" media-type="application/xhtml+xml"/>
    <item id="section-0010_html" href="Text/section-0010.html" media-type="application/xhtml+xml"/>
    <item id="section-0011_html" href="Text/section-0011.html" media-type="application/xhtml+xml"/>
    <item id="section-0012_html" href="Text/section-0012.html" media-type="application/xhtml+xml"/>
    <item id="section-0013_html" href="Text/section-0013.html" media-type="application/xhtml+xml"/>
    <item id="section-0014_html" href="Text/section-0014.html" media-type="application/xhtml+xml"/>
    <item id="section-0015_html" href="Text/section-0015.html" media-type="application/xhtml+xml"/>
    <item id="section-0016_html" href="Text/section-0016.html" media-type="application/xhtml+xml"/>
    <item id="section-0017_html" href="Text/section-0017.html" media-type="application/xhtml+xml"/>
    <item id="section-0018_html" href="Text/section-0018.html" media-type="application/xhtml+xml"/>
    <item id="section-0019_html" href="Text/section-0019.html" media-type="application/xhtml+xml"/>
    <item id="section-0020_html" href="Text/section-0020.html" media-type="application/xhtml+xml"/>
    <item id="section-0021_html" href="Text/section-0021.html" media-type="application/xhtml+xml"/>
    <item id="section-0022_html" href="Text/section-0022.html" media-type="application/xhtml+xml"/>
    <item id="section-0023_html" href="Text/section-0023.html" media-type="application/xhtml+xml"/>
    <item id="section-0024_html" href="Text/section-0024.html" media-type="application/xhtml+xml"/>
    <item id="section-0025_html" href="Text/section-0025.html" media-type="application/xhtml+xml"/>
    <item id="section-0026_html" href="Text/section-0026.html" media-type="application/xhtml+xml"/>
    <item id="section-0027_html" href="Text/section-0027.html" media-type="application/xhtml+xml"/>
    <item id="section-0028_html" href="Text/section-0028.html" media-type="application/xhtml+xml"/>
    <item id="section-0029_html" href="Text/section-0029.html" media-type="application/xhtml+xml"/>
    <item id="navid" href="Text/nav.xhtml" media-type="application/xhtml+xml" properties="nav"/>
    </manifest>
  <spine toc="ncx">
    <itemref idref="CoverPage_html"/>
    <itemref idref="toc"/>
    <itemref idref="section-0002_html"/>
    <itemref idref="section-0003_html"/>
    <itemref idref="section-0004_html"/>
    <itemref idref="section-0005_html"/>
    <itemref idref="section-0006_html"/>
    <itemref idref="section-0007_html"/>
    <itemref idref="section-0008_html"/>
    <itemref idref="section-0009_html"/>
    <itemref idref="section-0010_html"/>
    <itemref idref="section-0011_html"/>
    <itemref idref="section-0012_html"/>
    <itemref idref="section-0013_html"/>
    <itemref idref="section-0014_html"/>
    <itemref idref="section-0015_html"/>
    <itemref idref="section-0016_html"/>
    <itemref idref="section-0017_html"/>
    <itemref idref="section-0018_html"/>
    <itemref idref="section-0019_html"/>
    <itemref idref="section-0020_html"/>
    <itemref idref="section-0021_html"/>
    <itemref idref="section-0022_html"/>
    <itemref idref="section-0023_html"/>
    <itemref idref="section-0024_html"/>
    <itemref idref="section-0025_html"/>
    <itemref idref="section-0026_html"/>
    <itemref idref="section-0027_html"/>
    <itemref idref="section-0028_html"/>
    <itemref idref="section-0029_html"/>
    <itemref idref="navid"/>
    </spine>
  <guide>
    <reference type="cover" title="Cover Page" href="Text/CoverPage.html"/>
    <reference type="toc" title="Table of Contents" href="Text/section-0001.html"/>
  </guide>
</package>`,
	},
}

func TestFixUnreachableFile(t *testing.T) {
	t.Parallel()

	for name, args := range fixUnreachableFileTestCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			edit := rulefixes.FixUnreachableFile(args.line, args.column, args.opfContents)

			checkFinalOutputMatches(t, args.opfContents, args.expectedOutput, edit)
		})
	}
}
