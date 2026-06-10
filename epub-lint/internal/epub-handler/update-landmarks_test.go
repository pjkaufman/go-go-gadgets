//go:build unit

package epubhandler_test

import (
	"testing"

	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
	"github.com/stretchr/testify/assert"
)

type updateLandmarksTestCase struct {
	navText               string
	relativePathToReplace string
	relativePathToCover   string
	relativePathToToc     string
	expectedOutput        string
}

var updateLandmarksTestCases = map[string]updateLandmarksTestCase{
	"When the nav file has no landmarks, no change is made to the nav file": {
		navText: `<?xml version="1.0" encoding="utf-8"?>
<html xmlns="http://www.w3.org/1999/xhtml">
  <head><title>Test</title></head>
  <body>
    <nav id="somethingElse">
      <h2>Guide</h2>
      <ol>
        <li><a href="../Images/1.png">Cover Page</a></li>
      </ol>
    </nav>
  </body>
</html>`,
		relativePathToReplace: "../Images/1.png",
		relativePathToCover:   "../Text/cover.xhtml",
		relativePathToToc:     "../Text/nav.xhtml",
		expectedOutput: `<?xml version="1.0" encoding="utf-8"?>
<html xmlns="http://www.w3.org/1999/xhtml">
  <head><title>Test</title></head>
  <body>
    <nav id="somethingElse">
      <h2>Guide</h2>
      <ol>
        <li><a href="../Images/1.png">Cover Page</a></li>
      </ol>
    </nav>
  </body>
</html>`,
	},
	"When the nav file has landmarks, but no reference to the file to replace, no change is made": {
		navText: `<?xml version="1.0" encoding="utf-8"?>
<html xmlns="http://www.w3.org/1999/xhtml">
  <head><title>Test</title></head>
  <body>
    <nav epub:type="landmarks" id="landmarks" hidden="">
      <h2>Guide</h2>
      <ol>
        <li><a epub:type="cover" href="../Text/cover.xhtml">Cover Page</a></li>
        <li><a epub:type="toc" href="../Text/nav.xhtml">Table of Contents</a></li>
      </ol>
    </nav>
  </body>
</html>`,
		relativePathToReplace: "../Images/1.png",
		relativePathToCover:   "../Text/cover.xhtml",
		relativePathToToc:     "../Text/nav.xhtml",
		expectedOutput: `<?xml version="1.0" encoding="utf-8"?>
<html xmlns="http://www.w3.org/1999/xhtml">
  <head><title>Test</title></head>
  <body>
    <nav epub:type="landmarks" id="landmarks" hidden="">
      <h2>Guide</h2>
      <ol>
        <li><a epub:type="cover" href="../Text/cover.xhtml">Cover Page</a></li>
        <li><a epub:type="toc" href="../Text/nav.xhtml">Table of Contents</a></li>
      </ol>
    </nav>
  </body>
</html>`,
	},
	"When the nav file has landmarks and a cover reference to the file to replace, the cover href is swapped out for the relative path to the cover file": {
		navText: `<?xml version="1.0" encoding="utf-8"?>
<html xmlns="http://www.w3.org/1999/xhtml">
  <head><title>Test</title></head>
  <body>
    <nav epub:type="landmarks" id="landmarks" hidden="">
      <h2>Guide</h2>
      <ol>
        <li><a epub:type="cover" href="../Images/1.png">Cover Page</a></li>
      </ol>
    </nav>
  </body>
</html>`,
		relativePathToReplace: "../Images/1.png",
		relativePathToCover:   "../Text/cover.xhtml",
		relativePathToToc:     "../Text/nav.xhtml",
		expectedOutput: `<?xml version="1.0" encoding="utf-8"?>
<html xmlns="http://www.w3.org/1999/xhtml">
  <head><title>Test</title></head>
  <body>
    <nav epub:type="landmarks" id="landmarks" hidden="">
      <h2>Guide</h2>
      <ol>
        <li><a epub:type="cover" href="../Text/cover.xhtml">Cover Page</a></li>
      </ol>
    </nav>
  </body>
</html>`,
	},
	"When the nav file has landmarks and a toc reference to the file to replace, the toc href is swapped out for the relative path to the toc file": {
		navText: `<?xml version="1.0" encoding="utf-8"?>
<html xmlns="http://www.w3.org/1999/xhtml">
  <head><title>Test</title></head>
  <body>
    <nav epub:type="landmarks" id="landmarks" hidden="">
      <h2>Guide</h2>
      <ol>
        <li><a epub:type="toc" href="../Images/1.png">Table of Contents</a></li>
      </ol>
    </nav>
  </body>
</html>`,
		relativePathToReplace: "../Images/1.png",
		relativePathToCover:   "../Text/cover.xhtml",
		relativePathToToc:     "../Text/nav.xhtml",
		expectedOutput: `<?xml version="1.0" encoding="utf-8"?>
<html xmlns="http://www.w3.org/1999/xhtml">
  <head><title>Test</title></head>
  <body>
    <nav epub:type="landmarks" id="landmarks" hidden="">
      <h2>Guide</h2>
      <ol>
        <li><a epub:type="toc" href="../Text/nav.xhtml">Table of Contents</a></li>
      </ol>
    </nav>
  </body>
</html>`,
	},
	"When the nav file has landmarks and both a cover and toc reference to the file to replace, both the cover and toc hrefs are updated accordingly": {
		navText: `<?xml version="1.0" encoding="utf-8"?>
<html xmlns="http://www.w3.org/1999/xhtml">
  <head><title>Test</title></head>
  <body>
    <nav epub:type="landmarks" id="landmarks" hidden="">
      <h2>Guide</h2>
      <ol>
        <li><a epub:type="cover" href="../Images/1.png">Cover Page</a></li>
        <li><a epub:type="toc" href="../Images/1.png">Table of Contents</a></li>
      </ol>
    </nav>
  </body>
</html>`,
		relativePathToReplace: "../Images/1.png",
		relativePathToCover:   "../Text/cover.xhtml",
		relativePathToToc:     "../Text/nav.xhtml",
		expectedOutput: `<?xml version="1.0" encoding="utf-8"?>
<html xmlns="http://www.w3.org/1999/xhtml">
  <head><title>Test</title></head>
  <body>
    <nav epub:type="landmarks" id="landmarks" hidden="">
      <h2>Guide</h2>
      <ol>
        <li><a epub:type="cover" href="../Text/cover.xhtml">Cover Page</a></li>
        <li><a epub:type="toc" href="../Text/nav.xhtml">Table of Contents</a></li>
      </ol>
    </nav>
  </body>
</html>`,
	},
}

func TestUpdateLandmarks(t *testing.T) {
	t.Parallel()

	for name, args := range updateLandmarksTestCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actual := epubhandler.UpdateLandmarks(args.navText, args.relativePathToReplace, args.relativePathToCover, args.relativePathToToc)

			assert.Equal(t, args.expectedOutput, actual)
		})
	}
}
