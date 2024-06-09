//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

type AddCssPageBreakIfMissingTestCase struct {
	Input          string
	ExpectedOutput string
}

var cssFileWithHrBlankSpace = `p {
height: 10px,
}
hr.blankSpace {
border:0;
height:2em;
}`

var AddCssPageBreakIfMissingTestCases = map[string]AddCssPageBreakIfMissingTestCase{
	"make sure that an empty input becomes the hr blank space": {
		Input:          "",
		ExpectedOutput: linter.HrBlankSpace + "\n",
	},
	"make sure that a solely whitespace input becomes the hr blank space": {
		Input: `
				   `,
		ExpectedOutput: linter.HrBlankSpace + "\n",
	},
	"make sure that input that already contains blank space does not get it added": {
		Input:          cssFileWithHrBlankSpace,
		ExpectedOutput: cssFileWithHrBlankSpace,
	},
	"make sure that input that does not contain blank space does get it added": {
		Input: `p {
height: 10px,
}`,
		ExpectedOutput: cssFileWithHrBlankSpace,
	},
}

func TestAddCssPageBreakIfMissing(t *testing.T) {
	for name, args := range AddCssPageBreakIfMissingTestCases {
		t.Run(name, func(t *testing.T) {
			actual := linter.AddCssPageBreakIfMissing(args.Input)

			assert.Equal(t, args.ExpectedOutput, actual)
		})
	}
}
