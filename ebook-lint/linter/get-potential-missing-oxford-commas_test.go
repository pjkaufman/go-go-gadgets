//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/linter"
	"github.com/stretchr/testify/assert"
)

type GetPotentialMissingOxfordCommasTestCase struct {
	InputText           string
	ExpectedSuggestions map[string]string
}

var GetPotentialMissingOxfordCommasTestCases = map[string]GetPotentialMissingOxfordCommasTestCase{
	"make sure that a file with no missing and's or or's without a comma proceeding it gets no suggestions": {
		InputText: `<p>Here is some content.</p>
<p>Here is some more content</p>`,
		ExpectedSuggestions: map[string]string{},
	},
	"make sure that a file with a missing comma before an and gets a suggestion": {
		InputText: `<p>Here is some content.</p>
		<p>Here is a situation where I run, skip and jump for a long time.</p>`,
		ExpectedSuggestions: map[string]string{
			`
		<p>Here is a situation where I run, skip and jump for a long time.</p>`: `
		<p>Here is a situation where I run, skip, and jump for a long time.</p>`,
		},
	},
	"make sure that a file with a missing comma before an or gets a suggestion": {
		InputText: `<p>Here is some content.</p>
		<p>Here is a situation where I run, skip or jump for a long time.</p>`,
		ExpectedSuggestions: map[string]string{
			`
		<p>Here is a situation where I run, skip or jump for a long time.</p>`: `
		<p>Here is a situation where I run, skip, or jump for a long time.</p>`,
		},
	},
}

func TestGetPotentialMissingOxfordCommas(t *testing.T) {
	for name, args := range GetPotentialMissingOxfordCommasTestCases {
		t.Run(name, func(t *testing.T) {
			actual := linter.GetPotentialMissingOxfordCommas(args.InputText)

			assert.Equal(t, args.ExpectedSuggestions, actual)
		})
	}
}
