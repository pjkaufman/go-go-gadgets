//go:build unit

package rulefixes_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
	"github.com/stretchr/testify/assert"
)

type removeDuplicateManifestEntryTestCase struct {
	opfContents     string
	line, column    int
	expectedChanges []positions.TextEdit
}

var removeDuplicateUniqueIdentifierIdTestCases = map[string]removeDuplicateManifestEntryTestCase{
	"When there is a duplicate entry on the same line as another item, it should be properly removed": {
		opfContents: `<package version="3.0">
<?xml version='1.0' encoding='utf-8'?>
<package xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns="http://www.idpf.org/2007/opf" version="3.0" unique-identifier="BookId" prefix="rendition: http://www.idpf.org/vocab/rendition/#">
  <metadata xmlns:dcterms="http://purl.org/dc/terms/" xmlns:opf="http://www.idpf.org/2007/opf" xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:title id="title1">Some title</dc:title>
    <dc:identifier id="BookId">78436873456487534</dc:identifier>
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
    <item id="section-0030_html" href="Text/section-0030.html" media-type="application/xhtml+xml"/>
    <item id="section-0031_html" href="Text/section-0031.html" media-type="application/xhtml+xml"/>
    <item id="section-0032_html" href="Text/section-0032.html" media-type="application/xhtml+xml"/>
    <item id="section-0033_html" href="Text/section-0033.html" media-type="application/xhtml+xml"/>
    <item id="section-0034_html" href="Text/section-0034.html" media-type="application/xhtml+xml"/>
    <item id="section-0035_html" href="Text/section-0035.html" media-type="application/xhtml+xml"/>
    <item id="section-0036_html" href="Text/section-0036.html" media-type="application/xhtml+xml"/>
    <item id="section-0037_html" href="Text/section-0037.html" media-type="application/xhtml+xml"/>
    <item id="section-0038_html" href="Text/section-0038.html" media-type="application/xhtml+xml"/>
    <item id="section-0039_html" href="Text/section-0039.html" media-type="application/xhtml+xml"/>
    <item id="section-0040_html" href="Text/section-0040.html" media-type="application/xhtml+xml"/>
    <item id="section-0041_html" href="Text/section-0041.html" media-type="application/xhtml+xml"/>
    <item id="section-0042_html" href="Text/section-0042.html" media-type="application/xhtml+xml"/>
    <item id="section-0043_html" href="Text/section-0043.html" media-type="application/xhtml+xml"/>
    <item id="section-0044_html" href="Text/section-0044.html" media-type="application/xhtml+xml"/>
    <item id="section-0045_html" href="Text/section-0045.html" media-type="application/xhtml+xml"/>
    <item id="INTERIORIMAGES_10.jpg" href="Images/INTERIORIMAGES_10.jpg" media-type="image/jpeg"/>
    <item id="INTERIORIMAGES_09.jpg" href="Images/INTERIORIMAGES_09.jpg" media-type="image/jpeg"/>
    <item id="INTERIORIMAGES_08.jpg" href="Images/INTERIORIMAGES_08.jpg" media-type="image/jpeg"/>
    <item id="INTERIORIMAGES_07.jpg" href="Images/INTERIORIMAGES_07.jpg" media-type="image/jpeg"/>
    <item id="INTERIORIMAGES_06.jpg" href="Images/INTERIORIMAGES_06.jpg" media-type="image/jpeg"/>
    <item id="INTERIORIMAGES_05.jpg" href="Images/INTERIORIMAGES_05.jpg" media-type="image/jpeg"/>
    <item id="INTERIORIMAGES_04.jpg" href="Images/INTERIORIMAGES_04.jpg" media-type="image/jpeg"/>
    <item id="INTERIORIMAGES_03.jpg" href="Images/INTERIORIMAGES_03.jpg" media-type="image/jpeg"/>
    <item id="INTERIORIMAGES_02.jpg" href="Images/INTERIORIMAGES_02.jpg" media-type="image/jpeg"/>
    <item id="FRONTMATTER_03.jpg" href="Images/FRONTMATTER_03.jpg" media-type="image/jpeg"/>
    <item id="INTERIORIMAGES_01.jpg" href="Images/INTERIORIMAGES_01.jpg" media-type="image/jpeg"/>
    <item id="FRONTMATTER_01.jpg" href="Images/FRONTMATTER_01.jpg" media-type="image/jpeg"/>
    <item id="FRONTMATTER_02.jpg" href="Images/FRONTMATTER_02.jpg" media-type="image/jpeg"/>
    <item id="CoverDesign.jpg" href="Images/CoverDesign.jpg" media-type="image/jpeg" properties="cover-image"/>
    <item id="COLORGALLERY_02.jpg" href="Images/COLORGALLERY_02.jpg" media-type="image/jpeg"/>
    <item id="COLORGALLERY_01.jpg" href="Images/COLORGALLERY_01.jpg" media-type="image/jpeg"/>
    <item id="COLORGALLERY_03.jpg" href="Images/COLORGALLERY_03.jpg" media-type="image/jpeg"/>
    <item id="navid" href="Text/nav.xhtml" media-type="application/xhtml+xml" properties="nav"/>
  <item id="jnovels_xhtml" href="Text/jnovels.xhtml" media-type="application/xhtml+xml"/><item id="jnov_img_1" href="Images/1.png" media-type="image/png"/><item id="jnov_xhtml_jnovels" href="Text/jnovels.xhtml" media-type="application/xhtml+xml"/></manifest>
  <spine toc="ncx">
    <itemref idref="CoverPage_html"/>
    <itemref idref="toc"/>
    <itemref idref="section-0002_html"/>
    <itemref idref="section-0003_html"/>
    <itemref idref="section-0004_html"/>
    <itemref idref="section-0005_html"/>
    <itemref idref="section-0006_html"/>
<itemref idref="jnovels_xhtml"/><itemref idref="jnov_xhtml_jnovels"/><itemref idref="section-0007_html"/>
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
    <itemref idref="section-0030_html"/>
    <itemref idref="section-0031_html"/>
    <itemref idref="section-0032_html"/>
    <itemref idref="section-0033_html"/>
    <itemref idref="section-0034_html"/>
    <itemref idref="section-0035_html"/>
    <itemref idref="section-0036_html"/>
    <itemref idref="section-0037_html"/>
    <itemref idref="section-0038_html"/>
    <itemref idref="section-0039_html"/>
    <itemref idref="section-0040_html"/>
    <itemref idref="section-0041_html"/>
    <itemref idref="section-0042_html"/>
    <itemref idref="section-0043_html"/>
    <itemref idref="section-0044_html"/>
    <itemref idref="section-0045_html"/>
    </spine>
  <guide>
    <reference type="cover" title="Cover Page" href="Text/CoverPage.html"/>
    <reference type="toc" title="Table of Contents" href="Text/section-0001.html"/>
  </guide>
</package>`,
		line:   75,
		column: 90,
		expectedChanges: []positions.TextEdit{
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   75,
						Column: 3,
					},
					End: positions.Position{
						Line:   75,
						Column: 90,
					},
				},
			},
			{
				Range: positions.Range{
					Start: positions.Position{
						Line:   84,
						Column: 1,
					},
					End: positions.Position{
						Line:   84,
						Column: 33,
					},
				},
			},
		},
	},
}

func TestRemoveDuplicateManifestEntry(t *testing.T) {
	for name, args := range removeDuplicateUniqueIdentifierIdTestCases {
		t.Run(name, func(t *testing.T) {
			actual, err := rulefixes.RemoveDuplicateManifestEntry(args.line, args.column, args.opfContents)

			assert.Nil(t, err)
			assert.Equal(t, args.expectedChanges, actual)
		})
	}
}
