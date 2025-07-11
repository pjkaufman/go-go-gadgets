//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

type parseTextReplacementsTestCase struct {
	input    string
	expected map[string]string
}

var parseTextReplacementsTestCases = map[string]parseTextReplacementsTestCase{
	"make sure that an empty table results in an empty map": {
		input: `| Text to replace | Text replacement |
		| ---- | ---- |`,
		expected: map[string]string{},
	},
	"make sure that a non-empty table results in the appropriate amount of entries being placed in a map": {
		input: `| Text to replace | Text replacement |
		| ---- | ---- |
		| replace | with me |
		| "I am quoted" | 'I am single quoted' |`,
		expected: map[string]string{
			"replace":         "with me",
			"\"I am quoted\"": "'I am single quoted'",
		},
	},
	"make sure that values get trimmed before getting added to the map": {
		input: `| Text to replace | Text replacement |
		| ---- | ---- |
		| replace | with me |
		| "I am quoted" | 'I am single quoted' |
		|       I have lots of whitespace around me      | I have   wonky internal spacing |`,
		expected: map[string]string{
			"replace":                             "with me",
			"\"I am quoted\"":                     "'I am single quoted'",
			"I have lots of whitespace around me": "I have   wonky internal spacing",
		},
	},
	"make sure that lines without a pipe/table row get ignored": {
		input: `| Text to replace | Text replacement |
		| ---- | ---- |
		| replace | with me |
		| "I am quoted" | 'I am single quoted' |
		|       I have lots of whitespace around me      | I have   wonky internal spacing |
		Some text here
		Another line here`,
		expected: map[string]string{
			"replace":                             "with me",
			"\"I am quoted\"":                     "'I am single quoted'",
			"I have lots of whitespace around me": "I have   wonky internal spacing",
		},
	},
}

func TestParseTextReplacements(t *testing.T) {
	for name, args := range parseTextReplacementsTestCases {
		t.Run(name, func(t *testing.T) {
			actual, err := linter.ParseTextReplacements(args.input)

			assert.Nil(t, err)
			assert.Equal(t, args.expected, actual)
		})
	}
}
