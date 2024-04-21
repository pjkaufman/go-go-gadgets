//go:build unit

package linter_test

import (
	"strings"
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/linter"
	"github.com/stretchr/testify/assert"
)

type ParseOpfContentsTestCase struct {
	InputText        string
	ExpectedEpubInfo linter.EpubInfo
	ExpectedErr      error
	IsSyntaxError    bool
}

const (
	epub3PackageFile = `<?xml version="1.0" encoding="utf-8"?>
<package version="3.0" unique-identifier="BookId" xmlns="http://www.idpf.org/2007/opf">
  <metadata xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcterms="http://purl.org/dc/terms/" xmlns:opf="http://www.idpf.org/2007/opf" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
    <dc:title id="title1">Mushoku Tensei: Jobless Reincarnation Vol. 24</dc:title>
    <dc:creator id="id">Rifujin na Magonote</dc:creator>
    <dc:identifier>calibre:17</dc:identifier>
    <dc:identifier>uuid:d26524a6-6710-4c70-a8f1-4b95864f2eed</dc:identifier>
    <dc:identifier id="BookId">9798888439722</dc:identifier>
    <dc:relation>http://sevenseasentertainment.com</dc:relation>
    <dc:identifier>urn:calibre:9798888439722</dc:identifier>
    <dc:language>en</dc:language>
    <dc:publisher>Seven Seas Entertainment</dc:publisher>
    <dc:subject>light novel</dc:subject>
    <meta refines="#title1" property="title-type">main</meta>
    <meta refines="#title1" property="file-as">Mushoku Tensei: Jobless Reincarnation Vol. 24</meta>
    <meta name="cover" content="cover" />
    <meta content="1.9.30" name="Sigil version" />
    <meta property="dcterms:modified">2023-09-02T17:43:45Z</meta>
    <meta refines="#id" property="role" scheme="marc:relators">aut</meta>
    <meta refines="#id" property="file-as">Magonote, Rifujin na</meta>
  </metadata>
  <manifest>
    <item id="CoverPage_html" href="Text/CoverPage.html" media-type="application/xhtml+xml"/>
    <item id="toc" href="Text/TableOfContents.html" media-type="application/xhtml+xml"/>
    <item id="jnovels.xhtml" href="Text/jnovels.xhtml" media-type="application/xhtml+xml"/>
    <item id="section-0001_html" href="Text/section-0001.html" media-type="application/xhtml+xml"/>
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
    <item id="navid" href="nav.xhtml" media-type="application/xhtml+xml" properties="nav"/>
    <item id="ncx" href="toc.ncx" media-type="application/x-dtbncx+xml"/>
    <item id="styles_css" href="Styles/styles.css" media-type="text/css"/>
    <item id="x1.png" href="Images/1.png" media-type="image/png"/>
    <item id="COLORGALLERY__jpg" href="Images/COLORGALLERY_.jpg" media-type="image/jpeg"/>
    <item id="COLORGALLERY_1_jpg" href="Images/COLORGALLERY_1.jpg" media-type="image/jpeg"/>
    <item id="COLORGALLERY_2_jpg" href="Images/COLORGALLERY_2.jpg" media-type="image/jpeg"/>
    <item id="cover" href="Images/CoverDesign.jpg" media-type="image/jpeg" properties="cover-image"/>
    <item id="FRONTMATTER__jpg" href="Images/FRONTMATTER_.jpg" media-type="image/jpeg"/>
    <item id="FRONTMATTER_2_jpg" href="Images/FRONTMATTER_2.jpg" media-type="image/jpeg"/>
    <item id="FRONTMATTER_3_jpg" href="Images/FRONTMATTER_3.jpg" media-type="image/jpeg"/>
    <item id="FRONTMATTER_4_jpg" href="Images/FRONTMATTER_4.jpg" media-type="image/jpeg"/>
    <item id="INTERIORIMAGES__jpg" href="Images/INTERIORIMAGES_.jpg" media-type="image/jpeg"/>
    <item id="INTERIORIMAGES_2_jpg" href="Images/INTERIORIMAGES_2.jpg" media-type="image/jpeg"/>
    <item id="INTERIORIMAGES_3_jpg" href="Images/INTERIORIMAGES_3.jpg" media-type="image/jpeg"/>
    <item id="INTERIORIMAGES_4_jpg" href="Images/INTERIORIMAGES_4.jpg" media-type="image/jpeg"/>
    <item id="INTERIORIMAGES_5_jpg" href="Images/INTERIORIMAGES_5.jpg" media-type="image/jpeg"/>
    <item id="INTERIORIMAGES_6_jpg" href="Images/INTERIORIMAGES_6.jpg" media-type="image/jpeg"/>
    <item id="INTERIORIMAGES_7_jpg" href="Images/INTERIORIMAGES_7.jpg" media-type="image/jpeg"/>
    <item id="sevenseaslogo_jpg" href="Images/sevenseaslogo.jpg" media-type="image/jpeg"/>
  </manifest>
  <spine toc="ncx">
    <itemref idref="CoverPage_html"/>
    <itemref idref="toc"/>
    <itemref idref="jnovels.xhtml"/>
    <itemref idref="section-0001_html"/>
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
    <itemref idref="navid" linear="no"/>
  </spine>
  <guide>
    <reference type="cover" title="Cover Page" href="Text/CoverPage.html"/>
    <reference type="toc" title="Table of Contents" href="Text/TableOfContents.html#tableofcontents"/>
  </guide>
</package>
`
	noPackageFile = `<?xml version="1.0" encoding="utf-8"?>`
	noVersionFile = `<?xml version="1.0" encoding="utf-8"?>
<package unique-identifier="BookId" xmlns="http://www.idpf.org/2007/opf">
</package>
`
	noManifestFile = `<?xml version="1.0" encoding="utf-8"?>
<package version="3.0" unique-identifier="BookId" xmlns="http://www.idpf.org/2007/opf">
</package>
`
	noManifestEndFile = `<?xml version="1.0" encoding="utf-8"?>
<package version="3.0" unique-identifier="BookId" xmlns="http://www.idpf.org/2007/opf">
<manifest>
</package>
`
	noManifestContentsFile = `<?xml version="1.0" encoding="utf-8"?>
<package version="3.0" unique-identifier="BookId" xmlns="http://www.idpf.org/2007/opf">
<manifest></manifest>
</package>
`
	epub2PackageFile = `<?xml version="1.0" encoding="utf-8"?>
<package version="2.0" unique-identifier="uuid_id" xmlns="http://www.idpf.org/2007/opf">
  <metadata xmlns:calibre="http://calibre.kovidgoyal.net/2009/metadata" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcterms="http://purl.org/dc/terms/" xmlns:opf="http://www.idpf.org/2007/opf" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
    <dc:creator opf:file-as="Satou, Tsutomu" opf:role="aut">Tsutomu Satou</dc:creator>
    <dc:contributor opf:role="bkp">calibre (3.48.0) [https://calibre-ebook.com]</dc:contributor>
    <meta name="cover" content="cover" />
    <dc:date>0101-01-01T00:00:00+00:00</dc:date>
    <meta name="calibre:title_sort" content="MKnR 23 -(g)- Isolation" />
    <meta name="calibre:series" content="The Irregular at Magic High School" />
    <dc:title>MKnR 23 -(g)- Isolation</dc:title>
    <dc:language>en</dc:language>
    <meta name="calibre:timestamp" content="2019-09-18T22:51:46.874000+00:00" />
    <meta name="calibre:series_index" content="23" />
    <dc:identifier id="uuid_id" opf:scheme="uuid">50d7b0d3-304c-4c5c-9414-63ee4c15e9f6</dc:identifier>
    <dc:identifier opf:scheme="calibre">50d7b0d3-304c-4c5c-9414-63ee4c15e9f6</dc:identifier>
    <meta name="Sigil version" content="2.0.1" />
    <dc:date opf:event="modification" xmlns:opf="http://www.idpf.org/2007/opf">2023-10-07</dc:date>
  </metadata>
  <manifest>
    <item id="cover" href="cover.jpeg" media-type="image/jpeg"/>
    <item id="titlepage" href="titlepage.xhtml" media-type="application/xhtml+xml"/>
    <item id="chapter_1.html" href="chapter_1.html" media-type="application/xhtml+xml"/>
    <item id="chapter_2.html" href="chapter_2.html" media-type="application/xhtml+xml"/>
    <item id="chapter_3.html" href="chapter_3.html" media-type="application/xhtml+xml"/>
    <item id="chapter_4.html" href="chapter_4.html" media-type="application/xhtml+xml"/>
    <item id="chapter_6.html" href="chapter_6.html" media-type="application/xhtml+xml"/>
    <item id="chapter_7.html" href="chapter_7.html" media-type="application/xhtml+xml"/>
    <item id="chapter_8.html" href="chapter_8.html" media-type="application/xhtml+xml"/>
    <item id="chapter_9.html" href="chapter_9.html" media-type="application/xhtml+xml"/>
    <item id="chapter_10.html" href="chapter_10.html" media-type="application/xhtml+xml"/>
    <item id="id12" href="index-1_1.jpg" media-type="image/jpeg"/>
    <item id="id14" href="index-2_2.jpg" media-type="image/jpeg"/>
    <item id="id15" href="index-2_3.jpg" media-type="image/jpeg"/>
    <item id="id18" href="index-3_1.jpg" media-type="image/jpeg"/>
    <item id="id7" href="index-16_2.jpg" media-type="image/jpeg"/>
    <item id="id17" href="index-38_2.jpg" media-type="image/jpeg"/>
    <item id="id20" href="index-94_2.jpg" media-type="image/jpeg"/>
    <item id="id3" href="index-116_2.jpg" media-type="image/jpeg"/>
    <item id="id5" href="index-154_2.jpg" media-type="image/jpeg"/>
    <item id="id9" href="index-175_2.jpg" media-type="image/jpeg"/>
    <item id="id11" href="index-192_2.jpg" media-type="image/jpeg"/>
    <item id="page_css" href="page_styles.css" media-type="text/css"/>
    <item id="css" href="stylesheet.css" media-type="text/css"/>
    <item id="ncx" href="toc.ncx" media-type="application/x-dtbncx+xml"/>
    <item id="chapter_5.html" href="chapter_5.html" media-type="application/xhtml+xml"/>
    <item id="afterward.html" href="afterward.html" media-type="application/xhtml+xml"/>
    <item id="translators_notes.xhtml" href="translators_notes.xhtml" media-type="application/xhtml+xml"/>
    <item id="images.xhtml" href="images.xhtml" media-type="application/xhtml+xml"/>
  </manifest>
  <spine toc="ncx">
    <itemref idref="titlepage"/>
    <itemref idref="images.xhtml"/>
    <itemref idref="chapter_1.html"/>
    <itemref idref="chapter_2.html"/>
    <itemref idref="chapter_3.html"/>
    <itemref idref="chapter_4.html"/>
    <itemref idref="chapter_5.html"/>
    <itemref idref="chapter_6.html"/>
    <itemref idref="chapter_7.html"/>
    <itemref idref="chapter_8.html"/>
    <itemref idref="chapter_9.html"/>
    <itemref idref="chapter_10.html"/>
    <itemref idref="afterward.html"/>
    <itemref idref="translators_notes.xhtml"/>
  </spine>
  <guide>
    <reference type="cover" title="Cover" href="titlepage.xhtml"/>
  </guide>
</package>
`
)

