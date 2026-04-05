//go:build unit

package jnovels_test

import (
	"fmt"
	"testing"

	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/jnovels"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	opfBase = `<?xml version="1.0"?>
<package>
  <manifest>
    <item id="jnovels" href="Text/jnovels.xhtml" />
    <item id="img1" href="Images/1.png" />
    <item id="chap1" href="Text/ch1.xhtml" />
  </manifest>
  <spine><itemref idref="jnovels" />
    <itemref idref="chap1" />
  </spine>
</package>
`

	opfNoJNovels = `<?xml version="1.0"?>
<package>
  <manifest>
    <item id="img1" href="Images/1.png" />
    <item id="chap1" href="Text/ch1.xhtml" />
  </manifest>
  <spine>
    <itemref idref="chap1" />
  </spine>
</package>
`

	opfNoImage = `<?xml version="1.0"?>
<package>
  <manifest>
    <item id="jnovels" href="Text/jnovels.xhtml" />
    <item id="chap1" href="Text/ch1.xhtml" />
  </manifest>
  <spine><itemref idref="jnovels" />
    <itemref idref="chap1" />
  </spine>
</package>
`

	opfOnlyChapter = `<?xml version="1.0"?>
<package>
  <manifest>
    <item id="chap1" href="Text/ch1.xhtml" />
  </manifest>
  <spine>
    <itemref idref="chap1" />
  </spine>
</package>
`

	ncxWithJNovels = `<ncx>
<navMap>
  <navPoint id="navPoint-1" playOrder="1">
    <content src="Text/jnovels.xhtml"/>
  </navPoint>
  <navPoint id="navPoint-2" playOrder="2">
    <content src="Text/ch1.xhtml"/>
  </navPoint>
</navMap>
</ncx>
`

	ncxWithoutJNovels = `<ncx>
<navMap>
  <navPoint id="navPoint-2" playOrder="1">
    <content src="Text/ch1.xhtml"/>
  </navPoint>
</navMap>
</ncx>
`

	navTocWithJNovels = `<nav epub:type="toc">
  <ol>
    <li><a href="Text/cover.xhtml">Cover</a></li>
    <li><a href="Text/jnovels.xhtml">JNovels</a></li>
    <li><a href="Text/ch1.xhtml">Chapter 1</a></li>
  </ol>
</nav>
`

	navTocWithoutJNovels = `<nav epub:type="toc">
  <ol>
    <li><a href="Text/cover.xhtml">Cover</a></li>
    <li><a href="Text/ch1.xhtml">Chapter 1</a></li>
  </ol>
</nav>
`

	navLandmarksWithImage = `<nav epub:type="landmarks">
  <ol>
    <li><a epub:type="cover" href="Images/1.png">Cover</a></li>
    <li><a epub:type="toc" href="Images/1.png">TOC</a></li>
  </ol>
</nav>
`

	navLandmarksCoverUpdated = `<nav epub:type="landmarks">
  <ol>
    <li><a epub:type="cover" href="Text/cover.xhtml">Cover</a></li>
    <li><a epub:type="toc" href="Images/1.png">TOC</a></li>
  </ol>
</nav>
`

	navLandmarksTocUpdated = `<nav epub:type="landmarks">
  <ol>
    <li><a epub:type="cover" href="Images/1.png">Cover</a></li>
    <li><a epub:type="toc" href="Text/toc.xhtml">TOC</a></li>
  </ol>
</nav>
`

	navLandmarksCoverAndTocUpdated = `<nav epub:type="landmarks">
  <ol>
    <li><a epub:type="cover" href="Text/cover.xhtml">Cover</a></li>
    <li><a epub:type="toc" href="Text/toc.xhtml">TOC</a></li>
  </ol>
</nav>
`
)

type cleanupJNovelsFilesTestCase struct {
	input                jnovels.JNovelsCleanupContext
	expectedFileContent  map[string]string
	expectedHandledFiles []string
}

func makeCtx(
	updated map[string]string,
	fileBasenameMap map[string][]string,
	navFile, tocFile, coverFile string,
) jnovels.JNovelsCleanupContext {
	opfFolder := "OEBPS"
	return jnovels.JNovelsCleanupContext{
		EpubInfo: epubhandler.EpubInfo{
			NavFile:   navFile,
			TocFile:   tocFile,
			CoverFile: coverFile,
		},
		OpfFolder:           opfFolder,
		OpfFileName:         filehandler.JoinPath(opfFolder, "content.opf"),
		NcxFileName:         filehandler.JoinPath(opfFolder, "toc.ncx"),
		FileBasenameMap:     fileBasenameMap,
		UpdatedFileContents: updated,
		GetFileContents: func(path string) (string, error) {
			if c, ok := updated[path]; ok {
				return c, nil
			}
			return "", fmt.Errorf("file not found: %s", path)
		},
	}
}

var cleanupJNovelsFilesTestCases = map[string]cleanupJNovelsFilesTestCase{
	"When no JNovels files are present, no changes are made and no files are handled": {
		input: makeCtx(
			map[string]string{
				"OEBPS/content.opf": opfBase,
				"OEBPS/toc.ncx":     ncxWithJNovels,
			},
			map[string][]string{},
			"", "", "",
		),
		expectedFileContent: map[string]string{
			"OEBPS/content.opf": opfBase,
			"OEBPS/toc.ncx":     ncxWithJNovels,
		},
		expectedHandledFiles: nil,
	},
	"When there is a JNovels html file present, it should be handled and removed from the OPF file": {
		input: makeCtx(
			map[string]string{
				"OEBPS/content.opf": opfBase,
				"OEBPS/toc.ncx":     ncxWithJNovels,
			},
			map[string][]string{
				jnovels.JnovelsFile: {"OEBPS/Text/jnovels.xhtml"},
			},
			"", "", "",
		),
		expectedFileContent: map[string]string{
			"OEBPS/content.opf": opfNoJNovels,
			"OEBPS/toc.ncx":     ncxWithoutJNovels,
		},
		expectedHandledFiles: []string{"OEBPS/Text/jnovels.xhtml"},
	},
	"When there is a JNovels image file present, it should be handled and removed from the OPF file": {
		input: makeCtx(
			map[string]string{
				"OEBPS/content.opf": opfBase,
				"OEBPS/toc.ncx":     ncxWithJNovels,
			},
			map[string][]string{
				jnovels.JnovelsImage: {"OEBPS/Images/1.png"},
			},
			"", "", "",
		),
		expectedFileContent: map[string]string{
			"OEBPS/content.opf": opfNoImage,
			"OEBPS/toc.ncx":     ncxWithJNovels,
		},
		expectedHandledFiles: []string{"OEBPS/Images/1.png"},
	},
	"When there are the JNovels html and image files present, they should be handled and removed from the OPF file": {
		input: makeCtx(
			map[string]string{
				"OEBPS/content.opf": opfBase,
				"OEBPS/toc.ncx":     ncxWithJNovels,
			},
			map[string][]string{
				jnovels.JnovelsFile:  {"OEBPS/Text/jnovels.xhtml"},
				jnovels.JnovelsImage: {"OEBPS/Images/1.png"},
			},
			"", "", "",
		),
		expectedFileContent: map[string]string{
			"OEBPS/content.opf": opfOnlyChapter,
			"OEBPS/toc.ncx":     ncxWithoutJNovels,
		},
		expectedHandledFiles: []string{"OEBPS/Text/jnovels.xhtml", "OEBPS/Images/1.png"},
	},
	"When there is a JNovels html file present and it is referenced from the nav file, it should be handled and removed from the OPF and nav files": {
		input: makeCtx(
			map[string]string{
				"OEBPS/content.opf": opfBase,
				"OEBPS/toc.ncx":     ncxWithJNovels,
				"OEBPS/nav.xhtml":   navTocWithJNovels,
			},
			map[string][]string{
				jnovels.JnovelsFile: {"OEBPS/Text/jnovels.xhtml"},
			},
			"nav.xhtml", "", "",
		),
		expectedFileContent: map[string]string{
			"OEBPS/content.opf": opfNoJNovels,
			"OEBPS/toc.ncx":     ncxWithoutJNovels,
			"OEBPS/nav.xhtml":   navTocWithoutJNovels,
		},
		expectedHandledFiles: []string{"OEBPS/Text/jnovels.xhtml"},
	},
	"When there is a JNovels html file present and it is referenced from the nav and NCX files, it should be handled and removed from the OPF, NCX, and nav files": {
		input: makeCtx(
			map[string]string{
				"OEBPS/content.opf": opfBase,
				"OEBPS/toc.ncx":     ncxWithJNovels,
				"OEBPS/nav.xhtml":   navTocWithJNovels,
				"OEBPS/toc.xhtml":   navTocWithJNovels,
			},
			map[string][]string{
				jnovels.JnovelsFile: {"OEBPS/Text/jnovels.xhtml"},
			},
			"nav.xhtml", "toc.xhtml", "",
		),
		expectedFileContent: map[string]string{
			"OEBPS/content.opf": opfNoJNovels,
			"OEBPS/toc.ncx":     ncxWithoutJNovels,
			"OEBPS/nav.xhtml":   navTocWithoutJNovels,
			"OEBPS/toc.xhtml":   navTocWithoutJNovels,
		},
		expectedHandledFiles: []string{"OEBPS/Text/jnovels.xhtml"},
	},
	"When there is a JNovels image file present and it is referenced as a part of the landmarks as the cover and there is a cover file, it should be handled, removed from the OPF file, and update to the cover file in the nav file": {
		input: makeCtx(
			map[string]string{
				"OEBPS/content.opf": opfBase,
				"OEBPS/toc.ncx":     ncxWithJNovels,
				"OEBPS/nav.xhtml":   navLandmarksWithImage,
			},
			map[string][]string{
				jnovels.JnovelsImage: {"OEBPS/Images/1.png"},
			},
			"nav.xhtml", "", "Text/cover.xhtml",
		),
		expectedFileContent: map[string]string{
			"OEBPS/content.opf": opfNoImage,
			"OEBPS/toc.ncx":     ncxWithJNovels,
			"OEBPS/nav.xhtml":   navLandmarksCoverUpdated,
		},
		expectedHandledFiles: []string{"OEBPS/Images/1.png"},
	},
	"When there is a JNovels image file present and it is referenced as a part of the landmarks as the toc and there is a toc file, it should be handled, removed from the OPF file, and update to the toc file in the nav file": {
		input: makeCtx(
			map[string]string{
				"OEBPS/content.opf": opfBase,
				"OEBPS/toc.ncx":     ncxWithJNovels,
				"OEBPS/nav.xhtml":   navLandmarksWithImage,
			},
			map[string][]string{
				jnovels.JnovelsImage: {"OEBPS/Images/1.png"},
			},
			"nav.xhtml", "Text/toc.xhtml", "",
		),
		expectedFileContent: map[string]string{
			"OEBPS/content.opf": opfNoImage,
			"OEBPS/toc.ncx":     ncxWithJNovels,
			"OEBPS/nav.xhtml":   navLandmarksTocUpdated,
		},
		expectedHandledFiles: []string{"OEBPS/Images/1.png"},
	},
	"When there is a JNovels image file present and it is referenced as a part of the landmarks as the cover and toc and there is a cover and toc file, it should be handled, removed from the OPF file, and update to the cover and toc files in the nav file": {
		input: makeCtx(
			map[string]string{
				"OEBPS/content.opf": opfBase,
				"OEBPS/toc.ncx":     ncxWithJNovels,
				"OEBPS/nav.xhtml":   navLandmarksWithImage,
			},
			map[string][]string{
				jnovels.JnovelsImage: {"OEBPS/Images/1.png"},
			},
			"nav.xhtml", "Text/toc.xhtml", "Text/cover.xhtml",
		),
		expectedFileContent: map[string]string{
			"OEBPS/content.opf": opfNoImage,
			"OEBPS/toc.ncx":     ncxWithJNovels,
			"OEBPS/nav.xhtml":   navLandmarksCoverAndTocUpdated,
		},
		expectedHandledFiles: []string{"OEBPS/Images/1.png"},
	},
}

func TestCleanupJNovelsFiles(t *testing.T) {
	for name, tc := range cleanupJNovelsFilesTestCases {
		t.Run(name, func(t *testing.T) {
			handledFiles, err := jnovels.CleanupJNovelsFiles(tc.input)

			require.NoError(t, err)
			assert.Equal(t, tc.expectedHandledFiles, handledFiles)
			for name, fileContent := range tc.expectedFileContent {
				updatedContents, found := tc.input.UpdatedFileContents[name]
				if !found {
					assert.Fail(t, fmt.Sprintf("expected %q to be updated, but it was not", name))
				} else {
					assert.Equal(t, fileContent, updatedContents, "%q does not match expected content")
				}
			}
		})
	}
}
