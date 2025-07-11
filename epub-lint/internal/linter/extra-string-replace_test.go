//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

type extraStringReplaceTestCase struct {
	inputText             string
	inputFindsAndReplaces map[string]string
	inputHits             map[string]int
	expectedText          string
	expectedHits          map[string]int
}

var extraStringReplaceTestCases = map[string]extraStringReplaceTestCase{
	"make sure that when a replacement is made with an empty map, the number of hits in the map is updated accordingly": {
		inputText: `Here is some text that gets broken into
		multiple lines with a couple of words to be replaced`,
		inputFindsAndReplaces: map[string]string{
			"Here is":        "This was",
			"to be replaced": "that were replaced",
		},
		inputHits: map[string]int{},
		expectedHits: map[string]int{
			"Here is":        1,
			"to be replaced": 1,
		},
		expectedText: `This was some text that gets broken into
		multiple lines with a couple of words that were replaced`,
	},
	"make sure that when multiple instances of a value to replace in a string are present that all of them get replaced": {
		inputText: `I talk way too much as if I were not going to get another chance to talk to myself. I wonder why that is.`,
		inputFindsAndReplaces: map[string]string{
			"I": "You",
		},
		inputHits: map[string]int{},
		expectedHits: map[string]int{
			"I": 3,
		},
		expectedText: `You talk way too much as if You were not going to get another chance to talk to myself. You wonder why that is.`,
	},
	"make sure that not finding a value in a file when it does not already exist just sets the value for that search value to 0": {
		inputText: `Text not found`,
		inputFindsAndReplaces: map[string]string{
			"I": "You",
		},
		inputHits: map[string]int{},
		expectedHits: map[string]int{
			"I": 0,
		},
		expectedText: `Text not found`,
	},
	"make sure that not finding a value in a file when it does not already exist does not affect the resulting hit count": {
		inputText: `Text not found`,
		inputFindsAndReplaces: map[string]string{
			"I": "You",
		},
		inputHits: map[string]int{
			"I": 5,
		},
		expectedHits: map[string]int{
			"I": 5,
		},
		expectedText: `Text not found`,
	},
	"make sure that when a replacement is made and the value already exists in the hit count it gets incremented": {
		inputText: `This is not what I expected. This could get dangerous. This is not what I signed up for!`,
		inputFindsAndReplaces: map[string]string{
			"This": "That",
		},
		inputHits: map[string]int{
			"This": 2,
		},
		expectedHits: map[string]int{
			"This": 5,
		},
		expectedText: `That is not what I expected. That could get dangerous. That is not what I signed up for!`,
	},
}

func TestExtraStringReplace(t *testing.T) {
	for name, args := range extraStringReplaceTestCases {
		t.Run(name, func(t *testing.T) {
			actual := linter.ExtraStringReplace(args.inputText, args.inputFindsAndReplaces, args.inputHits)

			assert.Equal(t, args.expectedText, actual, "output text doesn't match")
			assert.Equal(t, args.expectedHits, args.inputHits, "output map doesn't match")
		})
	}
}
