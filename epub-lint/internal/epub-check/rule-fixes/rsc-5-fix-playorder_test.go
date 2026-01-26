//go:build unit

package rulefixes_test

import (
	"testing"

	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
)

type updatePlayOrderTestCase struct {
	input          string
	expectedOutput string
}

var updatePlayOrderTestCases = map[string]updatePlayOrderTestCase{
	"Updating the play order works when there are duplicate playOrder values": {
		input: `<ncx>
  <navMap>
    <navPoint id="navPoint-1" playOrder="1">
      <navLabel><text>Chapter 1</text></navLabel>
      <content src="chapter1.html" />
    </navPoint>
    <navPoint id="navPoint-2" playOrder="1">
      <navLabel><text>Chapter 2</text></navLabel>
      <content src="chapter2.html" />
    </navPoint>
  </navMap>
</ncx>`,
		expectedOutput: `<ncx>
  <navMap>
    <navPoint id="navPoint-1" playOrder="1">
      <navLabel><text>Chapter 1</text></navLabel>
      <content src="chapter1.html" />
    </navPoint>
    <navPoint id="navPoint-2" playOrder="2">
      <navLabel><text>Chapter 2</text></navLabel>
      <content src="chapter2.html" />
    </navPoint>
  </navMap>
</ncx>`,
	},
	"Updating the play order works when there is a missing playOrder attribute": {
		input: `<ncx>
  <navMap>
    <navPoint id="navPoint-1" playOrder="1">
      <navLabel><text>Chapter 1</text></navLabel>
      <content src="chapter1.html" />
    </navPoint>
    <navPoint id="navPoint-2">
      <navLabel><text>Chapter 2</text></navLabel>
      <content src="chapter2.html" />
    </navPoint>
  </navMap>
</ncx>`,
		expectedOutput: `<ncx>
  <navMap>
    <navPoint id="navPoint-1" playOrder="1">
      <navLabel><text>Chapter 1</text></navLabel>
      <content src="chapter1.html" />
    </navPoint>
    <navPoint id="navPoint-2" playOrder="2">
      <navLabel><text>Chapter 2</text></navLabel>
      <content src="chapter2.html" />
    </navPoint>
  </navMap>
</ncx>`,
	},
	"Updating the play order does nothing if all playOrders are in order": {
		input: `<ncx>
  <navMap>
    <navPoint id="navPoint-1" playOrder="1">
      <navLabel><text>Chapter 1</text></navLabel>
      <content src="chapter1.html" />
    </navPoint>
    <navPoint id="navPoint-2" playOrder="2">
      <navLabel><text>Chapter 2</text></navLabel>
      <content src="chapter2.html" />
    </navPoint>
  </navMap>
</ncx>`,
		expectedOutput: `<ncx>
  <navMap>
    <navPoint id="navPoint-1" playOrder="1">
      <navLabel><text>Chapter 1</text></navLabel>
      <content src="chapter1.html" />
    </navPoint>
    <navPoint id="navPoint-2" playOrder="2">
      <navLabel><text>Chapter 2</text></navLabel>
      <content src="chapter2.html" />
    </navPoint>
  </navMap>
</ncx>`,
	},
}

func TestFixPlayOrder(t *testing.T) {
	for name, tc := range updatePlayOrderTestCases {
		t.Run(name, func(t *testing.T) {
			edits := rulefixes.FixPlayOrder(tc.input)

			checkFinalOutputMatches(t, tc.input, tc.expectedOutput, edits...)
		})
	}
}
