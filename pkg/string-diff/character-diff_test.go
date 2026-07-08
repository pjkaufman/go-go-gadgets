//go:build unit

package stringdiff_test

import (
	"testing"

	"github.com/charmbracelet/x/ansi"
	stringdiff "github.com/pjkaufman/go-go-gadgets/pkg/string-diff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type characterDiffTestCase struct {
	inputOriginal  string
	inputNew       string
	expectedOutput string
}

var characterDiffTestCases = map[string]characterDiffTestCase{
	"simple character replace should have the expected output": {
		inputOriginal:  `<p class="calibre1">In addition to Tatsuya's group, Erika and Shizuku also came to visit Honoka. </p>`,
		inputNew:       `<p class="calibre1">In addition to Tatsuya's group, Erika, and Shizuku also came to visit Honoka. </p>`,
		expectedOutput: `<p class="calibre1">In addition to Tatsuya's group, Erika, and Shizuku also came to visit Honoka. </p>`,
	},
	"a simple removal should result in the expected output": {
		inputOriginal:  `<p class="calibre1">Having put the cup in the dishwasher (there was no need to press a button, since the <a id="p8"></a>sink started automatically when a sufficient amount of dirty dishes was reached), Tatsuya went to his room while walking past the living room. </p>`,
		inputNew:       `<p class="calibre1">Having put the cup in the dishwasher <i>there was no need to press a button, since the <a id="p8"></a>sink started automatically when a sufficient amount of dirty dishes was reached</i>, Tatsuya went to his room while walking past the living room. </p>`,
		expectedOutput: `<p class="calibre1">Having put the cup in the dishwasher (<i>there was no need to press a button, since the <a id="p8"></a>sink started automatically when a sufficient amount of dirty dishes was reached)</i>, Tatsuya went to his room while walking past the living room. </p>`,
	},
	"… being in the string should get displayed correctly": {
		inputOriginal:  `<p class="calibre1">Original… was here </p>`,
		inputNew:       `<p class="calibre1">Original… was here. </p>`,
		expectedOutput: `<p class="calibre1">Original… was here. </p>`,
	},
	"– being in the string should get displayed correctly": {
		inputOriginal:  `<p class="calibre1">Original– was here </p>`,
		inputNew:       `<p class="calibre1">Original– was here. </p>`,
		expectedOutput: `<p class="calibre1">Original– was here. </p>`,
	},
	"◇ being in the string should get displayed correctly": {
		inputOriginal:  `<p class="calibre1">◇◇◇</p>`,
		inputNew:       `<p class="calibre1">◇◇◇ </p>`,
		expectedOutput: `<p class="calibre1">◇◇◇ </p>`,
	},
}

func TestCharacterDiff(t *testing.T) {
	t.Parallel()

	for name, args := range characterDiffTestCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actual, err := stringdiff.GetPrettyDiffString(args.inputOriginal, args.inputNew)
			require.NoError(t, err)
			assert.Equal(t, args.expectedOutput, ansi.Strip(actual))
		})
	}
}