var ParseOpfContentsTestCases = map[string]ParseOpfContentsTestCase{
	"make sure that parsing an opf that does not have a package tag results in an error": {
		InputText:     noPackageFile,
		IsSyntaxError: true,
	},
	"make sure that parsing an opf that does have a package tag, but no version info results in an error": {
		InputText:   noVersionFile,
		ExpectedErr: linter.ErrNoPackageInfo,
	},
	"make sure that parsing an opf that does not have a manifest tag results in an error": {
		InputText:   noManifestFile,
		ExpectedErr: linter.ErrNoManifest,
	},
	"make sure that parsing an opf that does have a manifest tag, but no ending manifest tag results in an error": {
		InputText:     noManifestEndFile,
		IsSyntaxError: true,
	},
	"make sure that parsing an opf that does have a manifest tag, but no list items in it results in an error": {
		InputText:   noManifestContentsFile,
		ExpectedErr: linter.ErrNoItemEls,
	},
	"make sure that parsing an epub 2 has the proper version info and other package data": {
		InputText: epub2PackageFile,
		ExpectedEpubInfo: linter.EpubInfo{
			HtmlFiles: map[string]struct{}{
				"titlepage.xhtml":         {},
				"chapter_1.html":          {},
				"chapter_2.html":          {},
				"chapter_3.html":          {},
				"chapter_4.html":          {},
				"chapter_6.html":          {},
				"chapter_7.html":          {},
				"chapter_8.html":          {},
				"chapter_9.html":          {},
				"chapter_10.html":         {},
				"chapter_5.html":          {},
				"afterward.html":          {},
				"translators_notes.xhtml": {},
				"images.xhtml":            {},
			},
			ImagesFiles: map[string]struct{}{
				"cover.jpeg":      {},
				"index-1_1.jpg":   {},
				"index-2_2.jpg":   {},
				"index-2_3.jpg":   {},
				"index-3_1.jpg":   {},
				"index-16_2.jpg":  {},
				"index-38_2.jpg":  {},
				"index-94_2.jpg":  {},
				"index-116_2.jpg": {},
				"index-154_2.jpg": {},
				"index-175_2.jpg": {},
				"index-192_2.jpg": {},
			},
			CssFiles: map[string]struct{}{
				"page_styles.css": {},
				"stylesheet.css":  {},
			},
			OtherFiles: map[string]struct{}{
				"toc.ncx": {},
			},
			TocFile: "",
			NavFile: "",
			NcxFile: "toc.ncx",
			Version: 2,
		},
	},
	"make sure that parsing an epub 3 has the proper version info and other package data": {
		InputText: epub3PackageFile,
		ExpectedEpubInfo: linter.EpubInfo{
			HtmlFiles: map[string]struct{}{
				"Text/CoverPage.html":       {},
				"Text/TableOfContents.html": {},
				"Text/jnovels.xhtml":        {},
				"Text/section-0001.html":    {},
				"Text/section-0002.html":    {},
				"Text/section-0003.html":    {},
				"Text/section-0004.html":    {},
				"Text/section-0005.html":    {},
				"Text/section-0006.html":    {},
				"Text/section-0007.html":    {},
				"Text/section-0008.html":    {},
				"Text/section-0009.html":    {},
				"Text/section-0010.html":    {},
				"Text/section-0011.html":    {},
				"Text/section-0012.html":    {},
				"Text/section-0013.html":    {},
				"Text/section-0014.html":    {},
				"Text/section-0015.html":    {},
				"Text/section-0016.html":    {},
				"Text/section-0017.html":    {},
				"Text/section-0018.html":    {},
				"nav.xhtml":                 {},
			},
			ImagesFiles: map[string]struct{}{
				"Images/1.png":                {},
				"Images/COLORGALLERY_.jpg":    {},
				"Images/COLORGALLERY_1.jpg":   {},
				"Images/COLORGALLERY_2.jpg":   {},
				"Images/CoverDesign.jpg":      {},
				"Images/FRONTMATTER_.jpg":     {},
				"Images/FRONTMATTER_2.jpg":    {},
				"Images/FRONTMATTER_3.jpg":    {},
				"Images/FRONTMATTER_4.jpg":    {},
				"Images/INTERIORIMAGES_.jpg":  {},
				"Images/INTERIORIMAGES_2.jpg": {},
				"Images/INTERIORIMAGES_3.jpg": {},
				"Images/INTERIORIMAGES_4.jpg": {},
				"Images/INTERIORIMAGES_5.jpg": {},
				"Images/INTERIORIMAGES_6.jpg": {},
				"Images/INTERIORIMAGES_7.jpg": {},
				"Images/sevenseaslogo.jpg":    {},
			},
			CssFiles: map[string]struct{}{
				"Styles/styles.css": {},
			},
			OtherFiles: map[string]struct{}{
				"toc.ncx": {},
			},
			TocFile: "Text/TableOfContents.html",
			NavFile: "nav.xhtml",
			NcxFile: "toc.ncx",
			Version: 3,
		},
	},
}

func TestParseOpfContents(t *testing.T) {
	for name, args := range ParseOpfContentsTestCases {
		t.Run(name, func(t *testing.T) {
			actual, err := linter.ParseOpfFile(args.InputText)

			if !args.IsSyntaxError {
				assert.Equal(t, args.ExpectedErr, err)
			} else {
				assert.NotNil(t, err)
				assert.True(t, strings.HasPrefix(err.Error(), linter.ErrorParsingXmlMessageStart), "A syntax error must start with the syntax error prefix")
			}

			if err == nil {
				assert.Equal(t, args.ExpectedEpubInfo, actual)
			}
		})
	}
}
