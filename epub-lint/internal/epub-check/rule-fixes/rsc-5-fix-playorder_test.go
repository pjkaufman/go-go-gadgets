//go:build unit

package rulefixes_test

import (
	"testing"

	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
	"github.com/stretchr/testify/assert"
)

type updatePlayOrderTestCase struct {
	input           string
	expectedChanges []rulefixes.TextEdit
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
		expectedChanges: []rulefixes.TextEdit{
			{
				Range: rulefixes.Range{
					Start: rulefixes.Position{
						Line:   7,
						Column: 42,
					},
					End: rulefixes.Position{
						Line:   7,
						Column: 43,
					},
				},
				NewText: "2",
			},
		},
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
		expectedChanges: []rulefixes.TextEdit{
			{
				Range: rulefixes.Range{
					Start: rulefixes.Position{
						Line:   7,
						Column: 30,
					},
					End: rulefixes.Position{
						Line:   7,
						Column: 30,
					},
				},
				NewText: ` playOrder="2"`,
			},
		},
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
		expectedChanges: nil,
	},
}

func TestFixPlayOrder(t *testing.T) {
	for name, tc := range updatePlayOrderTestCases {
		t.Run(name, func(t *testing.T) {
			actual := rulefixes.FixPlayOrder(tc.input)
			assert.Equal(t, tc.expectedChanges, actual)
		})
	}
}
