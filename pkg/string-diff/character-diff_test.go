//go:build unit

package stringdiff_test

import (
	"testing"

	stringdiff "github.com/pjkaufman/go-go-gadgets/pkg/string-diff"
	"github.com/stretchr/testify/assert"
)

type CharacterDiffTestCase struct {
	InputOriginal  string
	InputNew       string
	ExpectedOutput string
}

var CharacterDiffTestCases = map[string]CharacterDiffTestCase{
	"simple character replace should have the expected output": {
		InputOriginal:  `<p class="calibre1">In addition to Tatsuya's group, Erika and Shizuku also came to visit Honoka. </p>`,
		InputNew:       `<p class="calibre1">In addition to Tatsuya's group, Erika, and Shizuku also came to visit Honoka. </p>`,
		ExpectedOutput: `<p class="calibre1">In addition to Tatsuya's group, Erika, and Shizuku also came to visit Honoka. </p>`,
	},
	"a simple removal should result in the expected output": {
		InputOriginal:  `<p class="calibre1">Having put the cup in the dishwasher (there was no need to press a button, since the <a id="p8"></a>sink started automatically when a sufficient amount of dirty dishes was reached), Tatsuya went to his room while walking past the living room. </p>`,
		InputNew:       `<p class="calibre1">Having put the cup in the dishwasher <i>there was no need to press a button, since the <a id="p8"></a>sink started automatically when a sufficient amount of dirty dishes was reached</i>, Tatsuya went to his room while walking past the living room. </p>`,
		ExpectedOutput: `<p class="calibre1">Having put the cup in the dishwasher (<i>there was no need to press a button, since the <a id="p8"></a>sink started automatically when a sufficient amount of dirty dishes was reached)</i>, Tatsuya went to his room while walking past the living room. </p>`,
	},
	"… being in the string should get displayed correctly": {
		InputOriginal:  `<p class="calibre1">Original… was here </p>`,
		InputNew:       `<p class="calibre1">Original… was here. </p>`,
		ExpectedOutput: `<p class="calibre1">Original… was here. </p>`,
	},
	"– being in the string should get displayed correctly": {
		InputOriginal:  `<p class="calibre1">Original– was here </p>`,
		InputNew:       `<p class="calibre1">Original– was here. </p>`,
		ExpectedOutput: `<p class="calibre1">Original– was here. </p>`,
	},
	"◇ being in the string should get displayed correctly": {
		InputOriginal:  `<p class="calibre1">◇◇◇</p>`,
		InputNew:       `<p class="calibre1">◇◇◇ </p>`,
		ExpectedOutput: `<p class="calibre1">◇◇◇ </p>`,
	},
}

func TestCharacterDiff(t *testing.T) {
	for name, args := range CharacterDiffTestCases {
		t.Run(name, func(t *testing.T) {

			actual, err := stringdiff.GetPrettyDiffString(args.InputOriginal, args.InputNew, true)
			assert.Nil(t, err)
			assert.Equal(t, args.ExpectedOutput, actual)
		})
	}
}
