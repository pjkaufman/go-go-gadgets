//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

type ParseTextReplacementsTestCase struct {
	Input    string
	Expected map[string]string
}

var ParseTextReplacementsTestCases = map[string]ParseTextReplacementsTestCase{
	"make sure that an empty table results in an empty map": {
		Input: `| Text to replace | Text replacement |
		| ---- | ---- |`,
		Expected: map[string]string{},
	},
	"make sure that a non-empty table results in the appropriate amount of entries being placed in a map": {
		Input: `| Text to replace | Text replacement |
		| ---- | ---- |
		| replace | with me |
		| "I am quoted" | 'I am single quoted' |`,
		Expected: map[string]string{
			"replace":         "with me",
			"\"I am quoted\"": "'I am single quoted'",
		},
	},
	"make sure that values get trimmed before getting added to the map": {
		Input: `| Text to replace | Text replacement |
		| ---- | ---- |
		| replace | with me |
		| "I am quoted" | 'I am single quoted' |
		|       I have lots of whitespace around me      | I have   wonky internal spacing |`,
		Expected: map[string]string{
			"replace":                             "with me",
			"\"I am quoted\"":                     "'I am single quoted'",
			"I have lots of whitespace around me": "I have   wonky internal spacing",
		},
	},
	"make sure that lines without a pipe/table row get ignored": {
		Input: `| Text to replace | Text replacement |
		| ---- | ---- |
		| replace | with me |
		| "I am quoted" | 'I am single quoted' |
		|       I have lots of whitespace around me      | I have   wonky internal spacing |
		Some text here
		Another line here`,
		Expected: map[string]string{
			"replace":                             "with me",
			"\"I am quoted\"":                     "'I am single quoted'",
			"I have lots of whitespace around me": "I have   wonky internal spacing",
		},
	},
}

func TestParseTextReplacements(t *testing.T) {
	for name, args := range ParseTextReplacementsTestCases {
		t.Run(name, func(t *testing.T) {
			actual, err := linter.ParseTextReplacements(args.Input)

			assert.Nil(t, err)
			assert.Equal(t, args.Expected, actual)
		})
	}
}
